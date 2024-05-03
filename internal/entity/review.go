package entity

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

type Review struct {
	ID            int       `json:"id"            example:"1"                    format:"int"`
	AuthorID      int       `json:"authorID"      example:"1"                    format:"int"`
	ContentID     int       `json:"contentID"     example:"1"                    format:"int"`
	ContentRating int       `json:"contentRating" example:"5"                    format:"int"`
	Title         string    `json:"title"         example:"Title"                format:"string"`
	Text          string    `json:"text"          example:"i like it"            format:"string"`
	CreatedAt     time.Time `json:"createdAt"     example:"2022-01-02T15:04:05Z" format:"int"`
	UpdatedAt     time.Time `json:"updatedAt"     example:"2022-01-02T15:04:05Z" format:"int"`
	Rating        int       `json:"rating"        example:"5"                    format:"int"`
	Likes         int       `json:"likes"         example:"5"                    format:"int"`
	Dislikes      int       `json:"dislikes"      example:"5"                    format:"int"`
}

// ValidateReviewRating проверяет, что рейтинг находится в диапазоне от 1 до 10
func ValidateReviewRating(rating int) error {
	if rating < 1 || rating > 10 {
		return errors.New("рейтинг должен быть в диапазоне от 1 до 10")
	}
	return nil
}

// ValidateReviewText проверяет, что длина текста находится в диапазоне от 1 до 10000 символов
func ValidateReviewText(text string) error {
	if utf8.RuneCountInString(strings.TrimSpace(text)) < 1 || utf8.RuneCountInString(strings.TrimSpace(text)) > 10000 {
		return errors.New("количество символов в тексте рецензии должно быть от 1 до 10000")
	}
	return nil
}

// ValidateReviewTitle проверяет, что длина заголовка находится в диапазоне от 1 до 50 символов
func ValidateReviewTitle(title string) error {
	if utf8.RuneCountInString(strings.TrimSpace(title)) < 1 || utf8.RuneCountInString(strings.TrimSpace(title)) > 50 {
		return errors.New("количество символов в заголовке рецензии должно быть от 1 до 50")
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
