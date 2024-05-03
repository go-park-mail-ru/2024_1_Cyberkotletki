package dto

type Compilation struct {
	ID                int    `json:"id"                  example:"1"                 format:"int"`
	Title             string `json:"title"               example:"The Best"          format:"string"`
	CompilationTypeID int    `json:"compilation_type_id" example:"1"                 format:"int"`
	PosterURL         string `json:"poster"              example:"static/poster.jpg" format:"string"`
}

type CompilationType struct {
	ID   int    `json:"id"   example:"1" format:"int"`
	Type string `json:"type" example:"C" format:"string"`
}

type CompilationTypeResponseList struct {
	CompilationTypes []CompilationType `json:"compilation_types"`
}

type CompilationResponseList struct {
	Compilations []Compilation `json:"compilations"`
}

type CompilationResponse struct {
	Compilation   Compilation `json:"compilation"`
	ContentIDs    []int       `json:"content_ids"`
	ContentLength int         `json:"content_length"`
	Page          int         `json:"page"`
	PerPage       int         `json:"per_page"`
	TotalPages    int         `json:"total_pages"`
}
