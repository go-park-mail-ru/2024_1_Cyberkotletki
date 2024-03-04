package country

// Country представляет страну.
type Country struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"` // Название страны
}

func (c *Country) Equals(other *Country) bool {
	return c.Id == other.Id
}

func (a *Country) NewCountryEmpty() *Country {
	return &Country{}
}

func (a *Country) NewCountryFull(name string, id int) *Country {
	return &Country{
		Id:   id,
		Name: name,
	}
}
