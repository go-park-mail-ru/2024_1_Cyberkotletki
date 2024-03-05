package content

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/audience"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/award"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/boxoffice"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/country"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/genre"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	"time"
)

// Content представляет основную структуру для хранения информации о контенте.
type Content struct {
	Id               int                   `json:"id"`               // Уникальный идентификатор
	Title            string                `json:"title"`            // Название
	OriginalTitle    string                `json:"original_title"`   // Название
	Country          []country.Country     `json:"country"`          // Страны, где был произведен контент
	Genres           []genre.Genre         `json:"genres"`           // Жанры
	Directors        []person.Person       `json:"directors"`        // Режиссеры
	Writers          []person.Person       `json:"writers"`          // Сценаристы
	Producers        []person.Person       `json:"producers"`        // Продюсеры
	Cinematographers []person.Person       `json:"cinematographers"` // Операторы
	Slogan           string                `json:"slogan"`           // Слоган
	Composers        []person.Person       `json:"composers"`        // Композиторы
	Artists          []person.Person       `json:"artists"`          // Художники
	Editors          []person.Person       `json:"editors"`          // Редакторы
	Budget           int                   `json:"budget"`           // Бюджет
	Marketing        int                   `json:"marketing"`        // Маркетинговые затраты
	BoxOffices       []boxoffice.BoxOffice `json:"box_offices"`      // Кассовые сборы
	Audiences        []audience.Audience   `json:"audiences"`        // Аудитория
	Premiere         time.Time             `json:"premiere"`         // Дата премьеры
	Release          time.Time             `json:"release"`          // Дата выпуска
	AgeRestriction   int                   `json:"age_restriction"`  // Возрастное ограничение
	Rating           float64               `json:"rating"`           // Рейтинг
	Actors           []person.Person       `json:"actors"`           // Актеры
	Dubbing          []person.Person       `json:"dubbing"`          // Дубляж
	Awards           []award.Award         `json:"awards,omitempty"` // Награды
	Description      string                `json:"description"`      // Описание
	Poster           string                `json:"poster"`           // Постер
	Playback         string                `json:"playback"`         // Воспроизведение на заднем плане небольшоко фрагмента видео
}
