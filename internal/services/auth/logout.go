package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/session"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"time"
)

func Logout(cookie string) *exc.Exception {
	if ok := session.SessionsDB.DeleteSession(cookie); ok {
		return nil
	}
	return &exc.Exception{
		When:  time.Now(),
		What:  "Не авторизован",
		Layer: exc.Service,
		Type:  exc.Forbidden,
	}
}
