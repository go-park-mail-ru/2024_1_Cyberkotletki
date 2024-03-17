package repository

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/content"

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_content.go
type Content interface {
	content.Film
	// todo: content.Series для сериалов
}
