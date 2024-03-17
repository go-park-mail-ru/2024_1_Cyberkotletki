package usecase

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/DTO"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_collections.go
type Collections interface {
	GetCompilation(genre string) (*DTO.Compilation, error)
	GetGenres() (*DTO.Genres, error)
}
