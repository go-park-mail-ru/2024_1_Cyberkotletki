package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
)

type OngoingContentEndpoints struct {
	usecase usecase.OngoingContent
}

func NewOngoingContentEndpoints(usecase usecase.OngoingContent) *OngoingContentEndpoints {
	return &OngoingContentEndpoints{
		usecase: usecase,
	}
}

func (h *OngoingContentEndpoints) Configure(server *echo.Group) {
	server.GET("/:id", h.GetOngoingContentByContentID)
	server.GET("/nearest", h.GetNearestOngoings)
	server.GET("/:year/:month", h.GetOngoingContentByMonthAndYear)
	server.GET("/years", h.GetAllReleaseYears)
}

// GetOngoingContentByContentID godoc
// @Summary Получить контент календаря релизов по id контента
// @Tags ongoing_content
// @Produce json
// @Param id path int true "ID контента календаря релизов"
// @Success 200 {object} dto.PreviewOngoingContentCardVertical
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /ongoing/{id} [get]
func (h *OngoingContentEndpoints) GetOngoingContentByContentID(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "невалидный id контента календаря релизов", err)
	}

	ongoingContent, err := h.usecase.GetOngoingContentByContentID(int(id))
	switch {
	case errors.Is(err, usecase.ErrOngoingContentNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "контент календаря релизов не найден", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "ошибка при получении контента календаря релизов", err)
	default:
		return utils.WriteJSON(ctx, ongoingContent)
	}
}

// GetNearestOngoings godoc
// @Summary Получить ближайшие релизы
// @Tags ongoing_content
// @Produce json
// @Param limit query int false "Количество релизов"
// @Success 200 {array} dto.PreviewOngoingContentCardVertical
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /ongoing/nearest [get]
func (h *OngoingContentEndpoints) GetNearestOngoings(ctx echo.Context) error {
	limit, err := strconv.ParseInt(ctx.QueryParam("limit"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "невалидное количество релизов", err)
	}

	ongoingContent, err := h.usecase.GetNearestOngoings(int(limit))
	switch {
	case errors.Is(err, usecase.ErrOngoingContentNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "контент календаря релизов не найден", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "ошибка при получении ближайших релизов", err)
	default:
		return utils.WriteJSON(ctx, ongoingContent)
	}
}

// GetOngoingContentByMonthAndYear godoc
// @Summary Получить релизы по месяцу и году
// @Tags ongoing_content
// @Produce json
// @Param month path int true "Месяц"
// @Param year path int true "Год"
// @Success 200 {array} dto.PreviewOngoingContentCardVertical
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /ongoing/{year}/{month} [get]
func (h *OngoingContentEndpoints) GetOngoingContentByMonthAndYear(ctx echo.Context) error {
	month, err := strconv.ParseInt(ctx.Param("month"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "невалидный месяц", err)
	}
	year, err := strconv.ParseInt(ctx.Param("year"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "невалидный год", err)
	}

	ongoingContent, err := h.usecase.GetOngoingContentByMonthAndYear(int(month), int(year))
	switch {
	case errors.Is(err, usecase.ErrOngoingContentNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "контент календаря релизов не найден", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "ошибка при получении релизов по месяцу и году", err)
	default:
		return utils.WriteJSON(ctx, ongoingContent)
	}
}

// GetAllReleaseYears godoc
// @Summary Получить все года релизов
// @Tags ongoing_content
// @Produce json
// @Success 200 {array} int
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /ongoing/years [get]
func (h *OngoingContentEndpoints) GetAllReleaseYears(ctx echo.Context) error {
	years, err := h.usecase.GetAllReleaseYears()
	switch {
	case errors.Is(err, usecase.ErrOngoingContentYearsNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "года релизов не найдены", err)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "ошибка при получении годов релизов", err)
	default:
		return utils.WriteJSON(ctx, years)
	}
}
