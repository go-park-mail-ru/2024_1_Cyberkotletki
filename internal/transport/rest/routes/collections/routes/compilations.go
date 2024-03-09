package routes

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/collections"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/httputil"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"github.com/gorilla/mux"
	"net/http"
)

// GetCompilation
// @Tags Collections
// @Description Возвращает актуальные подборки фильмов по указанному жанру. Если передать cookies с сессией, то подборка будет персонализированной
// @Accept json
// @Param 	Cookie header string  false "session"     default(session=xxx)
// @Param	genre	path	string	true	"Название жанра"
// @Success		200	{object}	collections.CompilationData	"Список с id фильмов указанного жанра"
// @Success		400	{object}	httputil.HTTPError	"Требуется указать жанр"
// @Success		404	{object}	httputil.HTTPError	"Такой жанр не найден"
// @Failure		500	{object}	httputil.HTTPError	"Внутренняя ошибка сервера"
// @Router /collections/compilation/genre/{genre} [get]
func GetCompilation(w http.ResponseWriter, r *http.Request) {
	genre, ok := mux.Vars(r)["genre"]
	if !ok {
		httputil.NewError(w, http.StatusBadRequest, exc.New(exc.Transport, exc.BadRequest, "требуется указать жанр"))
		return
	}
	compilation, err := collections.GetCompilation(genre)
	if err != nil {
		if exc.Is(err, exc.NotFoundErr) {
			httputil.NewError(w, http.StatusNotFound, err)
		} else {
			httputil.NewError(w, http.StatusInternalServerError, err)
		}
		return
	}
	httputil.WriteJSON(w, *compilation)
}
