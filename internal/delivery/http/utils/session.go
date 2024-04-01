package utils

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func SessionSet(ctx echo.Context, value string, expires time.Time) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    value,
		Expires:  expires,
		HttpOnly: true,
		Secure:   ctx.Get("params").(config.Config).HTTP.SecureCookies,
		Path:     "/",
	}
	ctx.SetCookie(&cookie)
}
