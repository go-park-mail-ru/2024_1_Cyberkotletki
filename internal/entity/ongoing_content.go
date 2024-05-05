package entity

import "time"

type OngoingContent struct {
	ID             int       // Уникальный идентификатор
	Title          string    // Название
	PosterStaticID int       // Постер
	Genres         []Genre   `json:"genres"`      // Жанры
	ReleaseDate    time.Time `json:"releaseDate"` // Дата выхода

	Type string `json:"type"` // Тип контента (movie / series)
	// Поля, которые есть только у фильмов
	Movie *OngoingMovie `json:"movie"`
	// Поля, которые есть только у сериалов
	Series *OngoingSeries `json:"series"`
}

type OngoingMovie struct {
	Premiere time.Time `json:"premiere"` // Дата премьеры
	Duration int       `json:"duration"` // Продолжительность
}

type OngoingSeries struct {
	YearStart int      `json:"year_start"` // Год начала сериала
	YearEnd   int      `json:"year_end"`   // Год окончания сериала
	Seasons   []Season `json:"seasons"`    // Сезоны в сериале
}

const (
	OngoingContentTypeMovie  = "movie"
	OngoingContentTypeSeries = "series"
)
