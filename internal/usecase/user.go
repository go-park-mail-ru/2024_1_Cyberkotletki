package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"io"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_user.go
type User interface {
	// Register регистрация пользователя.
	// Возможные ошибки:
	// ErrUserAlreadyExists - пользователь с такой почтой уже существует
	// UserIncorrectDataError - некорректные данные
	Register(*dto.Register) (int, error)
	// Login авторизация пользователя.
	// Возможные ошибки:
	// ErrUserNotFound - пользователь не найден
	// UserIncorrectDataError - некорректные данные
	Login(*dto.Login) (int, error)
	// UpdateAvatar обновление аватара пользователя.
	// Возможные ошибки:
	// ErrUserNotFound - пользователь не найден
	// UserIncorrectDataError - некорректные данные
	UpdateAvatar(userID int, reader io.ReadSeeker) error
	// UpdateInfo обновление информации о пользователе.
	// Возможные ошибки:
	// ErrUserNotFound - пользователь не найден
	// UserIncorrectDataError - некорректные данные
	UpdateInfo(userID int, update *dto.UserUpdate) error
	// UpdatePassword обновление пароля пользователя.
	// Возможные ошибки:
	// ErrUserNotFound - пользователь не найден
	// UserIncorrectDataError - некорректные данные
	UpdatePassword(userID int, update *dto.UpdatePassword) error
	// GetUser получение профиля пользователя.
	// Возможные ошибки:
	// ErrUserNotFound - пользователь не найден
	GetUser(userID int) (*dto.UserProfile, error)
}

// UserIncorrectDataError это ошибка некорректных данных пользователя
// Err содержит точную природу ошибки
type UserIncorrectDataError struct {
	Err error
}

func (u UserIncorrectDataError) Error() string {
	return u.Err.Error()
}

var (
	ErrUserAlreadyExists = errors.New("пользователь с такой почтой уже существует")
	ErrUserNotFound      = errors.New("пользователь не найден")
)
