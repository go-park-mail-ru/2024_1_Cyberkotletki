package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
)

const (
	compilationContentLimit = 10
)

type CompilationService struct {
	compilationRepo repository.Compilation
	staticRepo      repository.Static
	contentRepo     repository.Content
}

func NewCompilationService(
	compilationRepo repository.Compilation,
	staticRepo repository.Static,
	contentRepo repository.Content,
) usecase.Compilation {
	return &CompilationService{
		compilationRepo: compilationRepo,
		staticRepo:      staticRepo,
		contentRepo:     contentRepo,
	}
}

// compilationEntityToDTO конвертирует entity.Compilation в dto.Compilation добавляя поле длины контента в подборке
func (c *CompilationService) compilationEntityToDTO(compEntity entity.Compilation) (*dto.Compilation, error) {
	posterURL, err := c.staticRepo.GetStatic(compEntity.PosterUploadID)
	switch {
	case errors.Is(err, repository.ErrStaticNotFound):
		posterURL = ""
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении статики"), err)
	}
	compilationDTO := dto.Compilation{
		ID:                compEntity.ID,
		Title:             compEntity.Title,
		CompilationTypeID: compEntity.CompilationTypeID,
		PosterURL:         posterURL,
	}
	return &compilationDTO, nil
}

// compilationTypeDTOToEntity конвертирует entity.CompilationType в dto.CompilationType
func (c *CompilationService) compilationTypeDTOToEntity(
	compTypeEntity entity.CompilationType,
) dto.CompilationType {
	return dto.CompilationType{
		ID:   compTypeEntity.ID,
		Type: compTypeEntity.Name,
	}
}

// compilationEntitiesToDTO конвертирует массив entity.Compilation в массив dto.Compilation
func (c *CompilationService) compilationEntitiesToDTO(compEntities []entity.Compilation) (
	*dto.CompilationResponseList,
	error,
) {
	compDTOs := make([]dto.Compilation, len(compEntities))
	for i, compEntity := range compEntities {
		compDTO, err := c.compilationEntityToDTO(compEntity)
		if err != nil {
			return nil, err
		}
		compDTOs[i] = *compDTO
	}
	return &dto.CompilationResponseList{Compilations: compDTOs}, nil
}

func (c *CompilationService) compilationTypeEntitiesToDTO(
	compTypeEntities []entity.CompilationType,
) *dto.CompilationTypeResponseList {
	compTypeDTOs := make([]dto.CompilationType, len(compTypeEntities))
	for i, compTypeEntity := range compTypeEntities {
		compTypeDTOs[i] = c.compilationTypeDTOToEntity(compTypeEntity)
	}
	return &dto.CompilationTypeResponseList{CompilationTypes: compTypeDTOs}
}

// GetCompilationTypes возвращает список типов подборок
func (c *CompilationService) GetCompilationTypes() (*dto.CompilationTypeResponseList, error) {
	compTypes, err := c.compilationRepo.GetAllCompilationTypes()
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении типов подборок"), err)
	}
	return c.compilationTypeEntitiesToDTO(compTypes), nil
}

func (c *CompilationService) GetCompilationsByCompilationType(compTypeID int) (
	*dto.CompilationResponseList,
	error,
) {
	compEntities, err := c.compilationRepo.GetCompilationsByTypeID(compTypeID)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении подборок по типу"), err)
	}
	return c.compilationEntitiesToDTO(compEntities)
}

func (c *CompilationService) GetCompilationContent(compID, page int) (*dto.CompilationResponse, error) {
	compilation, err := c.compilationRepo.GetCompilation(compID)
	switch {
	case errors.Is(err, repository.ErrCompilationNotFound):
		return nil, usecase.ErrCompilationNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении подборки"), err)
	}
	contentIDs, err := c.compilationRepo.GetCompilationContent(compID, page, compilationContentLimit)
	switch {
	case errors.Is(err, repository.ErrCompilationNotFound):
		return nil, usecase.ErrCompilationNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении контента подборки"), err)
	}
	contentTotal, err := c.compilationRepo.GetCompilationContentLength(compID)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении количества контента подборки"), err)
	}
	compilationDTO, err := c.compilationEntityToDTO(compilation)
	if err != nil {
		return nil, err
	}
	return &dto.CompilationResponse{
		Compilation:   *compilationDTO,
		ContentIDs:    contentIDs,
		ContentLength: contentTotal,
		Page:          page,
		PerPage:       compilationContentLimit,
		TotalPages:    (contentTotal + compilationContentLimit - 1) / compilationContentLimit,
	}, nil
}
