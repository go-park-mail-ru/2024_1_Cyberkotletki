package usecase

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_review.go
type Review interface {
	GetLatestReviews(count int) (*[]dto.ReviewResponse, error)
	GetUserReviews(userID, count, page int) (*[]dto.ReviewResponse, error)
	GetContentReviews(contentID, count, page int) (*[]dto.ReviewResponse, error)
	GetReview(reviewID int) (*dto.ReviewResponse, error)
	CreateReview(create dto.ReviewCreate) (*dto.ReviewResponse, error)
	EditReview(update dto.ReviewUpdate) (*dto.ReviewResponse, error)
	DeleteReview(reviewID, userID int) error
	LikeReview(userID, reviewID int) error
	DislikeReview(userID, reviewID int) error
	UnlikeReview(userID, reviewID int) error
	IsLikedByUser(userID, reviewID int) (int, error)
	GetUserRating(userID int) (int, error)
	GetContentRating(contentID int) (int, error)
}
