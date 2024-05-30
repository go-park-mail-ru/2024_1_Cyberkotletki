package entity

import (
	"log"
	"strings"
	"time"
	"unicode/utf8"
)

// Content представляет основную структуру для хранения информации о контенте.
// В зависимости от типа контента, некоторые поля могут быть пустыми.
type Content struct {
	ID             int     `json:"id"`               // Уникальный идентификатор
	Title          string  `json:"title"`            // Название
	OriginalTitle  string  `json:"original_title"`   // Название
	Slogan         string  `json:"slogan"`           // Слоган
	Budget         int     `json:"budget"`           // Бюджет
	AgeRestriction int     `json:"age_restriction"`  // Возрастное ограничение
	Audience       int     `json:"audience"`         // Аудитория
	IMDBRating     float64 `json:"imdb_rating"`      // Рейтинг IMDB
	Description    string  `json:"description"`      // Описание
	PosterStaticID int     `json:"poster_static_id"` // Постер
	BoxOffice      int     `json:"box_office"`       // Кассовые сборы
	Marketing      int     `json:"marketing"`        // Маркетинговые затраты

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
	Release  time.Time `json:"release"`  // Дата выпуска
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
}

// Season представляет сезон сериала
type Season struct {
	ID        int       `json:"id"`         // Уникальный идентификатор
	YearStart int       `json:"year_start"` // Год начала сезона
	YearEnd   int       `json:"year_end"`   // Год окончания сезона
	Episodes  []Episode `json:"episodes"`   // Эпизоды в сезоне
}

const (
	ContentTypeMovie  = "movie"
	ContentTypeSeries = "series"
)

func ValidateAllContent(title, originalTitle, slogan, description string,
	budget, ageRestriction, audience, boxOffice, marketing int, imdbRating float64,
	typeContent string) error {
	if utf8.RuneCountInString(strings.TrimSpace(title)) > 150 || utf8.RuneCountInString(strings.TrimSpace(title)) < 1 {
		return NewClientError("Чило символов в названии фильма не должно превышать 150 символов", ErrBadRequest)
	}
	if utf8.RuneCountInString(strings.TrimSpace(originalTitle)) > 150 {
		return NewClientError("Чило символов в оригинальном названии фильма не должно превышать 150 символов", ErrBadRequest)
	}
	if utf8.RuneCountInString(strings.TrimSpace(slogan)) > 150 {
		return NewClientError("Чило символов в слогане не должно превышать 150 символов", ErrBadRequest)
	}
	if utf8.RuneCountInString(strings.TrimSpace(description)) > 10000 {
		return NewClientError("Чило символов в описании не должно превышать 10000 символов", ErrBadRequest)
	}
	if typeContent != ContentTypeMovie && typeContent != ContentTypeSeries {
		return NewClientError("Тип контента должен быть либо movie, либо series", ErrBadRequest)
	}
	if budget < 0 {
		return NewClientError("Бюджет не может быть отрицательным", ErrBadRequest)
	}
	if ageRestriction < 0 {
		return NewClientError("Возрастное ограничение не может быть отрицательным", ErrBadRequest)
	}
	if audience < 0 {
		return NewClientError("Аудитория не может быть отрицательной", ErrBadRequest)
	}
	if imdbRating < 0 || imdbRating > 10 {
		return NewClientError("Рейтинг IMDB не может быть отрицательным или больше 10", ErrBadRequest)
	}
	if boxOffice < 0 {
		return NewClientError("Кассовые сборы не могут быть отрицательными", ErrBadRequest)
	}
	if marketing < 0 {
		return NewClientError("Маркетинговые затраты не могут быть отрицательными", ErrBadRequest)
	}
	return nil
}

func ValidateContent(title, originalTitle, slogan, description string,
	budget, ageRestriction, audience, boxOffice, marketing int, imdbRating float64,
	typeContent string) error {
	if err := ValidateAllContent(title, originalTitle, slogan, description, budget, ageRestriction, audience, boxOffice, marketing, imdbRating, typeContent); err != nil {
		log.Println("Valid Error entity", err)
		return err
	}
	return nil
}
