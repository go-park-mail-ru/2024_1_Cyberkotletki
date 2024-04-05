package dto

type UserUpdate struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserProfile struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Rating int    `json:"rating"`
	Avatar int    `json:"avatar"`
}
