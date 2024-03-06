package award

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/nomination"

// Award представляет премию.
type Award struct {
	Id         int                     `json:"id"`         // Уникальный идентификатор
	Year       int                     `json:"year"`       // Год премии
	AwardType  string                  `json:"type"`       // Тип премии
	Nomination []nomination.Nomination `json:"nomination"` // Номинация
}
