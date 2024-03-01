package collections

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/collections/routes"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/compilations", routes.GetCompilations).Methods("GET")
}
