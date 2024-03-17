package entity

import (
	"time"
)

// Content представляет основную структуру для хранения информации о контенте.
type Content struct {
	Id               int         `json:"id"`               // Уникальный идентификатор
	Title            string      `json:"title"`            // Название
	OriginalTitle    string      `json:"original_title"`   // Название
	Country          []Country   `json:"country"`          // Страны, где был произведен контент
	Genres           []Genre     `json:"genres"`           // Жанры
	Directors        []Person    `json:"directors"`        // Режиссеры
	Writers          []Person    `json:"writers"`          // Сценаристы
	Producers        []Person    `json:"producers"`        // Продюсеры
	Cinematographers []Person    `json:"cinematographers"` // Операторы
	Slogan           string      `json:"slogan"`           // Слоган
	Composers        []Person    `json:"composers"`        // Композиторы
	Artists          []Person    `json:"artists"`          // Художники
	Editors          []Person    `json:"editors"`          // Редакторы
	Budget           int         `json:"budget"`           // Бюджет
	Marketing        int         `json:"marketing"`        // Маркетинговые затраты
	BoxOffices       []BoxOffice `json:"box_offices"`      // Кассовые сборы
	Audiences        []Audience  `json:"audiences"`        // Аудитория
	Premiere         time.Time   `json:"premiere"`         // Дата премьеры
	Release          time.Time   `json:"release"`          // Дата выпуска
	AgeRestriction   int         `json:"age_restriction"`  // Возрастное ограничение
	Rating           float64     `json:"rating"`           // Рейтинг
	Actors           []Person    `json:"actors"`           // Актеры
	Dubbing          []Person    `json:"dubbing"`          // Дубляж
	Awards           []Award     `json:"awards,omitempty"` // Награды
	Description      string      `json:"description"`      // Описание
	Poster           string      `json:"poster"`           // Постер
	Playback         string      `json:"playback"`         // Воспроизведение на заднем плане небольшоко фрагмента видео
}
