package small_models

// PlaceOfBirth представляет место рождения.
type PlaceOfBirth struct {
	city    string  `json:"city"`    // Город рождения
	region  string  `json:"region"`  // Регион рождения
	country Country `json:"country"` // Страна рождения
}

// Nomination представляет номинацию на премию.
type Nomination struct {
	title string `json:"title"` // Название номинации
	movie string `json:"movie"` // Фильм, за который дана номинация
}

// Award представляет премию.
type Award struct {
	year       int          `json:"year"`       // Год премии
	awardType  string       `json:"type"`       // Тип премии
	nomination []Nomination `json:"nomination"` // Номинация
}

// Country представляет страну.
type Country struct {
	name string `json:"name"` // Название страны
}

// Genre представляет жанр.
type Genre struct {
	name string `json:"name"` // Название жанра
}

// BoxOffice представляет кассовые сборы.
type BoxOffice struct {
	country Country `json:"country"` // Страна, в которой были сборы
	revenue int     `json:"revenue"` // Сумма сборов
}

// Audience представляет аудиторию.
type Audience struct {
	country   Country `json:"country"`    // Страна аудитории
	audienceT float64 `json:"audience_m"` // Размер аудитории в тысячах
}

func (p *PlaceOfBirth) Equals(other *PlaceOfBirth) bool {
	return p.city == other.city && p.region == other.region && p.country.Equals(&other.country)
}

func (n *Nomination) Equals(other *Nomination) bool {
	return n.title == other.title && n.movie == other.movie
}

func (a *Award) Equals(other *Award) bool {
	return a.year == other.year && a.awardType == other.awardType
}

func (c *Country) Equals(other *Country) bool {
	return c.name == other.name
}

func (g *Genre) Equals(other *Genre) bool {
	return g.name == other.name
}

func (b *BoxOffice) Equals(other *BoxOffice) bool {
	return b.revenue == other.revenue && b.country.Equals(&other.country)
}

func (a *Audience) Equals(other *Audience) bool {
	return a.audienceT == other.audienceT && a.country.Equals(&other.country)
}

// Геттеры и сеттеры для PlaceOfBirth

func (p *PlaceOfBirth) GetCity() string {
	return p.city
}

func (p *PlaceOfBirth) SetCity(city string) {
	p.city = city
}

func (p *PlaceOfBirth) GetRegion() string {
	return p.region
}

func (p *PlaceOfBirth) SetRegion(region string) {
	p.region = region
}

func (p *PlaceOfBirth) GetCountry() Country {
	return p.country
}

func (p *PlaceOfBirth) SetCountry(country Country) {
	p.country = country
}

// Геттеры и сеттеры для Nomination

func (n *Nomination) GetTitle() string {
	return n.title
}

func (n *Nomination) SetTitle(title string) {
	n.title = title
}

func (n *Nomination) GetMovie() string {
	return n.movie
}

func (n *Nomination) SetMovie(movie string) {
	n.movie = movie
}

// Геттеры и сеттеры для  Award

func (a *Award) GetYear() int {
	return a.year
}

func (a *Award) SetYear(year int) {
	a.year = year
}

func (a *Award) GetType() string {
	return a.awardType
}

func (a *Award) SetType(awardType string) {
	a.awardType = awardType
}

// Геттеры и сеттеры для  Country

func (c *Country) GetName() string {
	return c.name
}

func (c *Country) SetName(name string) {
	c.name = name
}

// Геттеры и сеттеры для Genre

func (g *Genre) GetName() string {
	return g.name
}

func (g *Genre) SetName(name string) {
	g.name = name
}

// Геттеры и сеттеры для BoxOffice

func (b *BoxOffice) GetCountry() Country {
	return b.country
}

func (b *BoxOffice) SetCountry(country Country) {
	b.country = country
}

func (b *BoxOffice) GetRevenue() int {
	return b.revenue
}

func (b *BoxOffice) SetRevenue(revenue int) {
	b.revenue = revenue
}

// Геттеры и сеттеры для  Audience

func (a *Audience) GetCountry() Country {
	return a.country
}

func (a *Audience) SetCountry(country Country) {
	a.country = country
}

func (a *Audience) GetAudienceM() float64 {
	return a.audienceT
}

func (a *Audience) SetAudienceM(audienceT float64) {
	a.audienceT = audienceT
}

// Методы для добавления и удаления элементов из слайсов

func (a *Award) AddNomination(nomination Nomination) {
	a.nomination = append(a.nomination, nomination)
}

func (a *Award) RemoveNomination(nomination Nomination) {
	for i, n := range a.nomination {
		if n.Equals(&nomination) {
			a.nomination = append(a.nomination[:i], a.nomination[i+1:]...)
			break
		}
	}
}

// Конструкторы

func NewPlaceOfBirthEmpty() *PlaceOfBirth {
	return &PlaceOfBirth{}
}

func NewPlaceOfBirthFull(city, region string, country Country) *PlaceOfBirth {
	return &PlaceOfBirth{
		city:    city,
		region:  region,
		country: country,
	}
}

func NewNominationEmpty() *Nomination {
	return &Nomination{}
}

func NewNominationFull(title, movie string) *Nomination {
	return &Nomination{
		title: title,
		movie: movie,
	}
}

func NewAwardEmpty() *Award {
	return &Award{}
}

func NewAwardFull(year int, awardType string, nomination []Nomination) *Award {
	return &Award{
		year:       year,
		awardType:  awardType,
		nomination: nomination,
	}
}

func NewCountryEmpty() *Country {
	return &Country{}
}

func NewCountryFull(name string) *Country {
	return &Country{
		name: name,
	}
}

func NewGenreEmpty() *Genre {
	return &Genre{}
}

func NewGenreFull(name string) *Genre {
	return &Genre{
		name: name,
	}
}

func NewBoxOfficeEmpty() *BoxOffice {
	return &BoxOffice{}
}

func NewBoxOfficeFull(country Country, revenue int) *BoxOffice {
	return &BoxOffice{
		country: country,
		revenue: revenue,
	}
}

func NewAudienceEmpty() *Audience {
	return &Audience{}
}

func NewAudienceFull(country Country, audienceT float64) *Audience {
	return &Audience{
		country:   country,
		audienceT: audienceT,
	}
}
