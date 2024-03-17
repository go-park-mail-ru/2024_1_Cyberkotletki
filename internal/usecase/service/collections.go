package service

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/DTO"
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

func (c CollectionsService) GetCompilation(genre string) (*DTO.Compilation, error) {
	var genreId int
	switch genre {
	case "drama":
		genreId = 1
	case "action":
		genreId = 2
	case "comedian":
		genreId = 3
	default:
		return nil, entity.NewClientError("такого жанра не существует", entity.ErrNotFound)
	}
	if films, err := c.contentRepo.GetFilmsByGenre(genreId); err != nil {
		return nil, err
	} else {
		var filmsIds []int
		for _, film := range films {
			filmsIds = append(filmsIds, film.Id)
		}
		return &DTO.Compilation{
			Genre:              genre,
			ContentIdentifiers: filmsIds,
		}, nil
	}
}

func (c CollectionsService) GetGenres() (*DTO.Genres, error) {
	return &DTO.Genres{
		Genres: []string{"action", "drama", "comedian"},
	}, nil
}
