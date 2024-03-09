package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/routes/auth/routes"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, params config.InitParams) {
	router.HandleFunc("/login", routes.Login)
	router.HandleFunc("/register", routes.Register)
	router.HandleFunc("/isAuth", routes.IsAuth)
	router.HandleFunc("/logout", routes.Logout)
}
