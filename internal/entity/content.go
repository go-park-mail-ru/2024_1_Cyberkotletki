package entity

import "time"

// Content представляет основную структуру для хранения информации о контенте.
// В зависимости от типа контента, некоторые поля могут быть пустыми.
type Content struct {
	ID               int     // Уникальный идентификатор
	Title            string  // Название
	OriginalTitle    string  // Название
	Slogan           string  // Слоган
	Budget           string  // Бюджет
	AgeRestriction   int     // Возрастное ограничение
	IMDBRating       float64 // Рейтинг IMDB
	Rating           float64 // Рейтинг
	Description      string  // Описание
	PosterStaticID   int     // Постер
	TrailerLink      string  // Ссылка на трейлер
	BackdropStaticID int     // Фоновое изображение

	PicturesStaticID []int    // Изображения
	Facts            []string // Интересные факты

	Country   []Country `json:"country"`   // Страны, где был произведен контент
	Genres    []Genre   `json:"genres"`    // Жанры
	Actors    []Person  `json:"actors"`    // Актеры
	Directors []Person  `json:"directors"` // Режиссеры
	Producers []Person  `json:"producers"` // Продюсеры
	Writers   []Person  `json:"writers"`   // Сценаристы
	Operators []Person  `json:"operators"` // Операторы
	Composers []Person  `json:"composers"` // Композиторы
	Editors   []Person  `json:"editors"`   // Редакторы

	Type string `json:"type"` // Тип контента (movie / series)
	// Поля, которые есть только у фильмов
	Movie *Movie `json:"movie"`
	// Поля, которые есть только у сериалов
	Series *Series `json:"series"`
}

type Movie struct {
	Premiere time.Time `json:"premiere"` // Дата премьеры
	Duration int       `json:"duration"` // Продолжительность
}

type Series struct {
	YearStart int      `json:"year_start"` // Год начала сериала
	YearEnd   int      `json:"year_end"`   // Год окончания сериала
	Seasons   []Season `json:"seasons"`    // Сезоны в сериале
}

// Episode представляет эпизод сериала
type Episode struct {
	ID            int    `json:"id"`             // Уникальный идентификатор
	EpisodeNumber int    `json:"episode_number"` // Номер эпизода
	Title         string `json:"title"`          // Название эпизода
	Duration      int    `json:"duration"`       // Продолжительность
}

// Season представляет сезон сериала
type Season struct {
	ID       int       `json:"id"`       // Уникальный идентификатор
	Title    string    `json:"title"`    // Название сезона
	Episodes []Episode `json:"episodes"` // Эпизоды в сезоне
}

const (
	ContentTypeMovie  = "movie"
	ContentTypeSeries = "series"
)
