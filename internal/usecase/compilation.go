package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_compilation.go
type Compilation interface {
	// GetCompilationTypes возвращает список типов подборок
	GetCompilationTypes() (*dto.CompilationTypeResponseList, error)
	// GetCompilationsByCompilationType возвращает список подборок по типу подборок
	GetCompilationsByCompilationType(compTypeID int) (*dto.CompilationResponseList, error)
	// GetCompilationContent возвращает контент подборки по ее ID
	// Если подборка не найдена, возвращает ErrCompilationNotFound
	GetCompilationContent(compID, page int) (*dto.CompilationResponse, error)
}

var (
	ErrCompilationNotFound = errors.New("подборка не найдена")
)
