package content

// испорт персоны
// геттеры

import (
	"time"
)

type Film struct {
	Content
}

func NewFilm(content Content) *Film {
	return &Film{Content: content}
}

// сеттеры
func (f *Film) SetTitle(title string) {
	f.Content.Title = title
}

func (f *Film) SetYear(year int) {
	f.Content.Year = year
}

func (f *Film) SetCountry(country Country) {
	f.Content.Country = country
}

func (f *Film) SetGenres(genres []string) {
	f.Content.Genres = genres
}

func (f *Film) SetDirectors(directors []Person) {
	f.Content.Directors = directors
}

func (f *Film) SetWriters(writers []Person) {
	f.Content.Writers = writers
}

func (f *Film) SetProducers(producers []Person) {
	f.Content.Producers = producers
}

func (f *Film) SetCinematographers(cinematographers []Person) {
	f.Content.Cinematographers = cinematographers
}

func (f *Film) SetSlogan(slogan string) {
	f.Content.Slogan = slogan
}

func (f *Film) SetComposers(composers []Person) {
	f.Content.Composers = composers
}

func (f *Film) SetArtists(artists []Person) {
	f.Content.Artists = artists
}

func (f *Film) SetEditors(editors []Person) {
	f.Content.Editors = editors
}

func (f *Film) SetBudget(budget int) {
	f.Content.Budget = budget
}

func (f *Film) SetMarketing(marketing int) {
	f.Content.Marketing = marketing
}

func (f *Film) SetBoxOffices(boxOffices []BoxOffice) {
	f.Content.BoxOffices = boxOffices
}

func (f *Film) SetAudiences(audiences []Audience) {
	f.Content.Audiences = audiences
}

func (f *Film) SetPremiere(premiere time.Time) {
	f.Content.Premiere = premiere
}

func (f *Film) SetRelease(release time.Time) {
	f.Content.Release = release
}

func (f *Film) SetAgeRestriction(ageRestriction int) {
	f.Content.AgeRestriction = ageRestriction
}

func (f *Film) SetDuration(duration int) {
	f.Content.Duration = duration
}

func (f *Film) SetRating(rating float64) {
	f.Content.Rating = rating
}

func (f *Film) SetActors(actors []Person) {
	f.Content.Actors = actors
}

func (f *Film) SetDubbing(dubbing []Person) {
	f.Content.Dubbing = dubbing
}

func (f *Film) SetDescription(description string) {
	f.Content.Description = description
}

func (f *Film) SetPoster(poster string) {
	f.Content.Poster = poster
}

func (f *Film) SetPlayback(playback string) {
	f.Content.Playback = playback
}

// функции

func (f *Film) CalculateProfit() int {
	totalBoxOffice := 0
	for _, boxOffice := range f.Content.BoxOffices {
		totalBoxOffice += boxOffice.Revenue
	}
	return totalBoxOffice - f.Content.Budget - f.Content.Marketing
}

// ROI (Return on Investment)
func (f *Film) CalculateROI() float64 {
	profit := float64(f.CalculateProfit())
	investment := float64(f.Content.Budget + f.Content.Marketing)
	if investment == 0 {
		return 0
	}
	return profit / investment
}

func (f *Film) IsPremiereThisYear() bool {
	currentYear := time.Now().Year()
	return f.Content.Premiere.Year() == currentYear
}

// добавление в слайсы

func (f *Film) AddDirector(director Person) {
	f.Content.Directors = append(f.Content.Directors, director)
}

func (f *Film) AddWriter(writer Person) {
	f.Content.Writers = append(f.Content.Writers, writer)
}

func (f *Film) AddProducer(producer Person) {
	f.Content.Producers = append(f.Content.Producers, producer)
}

func (f *Film) AddCinematographer(cinematographer Person) {
	f.Content.Cinematographers = append(f.Content.Cinematographers, cinematographer)
}

func (f *Film) AddComposer(composer Person) {
	f.Content.Composers = append(f.Content.Composers, composer)
}

func (f *Film) AddArtist(artist Person) {
	f.Content.Artists = append(f.Content.Artists, artist)
}

func (f *Film) AddEditor(editor Person) {
	f.Content.Editors = append(f.Content.Editors, editor)
}

func (f *Film) AddActor(actor Person) {
	f.Content.Actors = append(f.Content.Actors, actor)
}

func (f *Film) AddDubbing(dubbing Person) {
	f.Content.Dubbing = append(f.Content.Dubbing, dubbing)
}

func (f *Film) AddGenre(genre string) {
	f.Content.Genres = append(f.Content.Genres, genre)
}

func (f *Film) AddBoxOffice(boxOffice BoxOffice) {
	f.Content.BoxOffices = append(f.Content.BoxOffices, boxOffice)
}

func (f *Film) AddAudience(audience Audience) {
	f.Content.Audiences = append(f.Content.Audiences, audience)
}
