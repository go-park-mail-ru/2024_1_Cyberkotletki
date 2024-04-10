package dto

import "time"

type PreviewContentCard struct {
	ID            int      `json:"id"            example:"1"`
	Title         string   `json:"title"         example:"Бэтмен"`
	OriginalTitle string   `json:"originalTitle" example:"Batman"`
	ReleaseYear   int      `json:"releaseYear"   example:"2020"`
	Country       string   `json:"country"       example:"Россия"`
	Genre         string   `json:"genre"         example:"Боевик"`
	Director      string   `json:"director"      example:"Тарантино"`
	Actors        []string `json:"actors"        example:"Том Хэнкс,Сергей Бодров"`
	Poster        string   `json:"poster"        example:"/static/poster.jpg"`
	Rating        float64  `json:"rating"        example:"9.1"`

	Type string `json:"type" example:"movie"`
	// Поля, которые есть только у фильмов
	Duration int `json:"duration,omitempty" example:"134"`
	// Поля, которые есть только у сериалов
	SeasonsNumber int `json:"seasonsNumber,omitempty" example:"1"`
	YearStart     int `json:"yearStart,omitempty"     example:"2020"`
	YearEnd       int `json:"yearEnd,omitempty"       example:"2021"`
}

type MovieContent struct {
	Premiere time.Time `json:"premiere" example:"2020-01-01"`
	Release  time.Time `json:"release"  example:"2020-01-01"`
	Duration int       `json:"duration" example:"134"`
}

type SeriesContent struct {
	YearStart int      `json:"yearStart" example:"2020"`
	YearEnd   int      `json:"yearEnd"   example:"2020"`
	Seasons   []Season `json:"seasons"`
}

type Season struct {
	ID        int       `json:"id"        example:"1"`
	YearStart int       `json:"yearStart" example:"2020"`
	YearEnd   int       `json:"yearEnd"   example:"2020"`
	Episodes  []Episode `json:"episodes"`
}

type Episode struct {
	ID            int    `json:"id"            example:"1"`
	EpisodeNumber int    `json:"episodeNumber" example:"1"`
	Title         string `json:"title"         example:"Название серии"`
}

type PersonPreview struct {
	ID        int    `json:"id"        example:"1"`
	FirstName string `json:"firstName" example:"Киану"`
	LastName  string `json:"lastName"  example:"Ривз"`
}

type Person struct {
	ID          int                  `json:"id"                    example:"1"`
	FirstName   string               `json:"firstName"             example:"Киану"`
	LastName    string               `json:"lastName"              example:"Ривз"`
	BirthDate   time.Time            `json:"birthDate,omitempty"   example:"1964-09-02"`
	DeathDate   time.Time            `json:"deathDate,omitempty"   example:"2021-09-02"`
	StartCareer time.Time            `json:"startCareer,omitempty" example:"1984-09-02"`
	EndCareer   time.Time            `json:"endCareer,omitempty"   example:"2021-09-02"`
	Sex         string               `json:"sex"                   example:"M"`
	PhotoURL    string               `json:"photoURL,omitempty"    example:"/static/photo.jpg"`
	BirthPlace  string               `json:"birthPlace,omitempty"  example:"Бейрут"`
	Height      int                  `json:"height,omitempty"      example:"185"`
	Spouse      string               `json:"spouse,omitempty"      example:"Алисия Викандер"`
	Children    string               `json:"children,omitempty"    example:"Homer, Bart, Lisa, Maggie"`
	Roles       []PreviewContentCard `json:"roles"`
}

type Content struct {
	ID             int             `json:"id"                  example:"1"`
	Title          string          `json:"title"               example:"Бэтмен"`
	OriginalTitle  string          `json:"originalTitle"       example:"Batman"`
	Slogan         string          `json:"slogan,omitempty"    example:"I'm Batman"`
	Budget         int             `json:"budget,omitempty"    example:"1000000"`
	AgeRestriction int             `json:"ageRestriction"      example:"18"`
	Audience       int             `json:"audience,omitempty"  example:"1000000"`
	Rating         float64         `json:"rating"              example:"9.1"`
	IMDBRating     float64         `json:"imdbRating"          example:"9.1"`
	Description    string          `json:"description"         example:"Описание фильма или сериала"`
	PosterURL      string          `json:"posterURL"           example:"/static/poster.jpg"`
	BoxOffice      int             `json:"boxOffice,omitempty" example:"1000000"`
	Marketing      int             `json:"marketing,omitempty" example:"1000000"`
	Countries      []string        `json:"countries"           example:"Россия,США"`
	Genres         []string        `json:"genres"              example:"Боевик,Драма"`
	Actors         []PersonPreview `json:"actors"`
	Directors      []PersonPreview `json:"directors"`
	Producers      []PersonPreview `json:"producers"`
	Writers        []PersonPreview `json:"writers"`
	Operators      []PersonPreview `json:"operators"`
	Composers      []PersonPreview `json:"composers"`
	Editors        []PersonPreview `json:"editors"`
	Type           string          `json:"type"                example:"movie"`
	Movie          MovieContent    `json:"movie,omitempty"`
	Series         SeriesContent   `json:"series,omitempty"`
}
