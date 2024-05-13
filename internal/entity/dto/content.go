package dto

import "time"

type PreviewContent struct {
	ID            int      `json:"id"                      example:"1"`
	Title         string   `json:"title"                   example:"Бэтмен"`
	OriginalTitle string   `json:"originalTitle,omitempty" example:"Batman"`
	Country       string   `json:"country"                 example:"Россия"`
	Genre         string   `json:"genre"                   example:"Боевик"`
	Director      string   `json:"director"                example:"Тарантино"`
	Actors        []string `json:"actors"                  example:"Том Хэнкс,Сергей Бодров"`
	Poster        string   `json:"poster"                  example:"/static/poster.jpg"`
	Rating        float64  `json:"rating"                  example:"9.1"`

	Type string `json:"type" example:"movie"`
	// Поля, которые есть только у фильмов
	Duration    int `json:"duration,omitempty" example:"134"`
	ReleaseYear int `json:"release,omitempty"  example:"2020"`
	// Поля, которые есть только у сериалов
	SeasonsNumber int `json:"seasonsNumber,omitempty" example:"1"`
	YearStart     int `json:"yearStart,omitempty"     example:"2020"`
	YearEnd       int `json:"yearEnd,omitempty"       example:"2021"`
}

type PreviewContentCardVertical struct {
	ID     int      `json:"id"     example:"1"`
	Title  string   `json:"title"  example:"Бэтмен"`
	Genres []string `json:"genre"  example:"Боевик"`
	Poster string   `json:"poster" example:"/static/poster.jpg"`
	Rating float64  `json:"rating" example:"9.1"`

	Type string `json:"type" example:"movie"`
	// Поля, которые есть только у фильмов
	ReleaseYear int `json:"releaseYear,omitempty" example:"2020"`
	// Поля, которые есть только у сериалов
	YearStart int `json:"yearStart,omitempty" example:"2020"`
	YearEnd   int `json:"yearEnd,omitempty"   example:"2021"`
}

type MovieContent struct {
	Premiere *time.Time `json:"premiere,omitempty" example:"2020-01-01"`
	Duration int        `json:"duration"           example:"134"`
}

type SeriesContent struct {
	YearStart int      `json:"yearStart" example:"2020"`
	YearEnd   int      `json:"yearEnd"   example:"2020"`
	Seasons   []Season `json:"seasons"`
}

type Season struct {
	ID       int       `json:"id"       example:"1"`
	Episodes []Episode `json:"episodes"`
}

type Episode struct {
	ID            int    `json:"id"            example:"1"`
	EpisodeNumber int    `json:"episodeNumber" example:"1"`
	Title         string `json:"title"         example:"Название серии"`
	Duration      int    `json:"duration"      example:"45"`
}

type PersonPreview struct {
	ID     int    `json:"id"     example:"1"`
	Name   string `json:"name"   example:"Киану Ривз"`
	EnName string `json:"enName" example:"Keanu Reeves"`
}

type PersonPreviewWithPhoto struct {
	ID       int    `json:"id"       example:"1"`
	Name     string `json:"name"     example:"Киану Ривз"`
	EnName   string `json:"enName"   example:"Keanu Reeves"`
	PhotoURL string `json:"photoURL" example:"/static/photo.jpg"`
}

type Person struct {
	ID        int        `json:"id"                  example:"1"`
	Name      string     `json:"name"                example:"Киану Ривз"`
	EnName    string     `json:"enName"              example:"Keanu Reeves"`
	BirthDate *time.Time `json:"birthDate,omitempty" example:"1964-09-02"`
	DeathDate *time.Time `json:"deathDate,omitempty" example:"2021-09-02"`
	Sex       string     `json:"sex"                 example:"M"`
	PhotoURL  string     `json:"photoURL,omitempty"  example:"/static/photo.jpg"`
	Height    int        `json:"height,omitempty"    example:"185"`

	Roles map[string][]PreviewContentCardVertical `json:"roles"`
}

type Content struct {
	ID             int             `json:"id"               example:"1"`
	Title          string          `json:"title"            example:"Бэтмен"`
	OriginalTitle  string          `json:"originalTitle"    example:"Batman"`
	Slogan         string          `json:"slogan,omitempty" example:"I'm Batman"`
	Budget         string          `json:"budget,omitempty" example:"1000000"`
	AgeRestriction int             `json:"ageRestriction"   example:"18"`
	Rating         float64         `json:"rating"           example:"9.1"`
	IMDBRating     float64         `json:"imdbRating"       example:"9.1"`
	Description    string          `json:"description"      example:"Описание фильма или сериала"`
	Facts          []string        `json:"facts"            example:"Факты о фильме или сериале"`
	TrailerLink    string          `json:"trailerLink"      example:"https://www.youtube.com/watch?v=123456"`
	BackdropURL    string          `json:"backdropURL"      example:"/static/backdrop.jpg"`
	PicturesURL    []string        `json:"picturesURL"      example:"/static/picture1.jpg,/static/picture2.jpg"`
	PosterURL      string          `json:"posterURL"        example:"/static/poster.jpg"`
	Countries      []string        `json:"countries"        example:"Россия,США"`
	Genres         []string        `json:"genres"           example:"Боевик,Драма"`
	Actors         []PersonPreview `json:"actors"`
	Directors      []PersonPreview `json:"directors"`
	Producers      []PersonPreview `json:"producers"`
	Writers        []PersonPreview `json:"writers"`
	Operators      []PersonPreview `json:"operators"`
	Composers      []PersonPreview `json:"composers"`
	Editors        []PersonPreview `json:"editors"`
	Type           string          `json:"type"             example:"movie"`
	Movie          MovieContent    `json:"movie,omitempty"`
	Series         SeriesContent   `json:"series,omitempty"`
}
