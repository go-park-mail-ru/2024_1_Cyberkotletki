package usecase

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_review.go
type Review interface {
	GetGlobalReviews(int) (*[]dto.Review, error)
	GetUserReviews(int, int) (*[]dto.Review, error)
	GetContentReviews(int) (*[]dto.Review, error)
	GetReview(int) (*dto.Review, error)
	CreateReview(dto.Review) (*dto.Review, error)
	EditReview(dto.Review) (*dto.Review, error)
	DeleteReview(int) error
}
