package dto

import "time"

type PreviewOngoingContent struct {
	ID          int       `json:"id"          example:"1"`
	Title       string    `json:"title"       example:"Бэтмен"`
	Genres      []string  `json:"genre"       example:"Боевик"`
	Poster      string    `json:"poster"      example:"/static/poster.jpg"`
	ReleaseDate time.Time `json:"releaseDate" example:"2022-01-02T15:04:05Z"`

	Type string `json:"type" example:"movie"`
}

type PreviewOngoingContentList struct {
	OnGoingContentList []*PreviewOngoingContent `json:"ongoing_content_list"`
}

type ReleaseYearsResponse struct {
	Years []int `json:"years" example:"[2021, 2022]"`
}

type IsOngoingContentFinishedResponse struct {
	IsFinished bool `json:"is_finished" example:"true"`
}
