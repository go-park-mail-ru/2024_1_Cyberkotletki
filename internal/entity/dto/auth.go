package dto

type Register struct {
	Email    string `json:"email"    example:"email@email.com"  format:"string"`
	Password string `json:"password" example:"SecretPassword1!" format:"string"`
}

type Login struct {
	Login    string `json:"login"    example:"email@email.com"  format:"string"`
	Password string `json:"password" example:"SecretPassword1!" format:"string"`
}
