package entity

/*
TODO: тесты
*/

// Episode представляет эпизод сериала
type Episode struct {
	Id            int    `json:"id"`             // Уникальный идентификатор
	EpisodeNumber int    `json:"episode_number"` // Номер эпизода
	Duration      int    `json:"duration"`       // Продолжительность эпизода в минутах
	Description   string `json:"description"`    // Описание эпизода
}

// Season представляет сезон сериала
type Season struct {
	Content   Content   // Общий контент сезона
	YearStart int       `json:"year_start"` // Год начала сезона
	YearEnd   int       `json:"year_end"`   // Год окончания сезона
	Episodes  []Episode `json:"episodes"`   // Эпизоды в сезоне
}

// Series представляет сериал
type Series struct {
	Id        int      `json:"id"`         // Уникальный идентификатор
	Title     string   `json:"title"`      // Название
	YearStart int      `json:"year_start"` // Год начала сериала
	YearEnd   int      `json:"year_end"`   // Год окончания сериала
	Seasons   []Season `json:"seasons"`    // Сезоны в сериале
}
