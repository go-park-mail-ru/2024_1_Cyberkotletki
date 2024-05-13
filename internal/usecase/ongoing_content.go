package usecase

import (
	"errors"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_ongoing_content.go
type OngoingContent interface {
	// GetOngoingContentByContentID возвращает контент календаря релизов по id контента
	// Если контент не найден, возвращает ErrOngoingContentNotFound
	GetOngoingContentByContentID(contentID int) (*dto.PreviewOngoingContent, error)
	// GetNearestOngoings возвращает ближайшие релизы
	// Если контент не найден, возвращает ErrOngoingContentNotFound
	GetNearestOngoings(limit int) (*dto.PreviewOngoingContentList, error)
	// GetOngoingContentByMonthAndYear возвращает релизы по месяцу и году
	// Если контент не найден, возвращает ErrOngoingContentNotFound
	GetOngoingContentByMonthAndYear(month, year int) (*dto.PreviewOngoingContentList, error)
	// GetAllReleaseYears возвращает все года релизов
	// Если годы не найдены, возвращает ErrOngoingContentYearsNotFound
	GetAllReleaseYears() (*dto.ReleaseYearsResponse, error)
	// IsOngoingContentFinished возвращает true, если контент вышел.
	// Возвращает 1, если вышел, 0, если нет
	IsOngoingContentFinished(contentID int) (*dto.IsOngoingContentFinishedResponse, error)
}

var (
	ErrOngoingContentNotFound      = errors.New("контент календаря релизов не найден")
	ErrOngoingContentYearsNotFound = errors.New("года релизов не найдены")
)
