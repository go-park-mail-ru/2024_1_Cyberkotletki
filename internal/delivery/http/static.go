package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type StaticEndpoints struct {
	staticUC usecase.Static
}

func NewStaticEndpoints(staticUC usecase.Static) StaticEndpoints {
	return StaticEndpoints{staticUC: staticUC}
}

func (h *StaticEndpoints) Configure(e *echo.Group) {
	e.GET("/:id", h.GetStaticUrl)
}

// GetStaticUrl
// @Tags Static
// @Description Получение ссылки на статический файл по id. Возвращает ссылку подобного вида:
// avatars/uuid4.jpg. По умолчанию чтобы получить статику, нужно обратиться по
// адресу вида http://host:port/static/avatars/uuid4.jpg
// @Accept json
// @Param 	id	path	int	true	"ID статики"
// @Success     200
// @Failure		400	{object}	echo.HTTPError	"невалидный id статики"
// @Failure		404	{object}	echo.HTTPError	"файл не найден"
// @Failure		500	{object}	echo.HTTPError	"ошибка сервера"
// @Router /static/{id} [get]
func (h *StaticEndpoints) GetStaticUrl(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, entity.NewClientError("невалидный id статики"))
	}
	staticURL, err := h.staticUC.GetStaticUrl(int(id))
	if err != nil {
		switch {
		case entity.Contains(err, entity.ErrNotFound):
			return utils.NewError(ctx, http.StatusNotFound, err)
		default:
			return utils.NewError(ctx, http.StatusInternalServerError, err)
		}
	}
	return ctx.String(http.StatusOK, staticURL)
}
