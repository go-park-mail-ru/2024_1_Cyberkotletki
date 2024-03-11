package rest

import (
	"context"
	_ "github.com/go-park-mail-ru/2024_1_Cyberkotletki/docs"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/middlewares"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/auth"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/collections"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/playground"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"net/http"
)

type RouteRegister func(router *mux.Router, params config.InitParams)

// Routes это map с парами path - функция для регистрации путей
type Routes map[string]RouteRegister

var routes = Routes{
	"/auth":        auth.RegisterRoutes,
	"/collections": collections.RegisterRoutes,
	"/content":     content.RegisterRoutes,
	"/playground":  playground.RegisterRoutes,
}

func RegisterRoutes(router *mux.Router, params config.InitParams) {
	// CORS
	allowedOrigin := new(middlewares.CORSConfig)
	allowedOrigin.SetAllowedOriginsFromConfig(params)
	router.Use(allowedOrigin.SetCORS)

	// Проброс параметров
	router.Use(
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "params", params)))
			})
		},
	)

	// Статика
	fs := http.FileServer(http.Dir(params.StaticFolder))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Ручки
	for path, registerRoute := range routes {
		registerRoute(router.PathPrefix(path).Subrouter(), params)
	}

	// Swagger
	router.PathPrefix("/docs/").HandlerFunc(httpSwagger.WrapHandler)
}
