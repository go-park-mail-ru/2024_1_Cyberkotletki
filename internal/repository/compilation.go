package repository

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_compilation.go
type Compilation interface {
	// GetCompilation возвращает подборку по ее ID
	// Если подборка не найдена, возвращает ErrCompilationNotFound
	GetCompilation(id int) (entity.Compilation, error)
	// GetCompilationsByTypeID возвращает подборку по ее ID
	GetCompilationsByTypeID(compilationTypeID int) ([]entity.Compilation, error)
	// GetCompilationContentLength возвращает длину контента подборки по ее ID
	GetCompilationContentLength(id int) (int, error)
	// GetCompilationContent возвращает подборку по ее ID
	// Если подборка не найдена, возвращает ErrCompilationNotFound
	GetCompilationContent(id, page, limit int) ([]int, error)
	// GetAllCompilationTypes возвращает подборку по ее ID
	GetAllCompilationTypes() ([]entity.CompilationType, error)
}

var (
	ErrCompilationNotFound = errors.New("подборка не найдена")
)
