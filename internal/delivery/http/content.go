package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
)

type ContentEndpoints struct {
	useCase usecase.Content
}

func NewContentEndpoints(useCase usecase.Content) ContentEndpoints {
	return ContentEndpoints{useCase: useCase}
}
