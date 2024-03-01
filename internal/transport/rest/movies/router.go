package movies

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/movies/routes"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	// router.HandleFunc("/movieFullInfo", routes.MovieFullInfo)
	router.HandleFunc("/moviePreviewInfo", routes.GetMoviePreviewInfo)
}
