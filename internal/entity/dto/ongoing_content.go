package dto

import "time"

type PreviewOngoingContentCardVertical struct {
	ID          int       `json:"id"     example:"1"`
	Title       string    `json:"title"  example:"Бэтмен"`
	Genres      []string  `json:"genre"  example:"Боевик"`
	Poster      string    `json:"poster" example:"/static/poster.jpg"`
	ReleaseDate time.Time `json:"releaseDate" example:"2022-01-02T15:04:05Z"`

	Type string `json:"type" example:"movie"`
	// Поля, которые есть только у фильмов
	ReleaseYear int `json:"releaseYear" example:"2020"`
	// Поля, которые есть только у сериалов
	YearStart int `json:"yearStart,omitempty" example:"2020"`
	YearEnd   int `json:"yearEnd,omitempty"   example:"2021"`

	Movie  OngoingMovie  `json:"movie,omitempty"`
	Series OngoingSeries `json:"series,omitempty"`
}

type OngoingMovie struct {
	Premiere time.Time `json:"premiere" example:"2020-01-01"`
	Duration int       `json:"duration" example:"134"`
}

type OngoingSeries struct {
	YearStart int      `json:"yearStart" example:"2020"`
	YearEnd   int      `json:"yearEnd"   example:"2020"`
	Seasons   []Season `json:"seasons"`
}
