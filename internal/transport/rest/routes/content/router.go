package content

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/content/routes"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, params config.InitParams) {
	router.HandleFunc("/contentPreview", routes.GetContentPreview).Methods("GET")
}
