package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/session"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"time"
)

func IsAuth(cookie string) (int, *exc.Exception) {
	if userId, ok := session.SessionsDB.CheckSession(cookie); ok {
		return userId, nil
	}
	return -1, &exc.Exception{
		When:  time.Now(),
		What:  "Не авторизован",
		Layer: exc.Service,
		Type:  exc.Forbidden,
	}
}
