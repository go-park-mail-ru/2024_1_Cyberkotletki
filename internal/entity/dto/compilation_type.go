package dto

type CompilationType struct {
	ID   int    `json:"id" example:"1" format:"int"`
	Type string `json:"type" example:"C" format:"string"`
}

type CompilationTypeResponse struct {
	CompilationType
}

type CompilationTypeResponseList struct {
	CompilationTypes []CompilationTypeResponse `json:"compilation_types"`
}
