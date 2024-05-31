package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	_ "github.com/go-park-mail-ru/2024_1_Cyberkotletki/docs"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/profanity"
	delivery "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/redis"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/connector"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/prometheus/client_golang/prometheus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"time"
)

// @title API Киноскопа
// @version 1.0
// @Description сервис Киноскоп (аналог кинопоиска)
// @securityDefinitions.apikey _csrf
// @in header
// @name x-csrf
func main() {
	genCfg := flag.Bool("generate-example-config", false, "Генерирует пример конфига, с которым умеет работать сервер")
	flag.Parse()
	if *genCfg {
		config.GenerateExampleConfigs()
		return
	}

	logger := log.New("server: ")
	coreParams := config.ParseCoreServiceParams()
	authParams := config.ParseAuthServiceParams()
	staticParams := config.ParseStaticServiceParams()
	logger.Printf("Параметры запуска сервера: %v \n", coreParams)
	logger.Printf("Параметры сервиса авторизации: %v \n", authParams)
	logger.Printf("Параметры сервиса статики: %v \n", staticParams)
	logger.Printf(staticParams.S3.SecretAccessKey)

	echoServer := Init(logger, coreParams, authParams, staticParams)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	go Run(echoServer, coreParams)

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(coreParams.HTTP.Server.GracefulShutdownTimeout)*time.Second,
	)
	defer cancel()
	Shutdown(ctx, echoServer)
}

func Init(
	logger echo.Logger,
	coreParams config.Config,
	authParams config.AuthConfig,
	staticParams config.StaticConfig,
) *echo.Echo {
	// DBConn
	psqlConn, err := connector.GetPostgresConnector(coreParams.Postgres.GetConnectURL())
	if err != nil {
		logger.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	s3conn, err := connector.GetS3Connector(
		staticParams.S3.Endpoint, staticParams.S3.Region, staticParams.S3.AccessKeyID, staticParams.S3.SecretAccessKey,
	)
	if err != nil {
		logger.Fatal("Ошибка при подключении к S3: ", err)
	}
	redisConn, err := connector.GetRedisConnector(authParams.Redis.Addr, authParams.Redis.Password, authParams.Redis.DB)
	if err != nil {
		logger.Fatal("Ошибка при подключении к Redis: ", err)
	}

	// Repositories
	userRepo := postgres.NewUserRepository(psqlConn)
	contentRepo := postgres.NewContentRepository(psqlConn)
	reviewRepo := postgres.NewReviewRepository(psqlConn)
	compilationRepo := postgres.NewCompilationRepository(psqlConn)
	searchRepo := postgres.NewSearchRepository(psqlConn, contentRepo)
	favouriteRepo := postgres.NewFavouriteRepository(psqlConn)
	staticRepo := postgres.NewStaticRepository(psqlConn, s3conn, staticParams.S3.BucketName, staticParams.MaxFileSize)
	authRepository := redis.NewSessionRepository(redisConn, authParams.SessionAliveTime)

	// Use Cases
	profanityUseCase, err := profanity.NewGateway(coreParams.Microservices.ProfanityFilter.Addr)
	if err != nil {
		logger.Fatalf("Ошибка при подключении к сервису фильтрации сообщений: %v", err)
	}

	authUseCase := service.NewAuthService(authRepository)
	staticUseCase := service.NewStaticService(staticRepo)
	userUseCase := service.NewUserService(userRepo, staticUseCase)
	contentUseCase := service.NewContentService(contentRepo, staticUseCase, coreParams.ContentSecretKey)
	reviewUseCase := service.NewReviewService(reviewRepo, userRepo, contentRepo, staticUseCase, profanityUseCase)
	compilationUseCase := service.NewCompilationService(compilationRepo, staticUseCase, contentUseCase)
	searchUseCase := service.NewSearchService(searchRepo, contentUseCase)
	favouriteUseCase := service.NewFavouriteService(favouriteRepo, contentUseCase)

	sessionManager := utils.NewSessionManager(authUseCase,
		coreParams.Microservices.Auth.HTTPSessionAliveTime, coreParams.HTTP.SecureCookies)

	// Delivery
	staticDelivery := delivery.NewStaticEndpoints(staticUseCase)
	authDelivery := delivery.NewAuthEndpoints(authUseCase, sessionManager)
	userDelivery := delivery.NewUserEndpoints(userUseCase, authUseCase, staticUseCase, sessionManager)
	contentDelivery := delivery.NewContentEndpoints(contentUseCase)
	playgroundDelivery := delivery.NewPlaygroundEndpoints()
	reviewDelivery := delivery.NewReviewEndpoints(reviewUseCase, authUseCase)
	compilationDelivery := delivery.NewCompilationEndpoints(compilationUseCase)
	searchDelivery := delivery.NewSearchEndpoints(searchUseCase)
	ongoingDelivery := delivery.NewOngoingContentEndpoints(contentUseCase, authUseCase)
	favouriteDelivery := delivery.NewFavouriteEndpoints(favouriteUseCase, authUseCase)

	// REST API
	echoServer := echo.New()
	echoServer.Server.ReadTimeout = time.Duration(coreParams.HTTP.Server.ReadTimeout) * time.Second
	echoServer.Server.ReadHeaderTimeout = time.Duration(coreParams.HTTP.Server.ReadTimeout) * time.Second
	echoServer.Server.WriteTimeout = time.Duration(coreParams.HTTP.Server.WriteTimeout) * time.Second
	echoServer.Server.IdleTimeout = time.Duration(coreParams.HTTP.Server.ReadTimeout) * time.Second

	// static
	staticAPI := echoServer.Group("")
	staticDelivery.Configure(staticAPI)

	// middleware
	// metrics
	// echoServer.Use(echoprometheus.NewMiddleware("Kinoskop"))
	// echoServer.GET("/metrics", echoprometheus.NewHandler())
	// config
	echoServer.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set("params", coreParams)
			return next(ctx)
		}
	})
	// requestID middleware
	echoServer.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			reqID := uuid.New().String()
			ctx.Set(echo.HeaderXRequestID, reqID)
			ctx.Response().Header().Set(echo.HeaderXRequestID, reqID)
			return next(ctx)
		}
	})
	// CORS
	echoServer.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, coreParams.HTTP.CORSAllowedOrigins)
			ctx.Response().Header().Set(echo.HeaderAccessControlAllowMethods, strings.Join([]string{
				http.MethodGet,
				http.MethodPut,
				http.MethodPost,
				http.MethodDelete,
				http.MethodOptions,
			}, ","))
			ctx.Response().Header().Set(echo.HeaderAccessControlAllowHeaders, strings.Join([]string{
				echo.HeaderOrigin,
				echo.HeaderAccept,
				echo.HeaderXRequestedWith,
				echo.HeaderContentType,
				echo.HeaderAccessControlRequestMethod,
				echo.HeaderAccessControlRequestHeaders,
				echo.HeaderCookie,
				"X-Csrf",
			}, ","))
			ctx.Response().Header().Set(echo.HeaderAccessControlAllowCredentials, "true")
			ctx.Response().Header().Set(echo.HeaderAccessControlMaxAge, "86400")
			return next(ctx)
		}
	})
	// СSRF
	echoServer.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookieHTTPOnly: false,
		CookiePath:     "/",
		TokenLookup:    "header:X-Csrf",
	}))
	// recover
	echoServer.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			defer func() {
				if recErr := recover(); recErr != nil {
					reqID := ctx.Get(echo.HeaderXRequestID)
					if reqID == nil {
						reqID = "unknown"
					}
					logs := fmt.Errorf(
						"внутренняя ошибка сервера: %v\nRequestID: %v\nStack Trace:\n%s",
						recErr, reqID, debug.Stack(),
					)
					ctx.Logger().Error(logs)
					ctx.Error(entity.ErrInternal)
				}
			}()
			return next(ctx)
		}
	})
	// Endpoints
	api := echoServer.Group("/api")
	// docs
	api.GET("/docs*", echoSwagger.WrapHandler)
	// playground
	playgroundAPI := api.Group("/playground")
	playgroundAPI.GET("/ping", playgroundDelivery.Ping)
	// content
	contentAPI := api.Group("/content")
	contentDelivery.Configure(contentAPI)
	// user
	userAPI := api.Group("/user")
	userDelivery.Configure(userAPI)
	// auth
	authAPI := api.Group("/auth")
	authDelivery.Configure(authAPI)
	// reviews
	reviewAPI := api.Group("/review")
	reviewDelivery.Configure(reviewAPI)
	// compilations
	compilationAPI := api.Group("/compilation")
	compilationDelivery.Configure(compilationAPI)
	// search
	searchAPI := api.Group("/search")
	searchDelivery.Configure(searchAPI)
	// ongoing
	ongoingAPI := api.Group("/ongoing")
	ongoingDelivery.Configure(ongoingAPI)
	// favourite
	favouriteAPI := api.Group("/favourite")
	favouriteDelivery.Configure(favouriteAPI)
	return echoServer
}

func Run(server *echo.Echo, params config.Config) {
	if err := server.Start(params.GetServerAddr()); err != nil && !errors.Is(err, http.ErrServerClosed) {
		server.Logger.Fatalf("Сервер завершил свою работу по причине: %v\n", err)
	}
}

func Shutdown(ctx context.Context, server *echo.Echo) {
	if err := server.Shutdown(ctx); err != nil {
		server.Logger.Fatalf("Во время выключения сервера возникла ошибка: %s\n", err)
	}
}
