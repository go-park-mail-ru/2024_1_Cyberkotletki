package routes

import (
	"net/http"
)

// Ping
// @Tags Playground
// @Description Проверка соединения через классический ping pong
// @Success 200 {string} string "Pong"
// @Router /playground/ping [get]
func Ping(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("pong")); err != nil {
		http.Error(w, "что-то пошло не так...", 503)
	}
}
