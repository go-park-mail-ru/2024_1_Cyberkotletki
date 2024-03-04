package user

import (
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	"golang.org/x/crypto/bcrypt"
	"net/mail"
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

func (u *User) ValidatePassword(password string) *exc.Exception {
	if len(password) < 8 {
		return &exc.Exception{
			When:  time.Now(),
			What:  "Пароль должен содержать 8 символов или более",
			Layer: exc.Service,
			Type:  exc.Unprocessable,
		}
	}
	if len(password) > 72 {
		return &exc.Exception{
			When:  time.Now(),
			What:  "Слишком длинный пароль",
			Layer: exc.Service,
			Type:  exc.Unprocessable,
		}
	}
	return nil
}

func (u *User) ValidateEmail(email string) *exc.Exception {
	if _, err := mail.ParseAddress(email); err != nil {
		return &exc.Exception{
			When:  time.Now(),
			What:  "Невалидный Email",
			Layer: exc.Service,
			Type:  exc.Unprocessable,
		}
	}
	return nil
}

// CheckPassword сравнивает пароль с хешем и возвращает true, если они совпадают.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func NewUserEmpty() *User {
	return new(User)
}
