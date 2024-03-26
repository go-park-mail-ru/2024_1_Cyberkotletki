package service

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
)

type CollectionsService struct {
	contentRepo repository.Content
}

func NewCollectionsService(contentRepo repository.Content) usecase.Collections {
	return &CollectionsService{
		contentRepo: contentRepo,
	}
}

func (c CollectionsService) GetCompilation(genre string) (*dto.Compilation, error) {
	var genreID int
	genres := map[string]int{
		"drama":    1,
		"action":   2,
		"comedian": 3,
	}
	if id, ok := genres[genre]; ok {
		genreID = id
	} else {
		return nil, entity.NewClientError("такого жанра не существует", entity.ErrNotFound)
	}
	films, err := c.contentRepo.GetFilmsByGenre(genreID)
	if err != nil {
		return nil, err
	}
	filmsIDs := make([]int, 0, len(films))
	for _, film := range films {
		filmsIDs = append(filmsIDs, film.ID)
	}
	return &dto.Compilation{
		Genre:              genre,
		ContentIdentifiers: filmsIDs,
	}, nil
}

func (c CollectionsService) GetGenres() (*dto.Genres, error) {
	return &dto.Genres{
		Genres: []string{"action", "drama", "comedian"},
	}, nil
}
