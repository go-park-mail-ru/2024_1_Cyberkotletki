package audience

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/country"

// Audience представляет аудиторию.
type Audience struct {
	Id        int             `json:"id"`         // Уникальный идентификатор
	Country   country.Country `json:"country"`    // Страна аудитории
	AudienceT float64         `json:"audience_t"` // Размер аудитории в тысячах
}
