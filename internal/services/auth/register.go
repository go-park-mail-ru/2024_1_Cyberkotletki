package auth

import exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"

type RegisterData struct {
	Email    string `json:"email" example:"email@email.com" format:"string"`
	Password string `json:"password" example:"SecretPassword1!" format:"string"`
}

func Register(registerData RegisterData) (string, *exc.Exception) {
	// todo
	return "key", nil
}
