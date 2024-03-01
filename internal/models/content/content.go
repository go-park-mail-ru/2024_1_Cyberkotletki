package content

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/small_models"
	"time"
)

/*
TODO:
create ( empty, omitempty, full)
adder
remover
haser
interface
*/

// Content представляет основную структуру для хранения информации о контенте.
type Content struct {
	id               int                      `json:"id"`               // Уникальный идентификатор
	title            string                   `json:"title"`            // Название
	country          []small_models.Country   `json:"country"`          // Страны, где был произведен контент
	genres           []small_models.Genre     `json:"genres"`           // Жанры
	directors        []person.Person          `json:"directors"`        // Режиссеры
	writers          []person.Person          `json:"writers"`          // Сценаристы
	producers        []person.Person          `json:"producers"`        // Продюсеры
	cinematographers []person.Person          `json:"cinematographers"` // Операторы
	slogan           string                   `json:"slogan"`           // Слоган
	composers        []person.Person          `json:"composers"`        // Композиторы
	artists          []person.Person          `json:"artists"`          // Художники
	editors          []person.Person          `json:"editors"`          // Редакторы
	budget           int                      `json:"budget"`           // Бюджет
	marketing        int                      `json:"marketing"`        // Маркетинговые затраты
	boxOffices       []small_models.BoxOffice `json:"box_offices"`      // Кассовые сборы
	audiences        []small_models.Audience  `json:"audiences"`        // Аудитория
	premiere         time.Time                `json:"premiere"`         // Дата премьеры
	release          time.Time                `json:"release"`          // Дата выпуска
	ageRestriction   int                      `json:"age_restriction"`  // Возрастное ограничение
	rating           float64                  `json:"rating"`           // Рейтинг
	actors           []person.Person          `json:"actors"`           // Актеры
	dubbing          []person.Person          `json:"dubbing"`          // Дубляж
	description      string                   `json:"description"`      // Описание
	poster           string                   `json:"poster"`           // Постер
	playback         string                   `json:"playback"`         // Воспроизведение на заднем плане небольшоко фрагмента видео
}

// создает новый пустой объект Content
func NewContentEmpty() *Content {
	return &Content{}
}

// создает новый объект Content со всеми данными
func NewContentFull(id int, title string, country []small_models.Country, genres []small_models.Genre, directors []person.Person, writers []person.Person, producers []person.Person, cinematographers []person.Person, slogan string, composers []person.Person, artists []person.Person, editors []person.Person, budget int, marketing int, boxOffices []small_models.BoxOffice, audiences []small_models.Audience, premiere time.Time, release time.Time, ageRestriction int, rating float64, actors []person.Person, dubbing []person.Person, description string, poster string, playback string) *Content {
	return &Content{
		id:               id,
		title:            title,
		country:          country,
		genres:           genres,
		directors:        directors,
		writers:          writers,
		producers:        producers,
		cinematographers: cinematographers,
		slogan:           slogan,
		composers:        composers,
		artists:          artists,
		editors:          editors,
		budget:           budget,
		marketing:        marketing,
		boxOffices:       boxOffices,
		audiences:        audiences,
		premiere:         premiere,
		release:          release,
		ageRestriction:   ageRestriction,
		rating:           rating,
		actors:           actors,
		dubbing:          dubbing,
		description:      description,
		poster:           poster,
		playback:         playback,
	}
}

// Все геттеры для структуры Content

func (c *Content) GetID() int {
	return c.id
}

func (c *Content) GetTitle() string {
	return c.title
}

func (c *Content) GetCountry() []small_models.Country {
	return c.country
}

func (c *Content) GetGenres() []small_models.Genre {
	return c.genres
}

func (c *Content) GetDirectors() []person.Person {
	return c.directors
}

func (c *Content) GetWriters() []person.Person {
	return c.writers
}

func (c *Content) GetProducers() []person.Person {
	return c.producers
}

func (c *Content) GetCinematographers() []person.Person {
	return c.cinematographers
}

func (c *Content) GetSlogan() string {
	return c.slogan
}

func (c *Content) GetComposers() []person.Person {
	return c.composers
}

func (c *Content) GetArtists() []person.Person {
	return c.artists
}

func (c *Content) GetEditors() []person.Person {
	return c.editors
}

func (c *Content) GetBudget() int {
	return c.budget
}

func (c *Content) GetMarketing() int {
	return c.marketing
}

func (c *Content) GetBoxOffices() []small_models.BoxOffice {
	return c.boxOffices
}

func (c *Content) GetAudiences() []small_models.Audience {
	return c.audiences
}

func (c *Content) GetPremiere() time.Time {
	return c.premiere
}

func (c *Content) GetRelease() time.Time {
	return c.release
}

func (c *Content) GetAgeRestriction() int {
	return c.ageRestriction
}

func (c *Content) GetRating() float64 {
	return c.rating
}

func (c *Content) GetActors() []person.Person {
	return c.actors
}

func (c *Content) GetDubbing() []person.Person {
	return c.dubbing
}

func (c *Content) GetDescription() string {
	return c.description
}

func (c *Content) GetPoster() string {
	return c.poster
}

func (c *Content) GetPlayback() string {
	return c.playback
}

// Все сеттеры для структуры Content

func (c *Content) SetID(id int) {
	c.id = id
}

func (c *Content) SetTitle(title string) {
	c.title = title
}

func (c *Content) SetCountry(country []small_models.Country) {
	c.country = country
}

func (c *Content) SetGenres(genres []small_models.Genre) {
	c.genres = genres
}

func (c *Content) SetDirectors(directors []person.Person) {
	c.directors = directors
}

func (c *Content) SetWriters(writers []person.Person) {
	c.writers = writers
}

func (c *Content) SetProducers(producers []person.Person) {
	c.producers = producers
}

func (c *Content) SetCinematographers(cinematographers []person.Person) {
	c.cinematographers = cinematographers
}

func (c *Content) SetSlogan(slogan string) {
	c.slogan = slogan
}

func (c *Content) SetComposers(composers []person.Person) {
	c.composers = composers
}

func (c *Content) SetArtists(artists []person.Person) {
	c.artists = artists
}

func (c *Content) SetEditors(editors []person.Person) {
	c.editors = editors
}

func (c *Content) SetBudget(budget int) {
	c.budget = budget
}

func (c *Content) SetMarketing(marketing int) {
	c.marketing = marketing
}

func (c *Content) SetBoxOffices(boxOffices []small_models.BoxOffice) {
	c.boxOffices = boxOffices
}

func (c *Content) SetAudiences(audiences []small_models.Audience) {
	c.audiences = audiences
}

func (c *Content) SetPremiere(premiere time.Time) {
	c.premiere = premiere
}

func (c *Content) SetRelease(release time.Time) {
	c.release = release
}

func (c *Content) SetAgeRestriction(ageRestriction int) {
	c.ageRestriction = ageRestriction
}

func (c *Content) SetRating(rating float64) {
	c.rating = rating
}

func (c *Content) SetActors(actors []person.Person) {
	c.actors = actors
}

func (c *Content) SetDubbing(dubbing []person.Person) {
	c.dubbing = dubbing
}

func (c *Content) SetDescription(description string) {
	c.description = description
}

func (c *Content) SetPoster(poster string) {
	c.poster = poster
}

func (c *Content) SetPlayback(playback string) {
	c.playback = playback
}

// Методы для добавления и удаления элементов из слайсов

func (c *Content) AddCountry(country small_models.Country) {
	c.country = append(c.country, country)
}

func (c *Content) AddGenre(genre small_models.Genre) {
	c.genres = append(c.genres, genre)
}

func (c *Content) AddDirector(director person.Person) {
	c.directors = append(c.directors, director)
}

func (c *Content) AddWriter(writer person.Person) {
	c.writers = append(c.writers, writer)
}

func (c *Content) AddProduces(producer person.Person) {
	c.producers = append(c.producers, producer)
}

func (c *Content) AddCinematographer(cinematographer person.Person) {
	c.cinematographers = append(c.cinematographers, cinematographer)
}

func (c *Content) AddComposer(composer person.Person) {
	c.composers = append(c.composers, composer)
}

func (c *Content) AddArtists(artist person.Person) {
	c.artists = append(c.artists, artist)
}

func (c *Content) AddEditors(editor person.Person) {
	c.editors = append(c.editors, editor)
}

func (c *Content) AddActors(actor person.Person) {
	c.actors = append(c.actors, actor)
}

func (c *Content) AddDubbing(dubbing person.Person) {
	c.dubbing = append(c.dubbing, dubbing)
}

func (c *Content) AddBoxOffices(boxOffice small_models.BoxOffice) {
	c.boxOffices = append(c.boxOffices, boxOffice)
}

func (c *Content) AddAudiences(audience small_models.Audience) {
	c.audiences = append(c.audiences, audience)
}

func (c *Content) RemoveWriter(writer person.Person) {
	for i, w := range c.writers {
		if w.Equals(&writer) {
			c.writers = append(c.writers[:i], c.writers[i+1:]...)
			break
		}
	}
}

func (c *Content) RemoveProducer(producer person.Person) {
	for i, p := range c.producers {
		if p.Equals(&producer) {
			c.producers = append(c.producers[:i], c.producers[i+1:]...)
			break
		}
	}
}

func (c *Content) RemoveCinematographer(cinematographer person.Person) {
	for i, ci := range c.cinematographers {
		if ci.Equals(&cinematographer) {
			c.cinematographers = append(c.cinematographers[:i], c.cinematographers[i+1:]...)
			break
		}
	}
}

func (c *Content) RemoveComposer(composer person.Person) {
	for i, co := range c.composers {
		if co.Equals(&composer) {
			c.composers = append(c.composers[:i], c.composers[i+1:]...)
			break
		}
	}
}

func (c *Content) RemoveArtist(artist person.Person) {
	for i, a := range c.artists {
		if a.Equals(&artist) {
			c.artists = append(c.artists[:i], c.artists[i+1:]...)
			break
		}
	}
}

func (c *Content) RemoveEditor(editor person.Person) {
	for i, e := range c.editors {
		if e.Equals(&editor) {
			c.editors = append(c.editors[:i], c.editors[i+1:]...)
			break
		}
	}
}

func (c *Content) RemoveActor(actor person.Person) {
	for i, a := range c.actors {
		if a.Equals(&actor) {
			c.actors = append(c.actors[:i], c.actors[i+1:]...)
			break
		}
	}
}

func (c *Content) RemoveDubbing(dubbing person.Person) {
	for i, d := range c.dubbing {
		if d.Equals(&dubbing) {
			c.dubbing = append(c.dubbing[:i], c.dubbing[i+1:]...)
			break
		}
	}
}

func (c *Content) RemoveCountry(country small_models.Country) {
	for i, co := range c.country {
		if co.Equals(&country) {
			c.country = append(c.country[:i], c.country[i+1:]...)
			break
		}
	}
}

func (c *Content) RemoveGenre(genre small_models.Genre) {
	for i, g := range c.genres {
		if g.Equals(&genre) {
			c.genres = append(c.genres[:i], c.genres[i+1:]...)
			break
		}
	}
}

func (c *Content) RemoveBoxOffice(boxOffice small_models.BoxOffice) {
	for i, b := range c.boxOffices {
		if b.Equals(&boxOffice) {
			c.boxOffices = append(c.boxOffices[:i], c.boxOffices[i+1:]...)
			break
		}
	}
}

func (c *Content) RemoveAudience(audience small_models.Audience) {
	for i, a := range c.audiences {
		if a.Equals(&audience) {
			c.audiences = append(c.audiences[:i], c.audiences[i+1:]...)
			break
		}
	}
}

// функции has, которые проверяют наличие элемеента в слайсе

func (c *Content) HasCountry(country small_models.Country) bool {
	for _, co := range c.country {
		if co.Equals(&country) {
			return true
		}
	}
	return false
}

func (c *Content) HasGenre(genre small_models.Genre) bool {
	for _, g := range c.genres {
		if g.Equals(&genre) {
			return true
		}
	}
	return false
}

func (c *Content) HasDirector(director person.Person) bool {
	for _, d := range c.directors {
		if d.Equals(&director) {
			return true
		}
	}
	return false
}

func (c *Content) HasWriter(writer person.Person) bool {
	for _, w := range c.writers {
		if w.Equals(&writer) {
			return true
		}
	}
	return false
}

func (c *Content) HasProducer(producer person.Person) bool {
	for _, p := range c.producers {
		if p.Equals(&producer) {
			return true
		}
	}
	return false
}

func (c *Content) HasCinematographer(cinematographer person.Person) bool {
	for _, ci := range c.cinematographers {
		if ci.Equals(&cinematographer) {
			return true
		}
	}
	return false
}

func (c *Content) HasComposer(composer person.Person) bool {
	for _, co := range c.composers {
		if co.Equals(&composer) {
			return true
		}
	}
	return false
}

func (c *Content) HasArtist(artist person.Person) bool {
	for _, a := range c.artists {
		if a.Equals(&artist) {
			return true
		}
	}
	return false
}

func (c *Content) HasEditor(editor person.Person) bool {
	for _, e := range c.editors {
		if e.Equals(&editor) {
			return true
		}
	}
	return false
}

func (c *Content) HasActor(actor person.Person) bool {
	for _, a := range c.actors {
		if a.Equals(&actor) {
			return true
		}
	}
	return false
}

func (c *Content) HasDubbing(dubbing person.Person) bool {
	for _, d := range c.dubbing {
		if d.Equals(&dubbing) {
			return true
		}
	}
	return false
}

func (c *Content) HasBoxOffice(boxOffice small_models.BoxOffice) bool {
	for _, b := range c.boxOffices {
		if b.Equals(&boxOffice) {
			return true
		}
	}
	return false
}

func (c *Content) HasAudience(audience small_models.Audience) bool {
	for _, a := range c.audiences {
		if a.Equals(&audience) {
			return true
		}
	}
	return false
}
