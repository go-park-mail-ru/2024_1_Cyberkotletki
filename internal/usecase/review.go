package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_review.go
type Review interface {
	// GetLatestReviews получение последних рецензий
	GetLatestReviews(count int) (*dto.ReviewResponseList, error)
	// GetUserReviews получение рецензий пользователя
	GetUserReviews(userID, count, page int) (*dto.ReviewResponseList, error)
	// GetContentReviews получение рецензий на контент
	GetContentReviews(contentID, count, page int) (*dto.ReviewResponseList, error)
	// GetReview получение рецензии.
	// Возвращает ошибку ErrReviewNotFound, если рецензия не найдена
	GetReview(reviewID int) (*dto.ReviewResponse, error)
	// GetContentReviewByAuthor получение рецензии на контент от автора.
	// Возвращает ошибку ErrReviewNotFound, если рецензия не найдена
	GetContentReviewByAuthor(authorID, contentID int) (*dto.ReviewResponse, error)
	// CreateReview создание рецензии.
	// Возможные ошибки:
	// ErrReviewContentNotFound - контент не найден
	// ErrReviewAlreadyExists - рецензия уже существует
	// ReviewErrorIncorrectData - некорректные данные
	CreateReview(create dto.ReviewCreate) (*dto.ReviewResponse, error)
	// EditReview редактирование рецензии.
	// Возможные ошибки:
	// ErrReviewNotFound - рецензия не найдена
	// ErrReviewForbidden - недостаточно прав для выполнения операции
	// ReviewErrorIncorrectData - некорректные данные
	EditReview(update dto.ReviewUpdate) (*dto.ReviewResponse, error)
	// DeleteReview удаление рецензии.
	// Возможные ошибки:
	// ErrReviewNotFound - рецензия не найдена
	// ErrReviewForbidden - недостаточно прав для выполнения операции
	DeleteReview(reviewID, userID int) error
	// VoteReview голосование за рецензию.
	// Возможные ошибки:
	// ErrReviewNotFound - рецензия не найдена
	// ErrReviewVoteAlreadyExists - голос уже учтен
	VoteReview(userID, reviewID int, vote bool) error
	// UnVoteReview отмена голоса за рецензию.
	// Возможные ошибки:
	// ErrReviewVoteNotFound - рецензия не найдена
	UnVoteReview(userID, reviewID int) error
	// IsVotedByUser проверка, голосовал ли пользователь за рецензию
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
	ErrReviewNotFound        = errors.New("рецензия не найдена")
	ErrReviewContentNotFound = errors.New("контент не найден")
	ErrReviewAlreadyExists   = errors.New("рецензия уже существует")
	ErrReviewForbidden       = errors.New("недостаточно прав для выполнения операции")

	ErrReviewVoteNotFound      = errors.New("голос не найден")
	ErrReviewVoteAlreadyExists = errors.New("голос уже учтен")
)
