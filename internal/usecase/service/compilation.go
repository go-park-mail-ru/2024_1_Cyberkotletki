package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
)

type CompilationService struct {
	compilationRepo     repository.Compilation
	compilationTypeRepo repository.CompilationType
	staticRepo          repository.Static
	contentRepo         repository.Content
	reviewRepo          repository.Review
}

func NewCompilationService(compilationRepo repository.Compilation, compilationTypeRepo repository.CompilationType,
	staticRepo repository.Static, contentRepo repository.Content) *CompilationService {
	return &CompilationService{
		compilationRepo:     compilationRepo,
		compilationTypeRepo: compilationTypeRepo,
		staticRepo:          staticRepo,
		contentRepo:         contentRepo,
	}
}

// compilationEntityToDTO конвертирует entity.Compilation в dto.Compilation добавляя поле длины контента в подборке
func (c *CompilationService) compilationEntityToDTO(compEntity *entity.Compilation) (*dto.CompilationResponse, error) {
	contentLength, err := c.compilationRepo.GetCompilationContentLength(compEntity.ID)
	if err != nil {
		return nil, err
	}
	return &dto.CompilationResponse{
		Compilation: dto.Compilation{
			ID:                compEntity.ID,
			Title:             compEntity.Title,
			CompilationTypeID: compEntity.CompilationTypeID,
			PosterUploadID:    compEntity.PosterUploadID,
		},
		ContentLength: contentLength,
	}, nil
}

// compilationTypeDTOToEntity конвертирует entity.CompilationType в dto.CompilationType
func (c *CompilationService) compilationTypeDTOToEntity(
	compTypeEntity *entity.CompilationType,
) *dto.CompilationTypeResponse {
	return &dto.CompilationTypeResponse{
		CompilationType: dto.CompilationType{
			ID:   compTypeEntity.ID,
			Type: compTypeEntity.Type,
		},
	}
}

// compilationEntitiesToDTO конвертирует массив entity.Compilation в массив dto.Compilation
func (c *CompilationService) compilationEntitiesToDTO(compEntities []*entity.Compilation) (
	*dto.CompilationResponseList,
	error,
) {
	compDTOs := make([]dto.CompilationResponse, 0)
	for _, compEntity := range compEntities {
		compDTO, err := c.compilationEntityToDTO(compEntity)
		if err != nil {
			return nil, err
		}
		compDTOs = append(compDTOs, *compDTO)
	}
	return &dto.CompilationResponseList{Compilations: compDTOs}, nil
}

func (c *CompilationService) contentEntityToDTO(content *entity.Content) (*dto.PreviewContentCardResponse, error) {
	if len(content.Actors) < 2 ||
		len(content.Directors) == 0 ||
		len(content.Genres) == 0 ||
		len(content.Country) == 0 ||
		len(content.Country) == 0 ||
		len(content.Genres) == 0 {
		return nil, entity.NewClientError("недостаточно данных для создания PreviewContentCard")
	}

	preview, err := c.contentRepo.GetPreviewContent(content.ID)
	if err != nil {
		return nil, err
	}

	actors := make([]string, len(preview.Actors))
	for i, actor := range preview.Actors {
		actors[i] = actor.FirstName + " " + actor.LastName
	}

	poster, err := c.staticRepo.GetStatic(preview.PosterStaticID)
	if err != nil {
		return nil, err
	}

	rating, err := c.reviewRepo.GetContentRating(preview.ID)
	if err != nil {
		return nil, err
	}

	card := dto.PreviewContentCard{
		ID:            preview.ID,
		Title:         preview.Title,
		OriginalTitle: preview.OriginalTitle,
		Country:       preview.Country[0].Name,
		Genre:         preview.Genres[0].Name,
		Director:      preview.Directors[0].FirstName + " " + preview.Directors[0].LastName,
		Actors:        actors,
		Poster:        poster,
		Rating:        float64(rating),
	}

	if preview.Type == entity.ContentTypeMovie {
		card.Duration = preview.Movie.Duration
	} else if preview.Type == entity.ContentTypeSeries {
		card.SeasonsNumber = len(preview.Series.Seasons)
		card.YearStart = preview.Series.YearStart
		card.YearEnd = preview.Series.YearEnd
	}

	return &dto.PreviewContentCardResponse{
		PreviewContentCard: card,
		Type:               preview.Type,
	}, nil
}

func (c *CompilationService) compilationTypeEntitiesToDTO(compTypeEntities []*entity.CompilationType) (
	*dto.CompilationTypeResponseList,
	error,
) {
	compTypeDTOs := make([]dto.CompilationTypeResponse, 0)
	for _, compTypeEntity := range compTypeEntities {
		compTypeDTO := c.compilationTypeDTOToEntity(compTypeEntity)
		compTypeDTOs = append(compTypeDTOs, *compTypeDTO)
	}
	return &dto.CompilationTypeResponseList{CompilationTypes: compTypeDTOs}, nil

}

// GetCompilationTypes возвращает список типов подборок
func (c *CompilationService) GetCompilationTypes() (*dto.CompilationTypeResponseList, error) {
	compTypes, err := c.compilationTypeRepo.GetAllCompilationTypes()
	if err != nil {
		return nil, err
	}
	return c.compilationTypeEntitiesToDTO(compTypes)
}

func (c *CompilationService) GetCompilationsByCompilationType(compTypeID, count, page int) (
	*dto.CompilationResponseList,
	error,
) {
	compEntities, err := c.compilationRepo.GetCompilationsByCompilationTypeID(compTypeID, count, page)
	if err != nil {
		return nil, err
	}
	return c.compilationEntitiesToDTO(compEntities)
}

func (c *CompilationService) GetCompilationContent(compID int) ([]*dto.PreviewContentCardResponse, error) {
	contentIDs, err := c.compilationRepo.GetCompilationContent(compID, 0, 2)
	if err != nil {
		return nil, err
	}

	if len(contentIDs) == 0 {
		return nil, errors.New("no content found for this compilation")
	}

	contentCards := make([]*dto.PreviewContentCardResponse, 0)
	for _, contentID := range contentIDs {
		content, err := c.contentRepo.GetPreviewContent(contentID)
		if err != nil {
			return nil, err
		}

		contentCard, err := c.contentEntityToDTO(content)
		if err != nil {
			return nil, err
		}

		contentCards = append(contentCards, contentCard)
	}

	return contentCards, nil
}
