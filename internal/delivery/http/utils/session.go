package utils

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func CreateSession(ctx echo.Context, authUC usecase.Auth, userID int) error {
	session, err := authUC.CreateSession(userID)
	if err != nil {
		return NewError(ctx, http.StatusInternalServerError, err)
	}
	SessionSet(
		ctx,
		session,
		time.Now().Add(time.Duration(ctx.Get("params").(config.Config).Auth.SessionAliveTime)*time.Second),
	)
	return nil
}

func SessionSet(ctx echo.Context, value string, expires time.Time) {
	var cookie http.Cookie
	if params, ok := ctx.Get("params").(config.Config); ok {
		cookie = http.Cookie{
			Name:     "session",
			Value:    value,
			Expires:  expires,
			HttpOnly: true,
			Secure:   params.HTTP.SecureCookies,
			Path:     "/",
		}
	} else {
		cookie = http.Cookie{
			Name:     "session",
			Value:    value,
			Expires:  expires,
			HttpOnly: true,
			Path:     "/",
		}
	}
	ctx.SetCookie(&cookie)
}

func GetUserIDFromSession(ctx echo.Context, authUC usecase.Auth) (int, error) {
	session, err := ctx.Cookie("session")
	if err != nil {
		return -1, entity.NewClientError("отсутствует cookies с сессией", entity.ErrUnauthorized)
	}
	userID, err := authUC.GetUserIDBySession(session.Value)
	if err != nil {
		return -1, entity.NewClientError("необходима авторизация", entity.ErrUnauthorized)
	}
	return userID, nil
}
