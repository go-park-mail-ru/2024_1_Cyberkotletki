package repository

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_user.go
type User interface {
	AddUser(*entity.User) (*entity.User, error)
	GetUser(map[string]interface{}) (*entity.User, error)
	UpdateUser(params map[string]interface{}, values map[string]interface{}) error
}
