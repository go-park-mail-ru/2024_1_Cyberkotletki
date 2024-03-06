package collections

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/collections/routes"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/compilation/{genre:[a-z]+}", routes.GetCompilation).Methods("GET")
	router.HandleFunc("/genres", routes.GetGenres).Methods("GET")
}
