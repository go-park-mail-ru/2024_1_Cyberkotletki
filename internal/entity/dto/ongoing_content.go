package dto

type PreviewOngoingContentList struct {
	OnGoingContentList []*PreviewContent `json:"ongoing_content_list"`
}

type ReleaseYearsResponse struct {
	Years []int `json:"years" example:"[2021, 2022]"`
}

type IsOngoingContentFinishedResponse struct {
	IsReleased bool `json:"is_released" example:"true"`
}
