package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type PlaygroundEndpoints struct {
}

func NewPlaygroundEndpoints() PlaygroundEndpoints {
	return PlaygroundEndpoints{}
}

// Ping
// @Tags Playground
// @Description Проверка соединения через классический ping pong
// @Success 200 {string} string "Pong"
// @Router /playground/ping [get]
func (h *PlaygroundEndpoints) Ping(c echo.Context) error {
	if err := c.String(http.StatusOK, "pong"); err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return nil
}
