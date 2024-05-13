package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
)

type SearchService struct {
	searchRepo repository.Search
	staticUC   usecase.Static
}

func NewSearchService(searchRepo repository.Search, staticUC usecase.Static) usecase.Search {
	return &SearchService{
		searchRepo: searchRepo,
		staticUC:   staticUC,
	}
}

func (s SearchService) Search(query string) (*dto.SearchResult, error) {
	contents, err := s.searchRepo.SearchContent(query)
	if err != nil {
		return nil, entity.UsecaseWrap(err, errors.New("ошибка при поиске контента в SearchService"))
	}
	persons, err := s.searchRepo.SearchPerson(query)
	if err != nil {
		return nil, entity.UsecaseWrap(err, errors.New("ошибка при поиске персон в SearchService"))
	}

	result := dto.SearchResult{
		Content: make([]dto.PreviewContent, len(contents)),
		Persons: make([]dto.PersonPreviewWithPhoto, len(persons)),
	}
	for index, content := range contents {
		posterURL, err := s.staticUC.GetStatic(content.PosterStaticID)
		switch {
		case errors.Is(err, usecase.ErrStaticNotFound):
			posterURL = ""
		case err != nil:
			return nil, entity.UsecaseWrap(err, errors.New("ошибка при получении статики контента из Search"))
		}
		countries := content.Country
		var country string
		if len(countries) == 0 {
			country = ""
		} else {
			country = countries[0].Name
		}
		genres := content.Genres
		var genre string
		if len(genres) == 0 {
			genre = ""
		} else {
			genre = genres[0].Name
		}
		directors := content.Directors
		var director string
		if len(directors) == 0 {
			director = ""
		} else {
			director = directors[0].Name
		}
		actors := content.Actors
		var actorsList []string
		for _, actor := range actors {
			actorsList = append(actorsList, actor.Name)
		}
		contentDTO := dto.PreviewContent{
			ID:            content.ID,
			Title:         content.Title,
			OriginalTitle: content.OriginalTitle,
			Country:       country,
			Genre:         genre,
			Director:      director,
			Actors:        actorsList,
			Poster:        posterURL,
			Rating:        content.Rating,
			Type:          content.Type,
		}
		switch content.Type {
		case "movie":
			contentDTO.Duration = content.Movie.Duration
		case "series":
			contentDTO.SeasonsNumber = len(content.Series.Seasons)
			contentDTO.YearStart = content.Series.YearStart
			contentDTO.YearEnd = content.Series.YearEnd
		}
		result.Content[index] = contentDTO
	}
	for index, person := range persons {
		posterURL, err := s.staticUC.GetStatic(person.GetPhotoStaticID())
		switch {
		case errors.Is(err, usecase.ErrStaticNotFound):
			posterURL = ""
		case err != nil:
			return nil, entity.UsecaseWrap(err, errors.New("ошибка при получении статики персоны из Search"))
		}
		personDTO := dto.PersonPreviewWithPhoto{
			ID:       person.ID,
			Name:     person.Name,
			PhotoURL: posterURL,
		}
		result.Persons[index] = personDTO
	}
	return &result, nil
}
