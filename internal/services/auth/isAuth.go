package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/session"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
)

func IsAuth(cookie string) (int, error) {
	if userId, ok := session.SessionsDB.CheckSession(cookie); ok {
		return userId, nil
	}
	return -1, exc.New(exc.Service, exc.Forbidden, "не авторизован")
}
