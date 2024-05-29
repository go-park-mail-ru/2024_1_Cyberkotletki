package http

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
	"io"
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
	e.GET("/api/static/:id", h.GetStaticURL)
	e.GET("/static/:path", h.GetStaticFile)
}

// GetStaticURL
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
// @Router /api/static/{id} [get]
func (h *StaticEndpoints) GetStaticURL(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return utils.NewError(ctx, http.StatusBadRequest, "Невалидный id статики", nil)
	}
	staticURL, err := h.staticUC.GetStatic(int(id))
	switch {
	case errors.Is(err, usecase.ErrStaticNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "Статика не найдена", nil)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	}
	return ctx.String(http.StatusOK, staticURL)
}

// GetStaticFile
// @Tags Static
// @Description Получение статического файла по относительному пути. Возвращает файл в виде байтов.
// @Accept json
// @Param 	path	path	string	true	"Путь до статики"
// @Success     200
// @Failure		400	{object}	echo.HTTPError	"невалидный путь до статики"
// @Failure		404	{object}	echo.HTTPError	"файл не найден"
// @Failure		500	{object}	echo.HTTPError	"ошибка сервера"
// @Router /static/{path} [get]
func (h *StaticEndpoints) GetStaticFile(ctx echo.Context) error {
	path := ctx.Param("path")
	staticFile, err := h.staticUC.GetStaticFile(path)
	switch {
	case errors.Is(err, usecase.ErrStaticNotFound):
		return utils.NewError(ctx, http.StatusNotFound, "Статика не найдена", nil)
	case err != nil:
		return utils.NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	}
	// Читаем весь файл в буфер
	buf, err := io.ReadAll(staticFile)
	if err != nil {
		return utils.NewError(ctx, http.StatusInternalServerError, "Ошибка при чтении файла", err)
	}

	// Определение типа файла
	contentType := http.DetectContentType(buf)
	return ctx.Blob(http.StatusOK, contentType, buf)
}
