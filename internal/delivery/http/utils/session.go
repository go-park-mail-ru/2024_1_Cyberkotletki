package utils

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type SessionManager struct {
	authUC           usecase.Auth
	sessionAliveTime int
	secureCookies    bool
}

func NewSessionManager(authUC usecase.Auth, sessionAliveTime int, secureCookies bool) *SessionManager {
	return &SessionManager{
		authUC:           authUC,
		sessionAliveTime: sessionAliveTime,
		secureCookies:    secureCookies,
	}
}

func (s SessionManager) CreateSession(ctx echo.Context, authUC usecase.Auth, userID int) error {
	session, err := authUC.CreateSession(userID)
	if err != nil {
		return err
	}
	s.SessionSet(
		ctx,
		session,
		time.Now().Add(time.Duration(s.sessionAliveTime)*time.Second),
	)
	return nil
}

func (s SessionManager) SessionSet(ctx echo.Context, value string, expires time.Time) {
	cookie := http.Cookie{
		Name:     "session",
		Value:    value,
		Expires:  expires,
		HttpOnly: true,
		Secure:   s.secureCookies,
		Path:     "/",
	}
	ctx.SetCookie(&cookie)
}

// GetUserIDFromSession возвращает ID пользователя из сессии.
// В случае ошибки возвращает ErrUnauthorized
func GetUserIDFromSession(ctx echo.Context, authUC usecase.Auth) (int, error) {
	session, err := ctx.Cookie("session")
	if err != nil {
		return -1, ErrUnauthorized
	}
	userID, err := authUC.GetUserIDBySession(session.Value)
	if err != nil {
		return -1, ErrUnauthorized
	}
	return userID, nil
}

var (
	ErrUnauthorized = errors.New("необходима авторизация")
)
