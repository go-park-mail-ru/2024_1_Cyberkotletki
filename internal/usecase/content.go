package usecase

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_content.go
type Content interface {
	GetContentPreviewCard(int) (*dto.PreviewContentCard, error)
}
