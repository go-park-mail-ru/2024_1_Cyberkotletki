package tmpdb

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"sync"
)

type ContentDB struct {
	sync.RWMutex
	DB map[int]entity.Film
	// filmsLastId atomic.Int64
}

func NewContentRepository() repository.Film {
	f := &ContentDB{
		DB: make(map[int]entity.Film),
	}
	return f
}

func (f *ContentDB) GetFilm(id int) (*entity.Film, error) {
	f.Lock()
	defer f.Unlock()
	if filmObj, ok := f.DB[id]; ok {
		return &filmObj, nil
	}
	return nil, entity.NewClientError("фильм с таким id не найден", entity.ErrNotFound)
}

// GetFilmsByGenre возвращает фильмы определенного жанра
func (f *ContentDB) GetFilmsByGenre(genreID int) ([]entity.Film, error) {
	f.Lock()
	defer f.Unlock()

	var films []entity.Film
	for _, film := range f.DB {
		for _, genreObj := range film.Content.Genres {
			if genreObj.ID == genreID {
				films = append(films, film)
				break
			}
		}
	}
	if films == nil {
		return nil, entity.NewClientError("фильмы с таким жанром не найдены", entity.ErrNotFound)
	}
	return films, nil
}
