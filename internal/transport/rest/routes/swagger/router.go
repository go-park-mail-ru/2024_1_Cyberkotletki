package swagger

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func RegisterRoutes(router *mux.Router) {
	router.PathPrefix("/").HandlerFunc(httpSwagger.WrapHandler)
}
