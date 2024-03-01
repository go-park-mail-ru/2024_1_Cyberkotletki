package content

// нормально сделать контент

import (
	"time"
)

type Country struct {
	Name string `json:"name"`
}

// сборы
type BoxOffice struct {
	Country Country `json:"country"`
	Revenue int     `json:"revenue"`
}

// аудитория в миллионах
type Audience struct {
	Country   Country `json:"country"`
	AudienceM float64 `json:"audience_m"`
}

type Content struct {
	ID               int         `json:"id"`
	Title            string      `json:"title"`
	Year             int         `json:"year"`
	Country          Country     `json:"country"`
	Genres           []string    `json:"genres"`
	Directors        []Person    `json:"directors"`
	Writers          []Person    `json:"writers"`
	Producers        []Person    `json:"producers"`
	Cinematographers []Person    `json:"cinematographers"`
	Slogan           string      `json:"slogan"`
	Composers        []Person    `json:"composers"`
	Artists          []Person    `json:"artists"`
	Editors          []Person    `json:"editors"`
	Budget           int         `json:"budget"`
	Marketing        int         `json:"marketing"`
	BoxOffices       []BoxOffice `json:"box_offices"`
	Audiences        []Audience  `json:"audiences"`
	Premiere         time.Time   `json:"premiere"`
	Release          time.Time   `json:"release"`
	AgeRestriction   int         `json:"age_restriction"`
	Duration         int         `json:"duration"`
	Rating           float64     `json:"rating"`
	Actors           []Person    `json:"actors"`
	Dubbing          []Person    `json:"dubbing"`
	Description      string      `json:"description"`
	Poster           string      `json:"poster"`
	Playback         string      `json:"playback"`
}
