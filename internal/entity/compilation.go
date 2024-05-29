package entity

type Compilation struct {
	ID                int
	Title             string
	CompilationTypeID int
	PosterUploadID    int
}

type CompilationType struct {
	ID   int
	Name string
}
