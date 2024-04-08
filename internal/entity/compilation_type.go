package entity

import "unicode/utf8"

type CompilationType struct {
	ID   int    `json:"id" example:"1" format:"int"`
	Type string `json:"type" example:"C" format:"string"`
}

func ValidateCompilationTypeID(id int) error {
	if id < 1 {
		return NewClientError("Рейтинг должен быть в диапазоне от 1 до 10", ErrBadRequest)

	}
	return nil
}

func ValidateCompilationTypeType(compilationType string) error {
	if utf8.RuneCountInString(compilationType) == 0 {
		return NewClientError("Тип подборки не может быть пустым", ErrBadRequest)
	}
	return nil
}

func ValidateCompilationType(id int, compilationType string) error {
	if err := ValidateCompilationTypeID(id); err != nil {
		return err
	}
	if err := ValidateCompilationTypeType(compilationType); err != nil {
		return err
	}
	return nil
}
