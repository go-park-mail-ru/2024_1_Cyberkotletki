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
private
getter
setter
create ( empty, omitempty, full)
adder
remover
interface

Todo:
- пароли
- валидация
- регистрация
- добавление друга
- др друга
- фильмы друга
*/

type User struct {
	id               int              `json:"id"`                // Уникальный идентификатор
	name             string           `json:"name"`              // Имя пользователя
	email            string           `json:"email"`             // Электронная почта
	passwordHash     string           `json:"password_hash"`     // Хэш пароля пользователя
	birthDate        time.Time        `json:"birth_date"`        // День рождения
	savedFilms       []content.Film   `json:"saved_films"`       // Сохраненные фильмы
	savedSeries      []content.Series `json:"saved_series"`      // Сохраненные сериалы
	savedPersons     []person.Person  `json:"saved_persons"`     // Сохраненные персоны
	friends          []User           `json:"friends"`           // Друзья
	expectedFilms    []content.Film   `json:"expected_films"`    // Ожидаемые фильмы
	registrationDate time.Time        `json:"registration_date"` // Дата регистрации пользователя
}

// SetPassword хеширует пароль и сохраняет его в поле passwordHash.
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.passwordHash = string(hash)
	return nil
}

// CheckPassword сравнивает пароль с хешем и возвращает true, если они совпадают.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.passwordHash), []byte(password))
	return err == nil
}

// Validate проверяет, что все обязательные поля User заполнены и что электронная почта имеет правильный формат.
func (u *User) Validate() error {
	if u.id <= 0 {
		return errors.New("id is required")
	}
	if strings.TrimSpace(u.name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(u.email) == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(u.email); err != nil {
		return errors.New("invalid email format")
	}
	if u.registrationDate.IsZero() {
		return errors.New("registration date is required")
	}
	return nil
}

// NewEmptyUser создает новый пустой объект User.
func NewEmptyUser() *User {
	return &User{}
}

// NewUser создает новый объект User со всеми данными.
func NewUser(id int, name string, email string, passwordHash string, birthDate time.Time, savedFilms []content.Film, savedSeries []content.Series, savedPersons []person.Person, friends []User, expectedFilms []content.Film, registrationDate time.Time) *User {
	return &User{
		id:               id,
		name:             name,
		email:            email,
		passwordHash:     passwordHash,
		birthDate:        birthDate,
		savedFilms:       savedFilms,
		savedSeries:      savedSeries,
		savedPersons:     savedPersons,
		friends:          friends,
		expectedFilms:    expectedFilms,
		registrationDate: registrationDate,
	}
}

func (u *User) GetID() int {
	return u.id
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetPasswordHash() string {
	return u.passwordHash
}

func (u *User) GetBirthDate() time.Time {
	return u.birthDate
}

func (u *User) GetSavedFilms() []content.Film {
	return u.savedFilms
}

func (u *User) GetSavedSeries() []content.Series {
	return u.savedSeries
}

func (u *User) GetSavedPersons() []person.Person {
	return u.savedPersons
}

func (u *User) GetFriends() []User {
	return u.friends
}

func (u *User) GetExpectedFilms() []content.Film {
	return u.expectedFilms
}

func (u *User) GetRegistrationDate() time.Time {
	return u.registrationDate
}

func (u *User) SetID(id int) {
	u.id = id
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) SetEmail(email string) {
	u.email = email
}

func (u *User) SetPasswordHash(passwordHash string) {
	u.passwordHash = passwordHash
}

func (u *User) SetBirthDate(birthDate time.Time) {
	u.birthDate = birthDate
}

func (u *User) SetSavedFilms(savedFilms []content.Film) {
	u.savedFilms = savedFilms
}

func (u *User) SetSavedSeries(savedSeries []content.Series) {
	u.savedSeries = savedSeries
}

func (u *User) SetSavedPersons(savedPersons []person.Person) {
	u.savedPersons = savedPersons
}

func (u *User) SetFriends(friends []User) {
	u.friends = friends
}

func (u *User) SetExpectedFilms(expectedFilms []content.Film) {
	u.expectedFilms = expectedFilms
}

func (u *User) SetRegistrationDate(registrationDate time.Time) {
	u.registrationDate = registrationDate
}

func (u *User) AddSavedFilm(film content.Film) {
	u.savedFilms = append(u.savedFilms, film)
}

func (u *User) AddSavedSeries(series content.Series) {
	u.savedSeries = append(u.savedSeries, series)
}

func (u *User) AddSavedPerson(person person.Person) {
	u.savedPersons = append(u.savedPersons, person)
}

func (u *User) AddFriend(friend User) {
	u.friends = append(u.friends, friend)
}

func (u *User) AddExpectedFilm(film content.Film) {
	u.expectedFilms = append(u.expectedFilms, film)
}

func (u *User) RemoveSavedFilm(film content.Film) {
	for i, f := range u.savedFilms {
		if f.Equals(&film) {
			u.savedFilms = append(u.savedFilms[:i], u.savedFilms[i+1:]...)
			break
		}
	}
}

func (u *User) RemoveSavedSeries(series content.Series) {
	for i, s := range u.savedSeries {
		if s.Equals(&series) {
			u.savedSeries = append(u.savedSeries[:i], u.savedSeries[i+1:]...)
			break
		}
	}
}

func (u *User) RemoveSavedPerson(person person.Person) {
	for i, p := range u.savedPersons {
		if p.Equals(&person) {
			u.savedPersons = append(u.savedPersons[:i], u.savedPersons[i+1:]...)
			break
		}
	}
}

func (u *User) RemoveFriend(friend User) {
	for i, f := range u.friends {
		if f.Equals(&friend) {
			u.friends = append(u.friends[:i], u.friends[i+1:]...)
			break
		}
	}
}

func (u *User) RemoveExpectedFilm(film content.Film) {
	for i, f := range u.expectedFilms {
		if f.Equals(&film) {
			u.expectedFilms = append(u.expectedFilms[:i], u.expectedFilms[i+1:]...)
			break
		}
	}
}

// сохранен ли фильм пользователем.
func (u *User) HasSavedFilm(film content.Film) bool {
	for _, f := range u.savedFilms {
		if f.Equals(&film) {
			return true
		}
	}
	return false
}

func (u *User) HasSavedSeries(series content.Series) bool {
	for _, s := range u.savedSeries {
		if s.Equals(&series) {
			return true
		}
	}
	return false
}

func (u *User) HasSavedPerson(person person.Person) bool {
	for _, p := range u.savedPersons {
		if p.Equals(&person) {
			return true
		}
	}
	return false
}

// сравнивает двух пользователей на равенство
func (u *User) Equals(other *User) bool {
	return (u.id == other.id) || (u.email == u.email)
}

func (u *User) HasFriend(friend User) bool {
	for _, f := range u.friends {
		if f.Equals(&friend) {
			return true
		}
	}
	return false
}

func (u *User) HasExpectedFilm(film content.Film) bool {
	for _, f := range u.expectedFilms {
		if f.Equals(&film) {
			return true
		}
	}
	return false
}

// other

// возвращает день рождения друга.
func (u *User) FriendBirthday(friend *User) (time.Time, error) {
	for _, f := range u.friends {
		if f.id == friend.id {
			return f.birthDate, nil
		}
	}
	return time.Time{}, errors.New("friend not found")
}

// возвращает фильмы, сохраненные другом.
func (u *User) FriendFilms(friend *User) ([]content.Film, error) {
	for _, f := range u.friends {
		if f.id == friend.id {
			return f.savedFilms, nil
		}
	}
	return nil, errors.New("friend not found")
}
