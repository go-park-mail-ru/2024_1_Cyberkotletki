package repository

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_user.go
type User interface {
	AddUser(*entity.User) (*entity.User, error)
	GetUserByID(userID int) (*entity.User, error)
	GetUserByEmail(userEmail string) (*entity.User, error)
	UpdateUser(user *entity.User) error
}

var (
	ErrUserAlreadyExists = errors.New("пользователь с такой почтой уже существует")
	ErrUserNotFound      = errors.New("пользователь не найден")
	ErrUserIncorrectData = errors.New("поля заполнены некорректно")
)
