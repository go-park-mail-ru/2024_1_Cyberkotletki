package entity

import (
	"database/sql"
	"time"
)

type Person struct {
	ID            int           `db:"id"`
	Name          string        `db:"name"`
	EnName        string        `db:"en_name"`
	BirthDate     sql.NullTime  `db:"birth_date"`
	DeathDate     sql.NullTime  `db:"death_date"`
	Sex           string        `db:"sex"`
	Height        sql.NullInt64 `db:"height"`
	PhotoStaticID sql.NullInt64 `db:"photo_upload_id"`
}

// GetPhotoStaticID возвращает id статики с фотографией персоны. Если у персоны нет фото, то возвращает 0
func (p Person) GetPhotoStaticID() int {
	if p.PhotoStaticID.Valid {
		return int(p.PhotoStaticID.Int64)
	}
	return 0
}

// GetExamplePerson возвращает пример модели персоны
func GetExamplePerson() Person {
	return Person{
		ID:            1,
		Name:          "Имя",
		EnName:        "Name",
		BirthDate:     sql.NullTime{Time: time.Time{}, Valid: true},
		DeathDate:     sql.NullTime{Time: time.Time{}, Valid: true},
		Sex:           "M",
		Height:        sql.NullInt64{Int64: 175, Valid: true},
		PhotoStaticID: sql.NullInt64{Int64: 1, Valid: true},
	}
}

type PersonRole struct {
	PersonID  int
	Role      Role
	ContentID int
}

type Role struct {
	ID     int
	Name   string
	EnName string
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
