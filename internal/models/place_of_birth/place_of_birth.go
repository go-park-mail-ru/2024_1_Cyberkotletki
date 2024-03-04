package place_of_birth

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/country"

// PlaceOfBirth представляет место рождения.
type PlaceOfBirth struct {
	Id      int             `json:"Id"`      // Уникальный идентификатор
	City    string          `json:"City"`    // Город рождения
	Region  string          `json:"Region"`  // Регион рождения
	Country country.Country `json:"Country"` // Страна рождения
}

func (p *PlaceOfBirth) Equals(other *PlaceOfBirth) bool {
	return p.Id == other.Id
}

func (a *PlaceOfBirth) NewPlaceOfBirthEmpty() *PlaceOfBirth {
	return &PlaceOfBirth{}
}

func (a *PlaceOfBirth) NewPlaceOfBirthFull(city, region string, country country.Country) *PlaceOfBirth {
	return &PlaceOfBirth{
		City:    city,
		Region:  region,
		Country: country,
	}
}
