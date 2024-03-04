package audience

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/country"

// Audience представляет аудиторию.
type Audience struct {
	Id        int             `json:"Id"`         // Уникальный идентификатор
	Country   country.Country `json:"Country"`    // Страна аудитории
	AudienceT float64         `json:"audience_m"` // Размер аудитории в тысячах
}

func (a *Audience) Equals(other *Audience) bool {
	return a.Id == other.Id
}

func (a *Audience) NewAudienceEmpty() *Audience {
	return &Audience{}
}

func (a *Audience) NewAudienceFull(country country.Country, audienceT float64) *Audience {
	return &Audience{
		Country:   country,
		AudienceT: audienceT,
	}
}
