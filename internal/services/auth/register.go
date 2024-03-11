package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/session"
	userDB "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/user"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/user"
)

type RegisterData struct {
	Email    string `json:"email" example:"email@email.com" format:"string"`
	Password string `json:"password" example:"SecretPassword1!" format:"string"`
}

// Register создаёт пользователя в системе и сразу же возвращает сессию
func Register(registerData RegisterData) (string, error) {
	us := user.NewUserEmpty()
	if err := user.ValidateEmail(registerData.Email); err != nil {
		return "", err
	}
	if err := user.ValidatePassword(registerData.Password); err != nil {
		return "", err
	}
	salt, hash := user.HashPassword(registerData.Password)
	us.Email = registerData.Email
	us.PasswordHash = hash
	us.PasswordSalt = salt
	if us, err := userDB.UsersDatabase.AddUser(*us); err != nil {
		return "", err
	} else {
		return session.SessionsDB.NewSession(us.Id), nil
	}
}
