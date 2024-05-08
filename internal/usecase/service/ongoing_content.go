package service

import (
	"errors"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
)

type OngoingContentService struct {
	ongoingContentRepo repository.OngoingContent
	staticUC           usecase.Static
}

func NewOngoingContentService(
	ongoingContentRepo repository.OngoingContent,
	staticUC usecase.Static,
) usecase.OngoingContent {
	return &OngoingContentService{
		ongoingContentRepo: ongoingContentRepo,
		staticUC:           staticUC,
	}
}

// GetAllReleaseYears implements usecase.OngoingContent.
func (o *OngoingContentService) GetAllReleaseYears() ([]int, error) {
	return o.ongoingContentRepo.GetAllReleaseYears()
}

func (o *OngoingContentService) ongoingContentEntityToDTO(
	ongoingContentEntity *entity.OngoingContent,
) (*dto.PreviewOngoingContentCardVertical, error) {
	ongoingContent, err := o.ongoingContentRepo.GetOngoingContentByID(ongoingContentEntity.ID)
	switch {
	case errors.Is(err, repository.ErrOngoingContentNotFound):
		return nil, usecase.ErrOngoingContentNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении контента календаря релизов"), err)
	}
	posterURL, err := o.staticUC.GetStatic(ongoingContent.PosterStaticID)
	switch {
	case errors.Is(err, usecase.ErrStaticNotFound):
		posterURL = ""
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении постера"), err)
	}
	ongoingContentDTO := &dto.PreviewOngoingContentCardVertical{
		ID:          ongoingContent.ID,
		Title:       ongoingContent.Title,
		Poster:      posterURL,
		ReleaseDate: ongoingContent.ReleaseDate,
		Genres:      genreEntityToDTO(ongoingContent.Genres),
		Type:        ongoingContent.Type,
	}

	return ongoingContentDTO, nil
}

// GetNearestOngoings implements usecase.OngoingContent.
func (o *OngoingContentService) GetNearestOngoings(limit int) ([]*dto.PreviewOngoingContentCardVertical, error) {
	ongoingContentEntities, err := o.ongoingContentRepo.GetNearestOngoings(limit)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении ближайших релизов"), err)
	}
	ongoingContentDTOs := make([]*dto.PreviewOngoingContentCardVertical, 0, len(ongoingContentEntities))
	for _, ongoingContentEntity := range ongoingContentEntities {
		ongoingContentDTO, err := o.ongoingContentEntityToDTO(ongoingContentEntity)
		if err != nil {
			return nil, err
		}
		ongoingContentDTOs = append(ongoingContentDTOs, ongoingContentDTO)
	}
	return ongoingContentDTOs, nil
}

// GetOngoingContentByContentID implements usecase.OngoingContent.
func (o *OngoingContentService) GetOngoingContentByContentID(
	contentID int,
) (*dto.PreviewOngoingContentCardVertical, error) {
	ongoingContent, err := o.ongoingContentRepo.GetOngoingContentByID(contentID)
	switch {
	case errors.Is(err, repository.ErrOngoingContentNotFound):
		return nil, usecase.ErrOngoingContentNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении контента календаря релизов"), err)
	}
	posterURL, err := o.staticUC.GetStatic(ongoingContent.PosterStaticID)
	switch {
	case errors.Is(err, usecase.ErrStaticNotFound):
		posterURL = ""
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении постера"), err)
	}
	ongoingContentDTO := &dto.PreviewOngoingContentCardVertical{
		ID:          ongoingContent.ID,
		Title:       ongoingContent.Title,
		Poster:      posterURL,
		ReleaseDate: ongoingContent.ReleaseDate,
		Genres:      genreEntityToDTO(ongoingContent.Genres),
		Type:        ongoingContent.Type,
	}

	return ongoingContentDTO, nil
}

// GetOngoingContentByMonthAndYear implements usecase.OngoingContent.
func (o *OngoingContentService) GetOngoingContentByMonthAndYear(
	month int,
	year int,
) ([]*dto.PreviewOngoingContentCardVertical, error) {
	ongoingContentEntities, err := o.ongoingContentRepo.GetOngoingContentByMonthAndYear(month, year)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении релизов по месяцу и году"), err)
	}
	ongoingContentDTOs := make([]*dto.PreviewOngoingContentCardVertical, 0, len(ongoingContentEntities))
	for _, ongoingContentEntity := range ongoingContentEntities {
		ongoingContentDTO, err := o.ongoingContentEntityToDTO(ongoingContentEntity)
		if err != nil {
			return nil, err
		}
		ongoingContentDTOs = append(ongoingContentDTOs, ongoingContentDTO)
	}
	return ongoingContentDTOs, nil
}

// IsOngoingContentFinished implements usecase.OngoingContent.
func (o *OngoingContentService) IsOngoingContentFinished(contentID int) (bool, error) {
	return o.ongoingContentRepo.IsOngoingContentFinished(contentID)
}
