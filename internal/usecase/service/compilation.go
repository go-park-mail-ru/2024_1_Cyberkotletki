package service

//TODO Доделать сервис для компиляций

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
)

type CompilationService struct {
	compilationRepo     repository.Compilation
	compilationTypeRepo repository.CompilationType
	staticRepo          repository.Static
}

func NewCompilationService(compilationRepo repository.Compilation, compilationTypeRepo repository.CompilationType,
	staticRepo repository.Static) *CompilationService {
	return &CompilationService{
		compilationRepo:     compilationRepo,
		compilationTypeRepo: compilationTypeRepo,
		staticRepo:          staticRepo,
	}
}

func (compserv *CompilationService) compilationEntityToDTO(compEntity *entity.Compilation) (*dto.Compilation, error) {
	comp, err := compserv.compilationTypeRepo.GetCompilationType(compEntity.ID)
	if err != nil {
		return nil, err
	}
	var poster string
	poster, err = compserv.staticRepo.GetStatic(compEntity.PosterUploadID)
	if err != nil {
		return nil, err
	}
	contentLength, err := compserv.compilationRepo.GetCompilationContentLength(compEntity.ID)
	if err != nil {
		return nil, err
	}

}
