package repository

import (
	"errors"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_ongoing_content.go
type OngoingContent interface {
	// GetOngoingContentByID возвращает контент календаря релизов по id календаря релизов
	// Если контент не найден, возвращает ErrOngoingContentNotFound
	GetOngoingContentByID(id int) (*entity.OngoingContent, error)
	// GetNearestOngoings возвращает ближайшие релизы
	// Если контент не найден, возвращает ErrOngoingContentNotFound
	GetNearestOngoings(limit int) ([]*entity.OngoingContent, error)
	// GetOngoingContentByMonthAndYear возвращает релизы по месяцу и году
	// Если контент не найден, возвращает ErrOngoingContentNotFound
	GetOngoingContentByMonthAndYear(month, year int) ([]*entity.OngoingContent, error)
	// GetAllReleaseYears возвращает все года релизов
	// Если годы не найдены, возвращает ErrOngoingContentYearsNotFound
	GetAllReleaseYears() ([]int, error)
	// IsOngoingContentFinished возвращает true, если контент вышел.
	// Возвращает 1, если вышел, 0, если нет
	IsOngoingContentFinished(contentID int) (bool, error)
}

var (
	ErrOngoingContentNotFound      = errors.New("контент календаря релизов не найден")
	ErrOngoingContentYearsNotFound = errors.New("года релизов не найдены")
)
