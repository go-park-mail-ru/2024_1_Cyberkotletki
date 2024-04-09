package repository

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_compilation_type.go
type CompilationType interface {
	GetCompilationType(id int) (*entity.CompilationType, error)
	GetAllCompilationTypes() ([]*entity.CompilationType, error)
}
