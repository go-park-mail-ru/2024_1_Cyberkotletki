package content

/*
TODO: тесты
*/

// Episode представляет эпизод сериала
type Episode struct {
	Id            int    `json:"Id"`             // Уникальный идентификатор
	EpisodeNumber int    `json:"episode_number"` // Номер эпизода
	Duration      int    `json:"Duration"`       // Продолжительность эпизода в минутах
	Description   string `json:"Description"`    // Описание эпизода
}

// Season представляет сезон сериала
type Season struct {
	Content   Content   // Общий контент сезона
	YearStart int       `json:"year_start"` // Год начала сезона
	YearEnd   int       `json:"year_end"`   // Год окончания сезона
	Episodes  []Episode `json:"Episodes"`   // Эпизоды в сезоне
}

// Series представляет сериал
type Series struct {
	Id        int      `json:"Id"`         // Уникальный идентификатор
	Title     string   `json:"Title"`      // Название
	YearStart int      `json:"year_start"` // Год начала сериала
	YearEnd   int      `json:"year_end"`   // Год окончания сериала
	Seasons   []Season `json:"Seasons"`    // Сезоны в сериале
}

// конструкторы

func (e *Episode) NewEpisodeEmpty() *Episode {
	return &Episode{}
}

func (e *Episode) NewEpisodeFull(episodeNumber int, duration int, description string) *Episode {
	return &Episode{
		EpisodeNumber: episodeNumber,
		Duration:      duration,
		Description:   description,
	}
}

func (s *Season) NewSeasonEmpty() *Season {
	return &Season{}
}

func (s *Season) NewSeasonFull(content Content, yearStart int, yearEnd int, episodes []Episode) *Season {
	return &Season{
		Content:   content,
		YearStart: yearStart,
		YearEnd:   yearEnd,
		Episodes:  episodes,
	}
}

func (s *Series) NewSeriesEmpty() *Series {
	return &Series{}
}

func (s *Series) NewSeriesFull(id int, title string, yearStart int, yearEnd int, seasons []Season) *Series {
	return &Series{
		Id:        id,
		Title:     title,
		YearStart: yearStart,
		YearEnd:   yearEnd,
		Seasons:   seasons,
	}
}

// геттеры
func (p Season) GetEpisodes() []Episode {
	if p.Episodes == nil {
		return make([]Episode, 0)
	}
	return p.Episodes
}

func (p Series) Get() []Season {
	if p.Seasons == nil {
		return make([]Season, 0)
	}
	return p.Seasons
}

// проверка равенства
func (c *Episode) Equals(other *Episode) bool {
	return c.Id == other.Id
}

func (c *Season) Equals(other *Season) bool {
	return c.Content.Id == other.Content.Id
}

func (c *Series) Equals(other *Series) bool {
	return c.Id == other.Id
}

// добавление, удаление из слайсов и проверка на наличие элемента

// Season
func (s *Season) AddEpisode(episode Episode) {
	s.Episodes = append(s.Episodes, episode)
}

func (s *Season) RemoveEpisode(episode Episode) {
	for i, e := range s.Episodes {
		if e.Equals(&episode) {
			s.Episodes = append(s.Episodes[:i], s.Episodes[i+1:]...)
			break
		}
	}
}

func (s *Season) HasEpisode(episode Episode) bool {
	for _, e := range s.Episodes {
		if e.Equals(&episode) {
			return true
		}
	}
	return false
}

// Series
func (s *Series) AddSeason(season Season) {
	s.Seasons = append(s.Seasons, season)
}

func (s *Series) RemoveSeason(season Season) {
	for i, se := range s.Seasons {
		if se.Equals(&season) {
			s.Seasons = append(s.Seasons[:i], s.Seasons[i+1:]...)
			break
		}
	}
}

func (s *Series) HasSeason(season Season) bool {
	for _, se := range s.Seasons {
		if se.Equals(&season) {
			return true
		}
	}
	return false
}

// возвращает общую длительность всех эпизодов в сезоне
func (s *Season) TotalDuration() int {
	total := 0
	for _, e := range s.Episodes {
		total += e.GetDuration()
	}
	return total
}

// возвращает средний рейтинг всех сезонов в сериале
func (s *Series) AverageRating() float64 {
	if len(s.Seasons) == 0 {
		return 0
	}

	total := 0.0
	for _, season := range s.Seasons {
		total += season.Content.GetRating()
	}
	return total / float64(len(s.Seasons))
}

// общее количество эпизодов во всех сезонах сериала
func (s *Series) TotalEpisodes() int {
	total := 0
	for _, season := range s.Seasons {
		total += len(season.GetEpisodes())
	}
	return total
}
