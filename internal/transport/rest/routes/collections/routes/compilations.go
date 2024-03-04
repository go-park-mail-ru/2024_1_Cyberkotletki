package routes

import (
	"encoding/json"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/collections"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/httputil"
	"github.com/gorilla/mux"
	"net/http"
	"time"
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
// @Router /collections/compilation/{genre} [get]
func GetCompilation(w http.ResponseWriter, r *http.Request) {
	if genre, ok := mux.Vars(r)["genre"]; ok {
		if compilation, err := collections.GetCompilation(genre); err != nil {
			if err.Type == exc.NotFound {
				httputil.NewError(w, 404, *err)
			} else {
				httputil.NewError(w, 500, *err)
			}
		} else {
			j, _ := json.Marshal(compilation)
			_, _ = w.Write(j)
		}
	} else {
		httputil.NewError(w, 400, exc.Exception{
			When:  time.Now(),
			What:  "Требуется указать жанр",
			Layer: exc.Transport,
			Type:  exc.Unprocessable,
		})
	}
}
