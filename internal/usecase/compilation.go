package usecase

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_compilation.go
type Compilation interface {
	GetCompilation(compID int) (*dto.CompilationResponse, error)
	// GetCompilationTypes возвращает все типы подборок (кнопки сверху это и есть типы, которые должны лежать в бд)
	GetCompilationTypes() (*dto.CompilationTypeResponseList, error)
	// GetCompilationsByCompilationType получить список подборок по айди типа компиляции
	GetCompilationsByCompilationType(compTypeID int) (*dto.CompilationResponseList, error)
	// GetCompilation получить компиляцию по айди, должна возвращать все карточки контента
	GetCompilationContentCards(compID int) (*dto.PreviewContentCardResponse, error)
}
