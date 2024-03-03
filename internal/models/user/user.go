package user

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
	"strings"
	"time"
)

/*
TODO:
тесты
*/

type User struct {
	Id               int              `json:"Id"`                // Уникальный идентификатор
	Name             string           `json:"Name"`              // Имя пользователя
	Email            string           `json:"Email"`             // Электронная почта
	PasswordHash     string           `json:"password_hash"`     // Хэш пароля пользователя
	BirthDate        time.Time        `json:"birth_date"`        // День рождения
	SavedFilms       []content.Film   `json:"saved_films"`       // Сохраненные фильмы
	SavedSeries      []content.Series `json:"saved_series"`      // Сохраненные сериалы
	SavedPersons     []person.Person  `json:"saved_persons"`     // Сохраненные персоны
	Friends          []User           `json:"Friends"`           // Друзья
	ExpectedFilms    []content.Film   `json:"expected_films"`    // Ожидаемые фильмы
	RegistrationDate time.Time        `json:"registration_date"` // Дата регистрации пользователя
}

// SetPassword хеширует пароль и сохраняет его в поле PasswordHash.
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword сравнивает пароль с хешем и возвращает true, если они совпадают.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// Validate проверяет, что все обязательные поля User заполнены и что электронная почта имеет правильный формат.
func (u *User) Validate() error {
	if u.Id <= 0 {
		return errors.New("Id is required")
	}
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("Name is required")
	}
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("Email is required")
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("invalid Email format")
	}
	if u.RegistrationDate.IsZero() {
		return errors.New("registration date is required")
	}
	return nil
}

func (u *User) NewUserEmpty() *User {
	return &User{}
}

func (u *User) NewUserFull(id int, name string, email string, passwordHash string, birthDate time.Time, savedFilms []content.Film,
	savedSeries []content.Series, savedPersons []person.Person, friends []User, expectedFilms []content.Film,
	registrationDate time.Time) *User {
	return &User{
		Id:               id,
		Name:             name,
		Email:            email,
		PasswordHash:     passwordHash,
		BirthDate:        birthDate,
		SavedFilms:       savedFilms,
		SavedSeries:      savedSeries,
		SavedPersons:     savedPersons,
		Friends:          friends,
		ExpectedFilms:    expectedFilms,
		RegistrationDate: registrationDate,
	}
}

func (u *User) GetID() int {
	if u == nil {
		return 0
	}
	return u.Id
}

func (u *User) GetName() string {
	if u == nil {
		return ""
	}
	return u.Name
}

func (u *User) GetEmail() string {
	if u == nil {
		return ""
	}
	return u.Email
}

func (u *User) GetPasswordHash() string {
	if u == nil {
		return ""
	}
	return u.PasswordHash
}

func (u *User) GetBirthDate() time.Time {
	if u == nil {
		return time.Time{}
	}
	return u.BirthDate
}

func (u *User) GetSavedFilms() []content.Film {
	if u == nil {
		return nil
	}
	return u.SavedFilms
}

func (u *User) GetSavedSeries() []content.Series {
	if u == nil {
		return nil
	}
	return u.SavedSeries
}

func (u *User) GetSavedPersons() []person.Person {
	if u == nil {
		return nil
	}
	return u.SavedPersons
}

func (u *User) GetFriends() []User {
	if u == nil {
		return nil
	}
	return u.Friends
}

func (u *User) GetExpectedFilms() []content.Film {
	if u == nil {
		return nil
	}
	return u.ExpectedFilms
}

func (u *User) GetRegistrationDate() time.Time {
	if u == nil {
		return time.Time{}
	}
	return u.RegistrationDate
}

func (u *User) AddSavedFilm(film content.Film) {
	u.SavedFilms = append(u.SavedFilms, film)
}

func (u *User) AddSavedSeries(series content.Series) {
	u.SavedSeries = append(u.SavedSeries, series)
}

func (u *User) AddSavedPerson(person person.Person) {
	u.SavedPersons = append(u.SavedPersons, person)
}

func (u *User) AddFriend(friend User) {
	u.Friends = append(u.Friends, friend)
}

func (u *User) AddExpectedFilm(film content.Film) {
	u.ExpectedFilms = append(u.ExpectedFilms, film)
}

func (u *User) RemoveSavedFilm(film content.Film) {
	for i, f := range u.SavedFilms {
		if f.Equals(&film) {
			u.SavedFilms = append(u.SavedFilms[:i], u.SavedFilms[i+1:]...)
			break
		}
	}
}

func (u *User) RemoveSavedSeries(series content.Series) {
	for i, s := range u.SavedSeries {
		if s.Equals(&series) {
			u.SavedSeries = append(u.SavedSeries[:i], u.SavedSeries[i+1:]...)
			break
		}
	}
}

func (u *User) RemoveSavedPerson(person person.Person) {
	for i, p := range u.SavedPersons {
		if p.Equals(&person) {
			u.SavedPersons = append(u.SavedPersons[:i], u.SavedPersons[i+1:]...)
			break
		}
	}
}

func (u *User) RemoveFriend(friend User) {
	for i, f := range u.Friends {
		if f.Equals(&friend) {
			u.Friends = append(u.Friends[:i], u.Friends[i+1:]...)
			break
		}
	}
}

func (u *User) RemoveExpectedFilm(film content.Film) {
	for i, f := range u.ExpectedFilms {
		if f.Equals(&film) {
			u.ExpectedFilms = append(u.ExpectedFilms[:i], u.ExpectedFilms[i+1:]...)
			break
		}
	}
}

// сохранен ли фильм пользователем.
func (u *User) HasSavedFilm(film content.Film) bool {
	for _, f := range u.SavedFilms {
		if f.Equals(&film) {
			return true
		}
	}
	return false
}

func (u *User) HasSavedSeries(series content.Series) bool {
	for _, s := range u.SavedSeries {
		if s.Equals(&series) {
			return true
		}
	}
	return false
}

func (u *User) HasSavedPerson(person person.Person) bool {
	for _, p := range u.SavedPersons {
		if p.Equals(&person) {
			return true
		}
	}
	return false
}

// сравнивает двух пользователей на равенство
func (u *User) Equals(other *User) bool {
	return (u.Id == other.Id)
}

func (u *User) HasFriend(friend User) bool {
	for _, f := range u.Friends {
		if f.Equals(&friend) {
			return true
		}
	}
	return false
}

func (u *User) HasExpectedFilm(film content.Film) bool {
	for _, f := range u.ExpectedFilms {
		if f.Equals(&film) {
			return true
		}
	}
	return false
}

// возвращает день рождения друга.
func (u *User) FriendBirthday(friend *User) (time.Time, error) {
	for _, f := range u.Friends {
		if f.Id == friend.Id {
			return f.BirthDate, nil
		}
	}
	return time.Time{}, errors.New("friend not found")
}

// возвращает фильмы, сохраненные другом.
func (u *User) FriendFilms(friend *User) ([]content.Film, error) {
	for _, f := range u.Friends {
		if f.Id == friend.Id {
			return f.SavedFilms, nil
		}
	}
	return nil, errors.New("friend not found")
}
