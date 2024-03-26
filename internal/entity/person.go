package entity

import (
	"time"
)

type Person struct {
	ID          int          `json:"id"`
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	BirthDate   time.Time    `json:"birth_date"`
	Age         int          `json:"age"`
	DeathDate   time.Time    `json:"death_date,omitempty"`
	StartCareer time.Time    `json:"start_career"`
	EndCareer   time.Time    `json:"end_career,omitempty"`
	Photo       string       `json:"photo"`
	BirthPlace  PlaceOfBirth `json:"birth_place"`
	Genres      []Genre      `json:"genres"`
	Career      []string     `json:"career"`
	Height      int          `json:"height,omitempty"`
	// Жена/муж
	Spouse   string   `json:"spouse,omitempty"`
	Children []string `json:"children,omitempty"`
	Awards   []Award  `json:"awards,omitempty"`
}
