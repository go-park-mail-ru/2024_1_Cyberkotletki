package user

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	"time"
)

type User struct {
	ID               int              `json:"id"`
	Name             string           `json:"name"`
	Email            string           `json:"email"`
	PasswordHash     string           `json:"password_hash"`
	SavedFilms       []content.Film   `json:"saved_films"`
	SavedSeries      []content.Series `json:"saved_series"`
	SavedPersons     []content.Person `json:"saved_persons"`
	Friends          []User           `json:"friends"`
	ExpectedFilms    []content.Film   `json:"expected_films"`
	RegistrationDate time.Time        `json:"registration_date"`
}

// условная функция, надо доделать
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// доделать пароль, крайне подохрительно
func NewUser(id int, name string, email string, password string) *User {
	passwordHash := HashPassword(password)
	return &User{ID: id, Name: name, Email: email, PasswordHash: passwordHash}
}

// добавление в слайсы

func (u *User) AddSavedFilm(film content.Film) {
	u.SavedFilms = append(u.SavedFilms, film)
}

func (u *User) AddSavedSeries(series content.Series) {
	u.SavedSeries = append(u.SavedSeries, series)
}

func (u *User) AddSavedPerson(person content.Person) {
	u.SavedPersons = append(u.SavedPersons, person)
}

func (u *User) AddFriend(friend User) {
	u.Friends = append(u.Friends, friend)
}

func (u *User) AddExpectedFilm(film content.Film) {
	u.ExpectedFilms = append(u.ExpectedFilms, film)
}

func (u *User) RemoveSavedFilm(id int) {

}
func (u *User) RemoveSavedSeries(id int) {

}
func (u *User) RemoveSavedPerson(id int) {

}
func (u *User) RemoveFriend(id int) {

}

func (u *User) RemoveExpectedFilm(id int) {

}
