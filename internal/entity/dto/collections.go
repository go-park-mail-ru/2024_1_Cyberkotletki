package dto

type Compilation struct {
	Genre              string `json:"genre" example:"action"`
	ContentIdentifiers []int  `json:"ids"   example:"1,2,3"`
}

type Genres struct {
	Genres []string `json:"genres" example:"action,drama,comedian"`
}
