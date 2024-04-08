package usecase

//TODO добавить респонс

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"

type Compilation interface {
	CreateCompilation(dto.Compilation) (*dto.CompilationResponse, error)
	GetCompilation(int) (*dto.CompilationResponse, error)
	EditCompilation(dto.Compilation) (*dto.CompilationResponse, error)
	DeleteCompilation(int) error
	GetCompilationTypeCompilations(compilationTypeID, count, page int) (*dto.CompilationResponseList, error)
}
