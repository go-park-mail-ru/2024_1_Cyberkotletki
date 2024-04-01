package usecase

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_user.go
type User interface {
	Register(*dto.Register) (int, error)
	Login(*dto.Login) (int, error)
	UpdateAvatar(userID, uploadID int) error
	UpdateInfo(userID int, update *dto.UserUpdate) error
	UpdatePassword(userID int, update *dto.UpdatePassword) error
	GetUser(userID int) (*dto.UserProfile, error)
	// TODO см. таск 2.2
	// GetUserRating(userID int) (int, error)
	// TODO см. таск 2.2
	// GetUserReviews(userID int) (int, error)
}
