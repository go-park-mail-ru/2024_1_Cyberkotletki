package usecase

import (
	"errors"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_ongoing_content.go
type OngoingContent interface {
	// GetOngoingContentByContentID возвращает контент календаря релизов по id контента
	// Если контент не найден, возвращает ErrOngoingContentNotFound
	GetOngoingContentByContentID(contentID int) (*dto.PreviewOngoingContentCardVertical, error)
	// GetNearestOngoings возвращает ближайшие релизы
	// Если контент не найден, возвращает ErrOngoingContentNotFound
	GetNearestOngoings(limit int) ([]*dto.PreviewOngoingContentCardVertical, error)
	// GetOngoingContentByMonthAndYear возвращает релизы по месяцу и году
	// Если контент не найден, возвращает ErrOngoingContentNotFound
	GetOngoingContentByMonthAndYear(month, year int) ([]*dto.PreviewOngoingContentCardVertical, error)
	// GetAllReleaseYears возвращает все года релизов
	// Если годы не найдены, возвращает ErrOngoingContentYearsNotFound
	GetAllReleaseYears() ([]int, error)
	// IsConentFinished возвращает true, если контент вышел
	// Возвращает 1, если вышел, 0, если нет
	IsOngoingConentFinished(contentID int) (bool, error)
}

var (
	ErrOngoingContentNotFound      = errors.New("контент календаря релизов не найден")
	ErrOngoingContentYearsNotFound = errors.New("года релизов не найдены")
)
