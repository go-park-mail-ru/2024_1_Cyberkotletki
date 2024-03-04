package content

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/audience"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/award"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/boxoffice"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/country"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/genre"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	"time"
)

/*
TODO: тесты
*/

// Content представляет основную структуру для хранения информации о контенте.
type Content struct {
	Id               int                   `json:"Id"`               // Уникальный идентификатор
	Title            string                `json:"Title"`            // Название
	OriginalTitle    string                `json:"original_title"`   // Название
	Country          []country.Country     `json:"country"`          // Страны, где был произведен контент
	Genres           []genre.Genre         `json:"Genres"`           // Жанры
	Directors        []person.Person       `json:"Directors"`        // Режиссеры
	Writers          []person.Person       `json:"Writers"`          // Сценаристы
	Producers        []person.Person       `json:"Producers"`        // Продюсеры
	Cinematographers []person.Person       `json:"Cinematographers"` // Операторы
	Slogan           string                `json:"Slogan"`           // Слоган
	Composers        []person.Person       `json:"Composers"`        // Композиторы
	Artists          []person.Person       `json:"Artists"`          // Художники
	Editors          []person.Person       `json:"Editors"`          // Редакторы
	Budget           int                   `json:"Budget"`           // Бюджет
	Marketing        int                   `json:"Marketing"`        // Маркетинговые затраты
	BoxOffices       []boxoffice.BoxOffice `json:"box_offices"`      // Кассовые сборы
	Audiences        []audience.Audience   `json:"Audiences"`        // Аудитория
	Premiere         time.Time             `json:"Premiere"`         // Дата премьеры
	Release          time.Time             `json:"Release"`          // Дата выпуска
	AgeRestriction   int                   `json:"age_restriction"`  // Возрастное ограничение
	Rating           float64               `json:"Rating"`           // Рейтинг
	Actors           []person.Person       `json:"Actors"`           // Актеры
	Dubbing          []person.Person       `json:"Dubbing"`          // Дубляж
	Awards           []award.Award         `json:"awards,omitempty"` // Награды
	Description      string                `json:"Description"`      // Описание
	Poster           string                `json:"Poster"`           // Постер
	Playback         string                `json:"Playback"`         // Воспроизведение на заднем плане небольшоко фрагмента видео
}

// создает новый пустой объект Content
func (c *Content) NewContentEmpty() *Content {
	return &Content{}
}

// создает новый объект Content со всеми данными
func (c *Content) NewContentFull(id int, title string, country []country.Country, genres []genre.Genre, directors []person.Person,
	writers []person.Person, producers []person.Person, cinematographers []person.Person, slogan string, composers []person.Person,
	artists []person.Person, editors []person.Person, budget int, marketing int, boxOffices []boxoffice.BoxOffice, audiences []audience.Audience,
	premiere time.Time, release time.Time, ageRestriction int, rating float64, actors []person.Person, dubbing []person.Person, awards []award.Award, description string, poster string, playback string) *Content {
	return &Content{
		Id:               id,
		Title:            title,
		Country:          country,
		Genres:           genres,
		Directors:        directors,
		Writers:          writers,
		Producers:        producers,
		Cinematographers: cinematographers,
		Slogan:           slogan,
		Composers:        composers,
		Artists:          artists,
		Editors:          editors,
		Budget:           budget,
		Marketing:        marketing,
		BoxOffices:       boxOffices,
		Audiences:        audiences,
		Premiere:         premiere,
		Release:          release,
		AgeRestriction:   ageRestriction,
		Rating:           rating,
		Actors:           actors,
		Dubbing:          dubbing,
		Awards:           awards,
		Description:      description,
		Poster:           poster,
		Playback:         playback,
	}
}

// Все геттеры для структуры Content

func (c Content) GetCountry() []country.Country {
	if c.Country == nil {
		return make([]country.Country, 0)
	}
	return c.Country
}

func (c Content) GetAwards() []award.Award {
	if c.Awards == nil {
		return make([]award.Award, 0)
	}
	return c.Awards
}

func (c Content) GetGenres() []genre.Genre {
	if c.Genres == nil {
		return make([]genre.Genre, 0)
	}
	return c.Genres
}

func (c Content) GetDirectors() []person.Person {
	if c.Directors == nil {
		return make([]person.Person, 0)
	}
	return c.Directors
}

func (c Content) GetWriters() []person.Person {
	if c.Writers == nil {
		return make([]person.Person, 0)
	}
	return c.Writers
}

func (c Content) GetProducers() []person.Person {
	if c.Producers == nil {
		return make([]person.Person, 0)
	}
	return c.Producers
}

func (c Content) GetCinematographers() []person.Person {
	if c.Cinematographers == nil {
		return make([]person.Person, 0)
	}
	return c.Cinematographers
}

func (c Content) GetComposers() []person.Person {
	if c.Composers == nil {
		return make([]person.Person, 0)
	}
	return c.Composers
}

func (c Content) GetArtists() []person.Person {
	if c.Artists == nil {
		return make([]person.Person, 0)
	}
	return c.Artists
}

func (c Content) GetEditors() []person.Person {
	if c.Editors == nil {
		return make([]person.Person, 0)
	}
	return c.Editors
}

func (c Content) GetBoxOffices() []boxoffice.BoxOffice {
	if c.BoxOffices == nil {
		return make([]boxoffice.BoxOffice, 0)
	}
	return c.BoxOffices
}

func (c Content) GetAudiences() []audience.Audience {
	if c.Audiences == nil {
		return make([]audience.Audience, 0)
	}
	return c.Audiences
}

func (c Content) GetActors() []person.Person {
	if c.Actors == nil {
		return make([]person.Person, 0)
	}
	return c.Actors
}

func (c Content) GetDubbing() []person.Person {
	if c.Dubbing == nil {
		return make([]person.Person, 0)
	}
	return c.Dubbing
}

// Методы для добавления и удаления элементов из слайсов

func (c Content) AddCountry(country country.Country) {
	c.Country = append(c.Country, country)
}

func (c Content) AddGenre(genre genre.Genre) {
	c.Genres = append(c.Genres, genre)
}

func (c Content) AddDirector(director person.Person) {
	c.Directors = append(c.Directors, director)
}

func (c Content) AddAward(award award.Award) {
	c.Awards = append(c.Awards, award)
}

func (c Content) AddWriter(writer person.Person) {
	c.Writers = append(c.Writers, writer)
}

func (c Content) AddProduces(producer person.Person) {
	c.Producers = append(c.Producers, producer)
}

func (c Content) AddCinematographer(cinematographer person.Person) {
	c.Cinematographers = append(c.Cinematographers, cinematographer)
}

func (c Content) AddComposer(composer person.Person) {
	c.Composers = append(c.Composers, composer)
}

func (c Content) AddArtists(artist person.Person) {
	c.Artists = append(c.Artists, artist)
}

func (c Content) AddEditors(editor person.Person) {
	c.Editors = append(c.Editors, editor)
}

func (c Content) AddActors(actor person.Person) {
	c.Actors = append(c.Actors, actor)
}

func (c Content) AddDubbing(dubbing person.Person) {
	c.Dubbing = append(c.Dubbing, dubbing)
}

func (c Content) AddBoxOffices(boxOffice boxoffice.BoxOffice) {
	c.BoxOffices = append(c.BoxOffices, boxOffice)
}

func (c Content) AddAudiences(audience audience.Audience) {
	c.Audiences = append(c.Audiences, audience)
}

func (c Content) RemoveWriter(writer person.Person) {
	for i, w := range c.Writers {
		if w.Equals(&writer) {
			c.Writers = append(c.Writers[:i], c.Writers[i+1:]...)
			break
		}
	}
}

func (c Content) RemoveAward(award award.Award) {
	for i, w := range c.Awards {
		if w.Equals(&award) {
			c.Awards = append(c.Awards[:i], c.Awards[i+1:]...)
			break
		}
	}
}

func (c Content) RemoveProducer(producer person.Person) {
	for i, p := range c.Producers {
		if p.Equals(&producer) {
			c.Producers = append(c.Producers[:i], c.Producers[i+1:]...)
			break
		}
	}
}

func (c Content) RemoveCinematographer(cinematographer person.Person) {
	for i, ci := range c.Cinematographers {
		if ci.Equals(&cinematographer) {
			c.Cinematographers = append(c.Cinematographers[:i], c.Cinematographers[i+1:]...)
			break
		}
	}
}

func (c Content) RemoveComposer(composer person.Person) {
	for i, co := range c.Composers {
		if co.Equals(&composer) {
			c.Composers = append(c.Composers[:i], c.Composers[i+1:]...)
			break
		}
	}
}

func (c Content) RemoveArtist(artist person.Person) {
	for i, a := range c.Artists {
		if a.Equals(&artist) {
			c.Artists = append(c.Artists[:i], c.Artists[i+1:]...)
			break
		}
	}
}

func (c Content) RemoveEditor(editor person.Person) {
	for i, e := range c.Editors {
		if e.Equals(&editor) {
			c.Editors = append(c.Editors[:i], c.Editors[i+1:]...)
			break
		}
	}
}

func (c Content) RemoveActor(actor person.Person) {
	for i, a := range c.Actors {
		if a.Equals(&actor) {
			c.Actors = append(c.Actors[:i], c.Actors[i+1:]...)
			break
		}
	}
}

func (c Content) RemoveDubbing(dubbing person.Person) {
	for i, d := range c.Dubbing {
		if d.Equals(&dubbing) {
			c.Dubbing = append(c.Dubbing[:i], c.Dubbing[i+1:]...)
			break
		}
	}
}

func (c Content) RemoveCountry(country country.Country) {
	for i, co := range c.Country {
		if co.Equals(&country) {
			c.Country = append(c.Country[:i], c.Country[i+1:]...)
			break
		}
	}
}

func (c Content) RemoveGenre(genre genre.Genre) {
	for i, g := range c.Genres {
		if g.Equals(&genre) {
			c.Genres = append(c.Genres[:i], c.Genres[i+1:]...)
			break
		}
	}
}

func (c Content) RemoveBoxOffice(boxOffice boxoffice.BoxOffice) {
	for i, b := range c.BoxOffices {
		if b.Equals(&boxOffice) {
			c.BoxOffices = append(c.BoxOffices[:i], c.BoxOffices[i+1:]...)
			break
		}
	}
}

func (c Content) RemoveAudience(audience audience.Audience) {
	for i, a := range c.Audiences {
		if a.Equals(&audience) {
			c.Audiences = append(c.Audiences[:i], c.Audiences[i+1:]...)
			break
		}
	}
}

// функции has, которые проверяют наличие элемеента в слайсе

func (c Content) HasCountry(country country.Country) bool {
	for _, co := range c.Country {
		if co.Equals(&country) {
			return true
		}
	}
	return false
}

func (c Content) HasAward(award award.Award) bool {
	for _, co := range c.Awards {
		if co.Equals(&award) {
			return true
		}
	}
	return false
}

func (c Content) HasGenre(genre genre.Genre) bool {
	for _, g := range c.Genres {
		if g.Equals(&genre) {
			return true
		}
	}
	return false
}

func (c Content) HasDirector(director person.Person) bool {
	for _, d := range c.Directors {
		if d.Equals(&director) {
			return true
		}
	}
	return false
}

func (c Content) HasWriter(writer person.Person) bool {
	for _, w := range c.Writers {
		if w.Equals(&writer) {
			return true
		}
	}
	return false
}

func (c Content) HasProducer(producer person.Person) bool {
	for _, p := range c.Producers {
		if p.Equals(&producer) {
			return true
		}
	}
	return false
}

func (c Content) HasCinematographer(cinematographer person.Person) bool {
	for _, ci := range c.Cinematographers {
		if ci.Equals(&cinematographer) {
			return true
		}
	}
	return false
}

func (c Content) HasComposer(composer person.Person) bool {
	for _, co := range c.Composers {
		if co.Equals(&composer) {
			return true
		}
	}
	return false
}

func (c Content) HasArtist(artist person.Person) bool {
	for _, a := range c.Artists {
		if a.Equals(&artist) {
			return true
		}
	}
	return false
}

func (c Content) HasEditor(editor person.Person) bool {
	for _, e := range c.Editors {
		if e.Equals(&editor) {
			return true
		}
	}
	return false
}

func (c Content) HasActor(actor person.Person) bool {
	for _, a := range c.Actors {
		if a.Equals(&actor) {
			return true
		}
	}
	return false
}

func (c Content) HasDubbing(dubbing person.Person) bool {
	for _, d := range c.Dubbing {
		if d.Equals(&dubbing) {
			return true
		}
	}
	return false
}

func (c Content) HasBoxOffice(boxOffice boxoffice.BoxOffice) bool {
	for _, b := range c.BoxOffices {
		if b.Equals(&boxOffice) {
			return true
		}
	}
	return false
}

func (c Content) HasAudience(audience audience.Audience) bool {
	for _, a := range c.Audiences {
		if a.Equals(&audience) {
			return true
		}
	}
	return false
}
