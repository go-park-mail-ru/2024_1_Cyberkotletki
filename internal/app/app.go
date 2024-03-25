package app

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	delivery "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
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
	userRepo := tmpdb.NewUserRepository()
	contentRepo := tmpdb.NewContentRepository()
	sessionRepo := redis.NewSessionRepository(logger, params)

	// Use Cases
	authUseCase := service.NewAuthService(userRepo, sessionRepo)
	contentUseCase := service.NewContentService(contentRepo)
	collectionsUseCase := service.NewCollectionsService(contentRepo)

	// Delivery
	authDelivery := delivery.NewAuthEndpoints(authUseCase)
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
					ctx.Logger().Errorf(
						"Внутренняя ошибка сервера: %v\nRequestID: %v\nStack Trace:\n%s",
						recErr, reqID, debug.Stack(),
					)
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
	contentAPI.GET("/contentPreview", contentDelivery.GetContentPreview)
	// collections
	collectionsAPI := api.Group("/collections")
	collectionsAPI.GET("/genres", collectionsDelivery.GetGenres)
	collectionsAPI.GET("/compilation", collectionsDelivery.GetCompilationByGenre)
	// auth
	authAPI := api.Group("/auth")
	authAPI.POST("/register", authDelivery.Register)
	authAPI.POST("/login", authDelivery.Login)
	authAPI.GET("/isAuth", authDelivery.IsAuth)
	authAPI.POST("/logout", authDelivery.Logout)
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
