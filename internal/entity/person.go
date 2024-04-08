package entity

import (
	"time"
)

type Person struct {
	ID            int       `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	BirthDate     time.Time `json:"birth_date"`
	DeathDate     time.Time `json:"death_date"`
	StartCareer   time.Time `json:"start_career"`
	EndCareer     time.Time `json:"end_career"`
	Sex           string    `json:"sex"`
	PhotoStaticID int       `json:"photo_static_id"`
	BirthPlace    string    `json:"birth_place"`
	Height        int       `json:"height"`
	// Жена/муж
	Spouse   string   `json:"spouse"`
	Children []string `json:"children"`
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
