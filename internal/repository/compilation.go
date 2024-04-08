package repository

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"

type Compilation interface {
	AddCompilation(compilation *entity.Compilation) (*entity.Compilation, error)
	GetCompilation(id int) (*entity.Compilation, error)
	GetAllCompilationsByCompilationTypeID(id int, limit int, offset int) ([]*entity.Compilation, error)
	UpdateCompilation(compilation *entity.Compilation) (*entity.Compilation, error)
	DeleteCompilation(id int) error
	GetCompilationContentLength(id int) (int, error)
}
