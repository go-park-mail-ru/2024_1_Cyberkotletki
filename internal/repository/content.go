package repository

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_content.go
type Content interface {
	GetContent(id int) (*entity.Content, error)
	GetPreviewContent(id int) (*entity.Content, error)
	GetPerson(id int) (*entity.Person, error)
	GetPersonRoles(id int) ([]entity.Content, error)
}
