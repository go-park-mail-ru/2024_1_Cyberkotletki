package content

// импорт персоны
// геттеры
// сеттеры
// нормально сделать совмещение с контентом

type Episode struct {
	Content
	EpisodeNumber int `json:"episode_number"`
}

type Season struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Year     int       `json:"year"`
	Episodes []Episode `json:"episodes"`
}

type Series struct {
	ID               int         `json:"id"`
	Title            string      `json:"title"`
	Year             int         `json:"year"`
	Country          Country     `json:"country"`
	Genres           []string    `json:"genres"`
	Writers          []Person    `json:"writers"`
	Producers        []Person    `json:"producers"`
	Cinematographers []Person    `json:"cinematographers"`
	Slogan           string      `json:"slogan"`
	Composers        []Person    `json:"composers"`
	Artists          []Person    `json:"artists"`
	Editors          []Person    `json:"editors"`
	Budget           float64     `json:"budget"`
	Marketing        float64     `json:"marketing"`
	BoxOffices       []BoxOffice `json:"box_offices"`
	Audiences        []Audience  `json:"audiences"`
	Seasons          []Season    `json:"seasons"`
}

func NewEpisode(content Content, episodeNumber int) *Episode {
	return &Episode{Content: content, EpisodeNumber: episodeNumber}
}

func (e *Episode) GetEpisodeNumber() int {
	return e.EpisodeNumber
}

func (e *Episode) SetEpisodeNumber(newNumber int) {
	e.EpisodeNumber = newNumber
}

// эпизоды

// AddEpisode - добавляет эпизод в сезон
func (s *Season) AddEpisode(episode Episode) {
	s.Episodes = append(s.Episodes, episode)
}

// GetTotalDuration - возвращает общую длительность всех эпизодов в сезоне
func (s *Season) GetTotalDuration() int {
	totalDuration := 0
	for _, episode := range s.Episodes {
		totalDuration += episode.Content.Duration
	}
	return totalDuration
}

// GetAverageRating - возвращает средний рейтинг всех эпизодов в сезоне
func (s *Season) GetAverageRating() float64 {
	if len(s.Episodes) == 0 {
		return 0
	}
	totalRating := 0.0
	for _, episode := range s.Episodes {
		totalRating += episode.Content.Rating
	}
	return totalRating / float64(len(s.Episodes))
}

//сериалы

func NewSeries(id int, title string, year int, country Country, genres []string) *Series {
	return &Series{ID: id, Title: title, Year: year, Country: country, Genres: genres}
}

func (ser *Series) AddSeason(season Season) {
	ser.Seasons = append(ser.Seasons, season)
}

func (ser *Series) GetTotalEpisodes() int {
	totalEpisodes := 0
	for _, season := range ser.Seasons {
		totalEpisodes += len(season.Episodes)
	}
	return totalEpisodes
}

func (ser *Series) GetTotalDuration() int {
	totalDuration := 0
	for _, season := range ser.Seasons {
		totalDuration += season.GetTotalDuration()
	}
	return totalDuration
}

func (ser *Series) GetAverageRating() float64 {
	totalRating := 0.0
	totalEpisodes := 0
	for _, season := range ser.Seasons {
		totalEpisodes += len(season.Episodes)
		totalRating += season.GetAverageRating() * float64(len(season.Episodes))
	}
	if totalEpisodes == 0 {
		return 0
	}
	return totalRating / float64(totalEpisodes)
}
