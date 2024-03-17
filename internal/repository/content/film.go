package content

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
)

type Film interface {
	GetFilm(id int) (*entity.Film, error)
	GetFilmsByGenre(genreId int) ([]entity.Film, error)
}
