package routes

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest/httputil"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"net/http"
	"strconv"
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
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 32)
	if err != nil {
		httputil.NewError(w, http.StatusBadRequest, exc.New(exc.Transport, exc.BadRequest, "требуется указать валидный id контента"))
		return
	}
	contentPreview, err := content.GetContentPreviewInfo(int(id))
	if err != nil {
		if exc.Is(err, exc.NotFoundErr) {
			httputil.NewError(w, http.StatusNotFound, err)
		} else {
			httputil.NewError(w, http.StatusInternalServerError, err)
		}
		return
	}
	httputil.WriteJSON(w, *contentPreview)
}
