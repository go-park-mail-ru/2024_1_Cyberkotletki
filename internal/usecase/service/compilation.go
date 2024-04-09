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
func (compserv *CompilationService) compilationEntityToDTO(compEntity *entity.Compilation) (*dto.CompilationResponse, error) {
	contentLength, err := compserv.compilationRepo.GetCompilationContentLength(compEntity.ID)
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
func (compserv *CompilationService) compilationTypeDTOToEntity(compTypeEntity *entity.CompilationType) (*dto.CompilationTypeResponse, error) {
	return &dto.CompilationTypeResponse{
		CompilationType: dto.CompilationType{
			ID:   compTypeEntity.ID,
			Type: compTypeEntity.Type,
		},
	}, nil
}

// compilationEntitiesToDTO конвертирует массив entity.Compilation в массив dto.Compilation
func (compserv *CompilationService) compilationEntitiesToDTO(compEntities []*entity.Compilation) (*dto.CompilationResponseList, error) {
	compDTOs := make([]dto.CompilationResponse, 0)
	for _, compEntity := range compEntities {
		compDTO, err := compserv.compilationEntityToDTO(compEntity)
		if err != nil {
			return nil, err
		}
		compDTOs = append(compDTOs, *compDTO)
	}
	return &dto.CompilationResponseList{Compilations: compDTOs}, nil
}

// contentEntityToDTO преобразует entity.Content в dto.PreviewContentCard
func (compserv *CompilationService) contentEntityToDTO(content *entity.Content) (*dto.PreviewContentCardResponse, error) {
	if len(content.Actors) < 2 || len(content.Directors) == 0 || len(content.Genres) == 0 || len(content.Country) == 0 {
		return nil, entity.NewClientError("недостаточно данных для создания PreviewContentCard")
	}

	return &dto.PreviewContentCardResponse{
		PreviewContentCard: dto.PreviewContentCard{
			Title:         content.Title,
			OriginalTitle: content.OriginalTitle,
			ReleaseYear:   content.Movie.Release.Year(),
			Country:       content.Country[0].Name,
			Genre:         content.Genres[0].Name,
			Director:      content.Directors[0].FirstName + " " + content.Directors[0].LastName,
			Actors: []string{content.Actors[0].FirstName + " " + content.Actors[0].LastName,
				content.Actors[1].FirstName + " " + content.Actors[1].LastName},
			Poster:   content.PosterStaticID,
			Rating:   content.IMDBRating,
			Duration: content.Movie.Duration,
		},
		Type: content.Type,
	}, nil
}

func (compserv *CompilationService) compilationTypeEntitiesToDTO(compTypeEntities []*entity.CompilationType) (*dto.CompilationTypeResponseList, error) {
	compTypeDTOs := make([]dto.CompilationTypeResponse, 0)
	for _, compTypeEntity := range compTypeEntities {
		compTypeDTO, err := compserv.compilationTypeDTOToEntity(compTypeEntity)
		if err != nil {
			return nil, err
		}
		compTypeDTOs = append(compTypeDTOs, *compTypeDTO)
	}
	return &dto.CompilationTypeResponseList{CompilationTypes: compTypeDTOs}, nil

}

// GetCompilationTypes возвращает список типов подборок
func (compserv *CompilationService) GetCompilationTypes() (*dto.CompilationTypeResponseList, error) {
	compTypes, err := compserv.compilationTypeRepo.GetAllCompilationTypes()
	if err != nil {
		return nil, err
	}
	return compserv.compilationTypeEntitiesToDTO(compTypes)
}

func (compserv *CompilationService) GetCompilationsByCompilationType(compTypeID, count, page int) (*dto.CompilationResponseList, error) {
	compEntities, err := compserv.compilationRepo.GetCompilationsByCompilationTypeID(compTypeID, count, page)
	if err != nil {
		return nil, err
	}
	return compserv.compilationEntitiesToDTO(compEntities)
}

func (compserv *CompilationService) GetCompilationContent(compID int) ([]*dto.PreviewContentCardResponse, error) {
	contentIDs, err := compserv.compilationRepo.GetCompilationContent(compID, 0, 2)
	if err != nil {
		return nil, err
	}

	if len(contentIDs) == 0 {
		return nil, errors.New("no content found for this compilation")
	}

	var contentCards []*dto.PreviewContentCardResponse
	for _, contentID := range contentIDs {
		content, err := compserv.contentRepo.GetPreviewContent(contentID)
		if err != nil {
			return nil, err
		}

		contentCard, err := compserv.contentEntityToDTO(content)
		if err != nil {
			return nil, err
		}

		contentCards = append(contentCards, contentCard)
	}

	return contentCards, nil
}
