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
	// GetPersonByID возвращает персону по ее ID
	// Если персона не найдена, возвращает ErrPersonNotFound
	GetPersonByID(id int) (*dto.Person, error)
}

var (
	ErrContentNotFound = errors.New("контент не найден")
	ErrPersonNotFound  = errors.New("персона не найдена")
)
