package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/session"
	userDB "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/user"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/user"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
)

type RegisterData struct {
	Email    string `json:"email" example:"email@email.com" format:"string"`
	Password string `json:"password" example:"SecretPassword1!" format:"string"`
}

// Register создаёт пользователя в системе и сразу же возвращает сессию
func Register(registerData RegisterData) (string, error) {
	us := user.NewUserEmpty()
	if err := us.ValidateEmail(registerData.Email); err != nil {
		return "", err
	}
	if err := us.ValidatePassword(registerData.Password); err != nil {
		return "", err
	}
	if hash := user.HashPassword(registerData.Password); hash == "" {
		return "", exc.New(exc.Service, exc.Internal, "произошла непредвиденная ошибка", "не удалось получить хэш пароля")
	} else {
		us.Email = registerData.Email
		us.PasswordHash = hash
	}
	if us, err := userDB.UsersDatabase.AddUser(*us); err != nil {
		return "", err
	} else {
		return session.SessionsDB.NewSession(us.Id), nil
	}
}
