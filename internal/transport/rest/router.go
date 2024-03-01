package rest

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/auth"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/collections"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/movies"
	"github.com/gorilla/mux"
)

// Routes это map с парами path - функция для регистрации путей
type Routes map[string]func(*mux.Router)

var routes = Routes{
	"/auth":        auth.RegisterRoutes,
	"/collections": collections.RegisterRoutes,
	"/movies":      movies.RegisterRoutes,
}

func RegisterRoutes(router *mux.Router) {
	for path, registerRoute := range routes {
		r := router.PathPrefix(path).Subrouter()
		registerRoute(r)
	}
}
