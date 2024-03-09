package playground

import (
	_ "github.com/go-park-mail-ru/2024_1_Cyberkotletki/docs"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/playground/routes"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, params config.InitParams) {
	router.HandleFunc("/ping", routes.Ping).Methods("GET")
}
