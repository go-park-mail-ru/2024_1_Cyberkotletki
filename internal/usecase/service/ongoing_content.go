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
	staticRepo         repository.Static
}

func ongoindMovieEntityToDTO(ongoingMovieEntity entity.OngoingMovie) dto.OngoingMovie {
	return dto.OngoingMovie{
		Duration: ongoingMovieEntity.Duration,
		Premiere: ongoingMovieEntity.Premiere,
	}
}

func ongoingSeriesEntityToDTO(ongoingSeriesEntity entity.OngoingSeries) dto.OngoingSeries {
	return dto.OngoingSeries{
		YearStart: ongoingSeriesEntity.YearStart,
		YearEnd:   ongoingSeriesEntity.YearEnd,
	}
}

// GetAllReleaseYears implements usecase.OngoingContent.
func (o *OngoingContentService) GetAllReleaseYears() ([]int, error) {
	return o.ongoingContentRepo.GetAllReleaseYears()
}

func (o *OngoingContentService) ongoindContentEntityToDTO(ongoingContentEntity *entity.OngoingContent) (*dto.PreviewOngoingContentCardVertical, error) {
	ongoingContent, err := o.ongoingContentRepo.GetOngoingContentByID(ongoingContentEntity.ID)
	switch {
	case errors.Is(err, repository.ErrOngoingContentNotFound):
		return nil, usecase.ErrOngoingContentNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении контента календаря релизов"), err)
	}
	posterURL, err := o.staticRepo.GetStatic(ongoingContent.PosterStaticID)
	switch {
	case errors.Is(err, repository.ErrStaticNotFound):
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
	}
	switch ongoingContent.Type {
	case entity.OngoingContentTypeMovie:
		ongoingContentDTO.Movie = ongoindMovieEntityToDTO(*ongoingContent.Movie)
	case entity.OngoingContentTypeSeries:
		ongoingContentDTO.Series = ongoingSeriesEntityToDTO(*ongoingContent.Series)
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
		ongoingContentDTO, err := o.ongoindContentEntityToDTO(ongoingContentEntity)
		if err != nil {
			return nil, err
		}
		ongoingContentDTOs = append(ongoingContentDTOs, ongoingContentDTO)
	}
	return ongoingContentDTOs, nil
}

// GetOngoingContentByContentID implements usecase.OngoingContent.
func (o *OngoingContentService) GetOngoingContentByContentID(contentID int) (*dto.PreviewOngoingContentCardVertical, error) {
	ongoingContent, err := o.ongoingContentRepo.GetOngoingContentByID(contentID)
	switch {
	case errors.Is(err, repository.ErrOngoingContentNotFound):
		return nil, usecase.ErrOngoingContentNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении контента календаря релизов"), err)
	}
	posterURL, err := o.staticRepo.GetStatic(ongoingContent.PosterStaticID)
	switch {
	case errors.Is(err, repository.ErrStaticNotFound):
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
	}
	switch ongoingContent.Type {
	case entity.OngoingContentTypeMovie:
		ongoingContentDTO.Movie = ongoindMovieEntityToDTO(*ongoingContent.Movie)
	case entity.OngoingContentTypeSeries:
		ongoingContentDTO.Series = ongoingSeriesEntityToDTO(*ongoingContent.Series)
	}

	return ongoingContentDTO, nil
}

// GetOngoingContentByMonthAndYear implements usecase.OngoingContent.
func (o *OngoingContentService) GetOngoingContentByMonthAndYear(month int, year int) ([]*dto.PreviewOngoingContentCardVertical, error) {
	ongoingContentEntities, err := o.ongoingContentRepo.GetOngoingContentByMonthAndYear(month, year)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении релизов по месяцу и году"), err)
	}
	ongoingContentDTOs := make([]*dto.PreviewOngoingContentCardVertical, 0, len(ongoingContentEntities))
	for _, ongoingContentEntity := range ongoingContentEntities {
		ongoingContentDTO, err := o.ongoindContentEntityToDTO(ongoingContentEntity)
		if err != nil {
			return nil, err
		}
		ongoingContentDTOs = append(ongoingContentDTOs, ongoingContentDTO)
	}
	return ongoingContentDTOs, nil
}

// IsOngoingConentFinished implements usecase.OngoingContent.
func (o *OngoingContentService) IsOngoingConentFinished(contentID int) (bool, error) {
	return o.ongoingContentRepo.IsOngoingConentFinished(contentID)
}

func NewOngoingContentService(ongoingContentRepo repository.OngoingContent, staticRepo repository.Static) usecase.OngoingContent {
	return &OngoingContentService{
		ongoingContentRepo: ongoingContentRepo,
		staticRepo:         staticRepo,
	}
}
