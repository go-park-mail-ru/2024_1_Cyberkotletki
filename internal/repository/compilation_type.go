package repository

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"

type CompilationType interface {
	GetCompilationType(id int) (*entity.CompilationType, error)
	GetAllCompilationTypes() ([]*entity.CompilationType, error)
}
