package routes

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/httputil"
	"net/http"
)

// Ping
// @Tags Playground
// @Description Проверка соединения через классический ping pong
// @Success 200 {string} string "Pong"
// @Router /playground/ping [get]
func Ping(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("pong")); err != nil {
		httputil.NewError(w, http.StatusInternalServerError, err)
	}
}
