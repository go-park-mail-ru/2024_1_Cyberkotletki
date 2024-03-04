package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/session"
	userDB "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/user"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/user"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type RegisterData struct {
	Email    string `json:"email" example:"email@email.com" format:"string"`
	Password string `json:"password" example:"SecretPassword1!" format:"string"`
}

func Register(registerData RegisterData) (string, *exc.Exception) {
	us := user.NewUserEmpty()
	if err := us.ValidateEmail(registerData.Email); err != nil {
		return "", err
	}
	if err := us.ValidatePassword(registerData.Password); err != nil {
		return "", err
	}
	if hash, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), 14); err != nil {
		return "", &exc.Exception{
			When:  time.Now(),
			What:  "Внутренняя ошибка сервера",
			Layer: exc.Server,
			Type:  exc.Untyped,
		}
	} else {
		us.Email = registerData.Email
		us.PasswordHash = string(hash)
	}
	if us, err := userDB.UsersDatabase.AddUser(*us); err != nil {
		return "", err
	} else {
		return session.SessionsDB.NewSession(us.Id), nil
	}
}
