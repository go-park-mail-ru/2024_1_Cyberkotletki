package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/echoutil"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type ContentEndpoints struct {
	useCase usecase.Content
}

func NewContentEndpoints(useCase usecase.Content) ContentEndpoints {
	return ContentEndpoints{useCase: useCase}
}

// GetContentPreview
// @Tags Content
// @Description Возвращает краткую информацию о фильме или сериале
// @Param	id	query	int true	"ID искомого контента. Контентом может быть как фильм, так и сериал"
// @Success		200	{object}	dto.PreviewContentCard	"Список с id фильмов указанного жанра"
// @Failure		400	{object}	echo.HTTPError	"Требуется указать валидный id контента"
// @Failure		404	{object}	echo.HTTPError	"Контент с таким id не найден"
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /content/contentPreview [get]
func (h *ContentEndpoints) GetContentPreview(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.QueryParam("id"), 10, 64)
	if err != nil {
		return echoutil.NewError(ctx, http.StatusBadRequest, err)
	}
	contentPreview, err := h.useCase.GetContentPreviewCard(int(id))
	if err != nil {
		if entity.Contains(err, entity.ErrNotFound) {
			return echoutil.NewError(ctx, http.StatusNotFound, err)
		}
		return echoutil.NewError(ctx, http.StatusInternalServerError, err)
	}
	return echoutil.WriteJSON(ctx, *contentPreview)
}
