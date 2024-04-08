package entity

type CompilationContent struct {
	ContentID     int `json:"content_id" example:"1" format:"int"`
	CompilationID int `json:"compilation_id" example:"1" format:"int"`
}

func ValidateCompilationContent_ContentID(id int) error {
	if id < 0 {
		return NewClientError("ID контента не может быть отрицательным числом")
	}
	return nil
}

func ValidateCompilationContent_CompilationID(id int) error {
	if id < 0 {
		return NewClientError("ID компиляции не может быть отрицательным числом")
	}
	return nil
}

func ValidateCompilationContent(id_content, id_compilation int) error {
	if err := ValidateCompilationContent_ContentID(id_content); err != nil {
		return err
	}
	if err := ValidateCompilationContent_CompilationID(id_compilation); err != nil {
		return err
	}
	return nil
}
