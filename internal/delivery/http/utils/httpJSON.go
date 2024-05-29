package utils

import (
	"github.com/labstack/echo/v4"
	"github.com/mailru/easyjson"
	"net/http"
)

// WriteJSON отправляет json клиенту или ошибку, если по каким-либо причинам невозможно преобразовать
// response в json
func WriteJSON(ctx echo.Context, response easyjson.Marshaler) error {
	jsonData, err := easyjson.Marshal(response)
	if err != nil {
		return NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	}
	if err = ctx.JSONBlob(http.StatusOK, jsonData); err != nil {
		return NewError(ctx, http.StatusInternalServerError, "Внутренняя ошибка сервера", err)
	}
	return nil
}

// ReadJSON читает json из тела запроса и заполняет указанную структуру
func ReadJSON(ctx echo.Context, v easyjson.Unmarshaler) error {
	err := easyjson.UnmarshalFromReader(ctx.Request().Body, v)
	if err != nil {
		return NewError(ctx, http.StatusBadRequest, "Невалидный JSON", err)
	}
	return nil
}
