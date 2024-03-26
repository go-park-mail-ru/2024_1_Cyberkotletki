package dto

type PreviewContentCard struct {
	Title         string   `json:"title"         example:"Бэтмен"`
	OriginalTitle string   `json:"originalTitle" example:"Batman"`
	ReleaseYear   int      `json:"releaseYear"   example:"2020"`
	Country       string   `json:"country"       example:"Россия"`
	Genre         string   `json:"genre"         example:"Боевик"`
	Director      string   `json:"director"      example:"Тарантино"`
	Actors        []string `json:"actors"        example:"Том Хэнкс,Сергей Бодров"`
	Poster        string   `json:"poster"        example:"/static/poster.jpg"`
	Rating        float64  `json:"rating"        example:"9.1"`
	Duration      int      `json:"duration"      example:"134"`
}
