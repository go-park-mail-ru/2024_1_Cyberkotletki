package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	delivery "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/redis"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/tmpdb"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

func Init(logger echo.Logger, params config.Config) *echo.Echo {
	// Repositories
	userRepo, err := postgres.NewUserRepository(params.User.Postgres)
	if err != nil {
		logger.Fatalf("Ошибка при создании репозитория пользователей: %v", err)
	}
	contentRepo := tmpdb.NewContentRepository()
	sessionRepo, err := redis.NewSessionRepository(params)
	if err != nil {
		logger.Fatalf("Ошибка при создании репозитория сессий: %v", err)
	}
	staticRepo, err := postgres.NewStaticRepository(params.Static.Postgres, params.Static.Path, params.Static.MaxFileSize)
	if err != nil {
		logger.Fatalf("Ошибка при создании репозитория статики: %v", err)
	}

	// Use Cases
	staticUseCase := service.NewStaticService(staticRepo)
	authUseCase := service.NewAuthService(sessionRepo)
	userUseCase := service.NewUserService(userRepo)
	contentUseCase := service.NewContentService(contentRepo)
	collectionsUseCase := service.NewCollectionsService(contentRepo)

	// Delivery
	staticDelivery := delivery.NewStaticEndpoints(staticUseCase)
	authDelivery := delivery.NewAuthEndpoints(authUseCase)
	userDelivery := delivery.NewUserEndpoints(userUseCase, authUseCase, staticUseCase)
	contentDelivery := delivery.NewContentEndpoints(contentUseCase)
	collectionsDelivery := delivery.NewCollectionsEndpoints(collectionsUseCase)
	playgroundDelivery := delivery.NewPlaygroundEndpoints()

	// REST API
	echoServer := echo.New()
	echoServer.Server.ReadTimeout = time.Duration(params.HTTP.Server.ReadTimeout) * time.Second
	echoServer.Server.ReadHeaderTimeout = time.Duration(params.HTTP.Server.ReadTimeout) * time.Second
	echoServer.Server.WriteTimeout = time.Duration(params.HTTP.Server.WriteTimeout) * time.Second
	echoServer.Server.IdleTimeout = time.Duration(params.HTTP.Server.ReadTimeout) * time.Second
	// статика
	echoServer.Static("/static/", params.HTTP.StaticFolder)
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
			}, ","))
			ctx.Response().Header().Set(echo.HeaderAccessControlAllowCredentials, "true")
			ctx.Response().Header().Set(echo.HeaderAccessControlMaxAge, "86400")
			return next(ctx)
		}
	})
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
					fmt.Println(log)
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
	staticAPI := api.Group("/static")
	staticDelivery.Configure(staticAPI)
	// playground
	playgroundAPI := api.Group("/playground")
	playgroundAPI.GET("/ping", playgroundDelivery.Ping)
	// content
	contentAPI := api.Group("/content")
	contentAPI.GET("/contentPreview", contentDelivery.GetContentPreview)
	// collections
	collectionsAPI := api.Group("/collections")
	collectionsAPI.GET("/genres", collectionsDelivery.GetGenres)
	collectionsAPI.GET("/compilation", collectionsDelivery.GetCompilationByGenre)
	// user
	userAPI := api.Group("/user")
	userDelivery.Configure(userAPI)
	// auth
	authAPI := api.Group("/auth")
	authDelivery.Configure(authAPI)
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
