package content

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/content"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
)

type PreviewInfoData struct {
	Title         string   `json:"title" example:"Бэтмен"`
	OriginalTitle string   `json:"original_title" example:"Batman"`
	ReleaseYear   int      `json:"release_year" example:"2020"`
	Country       string   `json:"country" example:"Россия"`
	Genre         string   `json:"genre" example:"Боевик"`
	Director      string   `json:"director" example:"Тарантино"`
	Actors        []string `json:"actors" example:"Том Хэнкс,Сергей Бодров"`
	Poster        string   `json:"poster" example:"/static/poster.jpg"`
	Rating        float64  `json:"rating" example:"9.1"`
	Duration      int      `json:"duration" example:"134"`
}

func GetContentPreviewInfo(contentId int) (*PreviewInfoData, *exc.Exception) {
	if film, err := content.FilmsDatabase.GetFilm(contentId); err != nil {
		return nil, err
	} else {
		return &PreviewInfoData{
			Title:         film.Title,
			OriginalTitle: film.OriginalTitle,
			ReleaseYear:   film.Release.Year(),
			Country:       film.Country[0].Name,
			Genre:         film.Genres[0].Name,
			Director:      fmt.Sprintf("%s %s", film.Directors[0].FirstName, film.Directors[0].LastName),
			Actors: []string{
				fmt.Sprintf("%s %s", film.Actors[0].FirstName, film.Actors[0].LastName),
				fmt.Sprintf("%s %s", film.Actors[1].FirstName, film.Actors[1].LastName),
			},
			Poster:   film.Poster,
			Rating:   film.Rating,
			Duration: film.Duration,
		}, nil
	}
}
