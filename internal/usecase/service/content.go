package service

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/DTO"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
)

type ContentService struct {
	contentRepo repository.Content
}

func NewContentService(contentRepo repository.Content) usecase.Content {
	return &ContentService{
		contentRepo: contentRepo,
	}
}

func (c ContentService) GetContentPreviewCard(contentId int) (*DTO.PreviewContentCard, error) {
	film, err := c.contentRepo.GetFilm(contentId)
	if err != nil {
		return nil, err
	}
	var country string
	if len(film.Country) > 0 {
		country = film.Country[0].Name
	} else {
		country = ""
	}
	var genre string
	if len(film.Genres) > 0 {
		genre = film.Genres[0].Name
	} else {
		genre = ""
	}
	var director string
	if len(film.Directors) > 0 {
		director = fmt.Sprintf("%s %s", film.Directors[0].FirstName, film.Directors[0].LastName)
	} else {
		director = ""
	}
	var actors []string
	if len(film.Actors) > 0 {
		actors = []string{
			fmt.Sprintf("%s %s", film.Actors[0].FirstName, film.Actors[0].LastName),
			fmt.Sprintf("%s %s", film.Actors[1].FirstName, film.Actors[1].LastName),
		}
	} else {
		actors = []string{}
	}
	return &DTO.PreviewContentCard{
		Title:         film.Title,
		OriginalTitle: film.OriginalTitle,
		ReleaseYear:   film.Release.Year(),
		Country:       country,
		Genre:         genre,
		Director:      director,
		Actors:        actors,
		Poster:        film.Poster,
		Rating:        film.Rating,
		Duration:      film.Duration,
	}, nil
}
