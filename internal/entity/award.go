package entity

// Award представляет премию.
type Award struct {
	ID         int          `json:"id"`         // Уникальный идентификатор
	Year       int          `json:"year"`       // Год премии
	AwardType  string       `json:"type"`       // Тип премии
	Nomination []Nomination `json:"nomination"` // Номинация
}
