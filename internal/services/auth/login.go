package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/session"
	userDB "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/user"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
)

type LoginData struct {
	Login    string `json:"login" example:"email@email.com" format:"string"`
	Password string `json:"password" example:"SecretPassword1!" format:"string"`
}

func Login(loginData LoginData) (string, error) {
	us, err := userDB.UsersDatabase.GetUserByEmail(loginData.Login)
	if err != nil {
		return "", err
	}
	if us.CheckPassword(loginData.Password) {
		return session.SessionsDB.NewSession(us.Id), nil
	}
	return "", exc.New(exc.Service, exc.Forbidden, "неверный пароль")
}
