package nomination

// Nomination представляет номинацию на премию.
type Nomination struct {
	Id    int    `json:"Id"`    // Уникальный идентификатор
	Title string `json:"Title"` // Название номинации
	Movie string `json:"Movie"` // Фильм, за который дана номинация
}

func (n *Nomination) Equals(other *Nomination) bool {
	return n.Id == other.Id
}

func (a *Nomination) NewNominationEmpty() *Nomination {
	return &Nomination{}
}

func (a *Nomination) NewNominationFull(title, movie string) *Nomination {
	return &Nomination{
		Title: title,
		Movie: movie,
	}
}
