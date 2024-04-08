package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
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
// @Success		200	{object}	dto.Genres	"Список с id фильмов указанного жанра"
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /collections/genres [get]
func (h *CollectionsEndpoints) GetGenres(ctx echo.Context) error {
	genres, err := h.useCase.GetGenres()
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, err)
	}
	return utils.WriteJSON(ctx, genres)
}

// GetCompilationByGenre
// @Tags Collections
// @Description Возвращает список всех доступных жанров фильмов и сериалов
// @Param		genre	query	string	true	"Жанр для составления подборки"
// @Success		200	{object}	dto.Compilation	"Список с id фильмов указанного жанра"
// @Failure		404	{object}	echo.HTTPError	"Такого жанра не существует"
// @Failure		404	{object}	echo.HTTPError	"Такого жанра не существует"
// @Failure		500	{object}	echo.HTTPError	"Внутренняя ошибка сервера"
// @Router /collections/compilation [get]
func (h *CollectionsEndpoints) GetCompilationByGenre(ctx echo.Context) error {
	genre := ctx.QueryParam("genre")
	if genre == "" {
		return utils.NewError(
			ctx,
			http.StatusBadRequest,
			entity.NewClientError("Не указан жанр", entity.ErrBadRequest),
		)
	}
	compilation, err := h.useCase.GetCompilation(genre)
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		case entity.Contains(err, entity.ErrBadRequest):
			return utils.NewError(ctx, http.StatusBadRequest, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, compilation)
}
