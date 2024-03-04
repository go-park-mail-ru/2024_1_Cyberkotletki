package person

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/award"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/genre"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/place_of_birth"
	"time"
)

/*
TODO: тесты
*/

type Person struct {
	Id          int                         `json:"id"`
	FirstName   string                      `json:"first_name"`
	LastName    string                      `json:"last_name"`
	BirthDate   time.Time                   `json:"birth_date"`
	Age         int                         `json:"age"`
	DeathDate   time.Time                   `json:"death_date,omitempty"`
	StartCareer time.Time                   `json:"start_career"`
	EndCareer   time.Time                   `json:"end_career,omitempty"`
	Photo       string                      `json:"photo"`
	BirthPlace  place_of_birth.PlaceOfBirth `json:"birth_place"`
	Genres      []genre.Genre               `json:"genres"`
	Career      []string                    `json:"career"`
	Height      int                         `json:"height,omitempty"`
	// Жена/муж
	Spouse   string        `json:"spouse,omitempty"`
	Children []string      `json:"children,omitempty"`
	Awards   []award.Award `json:"awards,omitempty"`
}

// конструкторы
// создает пустую структуру
func (p *Person) NewPersonEmpty() *Person {
	return &Person{}
}

// создает структуру со всеми полями, кроме отмеченных как omitempty, те заполняет только обязательные к заполнению поля
func (p *Person) NewPersonOmitempty(id int, firstName string, lastName string, birthDate time.Time, age int, startCareer time.Time,
	photo string, birthPlace place_of_birth.PlaceOfBirth, genres []genre.Genre, career []string) *Person {
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
func (p *Person) NewPersonFull(id int, firstName string, lastName string, birthDate time.Time, age int, deathDate time.Time, startCareer time.Time, endCareer time.Time, photo string,
	birthPlace place_of_birth.PlaceOfBirth, genres []genre.Genre, career []string, height int, spouse string, children []string, awards []award.Award) *Person {
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

func (p Person) GetGenres() []genre.Genre {
	if p.Genres == nil {
		return make([]genre.Genre, 0)
	}
	return p.Genres
}

func (p Person) GetCareer() []string {
	if p.Career == nil {
		return make([]string, 0)
	}
	return p.Career
}

func (p Person) Getstring() []string {
	if p.Children == nil {
		return make([]string, 0)
	}
	return p.Children
}

func (p Person) GetAwards() []award.Award {
	if p.Genres == nil {
		return make([]award.Award, 0)
	}
	return p.Awards
}

// добавление, удаление и проверка наличия элемента в слайсах

func (p *Person) AddGenre(genre genre.Genre) {
	p.Genres = append(p.Genres, genre)
}

func (p *Person) RemoveGenre(genre genre.Genre) {
	for i, g := range p.Genres {
		if g.Equals(&genre) {
			p.Genres = append(p.Genres[:i], p.Genres[i+1:]...)
			break
		}
	}
}

func (p *Person) HasGenre(genre genre.Genre) bool {
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

func (p *Person) AddAward(award award.Award) {
	p.Awards = append(p.Awards, award)
}

func (p *Person) RemoveAward(award award.Award) {
	for i, a := range p.Awards {
		if a.Equals(&award) {
			p.Awards = append(p.Awards[:i], p.Awards[i+1:]...)
			break
		}
	}
}

func (p *Person) HasAward(award award.Award) bool {
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
	return !p.DeathDate.IsZero()
}

// на пенсии или нет
func (p *Person) IsRetired() bool {
	return p.EndCareer.IsZero()
}

// в браке или нет
func (p *Person) IsMarried() bool {
	return p.Spouse != ""
}

func (p *Person) Equals(other *Person) bool {
	return p.Id == other.Id
}
