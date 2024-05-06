package repository

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_search.go
type Search interface {
	SearchContent(query string) ([]entity.Content, error)
	SearchPerson(query string) ([]entity.Person, error)
}
