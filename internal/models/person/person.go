package person

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/small_models"
	"time"
)

type Person struct {
	ID          int                       `json:"id"`
	FirstName   string                    `json:"first_name"`
	LastName    string                    `json:"last_name"`
	BirthDate   time.Time                 `json:"birth_date"`
	Age         int                       `json:"age"`
	DeathDate   *time.Time                `json:"death_date,omitempty"`
	StartCareer time.Time                 `json:"start_career"`
	EndCareer   *time.Time                `json:"end_career,omitempty"`
	Photo       string                    `json:"photo"`
	BirthPlace  small_models.PlaceOfBirth `json:"birth_place"`
	Genres      []small_models.Genre      `json:"genres"`
	Career      []string                  `json:"career"`
	Height      *int                      `json:"height,omitempty"`
	// Жена/муж
	Spouse   *Person              `json:"spouse,omitempty"`
	Children []*Person            `json:"children,omitempty"`
	Awards   []small_models.Award `json:"awards,omitempty"`
}

// конструкторы
// создает пустую структуру
func NewPersonEmpty() *Person {
	return &Person{}
}

// создает структуру со всеми полями, кроме отмеченных как omitempty, те заполняет только обязательные к заполнению поля
func NewPersonOmitempty(id int, firstName string, lastName string, birthDate time.Time, age int, startCareer time.Time, photo string, birthPlace small_models.PlaceOfBirth,
	genres []small_models.Genre, career []string) *Person {
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

// создает новый объект Person со всеми данными
func NewPersonFull(id int, firstName string, lastName string, birthDate time.Time, age int, deathDate *time.Time, startCareer time.Time, endCareer *time.Time, photo string,
	birthPlace small_models.PlaceOfBirth, genres []small_models.Genre, career []string, height *int, spouse *Person, children []*Person, awards []small_models.Award) *Person {
	return &Person{
		ID:          id,
		FirstName:   firstName,
		LastName:    lastName,
		BirthDate:   birthDate,
		Age:         age,
		DeathDate:   deathDate,
		StartCareer: startCareer,
		EndCareer:   endCareer,
		Photo:       photo,
		BirthPlace:  birthPlace,
		Genres:      genres,
		Career:      career,
		Height:      height,
		Spouse:      spouse,
		Children:    children,
		Awards:      awards,
	}
}

//сеттеры

func (p *Person) SetID(id int) {
	p.ID = id
}

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

func (p *Person) SetBirthPlace(birthPlace small_models.PlaceOfBirth) {
	p.BirthPlace = birthPlace
}

func (p *Person) SetGenres(genres []small_models.Genre) {
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

func (p *Person) GetBirthPlace() small_models.PlaceOfBirth {
	return p.BirthPlace
}

func (p *Person) GetGenres() []small_models.Genre {
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

// добавление, удаление и проверка наличия элемента в слайсах

func (p *Person) AddGenre(genre small_models.Genre) {
	p.Genres = append(p.Genres, genre)
}

func (p *Person) RemoveGenre(genre small_models.Genre) {
	for i, g := range p.Genres {
		if g.Equals(&genre) {
			p.Genres = append(p.Genres[:i], p.Genres[i+1:]...)
			break
		}
	}
}

func (p *Person) HasGenre(genre small_models.Genre) bool {
	for _, g := range p.Genres {
		if g.Equals(&genre) {
			return true
		}
	}
	return false
}

func (p *Person) AddCareer(career string) {
	p.Career = append(p.Career, career)
}

func (p *Person) RemoveCareer(career string) {
	for i, c := range p.Career {
		if c == career {
			p.Career = append(p.Career[:i], p.Career[i+1:]...)
			break
		}
	}
}

func (p *Person) HasCareer(career string) bool {
	for _, c := range p.Career {
		if c == career {
			return true
		}
	}
	return false
}

func (p *Person) AddChild(child *Person) {
	p.Children = append(p.Children, child)
}

func (p *Person) RemoveChild(child *Person) {
	for i, ch := range p.Children {
		if ch == child {
			p.Children = append(p.Children[:i], p.Children[i+1:]...)
			break
		}
	}
}

func (p *Person) HasChild(child *Person) bool {
	for _, ch := range p.Children {
		if ch == child {
			return true
		}
	}
	return false
}

func (p *Person) AddAward(award small_models.Award) {
	p.Awards = append(p.Awards, award)
}

func (p *Person) RemoveAward(award small_models.Award) {
	for i, a := range p.Awards {
		if a.Equals(&award) {
			p.Awards = append(p.Awards[:i], p.Awards[i+1:]...)
			break
		}
	}
}

func (p *Person) HasAward(award small_models.Award) bool {
	for _, a := range p.Awards {
		if a.Equals(&award) {
			return true
		}
	}
	return false
}

// функции

// число детей
func (p *Person) NumberOfChildren() int {
	return len(p.Children)
}

// жив или нет
func (p *Person) IsAlive() bool {
	return p.DeathDate == nil
}

// на пенсии или нет
func (p *Person) IsRetired() bool {
	return p.EndCareer != nil
}

// в браке или нет
func (p *Person) IsMarried() bool {
	return p.Spouse != nil
}

// сравнивает персон по имени и дате рождения. Если они совпадают, то возвращается true
func (p *Person) Equals(other *Person) bool {
	return p.FirstName == other.FirstName &&
		p.LastName == other.LastName &&
		p.BirthDate.Equal(other.BirthDate)
}
