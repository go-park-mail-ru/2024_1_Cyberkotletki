package usecase

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_compilation.go
type Compilation interface {
	GetCompilationTypes() (*dto.CompilationTypeResponseList, error)
	GetCompilationsByCompilationType(compTypeID int) (*dto.CompilationResponseList, error)
	GetCompilationContent(compID, page, limit int) ([]*dto.PreviewContentCard, error)
}
