package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type CompilationEndpoints struct {
	сompilationUC usecase.Compilation
}

func NewCompilationEndpoints(сompilationUC usecase.Compilation) CompilationEndpoints {
	return CompilationEndpoints{сompilationUC: сompilationUC}
}

func (h *CompilationEndpoints) Configure(server *echo.Group) {
	server.GET("/compilation/:id", h.GetCompilation)
	server.GET("/compilation/types", h.GetCompilationTypes)
	server.GET("/compilation/type/:compilationType", h.GetCompilationsByCompilationType)

}

// GetCompilationTypes
// @Summary Получение списка подборок
// @Tags compilation
// @Description Получение списка подборок по id
// @Accept json
// @Produce json
// @Success 200 {object} dto.CompilationTypeResponseList
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /compilation/types [get]
func (h *CompilationEndpoints) GetCompilationTypes(ctx echo.Context) error {
	compType, err := h.сompilationUC.GetCompilationTypes()
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, compType)
}

// GetCompilation
// @Summary Получение подборки
// @Tags compilation
// @Description Получение подборки по id
// @Accept json
// @Produce json
// @Param id path string true "id подборки"
// @Success 200 {object} dto.CompilationResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /compilation/{id} [get]
func (h *CompilationEndpoints) GetCompilation(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id подборки"))
	}
	compilation, err := h.сompilationUC.GetCompilation(int(id))
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, compilation)
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
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /compilation/type/{compilationType} [get]
func (h *CompilationEndpoints) GetCompilationsByCompilationType(ctx echo.Context) error {
	compilationType, err := strconv.ParseInt(ctx.Param("compilationType"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id типа подборки"))
	}
	compilations, err := h.сompilationUC.GetCompilationsByCompilationType(int(compilationType))
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, compilations)
}
