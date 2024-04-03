package utils

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
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

func GetUserIDFromSession(ctx echo.Context, authUC usecase.Auth) (int, error) {
	session, err := ctx.Cookie("session")
	if err != nil {
		return -1, entity.NewClientError("необходима авторизация", entity.ErrUnauthorized)
	}
	userID, err := authUC.GetUserIDBySession(session.Value)
	if err != nil {
		return -1, entity.NewClientError("необходима авторизация", entity.ErrUnauthorized)
	}
	return userID, nil
}
