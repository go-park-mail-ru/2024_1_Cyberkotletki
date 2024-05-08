package http

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CompilationEndpoints struct {
	compilationUC usecase.Compilation
}

func NewCompilationEndpoints(compilationUC usecase.Compilation) CompilationEndpoints {
	return CompilationEndpoints{compilationUC: compilationUC}
}

func (h *CompilationEndpoints) Configure(server *echo.Group) {
	server.GET("/types", h.GetCompilationTypes)
	server.GET("/type/:compilationType", h.GetCompilationsByCompilationType)
	server.GET("/:id/:page", h.GetCompilationContent)
}

// GetCompilationTypes
// @Summary Получение списка подборок
// @Tags compilation
// @Description Получение списка подборок по id
// @Accept json
// @Produce json
// @Success 200 {object} dto.CompilationTypeResponseList
// @Failure 500 {object} echo.HTTPError
// @Router /api/compilation/types [get]
func (h *CompilationEndpoints) GetCompilationTypes(ctx echo.Context) error {
	compType, err := h.compilationUC.GetCompilationTypes()
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	}
	return utils.WriteJSON(ctx, compType)
}

// GetCompilationsByCompilationType
// @Summary Получение списка подборок по типу подборок
// @Tags compilation
// @Description Получение списка подборок по id типа подборки
// @Accept json
// @Produce json
// @Param compilationType path string true "id типа подборки"
// @Success 200 {object} dto.CompilationResponseList
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/compilation/type/{compilationType} [get]
func (h *CompilationEndpoints) GetCompilationsByCompilationType(ctx echo.Context) error {
	compilationType, err := strconv.ParseInt(ctx.Param("compilationType"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Невалидный id типа подборки", nil)
	}
	compilations, err := h.compilationUC.GetCompilationsByCompilationType(int(compilationType))
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	}
	return utils.WriteJSON(ctx, compilations)
}

// GetCompilationContent
// @Summary Получение карточек контента подборки
// @Tags compilation
// @Description Получение карточек контента подборки по id
// @Accept json
// @Produce json
// @Param id path int true "id подборки"
// @Param page path int true "номер страницы"
// @Success 200 {object} dto.Compilation
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /api/compilation/{id}/{page} [get]
func (h *CompilationEndpoints) GetCompilationContent(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Невалидный id подборки", nil)
	}
	page, err := strconv.ParseInt(ctx.Param("page"), 10, 64)
	if err != nil {
		page = 1
	}
	compilation, err := h.compilationUC.GetCompilationContent(int(id), int(page))
	switch {
	case errors.Is(err, usecase.ErrCompilationNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "Подборка не найдена", nil)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	}
	return utils.WriteJSON(ctx, compilation)
}
