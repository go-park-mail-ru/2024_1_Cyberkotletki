package service

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
)

type CompilationService struct {
	compilationRepo repository.Compilation
	staticRepo      repository.Static
	contentRepo     repository.Content
	reviewRepo      repository.Review
}

func NewCompilationService(
	compilationRepo repository.Compilation,
	staticRepo repository.Static,
	contentRepo repository.Content,
	reviewRepo repository.Review,
) *CompilationService {
	return &CompilationService{
		compilationRepo: compilationRepo,
		staticRepo:      staticRepo,
		contentRepo:     contentRepo,
		reviewRepo:      reviewRepo,
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

func (c *CompilationService) contentEntityToDTO(content *entity.Content) (*dto.PreviewContentCard, error) {
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
	var country string
	if len(preview.Country) > 0 {
		country = preview.Country[0].Name
	}
	var genre string
	if len(preview.Genres) > 0 {
		genre = preview.Genres[0].Name
	}
	var director string
	if len(preview.Directors) > 0 {
		director = preview.Directors[0].FirstName + " " + preview.Directors[0].LastName
	}
	card := dto.PreviewContentCard{
		ID:            preview.ID,
		Title:         preview.Title,
		OriginalTitle: preview.OriginalTitle,
		Country:       country,
		Genre:         genre,
		Director:      director,
		Actors:        actors,
		Poster:        poster,
		Rating:        rating,
	}
	if preview.Type == entity.ContentTypeMovie {
		card.Duration = preview.Movie.Duration
	} else if preview.Type == entity.ContentTypeSeries {
		card.SeasonsNumber = len(preview.Series.Seasons)
		card.YearStart = preview.Series.YearStart
		card.YearEnd = preview.Series.YearEnd
	}

	return &card, nil
}

func (c *CompilationService) compilationTypeEntitiesToDTO(compTypeEntities []*entity.CompilationType) (
	*dto.CompilationTypeResponseList,
	error,
) {
	compTypeDTOs := make([]dto.CompilationTypeResponse, 0, len(compTypeEntities))
	for index, compTypeEntity := range compTypeEntities {
		compTypeDTO := c.compilationTypeDTOToEntity(compTypeEntity)
		compTypeDTOs[index] = *compTypeDTO
	}
	return &dto.CompilationTypeResponseList{CompilationTypes: compTypeDTOs}, nil
}

// GetCompilationTypes возвращает список типов подборок
func (c *CompilationService) GetCompilationTypes() (*dto.CompilationTypeResponseList, error) {
	compTypes, err := c.compilationRepo.GetAllCompilationTypes()
	if err != nil {
		return nil, err
	}
	return c.compilationTypeEntitiesToDTO(compTypes)
}

func (c *CompilationService) GetCompilationsByCompilationType(compTypeID int) (
	*dto.CompilationResponseList,
	error,
) {
	compEntities, err := c.compilationRepo.GetCompilationsByTypeID(compTypeID)
	if err != nil {
		return nil, err
	}
	return c.compilationEntitiesToDTO(compEntities)
}

func (c *CompilationService) GetCompilationContent(compID, page, limit int) ([]*dto.PreviewContentCard, error) {
	contentIDs, err := c.compilationRepo.GetCompilationContent(compID, page, limit)
	if err != nil {
		return nil, err
	}
	contentCards := make([]*dto.PreviewContentCard, 0, limit)
	for index, contentID := range contentIDs {
		content, err := c.contentRepo.GetPreviewContent(contentID)
		if err != nil {
			return nil, err
		}
		contentCard, err := c.contentEntityToDTO(content)
		if err != nil {
			return nil, err
		}
		contentCards[index] = contentCard
	}

	return contentCards, nil
}
