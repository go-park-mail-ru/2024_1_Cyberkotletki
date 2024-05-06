package entity

import (
	"errors"
	"strings"
	"time"
	"unicode/utf8"
)

type Review struct {
	ID            int       `db:"id"`
	AuthorID      int       `db:"user_id"`
	ContentID     int       `db:"content_id"`
	ContentRating int       `db:"content_rating"`
	Title         string    `db:"title"`
	Text          string    `db:"text"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	Rating        int       `db:"rating"`
	Likes         int       `db:"likes"`
	Dislikes      int       `db:"dislikes"`
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
