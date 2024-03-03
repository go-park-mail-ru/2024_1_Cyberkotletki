package small_models

/*
TODO: тесты
*/

// PlaceOfBirth представляет место рождения.
type PlaceOfBirth struct {
	Id      int     `json:"Id"`      // Уникальный идентификатор
	City    string  `json:"City"`    // Город рождения
	Region  string  `json:"Region"`  // Регион рождения
	Country Country `json:"Country"` // Страна рождения
}

// Nomination представляет номинацию на премию.
type Nomination struct {
	Id    int    `json:"Id"`    // Уникальный идентификатор
	Title string `json:"Title"` // Название номинации
	Movie string `json:"Movie"` // Фильм, за который дана номинация
}

// Award представляет премию.
type Award struct {
	Id         int          `json:"Id"`         // Уникальный идентификатор
	Year       int          `json:"Year"`       // Год премии
	AwardType  string       `json:"type"`       // Тип премии
	Nomination []Nomination `json:"Nomination"` // Номинация
}

// Country представляет страну.
type Country struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"` // Название страны
}

// Genre представляет жанр.
type Genre struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"` // Название жанра
}

// BoxOffice представляет кассовые сборы.
type BoxOffice struct {
	Id      int     `json:"Id"`      // Уникальный идентификатор
	Country Country `json:"Country"` // Страна, в которой были сборы
	Revenue int     `json:"Revenue"` // Сумма сборов
}

// Audience представляет аудиторию.
type Audience struct {
	Id        int     `json:"Id"`         // Уникальный идентификатор
	Country   Country `json:"Country"`    // Страна аудитории
	AudienceT float64 `json:"audience_m"` // Размер аудитории в тысячах
}

func (p *PlaceOfBirth) Equals(other *PlaceOfBirth) bool {
	return p.Id == other.Id
}

func (n *Nomination) Equals(other *Nomination) bool {
	return n.Id == other.Id
}

func (a *Award) Equals(other *Award) bool {
	return a.Id == other.Id
}

func (c *Country) Equals(other *Country) bool {
	return c.Id == other.Id
}

func (g *Genre) Equals(other *Genre) bool {
	return g.Id == other.Id
}

func (b *BoxOffice) Equals(other *BoxOffice) bool {
	return b.Id == other.Id
}

func (a *Audience) Equals(other *Audience) bool {
	return a.Id == other.Id
}

func (p Award) GetGenres() []Nomination {
	if p.Nomination == nil {
		return make([]Nomination, 0)
	}
	return p.Nomination
}

// Методы для добавления и удаления элементов из слайсов

func (a *Award) AddNomination(nomination Nomination) {
	a.Nomination = append(a.Nomination, nomination)
}

func (a *Award) RemoveNomination(nomination Nomination) {
	for i, n := range a.Nomination {
		if n.Equals(&nomination) {
			a.Nomination = append(a.Nomination[:i], a.Nomination[i+1:]...)
			break
		}
	}
}

// Конструкторы

func (a *PlaceOfBirth) NewPlaceOfBirthEmpty() *PlaceOfBirth {
	return &PlaceOfBirth{}
}

func (a *PlaceOfBirth) NewPlaceOfBirthFull(city, region string, country Country) *PlaceOfBirth {
	return &PlaceOfBirth{
		City:    city,
		Region:  region,
		Country: country,
	}
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

func (a *Award) NewAwardEmpty() *Award {
	return &Award{}
}

func (a *Award) NewAwardFull(year int, awardType string, nomination []Nomination) *Award {
	return &Award{
		Year:       year,
		AwardType:  awardType,
		Nomination: nomination,
	}
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

func (a *Genre) NewGenreEmpty() *Genre {
	return &Genre{}
}

func (a *Genre) NewGenreFull(name string, id int) *Genre {
	return &Genre{
		Id:   id,
		Name: name,
	}
}

func (a *BoxOffice) NewBoxOfficeEmpty() *BoxOffice {
	return &BoxOffice{}
}

func (a *BoxOffice) NewBoxOfficeFull(country Country, revenue int) *BoxOffice {
	return &BoxOffice{
		Country: country,
		Revenue: revenue,
	}
}

func (a *Audience) NewAudienceEmpty() *Audience {
	return &Audience{}
}

func (a *Audience) NewAudienceFull(country Country, audienceT float64) *Audience {
	return &Audience{
		Country:   country,
		AudienceT: audienceT,
	}
}
