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
