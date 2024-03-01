package content

/*
TODO:
private
getter
setter
create ( empty, omitempty, full)
adder
remover
has-er
*/

// Episode представляет эпизод сериала
type Episode struct {
	id            int    `json:"id"`             // Уникальный идентификатор
	episodeNumber int    `json:"episode_number"` // Номер эпизода
	duration      int    `json:"duration"`       // Продолжительность эпизода в минутах
	description   string `json:"description"`    // Описание эпизода
}

// Season представляет сезон сериала
type Season struct {
	content   Content   // Общий контент сезона
	yearStart int       `json:"year_start"` // Год начала сезона
	yearEnd   int       `json:"year_end"`   // Год окончания сезона
	episodes  []Episode `json:"episodes"`   // Эпизоды в сезоне
}

// Series представляет сериал
type Series struct {
	id        int      `json:"id"`         // Уникальный идентификатор
	title     string   `json:"title"`      // Название
	yearStart int      `json:"year_start"` // Год начала сериала
	yearEnd   int      `json:"year_end"`   // Год окончания сериала
	seasons   []Season `json:"seasons"`    // Сезоны в сериале
}

// конструкторы

func NewEpisodeEmpty() *Episode {
	return &Episode{}
}

func NewEpisodeFull(episodeNumber int, duration int, description string) *Episode {
	return &Episode{
		episodeNumber: episodeNumber,
		duration:      duration,
		description:   description,
	}
}

func NewSeasonEmpty() *Season {
	return &Season{}
}

func NewSeasonFull(content Content, yearStart int, yearEnd int, episodes []Episode) *Season {
	return &Season{
		content:   content,
		yearStart: yearStart,
		yearEnd:   yearEnd,
		episodes:  episodes,
	}
}

func NewSeriesEmpty() *Series {
	return &Series{}
}

func NewSeriesFull(id int, title string, yearStart int, yearEnd int, seasons []Season) *Series {
	return &Series{
		id:        id,
		title:     title,
		yearStart: yearStart,
		yearEnd:   yearEnd,
		seasons:   seasons,
	}
}

// геттеры и сеттеры

// Episode
func (e *Episode) GetID() int {
	return e.id
}

func (e *Episode) SetID(id int) {
	e.id = id
}

func (e *Episode) GetEpisodeNumber() int {
	return e.episodeNumber
}

func (e *Episode) SetEpisodeNumber(episodeNumber int) {
	e.episodeNumber = episodeNumber
}

func (e *Episode) GetDuration() int {
	return e.duration
}

func (e *Episode) SetDuration(duration int) {
	e.duration = duration
}

func (e *Episode) GetDescription() string {
	return e.description
}

func (e *Episode) SetDescription(description string) {
	e.description = description
}

// Season
func (s *Season) GetContent() Content {
	return s.content
}

func (s *Season) SetContent(content Content) {
	s.content = content
}

func (s *Season) GetYearStart() int {
	return s.yearStart
}

func (s *Season) SetYearStart(yearStart int) {
	s.yearStart = yearStart
}

func (s *Season) GetYearEnd() int {
	return s.yearEnd
}

func (s *Season) SetYearEnd(yearEnd int) {
	s.yearEnd = yearEnd
}

func (s *Season) GetEpisodes() []Episode {
	return s.episodes
}

func (s *Season) SetEpisodes(episodes []Episode) {
	s.episodes = episodes
}

// Series
func (s *Series) GetID() int {
	return s.id
}

func (s *Series) SetID(id int) {
	s.id = id
}

func (s *Series) GetTitle() string {
	return s.title
}

func (s *Series) SetTitle(title string) {
	s.title = title
}

func (s *Series) GetYearStart() int {
	return s.yearStart
}

func (s *Series) SetYearStart(yearStart int) {
	s.yearStart = yearStart
}

func (s *Series) GetYearEnd() int {
	return s.yearEnd
}

func (s *Series) SetYearEnd(yearEnd int) {
	s.yearEnd = yearEnd
}

func (s *Series) GetSeasons() []Season {
	return s.seasons
}

func (s *Series) SetSeasons(seasons []Season) {
	s.seasons = seasons
}

// проверка равенства
func (c *Episode) Equals(other *Episode) bool {
	return c.episodeNumber == other.episodeNumber &&
		c.duration == other.duration
}

func (c *Season) Equals(other *Season) bool {
	return c.content.title == other.content.title &&
		c.yearStart == other.yearEnd
}

func (c *Series) Equals(other *Series) bool {
	return c.title == other.title &&
		c.yearStart == other.yearStart &&
		c.yearEnd == c.yearEnd
}

// добавление, удаление из слайсов и проверка на наличие элемента

// Season
func (s *Season) AddEpisode(episode Episode) {
	s.episodes = append(s.episodes, episode)
}

func (s *Season) RemoveEpisode(episode Episode) {
	for i, e := range s.episodes {
		if e.Equals(&episode) {
			s.episodes = append(s.episodes[:i], s.episodes[i+1:]...)
			break
		}
	}
}

func (s *Season) HasEpisode(episode Episode) bool {
	for _, e := range s.episodes {
		if e.Equals(&episode) {
			return true
		}
	}
	return false
}

// Series
func (s *Series) AddSeason(season Season) {
	s.seasons = append(s.seasons, season)
}

func (s *Series) RemoveSeason(season Season) {
	for i, se := range s.seasons {
		if se.Equals(&season) {
			s.seasons = append(s.seasons[:i], s.seasons[i+1:]...)
			break
		}
	}
}

func (s *Series) HasSeason(season Season) bool {
	for _, se := range s.seasons {
		if se.Equals(&season) {
			return true
		}
	}
	return false
}

// возвращает общую длительность всех эпизодов в сезоне
func (s *Season) TotalDuration() int {
	total := 0
	for _, e := range s.episodes {
		total += e.GetDuration()
	}
	return total
}

// возвращает средний рейтинг всех сезонов в сериале
func (s *Series) AverageRating() float64 {
	if len(s.seasons) == 0 {
		return 0
	}

	total := 0.0
	for _, season := range s.seasons {
		total += season.content.GetRating()
	}
	return total / float64(len(s.seasons))
}

// общее количество эпизодов во всех сезонах сериала
func (s *Series) TotalEpisodes() int {
	total := 0
	for _, season := range s.seasons {
		total += len(season.GetEpisodes())
	}
	return total
}
