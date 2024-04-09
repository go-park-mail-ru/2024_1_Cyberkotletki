package entity

type CompilationContent struct {
	ContentID     int `json:"content_id"     example:"1" format:"int"`
	CompilationID int `json:"compilation_id" example:"1" format:"int"`
}

func ValidateCompilationContentContentID(id int) error {
	if id < 0 {
		return NewClientError("ID контента не может быть отрицательным числом")
	}
	return nil
}

func ValidateCompilationContentCompilationID(id int) error {
	if id < 0 {
		return NewClientError("ID компиляции не может быть отрицательным числом")
	}
	return nil
}

func ValidateCompilationContent(idContent, idCompilation int) error {
	if err := ValidateCompilationContentContentID(idContent); err != nil {
		return err
	}
	if err := ValidateCompilationContentCompilationID(idCompilation); err != nil {
		return err
	}
	return nil
}
