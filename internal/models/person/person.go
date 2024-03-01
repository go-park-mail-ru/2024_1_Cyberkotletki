package person

import (
	"encoding/json"
	"time"
)

type PlaceOfBirth struct {
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
}

type Nomination struct {
	Title string `json:"title"`
	Movie string `json:"movie"`
}

type Award struct {
	Year       int        `json:"year"`
	Type       string     `json:"type"`
	Nomination Nomination `json:"nomination"`
}

type Person struct {
	ID          int          `json:"id"`
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	BirthDate   time.Time    `json:"birth_date"`
	Age         int          `json:"age"`
	DeathDate   *time.Time   `json:"death_date,omitempty"`
	StartCareer time.Time    `json:"start_career"`
	EndCareer   *time.Time   `json:"end_career,omitempty"`
	Photo       string       `json:"photo"`
	BirthPlace  PlaceOfBirth `json:"birth_place"`
	Genres      []string     `json:"genres"`
	Career      []string     `json:"career"`
	Height      *int         `json:"height,omitempty"`
	// Жена/муж
	Spouse   *Person   `json:"spouse,omitempty"`
	Children []*Person `json:"children,omitempty"`
	Awards   []Award   `json:"awards,omitempty"`
}

func (p *Person) UpdateName(firstName, lastName string) {
	p.FirstName = firstName
	p.LastName = lastName
}

func (p *Person) UpdateAge(age int) {
	p.Age = age
}

//сеттеры

/*func (p *Person) SetID(id int) {
	p.ID = id
}*/

func (p *Person) SetName(firstName, lastName string) {
	p.FirstName = firstName
	p.LastName = lastName
}

func (p *Person) SetBirthDate(birthDate time.Time) {
	p.BirthDate = birthDate
}

func (p *Person) SetAge(age int) {
	p.Age = age
}

func (p *Person) SetDeathDate(deathDate *time.Time) {
	p.DeathDate = deathDate
}

func (p *Person) SetStartCareer(startCareer time.Time) {
	p.StartCareer = startCareer
}

func (p *Person) SetEndCareer(endCareer *time.Time) {
	p.EndCareer = endCareer
}

func (p *Person) SetPhoto(photo string) {
	p.Photo = photo
}

func (p *Person) SetBirthPlace(birthPlace PlaceOfBirth) {
	p.BirthPlace = birthPlace
}

func (p *Person) SetGenres(genres []string) {
	p.Genres = genres
}

func (p *Person) SetCareer(career []string) {
	p.Career = career
}

func (p *Person) SetHeight(height *int) {
	p.Height = height
}

func (p *Person) SetSpouse(spouse *Person) {
	p.Spouse = spouse
}

func (p *Person) SetChildren(children []*Person) {
	p.Children = children
}

// геттеры

func (p *Person) GetID() int {
	return p.ID
}

func (p *Person) GetFirstName() string {
	return p.FirstName
}

func (p *Person) GetLastName() string {
	return p.LastName
}

func (p *Person) GetBirthDate() time.Time {
	return p.BirthDate
}

func (p *Person) GetAge() int {
	return p.Age
}

func (p *Person) GetDeathDate() *time.Time {
	return p.DeathDate
}

func (p *Person) GetStartCareer() time.Time {
	return p.StartCareer
}

func (p *Person) GetEndCareer() *time.Time {
	return p.EndCareer
}

func (p *Person) GetPhoto() string {
	return p.Photo
}

func (p *Person) GetBirthPlace() PlaceOfBirth {
	return p.BirthPlace
}

func (p *Person) GetGenres() []string {
	return p.Genres
}

func (p *Person) GetCareer() []string {
	return p.Career
}

func (p *Person) GetHeight() *int {
	return p.Height
}

func (p *Person) GetSpouse() *Person {
	return p.Spouse
}

func (p *Person) GetChildren() []*Person {
	return p.Children
}

// функции

// число детей
func (p *Person) NumberOfChildren() int {
	return len(p.Children)
}

// имеет ли кнокретную роль или нет
func (p *Person) HasRole(role string) bool {
	for _, r := range p.Career {
		if r == role {
			return true
		}
	}
	return false
}

// имеет ли отношение к этому жанру
func (p *Person) HasGenre(genre string) bool {
	for _, g := range p.Genres {
		if g == genre {
			return true
		}
	}
	return false
}

// жив или нет
func (p *Person) IsAlive() bool {
	return p.DeathDate == nil
}

func (p *Person) AddChild(child *Person) {
	p.Children = append(p.Children, child)
}

func (p *Person) AddCareerRole(role string) {
	p.Career = append(p.Career, role)
}

func (p *Person) AddGenre(genre string) {
	p.Genres = append(p.Genres, genre)
}

// на пенсии или нет
func (p *Person) IsRetired() bool {
	return p.EndCareer != nil
}

// в браке или нет
func (p *Person) IsMarried() bool {
	return p.Spouse != nil
}

// создает только обязательные параметры у персоны
func CreatePerson(id int, firstName, lastName string, birthDate time.Time, age int, startCareer time.Time, photo string, birthPlace PlaceOfBirth, genres, career []string) *Person {
	return &Person{
		ID:          id,
		FirstName:   firstName,
		LastName:    lastName,
		BirthDate:   birthDate,
		Age:         age,
		StartCareer: startCareer,
		Photo:       photo,
		BirthPlace:  birthPlace,
		Genres:      genres,
		Career:      career,
	}
}

// возвращает ошибку при неверном декодировании
func NewPersonFromJSON(input string) (*Person, error) {
	var person Person
	err := json.Unmarshal([]byte(input), &person)
	if err != nil {
		return nil, err
	}
	return &person, nil
}
