package content

import (
	"time"
)

// Film представляет фильм.
type Film struct {
	Content
	year     int `json:"year"`     // год выпуска
	duration int `json:"duration"` // Продолжительность
}

func NewFilmEmpty() *Film {
	return &Film{}
}

func NewFilmFull(content Content, year int, duration int) *Film {
	return &Film{
		Content:  content,
		year:     year,
		duration: duration,
	}
}

func (f *Film) Equals(other *Film) bool {
	return f.Content.title == other.title &&
		f.year == other.year
}

func (f *Film) GetYear() int {
	return f.year
}

func (f *Film) SetYear(year int) {
	f.year = year
}

func (f *Film) GetDuration() int {
	return f.duration
}

func (f *Film) SetDuration(duration int) {
	f.duration = duration
}

// функции

// IsNewRelease проверяет, является ли фильм новым релизом.
func (f *Film) IsNewRelease(currentYear int) bool {
	return f.year == currentYear
}

// является ли фильм старым (выпущен более 30 лет назад)
func (f *Film) IsOld(currentYear int) bool {
	return currentYear-f.year > 30
}

func (f *Film) CalculateProfit() int {
	totalBoxOffice := 0
	for _, boxOffice := range f.Content.GetBoxOffices() {
		totalBoxOffice += boxOffice.GetRevenue()
	}
	return totalBoxOffice - f.Content.GetBudget() - f.Content.GetMarketing()
}

// ROI (Return on Investment)
func (f *Film) CalculateROI() float64 {
	profit := float64(f.CalculateProfit())
	investment := float64(f.Content.GetBudget() + f.Content.GetMarketing())
	if investment == 0 {
		return 0
	}
	return profit / investment
}

// в этом году премьера или нет
func (f *Film) IsPremiereThisYear() bool {
	currentYear := time.Now().Year()
	return f.Content.GetPremiere().Year() == currentYear
}
