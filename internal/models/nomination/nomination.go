package nomination

// Nomination представляет номинацию на премию.
type Nomination struct {
	Id    int    `json:"id"`    // Уникальный идентификатор
	Title string `json:"title"` // Название номинации
	Movie string `json:"movie"` // Фильм, за который дана номинация
}
