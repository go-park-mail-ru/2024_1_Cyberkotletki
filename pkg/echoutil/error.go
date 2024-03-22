package echoutil

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/url"
)

// NewError возвращает *echo.HTTPError, чтобы middleware echo автоматически конвертировал её в JSON-ошибку. Если
// ошибка содержит в себе entity.ClientError, то использует его в качестве сообщения об ошибке, в противном
// случае приводит стандартное описание ошибки по её коду для избежания возможных утечек.
// Попутно логируются все пятисотки
func NewError(ctx echo.Context, status int, err error) *echo.HTTPError {
	httpError := &echo.HTTPError{Code: status}
	var e entity.ClientError
	if errors.As(err, &e) {
		httpError.Message = e.Error()
	} else {
		// клиентской ошибки нет, поэтому отобразим стандартное описание
		httpError.Message = http.StatusText(status)
	}
	if status >= 500 {
		ctx.Logger().Error(GetErrMsgFromContext(ctx))
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
}

func GetErrMsgFromContext(ctx echo.Context) ServerErrorMsg {
	// todo: гипотетически тело запроса может быть большим, лучше заасинхронить чтение
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
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
	}
}

var ErrBadJSON = errors.New("невалидный JSON")
