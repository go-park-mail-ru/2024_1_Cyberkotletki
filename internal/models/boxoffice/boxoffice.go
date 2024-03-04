package boxoffice

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/country"

// BoxOffice представляет кассовые сборы.
type BoxOffice struct {
	Id      int             `json:"Id"`      // Уникальный идентификатор
	Country country.Country `json:"Country"` // Страна, в которой были сборы
	Revenue int             `json:"Revenue"` // Сумма сборов
}

func (b *BoxOffice) Equals(other *BoxOffice) bool {
	return b.Id == other.Id
}

func (a *BoxOffice) NewBoxOfficeEmpty() *BoxOffice {
	return &BoxOffice{}
}

func (a *BoxOffice) NewBoxOfficeFull(country country.Country, revenue int) *BoxOffice {
	return &BoxOffice{
		Country: country,
		Revenue: revenue,
	}
}
