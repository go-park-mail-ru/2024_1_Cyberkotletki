package routes

import (
	"encoding/json"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/collections"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/httputil"
	"net/http"
	"time"
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
		httputil.NewError(w, 500, exc.Exception{
			When:  time.Now(),
			What:  "Внутреняя ошибка сервера",
			Layer: exc.Server,
			Type:  exc.Untyped,
		})
	} else {
		j, _ := json.Marshal(genres)
		_, _ = w.Write(j)
	}
}
