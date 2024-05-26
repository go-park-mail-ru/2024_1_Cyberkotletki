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
	contentUC  usecase.Content
}

func NewSearchService(searchRepo repository.Search, contentUC usecase.Content) usecase.Search {
	return &SearchService{
		searchRepo: searchRepo,
		contentUC:  contentUC,
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
		Content: make([]*dto.PreviewContent, len(contents)),
		Persons: make([]*dto.PersonPreviewWithPhoto, len(persons)),
	}
	for index, content := range contents {
		contentDTO, err := s.contentUC.GetPreviewContentByID(content)
		if err != nil {
			return nil, errors.Join(err, errors.New("ошибка при получении контента из Search"))
		}
		result.Content[index] = contentDTO
	}
	for index, person := range persons {
		personDTO, err := s.contentUC.GetPreviewPersonByID(person)
		if err != nil {
			return nil, errors.Join(err, errors.New("ошибка при получении персоны из Search"))
		}
		result.Persons[index] = personDTO
	}
	return &result, nil
}
