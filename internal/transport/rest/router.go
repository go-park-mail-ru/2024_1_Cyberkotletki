package rest

import (
	_ "github.com/go-park-mail-ru/2024_1_Cyberkotletki/docs"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/middlewares"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/auth"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/collections"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/playground"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/swagger"
	"github.com/gorilla/mux"
)

// Routes это map с парами path - функция для регистрации путей
type Routes map[string]func(*mux.Router)

var routes = Routes{
	"/auth":        auth.RegisterRoutes,
	"/collections": collections.RegisterRoutes,
	"/content":     content.RegisterRoutes,
	"/playground":  playground.RegisterRoutes,
	"/swagger":     swagger.RegisterRoutes,
}

func RegisterRoutes(router *mux.Router, params config.InitParams) {
	// CORS
	// router.Use(mux.CORSMethodMiddleware(router))
	allowedOrigin := new(middlewares.CORSConfig)
	allowedOrigin.SetAllowedOriginsFromConfig(params)
	router.Use(allowedOrigin.SetCORS)

	for path, registerRoute := range routes {
		r := router.PathPrefix(path).Subrouter()
		registerRoute(r)
	}
}
