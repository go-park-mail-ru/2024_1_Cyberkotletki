package repository

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_user.go
type User interface {
	// AddUser добавляет пользователя в базу данных.
	// Возможные ошибки:
	// ErrUserAlreadyExists - пользователь с такой почтой уже существует
	// ErrUserIncorrectData - поля заполнены некорректно
	AddUser(email string, passwordHash, passwordSalt []byte) (*entity.User, error)
	// GetUserByID возвращает пользователя по его id.
	// Возможные ошибки:
	// ErrUserNotFound - пользователь не найден
	GetUserByID(userID int) (*entity.User, error)
	// GetUserByEmail возвращает пользователя по его почте.
	// Возможные ошибки:
	// ErrUserNotFound - пользователь не найден
	GetUserByEmail(userEmail string) (*entity.User, error)
	// UpdateUser обновляет данные пользователя.
	UpdateUser(user *entity.User) error
}

var (
	ErrUserAlreadyExists = errors.New("пользователь с такой почтой уже существует")
	ErrUserNotFound      = errors.New("пользователь не найден")
	ErrUserIncorrectData = errors.New("поля заполнены некорректно")
)
