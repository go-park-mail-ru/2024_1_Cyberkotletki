package entity

import (
	"time"
	"unicode/utf8"
)

type Review struct {
	ID        int       `json:"id"        example:"1"                    format:"int"`
	AuthorID  int       `json:"authorID"  example:"1"                    format:"int"`
	ContentID int       `json:"contentID" example:"1"                    format:"int"`
	Rating    int       `json:"rating"    example:"5"                    format:"int"`
	Title     string    `json:"title"     example:"Title"                format:"string"`
	Text      string    `json:"text"      example:"i like it"            format:"string"`
	CreatedAt time.Time `json:"createdAt" example:"2022-01-02T15:04:05Z" format:"int"`
	UpdatedAt time.Time `json:"updatedAt" example:"2022-01-02T15:04:05Z" format:"int"`
	Likes     int       `json:"likes"     example:"5"                    format:"int"`
	Dislikes  int       `json:"dislikes"  example:"5"                    format:"int"`
}

// ValidateReviewRating проверяет, что рейтинг находится в диапазоне от 1 до 10
func ValidateReviewRating(rating int) error {
	if rating < 1 || rating > 10 {
		return NewClientError("Рейтинг должен быть в диапазоне от 1 до 10", ErrBadRequest)
	}
	return nil
}

// ValidateReviewText проверяет, что длина текста находится в диапазоне от 1 до 10000 символов
func ValidateReviewText(text string) error {
	if utf8.RuneCountInString(text) < 1 || utf8.RuneCountInString(text) > 10000 {
		return NewClientError("Количество символов в тексте рецензии должно быть от 1 до 10000", ErrBadRequest)
	}
	return nil
}

// ValidateReviewTitle проверяет, что длина заголовка находится в диапазоне от 1 до 50 символов
func ValidateReviewTitle(title string) error {
	if utf8.RuneCountInString(title) < 1 || utf8.RuneCountInString(title) > 50 {
		return NewClientError("Количество символов в заголовке рецензии должно быть от 1 до 50", ErrBadRequest)
	}
	return nil
}

func ValidateReview(rating int, title, text string) error {
	if err := ValidateReviewRating(rating); err != nil {
		return err
	}
	if err := ValidateReviewTitle(title); err != nil {
		return err
	}
	if err := ValidateReviewText(text); err != nil {
		return err
	}
	return nil
}
