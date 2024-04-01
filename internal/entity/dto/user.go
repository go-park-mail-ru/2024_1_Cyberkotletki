package dto

type UserUpdate struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserProfile struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar int    `json:"avatar"`
}
