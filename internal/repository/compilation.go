package repository

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_compilation.go
type Compilation interface {
	GetCompilationsByTypeID(compilationTypeID int) ([]*entity.Compilation, error)
	GetCompilationContentLength(id int) (int, error)
	GetCompilationContent(id, page, limit int) ([]int, error)
	GetAllCompilationTypes() ([]*entity.CompilationType, error)
}
