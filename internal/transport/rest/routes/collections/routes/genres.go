package routes

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/collections"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/httputil"
	"net/http"
)

// GetGenres
// @Tags Collections
// @Description Возвращает список всех доступных жанров фильмов и сериалов
// @Success		200	{object}	collections.GenresData	"Список с id фильмов указанного жанра"
// @Failure		500	{object}	httputil.HTTPError	"Внутренняя ошибка сервера"
// @Router /collections/genres [get]
func GetGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := collections.GetGenres()
	if err != nil {
		httputil.NewError(w, http.StatusInternalServerError, err)
		return
	}
	httputil.WriteJSON(w, genres)
}
