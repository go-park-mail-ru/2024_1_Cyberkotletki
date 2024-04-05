package repository

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_review.go
type Review interface {
	AddReview(review *entity.Review) (*entity.Review, error)
	GetReviewByID(id int) (*entity.Review, error)
	GetReviewsByContentID(contentID, page, limit int) ([]*entity.Review, error)
	UpdateReview(review *entity.Review) (*entity.Review, error)
	DeleteReviewByID(id int) error
	GetReviewsByAuthorID(authorID, page, limit int) ([]*entity.Review, error)
	GetAuthorRating(authorID int) (int, error)
	GetLatestReviews(limit int) ([]*entity.Review, error)
	LikeReview(reviewID, userID int, like bool) error
	UnlikeReview(reviewID, userID int) error
	IsLikedByUser(reviewID, userID int) (int, error)
	DeleteReview(reviewID int) error
	GetContentRating(reviewID int) (int, error)
}
