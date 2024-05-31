package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_content.go
type Content interface {
	// GetContentByID возвращает контент по его ID
	// Если контент не найден, возвращает ErrContentNotFound
	GetContentByID(id int) (*dto.Content, error)
	// GetPreviewPersonByID возвращает персону по ее ID, но только с минимальным набором полей
	// Если персона не найдена, возвращает ErrPersonNotFound
	GetPreviewPersonByID(id int) (*dto.PersonPreviewWithPhoto, error)
	// GetPersonByID возвращает персону по ее ID
	// Если персона не найдена, возвращает ErrPersonNotFound
	GetPersonByID(id int) (*dto.Person, error)
	// GetPreviewContentByID возвращает контент по его ID, но только с минимальным набором полей
	// Если контент не найден, возвращает ErrContentNotFound
	GetPreviewContentByID(id int) (*dto.PreviewContent, error)
	// GetNearestOngoings возвращает 10 ближайших релизов
	GetNearestOngoings() (*dto.PreviewOngoingContentList, error)
	// GetOngoingContentByMonthAndYear возвращает релизы по месяцу и году
	GetOngoingContentByMonthAndYear(month, year int) (*dto.PreviewOngoingContentList, error)
	// GetAllOngoingsYears возвращает все года релизов
	GetAllOngoingsYears() (*dto.ReleaseYearsResponse, error)
	// IsOngoingContentReleased возвращает true, если контент вышел
	// Если контент не найден, возвращает ErrContentNotFound
	IsOngoingContentReleased(contentID int, releasedCh chan<- bool, errCh chan<- error)
	// SetReleasedState устанавливает состояние релиза
	SetReleasedState(secretKey string, contentID int, isReleased bool) error
	// SubscribeOnContent подписывает пользователя на контент
	// Если контент не найден, возвращает ErrContentNotFound
	// Если пользователь не найден, возвращает ErrUserNotFound
	SubscribeOnContent(userID, contentID int) error
	// UnsubscribeFromContent отписывает пользователя от контента
	// Если контент не найден, возвращает ErrContentNotFound
	// Если пользователь не найден, возвращает ErrUserNotFound
	UnsubscribeFromContent(userID, contentID int) error
	// GetSubscribedContentIDs возвращает id контентов, на которые подписан пользователь
	// Если пользователь не найден, возвращает ErrUserNotFound
	GetSubscribedContentIDs(userID int) (*dto.SubscriptionsResponse, error)
	// GetAvailableToWatch возвращает контент, который доступен для просмотра
	GetAvailableToWatch(page, limit int) (*dto.ContentPreviewList, error)
}

var (
	ErrContentNotFound         = errors.New("контент не найден")
	ErrPersonNotFound          = errors.New("персона не найдена")
	ErrContentInvalidSecretKey = errors.New("неверный секретный ключ")
)
