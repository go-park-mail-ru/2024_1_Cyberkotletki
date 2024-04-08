package dto

type CompilationContent struct {
	ContentID     int `json:"content_id" example:"1" format:"int"`
	CompilationID int `json:"compilation_id" example:"1" format:"int"`
}
