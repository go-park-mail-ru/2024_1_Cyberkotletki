package repository

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_user.go
type User interface {
	HasUser(entity.User) bool
	AddUser(entity.User) (*entity.User, error)
	GetUserByEmail(string) (*entity.User, error)
}
