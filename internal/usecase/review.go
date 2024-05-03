package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_review.go
type Review interface {
	GetLatestReviews(count int) (*dto.ReviewResponseList, error)
	GetUserReviews(userID, count, page int) (*dto.ReviewResponseList, error)
	GetContentReviews(contentID, count, page int) (*dto.ReviewResponseList, error)
	GetReview(reviewID int) (*dto.ReviewResponse, error)
	GetContentReviewByAuthor(authorID, contentID int) (*dto.ReviewResponse, error)
	CreateReview(create dto.ReviewCreate) (*dto.ReviewResponse, error)
	EditReview(update dto.ReviewUpdate) (*dto.ReviewResponse, error)
	DeleteReview(reviewID, userID int) error
	VoteReview(userID, reviewID int, vote bool) error
	UnVoteReview(userID, reviewID int) error
	IsVotedByUser(userID, reviewID int) (int, error)
}

// ReviewErrorIncorrectData ошибка некорректных данных
type ReviewErrorIncorrectData struct {
	Err error
}

func (e ReviewErrorIncorrectData) Error() string {
	return e.Err.Error()
}

var (
	ErrReviewNotFound      = errors.New("рецензия не найдена")
	ErrReviewAlreadyExists = errors.New("рецензия уже существует")
	ErrReviewForbidden     = errors.New("недостаточно прав для выполнения операции")
)
