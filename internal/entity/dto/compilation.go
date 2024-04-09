package dto

//TODO^ сделать респонсы

type Compilation struct {
	ID                int    `json:"id" example:"1" format:"int"`
	Title             string `json:"title" example:"The Best" format:"string"`
	CompilationTypeID int    `json:"compilation_type_id" example:"1" format:"int"`
	PosterUploadID    int    `json:"poster_uploadId" example:"1" format:"int"`
}

type CompilationResponse struct {
	Compilation
	ContentLength int `json:"content_length" example:"1" format:"int"`
}

type CompilationResponseList struct {
	Compilations []CompilationResponse `json:"compilations"`
}
