package entity

import (
	"time"
)

type Person struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	EnName        string    `json:"en_name"`
	BirthDate     time.Time `json:"birth_date,omitempty"`
	DeathDate     time.Time `json:"death_date,omitempty"`
	Sex           string    `json:"sex"`
	Height        int       `json:"height,omitempty"`
	PhotoStaticID int       `json:"photo_static_id,omitempty"`
}

type PersonRole struct {
	PersonID  int  `json:"person_id"`
	Role      Role `json:"role"`
	ContentID int  `json:"content_id"`
}

type Role struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	EnName string `json:"en_name"`
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
