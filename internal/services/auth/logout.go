package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/session"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
)

func Logout(cookie string) error {
	if ok := session.SessionsDB.DeleteSession(cookie); ok {
		return nil
	}
	return exc.New(exc.Service, exc.Untyped, "сессия недействительна")
}
