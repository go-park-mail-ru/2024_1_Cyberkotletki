package dto

type CompilationType struct {
	ID   int    `json:"id" example:"1" format:"int"`
	Type string `json:"type" example:"C" format:"string"`
}
