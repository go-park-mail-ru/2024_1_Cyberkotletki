package entity

// Nomination представляет номинацию на премию.
type Nomination struct {
	ID    int    `json:"id"`    // Уникальный идентификатор
	Title string `json:"title"` // Название номинации
	Movie string `json:"movie"` // Фильм, за который дана номинация
}
