package dto

type UserUpdate struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserProfile struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Rating int    `json:"rating"`
	Avatar string `json:"avatar"`
}
