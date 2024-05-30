package http

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
)

type ContentEndpoints struct {
	useCase usecase.Content
}

func NewContentEndpoints(useCase usecase.Content) ContentEndpoints {
	return ContentEndpoints{useCase: useCase}
}

func (h *ContentEndpoints) Configure(server *echo.Group) {
	server.GET("/:id", h.GetContent)
	server.GET("/person/:id", h.GetPerson)
	server.POST("", h.CreateContent)
}

// CreateContent
// @Summary Создание контента
// @Tags content
// @Description Создание контента
// @Accept json
// @Produce json
// @Param content body dto.Content true "Контент"
// @Success 200 {object} dto.Content
// @Failure 400 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /content [post]
// @Security _csrf
func (h *ContentEndpoints) CreateContent(ctx echo.Context) error {
	log.Println("CreateContent")
	var contentCreate dto.Content
	err := ctx.Bind(&contentCreate)
	if err != nil {
		log.Println("Bind error", err)
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный контент"))
	}

	contentDTO, err := h.useCase.CreateContent(contentCreate)
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrBadRequest):
			return utils.NewError(ctx, http.StatusBadRequest, err)
		case entity.Contains(err, entity.ErrAlreadyExists):
			return utils.NewError(ctx, http.StatusConflict, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, contentDTO)

}

// GetContent
// @Summary Получение контента по id
// @Tags content
// @Description Получение контента по id
// @Produce json
// @Param id path int true "ID контента"
// @Success 200 {object} dto.Content
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /content/{id} [get]
func (h *ContentEndpoints) GetContent(ctx echo.Context) error {
	log.Println("GetContent")

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id контента"))
	}
	content, err := h.useCase.GetContentByID(int(id))
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, content)
}

// GetPerson
// @Summary Получение персоны по id
// @Tags content
// @Description Получение персоны по id
// @Produce json
// @Param id path int true "ID персоны"
// @Success 200 {object} dto.Person
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /content/person/{id} [get]
func (h *ContentEndpoints) GetPerson(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id персоны"))
	}
	person, err := h.useCase.GetPersonByID(int(id))
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return utils.WriteJSON(ctx, person)
}
