package usecase

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_search.go
type Search interface {
	Search(query string) (*dto.SearchResult, error)
}
