package repository

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_content.go
type Content interface {
	// GetContent возвращает контент по его id
	// Если контент не найден, возвращает ErrContentNotFound
	GetContent(id int) (*entity.Content, error)
	// GetPreviewContent возвращает контент по его id, но только с минимальным набором полей
	// Если контент не найден, возвращает ErrContentNotFound
	GetPreviewContent(id int) (*entity.Content, error)
	// GetPerson возвращает роли контента
	// Если контент не найден, возвращает ErrContentNotFound
	GetPerson(id int) (*entity.Person, error)
	// GetPersonRoles возвращает роли персоны
	GetPersonRoles(personID int) ([]entity.PersonRole, error)
	// GetSimilarContent возвращает похожий контент
	GetSimilarContent(id int) ([]entity.Content, error)
	// GetNearestOngoings возвращает ближайшие релизы
	GetNearestOngoings(limit int) ([]int, error)
	// GetOngoingContentByMonthAndYear возвращает релизы по месяцу и году
	GetOngoingContentByMonthAndYear(month, year int) ([]int, error)
	// GetAllOngoingsYears возвращает все года релизов
	GetAllOngoingsYears() ([]int, error)
	// IsOngoingContentReleased возвращает true, если контент вышел
	// Если контент не найден, возвращает ErrContentNotFound
	IsOngoingContentReleased(contentID int) (bool, error)
	// SetReleasedState устанавливает состояние релиза
	// Если контент не найден, возвращает ErrContentNotFound
	SetReleasedState(contentID int, isReleased bool) error
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
	GetSubscribedContentIDs(userID int) ([]int, error)
}

var (
	ErrContentNotFound = errors.New("контент не найден")
	ErrPersonNotFound  = errors.New("персона не найдена")
)
