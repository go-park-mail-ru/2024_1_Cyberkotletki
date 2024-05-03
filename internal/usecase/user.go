package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"io"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_user.go
type User interface {
	Register(*dto.Register) (int, error)
	Login(*dto.Login) (int, error)
	UpdateAvatar(userID int, reader io.Reader) error
	UpdateInfo(userID int, update *dto.UserUpdate) error
	UpdatePassword(userID int, update *dto.UpdatePassword) error
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
