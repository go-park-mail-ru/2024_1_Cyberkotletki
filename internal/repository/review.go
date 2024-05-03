package repository

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_review.go
type Review interface {
	// GetLatestReviews возвращает последние рецензии
	GetLatestReviews(limit int) ([]*entity.Review, error)
	// AddReview добавляет рецензию
	// Возможные ошибки:
	// ErrReviewViolation - контент с таким id не существует, либо такого пользователя не существует
	// ErrReviewBadRequest - некорректные данные для создания рецензии
	// ErrReviewAlreadyExists - рецензия уже существует
	AddReview(review *entity.Review) (*entity.Review, error)
	// GetReviewByID возвращает рецензию по id
	// Возможные ошибки:
	// ErrReviewNotFound - рецензия не найдена
	GetReviewByID(id int) (*entity.Review, error)
	// GetReviewsCountByContentID возвращает количество рецензий по id контента
	GetReviewsCountByContentID(contentID int) (int, error)
	// GetReviewsByContentID возвращает рецензии по id контента
	GetReviewsByContentID(contentID, page, limit int) ([]*entity.Review, error)
	// UpdateReview обновляет рецензию
	// Возможные ошибки:
	// ErrReviewNotFound - рецензия не найдена
	// ErrReviewBadRequest - некорректные данные для обновления рецензии
	UpdateReview(review *entity.Review) error
	// DeleteReviewByID удаляет рецензию по id
	// Возможные ошибки:
	// ErrReviewNotFound - рецензия не найдена
	DeleteReviewByID(id int) error
	// GetReviewsCountByAuthorID возвращает количество рецензий по id автора
	GetReviewsCountByAuthorID(authorID int) (int, error)
	// GetReviewsByAuthorID возвращает рецензии по id автора
	GetReviewsByAuthorID(authorID, page, limit int) ([]*entity.Review, error)
	// GetContentReviewByAuthor возвращает рецензию по id автора и id контента
	// Возможные ошибки:
	// ErrReviewNotFound - рецензия не найдена
	GetContentReviewByAuthor(authorID, contentID int) (*entity.Review, error)
	// VoteReview оценивает рецензию
	// Возможные ошибки:
	// ErrReviewNotFound - рецензия не найдена
	// ErrReviewVoteAlreadyExists - пользователь уже оценил рецензию
	VoteReview(reviewID, userID int, like bool) error
	// UnVoteReview убирает оценку рецензии
	// Возможные ошибки:
	// ErrReviewVoteNotFound - оценка не найдена
	UnVoteReview(reviewID, userID int) error
	// IsVotedByUser проверяет, оценил ли пользователь рецензию.
	// Возвращает 1, если оценил положительно, -1, если отрицательно, 0, если не оценил
	IsVotedByUser(reviewID, userID int) (int, error)
}

var (
	ErrReviewViolation         = errors.New("контент с таким id не существует, либо такого пользователя не существует")
	ErrReviewBadRequest        = errors.New("некорректные данные для создания рецензии")
	ErrReviewAlreadyExists     = errors.New("рецензия уже существует")
	ErrReviewNotFound          = errors.New("рецензия не найдена")
	ErrReviewVoteAlreadyExists = errors.New("пользователь уже оценил рецензию")
	ErrReviewVoteNotFound      = errors.New("оценка не найдена")
)
