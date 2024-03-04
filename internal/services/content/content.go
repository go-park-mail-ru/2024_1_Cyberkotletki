package content

import exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"

type PreviewInfoData struct {
	Title         string   `json:"title" example:"Бэтмен"`
	OriginalTitle string   `json:"original_title" example:"Batman"`
	ReleaseYear   int16    `json:"release_year" example:"2020"`
	Country       string   `json:"country" example:"Россия"`
	Genre         string   `json:"genre" example:"Боевик"`
	Director      string   `json:"director" example:"Тарантино"`
	Actors        []string `json:"actors" example:"Том Хэнкс,Сергей Бодров"`
	Poster        string   `json:"poster" example:"/static/poster.jpg"`
}

func GetContentPreviewInfo(contentId int64) (PreviewInfoData, *exc.Exception) {
	return PreviewInfoData{
		Title:         "Бэтмен",
		OriginalTitle: "Batman",
		ReleaseYear:   2020,
		Country:       "Россия",
		Genre:         "Боевик",
		Director:      "Тарантино",
		Actors:        []string{"Том Хэнкс", "Сергей Бодров"},
		Poster:        "/static/poster.jpg",
	}, nil
}
