package collections

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/collections/routes"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, params config.InitParams) {
	router.HandleFunc("/compilation/genre/{genre:[a-z]+}", routes.GetCompilation).Methods("GET")
	router.HandleFunc("/genres", routes.GetGenres).Methods("GET")
}
