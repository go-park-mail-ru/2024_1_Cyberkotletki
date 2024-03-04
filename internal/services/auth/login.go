package auth

import exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"

type LoginData struct {
	Login    string `json:"login" example:"email@email.com" format:"string"`
	Password string `json:"password" example:"SecretPassword1!" format:"string"`
}

func Login(loginData LoginData) (string, *exc.Exception) {
	// todo
	return "key", nil
}
