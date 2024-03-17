package entity

// Audience представляет аудиторию.
type Audience struct {
	Id        int     `json:"id"`         // Уникальный идентификатор
	Country   Country `json:"country"`    // Страна аудитории
	AudienceT float64 `json:"audience_t"` // Размер аудитории в тысячах
}
