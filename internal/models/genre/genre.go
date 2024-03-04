package genre

// Genre представляет жанр.
type Genre struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"` // Название жанра
}

func (g *Genre) Equals(other *Genre) bool {
	return g.Id == other.Id
}
func (a *Genre) NewGenreEmpty() *Genre {
	return &Genre{}
}

func (a *Genre) NewGenreFull(name string, id int) *Genre {
	return &Genre{
		Id:   id,
		Name: name,
	}
}
