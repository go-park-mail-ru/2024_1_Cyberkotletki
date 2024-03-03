package content

import (
	"time"
)

/*
TODO: тесты
*/

type Film struct {
	Content
	Duration int `json:"Duration"` // Продолжительность
}

func (f *Film) NewFilmEmpty() *Film {
	return &Film{}
}

func (f *Film) NewFilmFull(content Content, year int, duration int) *Film {
	return &Film{
		Content:  content,
		Duration: duration,
	}
}

func (f *Film) Equals(other *Film) bool {
	return f.Id == other.Id
}

func (f *Film) GetDuration() int {
	if f == nil {
		return 0
	}
	return f.Duration
}

// функции

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
