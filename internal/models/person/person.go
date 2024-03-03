package person

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/small_models"
	"time"
)

/*
TODO: тесты
*/

type Person struct {
	Id          int                       `json:"id"`
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
	Spouse   string               `json:"spouse,omitempty"`
	Children []string             `json:"children,omitempty"`
	Awards   []small_models.Award `json:"awards,omitempty"`
}

// конструкторы
// создает пустую структуру
func (p *Person) NewPersonEmpty() *Person {
	return &Person{}
}

// создает структуру со всеми полями, кроме отмеченных как omitempty, те заполняет только обязательные к заполнению поля
func (p *Person) NewPersonOmitempty(id int, firstName string, lastName string, birthDate time.Time, age int, startCareer time.Time, photo string, birthPlace small_models.PlaceOfBirth,
	genres []small_models.Genre, career []string) *Person {
	return &Person{
		Id:          id,
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
func (p *Person) NewPersonFull(id int, firstName string, lastName string, birthDate time.Time, age int, deathDate *time.Time, startCareer time.Time, endCareer *time.Time, photo string,
	birthPlace small_models.PlaceOfBirth, genres []small_models.Genre, career []string, height *int, spouse string, children []string, awards []small_models.Award) *Person {
	return &Person{
		Id:          id,
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

// геттеры

func (p *Person) GetID() int {
	if p == nil {
		return 0
	}
	return p.Id
}

func (p *Person) GetFirstName() string {
	if p == nil {
		return ""
	}
	return p.FirstName
}

func (p *Person) GetLastName() string {
	if p == nil {
		return ""
	}
	return p.LastName
}

func (p *Person) GetBirthDate() time.Time {
	if p == nil {
		return time.Time{}
	}
	return p.BirthDate
}

func (p *Person) GetAge() int {
	if p == nil {
		return 0
	}
	return p.Age
}

func (p *Person) GetDeathDate() *time.Time {
	if p == nil {
		return nil
	}
	return p.DeathDate
}

func (p *Person) GetStartCareer() time.Time {
	if p == nil {
		return time.Time{}
	}
	return p.StartCareer
}

func (p *Person) GetEndCareer() *time.Time {
	if p == nil {
		return nil
	}
	return p.EndCareer
}

func (p *Person) GetPhoto() string {
	if p == nil {
		return ""
	}
	return p.Photo
}

func (p *Person) GetBirthPlace() small_models.PlaceOfBirth {
	if p == nil {
		return small_models.PlaceOfBirth{}
	}
	return p.BirthPlace
}

func (p *Person) GetGenres() []small_models.Genre {
	if p == nil {
		return nil
	}
	return p.Genres
}

func (p *Person) GetCareer() []string {
	if p == nil {
		return nil
	}
	return p.Career
}

func (p *Person) GetHeight() *int {
	if p == nil {
		return nil
	}
	return p.Height
}

func (p *Person) GetSpouse() string {
	if p == nil {
		return ""
	}
	return p.Spouse
}

func (p *Person) GetChildren() []string {
	if p == nil {
		return nil
	}
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

func (p *Person) AddChild(child string) {
	p.Children = append(p.Children, child)
}

func (p *Person) RemoveChild(child string) {
	for i, ch := range p.Children {
		if ch == child {
			p.Children = append(p.Children[:i], p.Children[i+1:]...)
			break
		}
	}
}

func (p *Person) HasChild(child string) bool {
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
	return p.Spouse != ""
}

// сравнивает персон по имени и дате рождения. Если они совпадают, то возвращается true
func (p *Person) Equals(other *Person) bool {
	return p.Id == other.Id
}
