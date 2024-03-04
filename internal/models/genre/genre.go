package genre

// Genre представляет жанр.
type Genre struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`    // Название жанра
	RuName string `json:"ru_name"` // Название жанра на русском
}
