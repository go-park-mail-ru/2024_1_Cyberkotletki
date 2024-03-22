package entity

// PlaceOfBirth представляет место рождения.
type PlaceOfBirth struct {
	ID      int     `json:"id"`      // Уникальный идентификатор
	City    string  `json:"city"`    // Город рождения
	Region  string  `json:"region"`  // Регион рождения
	Country Country `json:"country"` // Страна рождения
}
