package utils

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// WriteJSON отправляет json клиенту или ошибку, если по каким-либо причинам невозможно преобразовать
// response в json
func WriteJSON(c echo.Context, response any) error {
	err := c.JSON(http.StatusOK, response)
	if err != nil {
		return NewError(c, http.StatusInternalServerError, err)
	}
	return nil
}
