package repository

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
)

type Film interface {
	GetFilm(id int) (*entity.Film, error)
	GetFilmsByGenre(genreID int) ([]entity.Film, error)
}

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_content.go
type Content interface {
	Film
	// todo: content.Series для сериалов
}
