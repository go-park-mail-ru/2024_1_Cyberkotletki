package entity

import "time"

type OngoingContent struct {
	ID             int       // Уникальный идентификатор
	Title          string    // Название
	PosterStaticID int       // Постер
	Genres         []Genre   `json:"genres"`      // Жанры
	ReleaseDate    time.Time `json:"releaseDate"` // Дата выхода

	Type string `json:"type"` // Тип контента (movie / series)
}
