package place_of_birth

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/country"

// PlaceOfBirth представляет место рождения.
type PlaceOfBirth struct {
	Id      int             `json:"id"`      // Уникальный идентификатор
	City    string          `json:"city"`    // Город рождения
	Region  string          `json:"region"`  // Регион рождения
	Country country.Country `json:"country"` // Страна рождения
}
