package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	_ "github.com/go-park-mail-ru/2024_1_Cyberkotletki/docs"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/auth"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/static"
	delivery "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/connector"
	"github.com/google/uuid"
	_ "github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/mcuadros/go-defaults"
	_ "github.com/prometheus/client_golang/prometheus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"time"
)

func GenerateExampleConfig() {
	// Создание экземпляра структуры с использованием значений по умолчанию из тегов default
	var conf config.Config
	defaults.SetDefaults(&conf)

	yamlData, err := yaml.Marshal(&conf)
	if err != nil {
		fmt.Printf("Ошибка при маршализации YAML: %v\n", err)
		return
	}
	file, err := os.Create("config.yaml")
	if err != nil {
		fmt.Printf("Ошибка при создании файла: %v\n", err)
		return
	}
	_, err = file.Write(yamlData)
	if err != nil {
		fmt.Printf("Ошибка при записи в файл: %v\n", err)
		return
	}
	err = file.Close()
	if err != nil {
		fmt.Printf("Ошибка при записи в файл: %v\n", err)
		return
	}

	fmt.Println("Конфигурационный файл успешно создан.")
}

func ParseParams() config.Config {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка при чтении конфига сервера: %v", err)
	}
	var cfg config.Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Ошибка при парсинге конфига сервера: %v", err)
	}

	return cfg
}

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
		GenerateExampleConfig()
		return
	}

	logger := log.New("server: ")
	params := ParseParams()
	logger.Printf("Параметры запуска сервера: %v \n", params)

	echoServer := Init(logger, params)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	go Run(echoServer, params)

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(params.HTTP.Server.GracefulShutdownTimeout)*time.Second,
	)
	defer cancel()
	Shutdown(ctx, echoServer)
}

func Init(logger echo.Logger, params config.Config) *echo.Echo {
	// DBConn
	psqlConn, err := connector.GetPostgresConnector(params.Postgres.ConnectURL)
	if err != nil {
		logger.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	// Repositories
	userRepo := postgres.NewUserRepository(psqlConn)
	contentRepo := postgres.NewContentRepository(psqlConn)
	reviewRepo := postgres.NewReviewRepository(psqlConn)
	compilationRepo := postgres.NewCompilationRepository(psqlConn)
	searchRepo := postgres.NewSearchRepository(psqlConn, contentRepo)
	ongoingRepo := postgres.NewOngoingContentRepository(psqlConn)
	favouriteRepo := postgres.NewFavouriteRepository(psqlConn)

	// Use Cases
	staticUseCase, err := static.NewGateway(params.Microservices.Static.Addr)
	if err != nil {
		logger.Fatalf("Ошибка при подключении к сервису статики: %v", err)
	}
	authUseCase, err := auth.NewGateway(params.Microservices.Auth.Addr)
	if err != nil {
		logger.Fatalf("Ошибка при подключении к сервису авторизации: %v", err)
	}
	userUseCase := service.NewUserService(userRepo, staticUseCase)
	contentUseCase := service.NewContentService(contentRepo, staticUseCase)
	reviewUseCase := service.NewReviewService(reviewRepo, userRepo, contentRepo, staticUseCase)
	compilationUseCase := service.NewCompilationService(compilationRepo, staticUseCase, contentRepo)
	searchUseCase := service.NewSearchService(searchRepo, staticUseCase)
	ongoingUseCase := service.NewOngoingContentService(ongoingRepo, staticUseCase)
	favouriteUseCase := service.NewFavouriteService(favouriteRepo)

	sessionManager := utils.NewSessionManager(authUseCase,
		params.Microservices.Auth.HTTPSessionAliveTime, params.HTTP.SecureCookies)

	// Delivery
	staticDelivery := delivery.NewStaticEndpoints(staticUseCase)
	authDelivery := delivery.NewAuthEndpoints(authUseCase, sessionManager)
	userDelivery := delivery.NewUserEndpoints(userUseCase, authUseCase, staticUseCase, sessionManager)
	contentDelivery := delivery.NewContentEndpoints(contentUseCase)
	playgroundDelivery := delivery.NewPlaygroundEndpoints()
	reviewDelivery := delivery.NewReviewEndpoints(reviewUseCase, authUseCase)
	compilationDelivery := delivery.NewCompilationEndpoints(compilationUseCase)
	searchDelivery := delivery.NewSearchEndpoints(searchUseCase)
	ongoingDelivery := delivery.NewOngoingContentEndpoints(ongoingUseCase)
	favouriteDelivery := delivery.NewFavouriteEndpoints(favouriteUseCase, authUseCase)

	// REST API
	echoServer := echo.New()
	echoServer.Server.ReadTimeout = time.Duration(params.HTTP.Server.ReadTimeout) * time.Second
	echoServer.Server.ReadHeaderTimeout = time.Duration(params.HTTP.Server.ReadTimeout) * time.Second
	echoServer.Server.WriteTimeout = time.Duration(params.HTTP.Server.WriteTimeout) * time.Second
	echoServer.Server.IdleTimeout = time.Duration(params.HTTP.Server.ReadTimeout) * time.Second

	// static
	staticAPI := echoServer.Group("")
	staticDelivery.Configure(staticAPI)

	// middleware
	// metrics
	echoServer.Use(echoprometheus.NewMiddleware("Kinoskop"))
	echoServer.GET("/metrics", echoprometheus.NewHandler())
	// config
	echoServer.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set("params", params)
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
			ctx.Response().Header().Set(echo.HeaderAccessControlAllowOrigin, params.HTTP.CORSAllowedOrigins)
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
