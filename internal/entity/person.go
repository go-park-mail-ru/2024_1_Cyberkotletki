package entity

import (
	"time"
)

type Person struct {
	ID            int       `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	BirthDate     time.Time `json:"birth_date,omitempty"`
	DeathDate     time.Time `json:"death_date,omitempty"`
	StartCareer   time.Time `json:"start_career,omitempty"`
	EndCareer     time.Time `json:"end_career,omitempty"`
	Sex           string    `json:"sex"`
	PhotoStaticID int       `json:"photo_static_id,omitempty"`
	BirthPlace    string    `json:"birth_place,omitempty"`
	Height        int       `json:"height,omitempty"`
	// Жена/муж
	Spouse   string `json:"spouse,omitempty"`
	Children string `json:"children,omitempty"`
}

const (
	RoleActor    = "actor"
	RoleDirector = "director"
	RoleProducer = "producer"
	RoleWriter   = "writer"
	RoleOperator = "operator"
	RoleComposer = "composer"
	RoleEditor   = "editor"
)
