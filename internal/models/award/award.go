package award

import "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/nomination"

// Award представляет премию.
type Award struct {
	Id         int                     `json:"Id"`         // Уникальный идентификатор
	Year       int                     `json:"Year"`       // Год премии
	AwardType  string                  `json:"type"`       // Тип премии
	Nomination []nomination.Nomination `json:"Nomination"` // Номинация
}

func (p Award) GetGenres() []nomination.Nomination {
	if p.Nomination == nil {
		return make([]nomination.Nomination, 0)
	}
	return p.Nomination
}
func (a *Award) Equals(other *Award) bool {
	return a.Id == other.Id
}

func (a *Award) AddNomination(nomination nomination.Nomination) {
	a.Nomination = append(a.Nomination, nomination)
}

func (a *Award) RemoveNomination(nomination nomination.Nomination) {
	for i, n := range a.Nomination {
		if n.Equals(&nomination) {
			a.Nomination = append(a.Nomination[:i], a.Nomination[i+1:]...)
			break
		}
	}
}

func (a *Award) NewAwardEmpty() *Award {
	return &Award{}
}

func (a *Award) NewAwardFull(year int, awardType string, nomination []nomination.Nomination) *Award {
	return &Award{
		Year:       year,
		AwardType:  awardType,
		Nomination: nomination,
	}
}
