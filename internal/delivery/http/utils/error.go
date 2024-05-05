package utils

import (
	"errors"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/url"
)

// NewError возвращает *echo.HTTPError, чтобы middleware echo автоматически конвертировал её в JSON-ошибку. Если
// ошибка содержит в себе entity.ClientError, то использует его в качестве сообщения об ошибке, в противном
// случае приводит стандартное описание ошибки по её коду для избежания возможных утечек.
// Попутно логируются все пятисотки
func NewError(ctx echo.Context, status int, message string, err error) *echo.HTTPError {
	httpError := &echo.HTTPError{Code: status}
	httpError.Message = message
	if status >= 500 {
		httpError.Internal = err
		ctx.Logger().Error(GetErrMsgFromContext(ctx, err))
	}
	return httpError
}

// ServerErrorMsg описывает структуру ошибки для логирования
type ServerErrorMsg struct {
	URL         url.URL
	UserAgent   string
	Header      http.Header
	Method      string
	QueryParams url.Values
	ClientIP    string
	RequestBody string
	Error       error
}

func GetErrMsgFromContext(ctx echo.Context, err error) ServerErrorMsg {
	body, e := io.ReadAll(ctx.Request().Body)
	if e != nil {
		body = []byte("не удалось получить тело запроса")
	}

	return ServerErrorMsg{
		URL:         *ctx.Request().URL,
		UserAgent:   ctx.Request().UserAgent(),
		Header:      ctx.Request().Header,
		Method:      ctx.Request().Method,
		QueryParams: ctx.QueryParams(),
		ClientIP:    ctx.RealIP(),
		RequestBody: string(body),
		Error:       err,
	}
}

var ErrBadJSON = errors.New("невалидный JSON")
