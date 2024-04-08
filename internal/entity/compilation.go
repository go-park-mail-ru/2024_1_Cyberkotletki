package entity

//TODO убрать коллекции

import "unicode/utf8"

type Compilation struct {
	ID                int    `json:"id" example:"1" format:"int"`
	Title             string `json:"title" example:"The Best" format:"string"`
	CompilationTypeID int    `json:"compilation_type_id" example:"1" format:"int"`
	PosterUploadID    int    `json:"poster_uploadId" example:"1" format:"int"`
}

func ValidateCompilationID(id int) error {
	if id <= 0 {
		return NewClientError("ID контента не может быть отрицательным числом", ErrBadRequest)
	}
	return nil
}

func ValidateCompilationTitle(title string) error {
	if utf8.RuneCountInString(title) == 0 {
		return NewClientError("Название компиляции не может быть пустым", ErrBadRequest)
	}
	return nil
}
func ValidateCompilationCompilationTypeID(id int) error {
	if id <= 0 {
		return NewClientError("ID типа компиляции не может быть отрицательным числом", ErrBadRequest)
	}
	return nil
}

func ValidateCompilationPosterUploadID(id int) error {
	if id <= 0 {
		return NewClientError("ID постера не может быть отрицательным числом", ErrBadRequest)
	}
	return nil
}

func ValidateCompilation(id int, title string, compilationTypeID int, posterUploadID int) error {
	if err := ValidateCompilationID(id); err != nil {
		return err
	}
	if err := ValidateCompilationTitle(title); err != nil {
		return err
	}
	if err := ValidateCompilationTypeID(compilationTypeID); err != nil {
		return err
	}
	if err := ValidateCompilationPosterUploadID(posterUploadID); err != nil {
		return err
	}
	return nil
}
