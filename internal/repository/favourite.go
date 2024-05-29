package repository

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_favourite.go
type Favourite interface {
	// CreateFavourite добавление в избранное. Если уже в избранном, то ошибка не возвращается (идемпотентный метод).
	// Возможные ошибки:
	// ErrFavouriteContentNotFound - контент не найден
	CreateFavourite(userID, contentID int, category string) error
	// DeleteFavourite удаление из избранного.
	// Возвращает ошибку ErrFavouriteNotFound, если контент не найден в избранном
	DeleteFavourite(userID, contentID int) error
	// GetFavourites получение избранного контента пользователя.
	// Возвращает ошибку ErrFavouriteUserNotFound, если пользователь не найден
	GetFavourites(userID int) ([]*entity.Favourite, error)
	// GetFavourite получение статуса контента в избранном
	GetFavourite(userID, contentID int) (*entity.Favourite, error)
}

var (
	ErrFavouriteNotFound        = errors.New("избранное не найдено")
	ErrFavouriteContentNotFound = errors.New("контент не найден")
	ErrFavouriteUserNotFound    = errors.New("пользователь не найден")
)
