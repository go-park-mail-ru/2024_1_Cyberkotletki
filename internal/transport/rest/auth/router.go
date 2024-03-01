package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/auth/routes"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", routes.Login)
	router.HandleFunc("/register", routes.Register)
}
