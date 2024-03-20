package app

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	delivery "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/redis"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/tmpDB"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"strings"
	"time"
)

func Init(logger echo.Logger, params config.Config) *echo.Echo {
	// Repositories
	userRepo := tmpDB.NewUserRepository()
	contentRepo := tmpDB.NewContentRepository()
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
	e := echo.New()
	e.Server.ReadTimeout = time.Duration(params.HTTP.Server.ReadTimeout) * time.Second
	e.Server.ReadHeaderTimeout = time.Duration(params.HTTP.Server.ReadTimeout) * time.Second
	e.Server.WriteTimeout = time.Duration(params.HTTP.Server.WriteTimeout) * time.Second
	e.Server.IdleTimeout = time.Duration(params.HTTP.Server.ReadTimeout) * time.Second
	// статика
	e.Static("/static/", params.HTTP.StaticFolder)
	// middleware
	// config
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("params", params)
			return next(c)
		}
	})
	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(params.HTTP.CORSAllowedOrigins, ","),
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderAccept,
			echo.HeaderXRequestedWith,
			echo.HeaderAccessControlRequestMethod,
			echo.HeaderAccessControlRequestHeaders,
			echo.HeaderCookie,
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}))
	// Endpoints
	api := e.Group("/api")
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
	// auth
	authAPI := api.Group("/auth")
	authAPI.POST("/register", authDelivery.Register)
	authAPI.POST("/login", authDelivery.Login)
	authAPI.GET("/isAuth", authDelivery.IsAuth)
	authAPI.POST("/logout", authDelivery.Logout)
	return e
}

func Run(server *echo.Echo, params config.Config) {
	if err := server.Start(params.GetServerAddr()); err != nil && !errors.Is(err, http.ErrServerClosed) {
		server.Logger.Fatalf("Сервер завершил свою работу по причине: %v\n", err)
	}
}

func Shutdown(server *echo.Echo, ctx context.Context) {
	if err := server.Shutdown(ctx); err != nil {
		server.Logger.Fatalf("Во время выключения сервера возникла ошибка: %s\n", err)
	}
}
