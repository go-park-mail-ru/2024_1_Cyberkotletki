package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/echoutil"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CollectionsEndpoints struct {
	useCase usecase.Collections
}

func NewCollectionsEndpoints(useCase usecase.Collections) CollectionsEndpoints {
	return CollectionsEndpoints{useCase: useCase}
}

// GetGenres
// @Tags Collections
// @Description Возвращает список всех доступных жанров фильмов и сериалов
// @Success		200	{object}	DTO.Genres	"Список с id фильмов указанного жанра"
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /collections/genres [get]
func (h *CollectionsEndpoints) GetGenres(c echo.Context) error {
	genres, err := h.useCase.GetGenres()
	if err != nil {
		return echoutil.NewError(c, http.StatusInternalServerError, err)
	}
	return echoutil.WriteJSON(c, genres)
}
