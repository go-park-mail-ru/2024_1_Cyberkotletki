package repository

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"

type Compilation interface {
	GetCompilationsByCompilationTypeID(compilationTypeID, page, limit int) ([]*entity.Compilation, error)
	GetCompilation(id int) (*entity.Compilation, error)
	GetCompilationContentLength(id int) (int, error)
	GetCompilationContent(id, page, limit int) ([]int, error)
}
