package routes

import (
	"encoding/json"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/httputil"
	"net/http"
	"strconv"
	"time"
)

// GetContentPreview
// @Tags Content
// @Description Возвращает краткую информацию о фильме или сериале
// @Param	id	query	int true	"ID искомого контента. Контентом может быть как фильм, так и сериал"
// @Success		200	{object}	content.PreviewInfoData	"Список с id фильмов указанного жанра"
// @Failure		400	{object}	httputil.HTTPError	"Требуется указать валидный id контента"
// @Failure		404	{object}	httputil.HTTPError	"Контент с таким id не найден"
// @Failure		500	{object}	httputil.HTTPError	"Внутренняя ошибка сервера"
// @Router /content/contentPreview [get]
func GetContentPreview(w http.ResponseWriter, r *http.Request) {
	if id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 32); err == nil {
		if compilation, err := content.GetContentPreviewInfo(int(id)); err != nil {
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
			What:  "Требуется указать валидный id контента",
			Layer: exc.Transport,
			Type:  exc.Unprocessable,
		})
	}
}
