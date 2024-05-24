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
}

var (
	ErrContentNotFound = errors.New("контент не найден")
	ErrPersonNotFound  = errors.New("персона не найдена")
)
