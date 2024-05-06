package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
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
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

func Init(logger echo.Logger, params config.Config) *echo.Echo {
	// DBConn
	psqlConn, err := connector.GetPostgresConnector(params.Postgres.ConnectURL)
	if err != nil {
		logger.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}
	// RedisConn
	redisConn, err := connector.GetRedisConnector(params.Redis.Addr, params.Redis.Password, params.Redis.DB)
	if err != nil {
		logger.Fatalf("Ошибка при подключении к Redis: %v", err)
	}

	// Repositories
	userRepo := postgres.NewUserRepository(psqlConn)
	contentRepo := postgres.NewContentRepository(psqlConn)
	sessionRepo := redis.NewSessionRepository(redisConn, params.Auth.SessionAliveTime)
	staticRepo := postgres.NewStaticRepository(psqlConn, params.Static.Path, params.Static.MaxFileSize)
	reviewRepo := postgres.NewReviewRepository(psqlConn)
	compilationRepo := postgres.NewCompilationRepository(psqlConn)
	ongoingRepo := postgres.NewOngoingContentRepository(psqlConn)

	// Use Cases
	staticUseCase := service.NewStaticService(staticRepo)
	authUseCase := service.NewAuthService(sessionRepo)
	userUseCase := service.NewUserService(userRepo, staticUseCase)
	contentUseCase := service.NewContentService(contentRepo, staticRepo)
	reviewUseCase := service.NewReviewService(reviewRepo, userRepo, contentRepo, staticRepo)
	compilationUseCase := service.NewCompilationService(compilationRepo, staticRepo, contentRepo)
	ongoingUseCase := service.NewOngoingContentService(ongoingRepo, staticRepo)

	sessionManager := utils.NewSessionManager(authUseCase, params.Auth.SessionAliveTime, params.HTTP.SecureCookies)

	// Delivery
	staticDelivery := delivery.NewStaticEndpoints(staticUseCase)
	authDelivery := delivery.NewAuthEndpoints(authUseCase, sessionManager)
	userDelivery := delivery.NewUserEndpoints(userUseCase, authUseCase, staticUseCase, sessionManager)
	contentDelivery := delivery.NewContentEndpoints(contentUseCase)
	playgroundDelivery := delivery.NewPlaygroundEndpoints()
	reviewDelivery := delivery.NewReviewEndpoints(reviewUseCase, authUseCase)
	compilationDelivery := delivery.NewCompilationEndpoints(compilationUseCase)
	ongoingDelivery := delivery.NewOngoingContentEndpoints(ongoingUseCase)

	// REST API
	echoServer := echo.New()
	echoServer.Server.ReadTimeout = time.Duration(params.HTTP.Server.ReadTimeout) * time.Second
	echoServer.Server.ReadHeaderTimeout = time.Duration(params.HTTP.Server.ReadTimeout) * time.Second
	echoServer.Server.WriteTimeout = time.Duration(params.HTTP.Server.WriteTimeout) * time.Second
	echoServer.Server.IdleTimeout = time.Duration(params.HTTP.Server.ReadTimeout) * time.Second
	// статика
	echoServer.Static("/static/", params.Static.Path)
	// middleware
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
					log := fmt.Errorf(
						"внутренняя ошибка сервера: %v\nRequestID: %v\nStack Trace:\n%s",
						recErr, reqID, debug.Stack(),
					)
					ctx.Logger().Error(log)
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
	// static
	staticAPI := api.Group(params.Static.Path)
	staticDelivery.Configure(staticAPI)
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
	// ongoing
	ongoingAPI := api.Group("/ongoing")
	ongoingDelivery.Configure(ongoingAPI)
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
