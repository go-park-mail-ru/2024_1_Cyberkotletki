package user

import (
	"crypto/sha256"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/random"
	"golang.org/x/crypto/pbkdf2"
	"regexp"
	"time"
)

const (
	PasswordHashIterations = 4096
	PasswordHashKeyLength  = 128
)

var PasswordHashFunction = sha256.New

type User struct {
	Id               int              `json:"id"`    // Уникальный идентификатор
	Name             string           `json:"name"`  // Имя пользователя
	Email            string           `json:"email"` // Электронная почта
	PasswordHash     string           // Хэш пароля пользователя
	PasswordSalt     string           // Соль для генерации хэша пароля
	BirthDate        time.Time        `json:"birth_date"`        // День рождения
	SavedFilms       []content.Film   `json:"saved_films"`       // Сохраненные фильмы
	SavedSeries      []content.Series `json:"saved_series"`      // Сохраненные сериалы
	SavedPersons     []person.Person  `json:"saved_persons"`     // Сохраненные персоны
	Friends          []User           `json:"friends"`           // Друзья
	ExpectedFilms    []content.Film   `json:"expected_films"`    // Ожидаемые фильмы
	RegistrationDate time.Time        `json:"registration_date"` // Дата регистрации пользователя
}

// ValidatePassword проверяет валидность пароля.
//
// Критерии валидности пароля:
//
// 1) пароль содержит от 8 до 32 символов включительно
//
// 2) пароль не содержит ничего, кроме латинских букв, цифр и символов !@#$%^&*
//
// 3) пароль пароль содержит как минимум одну заглавную, одну строчную букву, одну цифру и один из символов !@#$%^&*
func ValidatePassword(password string) error {
	switch {
	case len(password) < 8:
		return exc.New(exc.Service, exc.BadRequest, "пароль должен содержать не менее 8 символов")
	case len(password) > 32:
		return exc.New(exc.Service, exc.BadRequest, "пароль должен содержать не менее 8 символов")
	case !regexp.MustCompile("^[!@#$%^&*\\w]+$").MatchString(password):
		return exc.New(exc.Service, exc.BadRequest, "пароль должен состоять из латинских букв, цифр и специальных символов !@#$%^&*")
	case !regexp.MustCompile("[A-Z]").MatchString(password):
		return exc.New(exc.Service, exc.BadRequest, "пароль должен содержать как минимум одну заглавную букву")
	case !regexp.MustCompile("[a-z]").MatchString(password):
		return exc.New(exc.Service, exc.BadRequest, "пароль должен содержать как минимум одну строчную букву")
	case !regexp.MustCompile("\\d").MatchString(password):
		return exc.New(exc.Service, exc.BadRequest, "пароль должен содержать как минимум одну цифру")
	case !regexp.MustCompile("[!@#$%^&*]").MatchString(password):
		return exc.New(exc.Service, exc.BadRequest, "пароль должен содержать как минимум один из специальных символов !@#$%^&*")
	default:
		return nil
	}
}

// ValidateEmail проверяет валидность Email.
//
// Используется приближенная (т.е. без учёта исключительных случаев, которые почтовые сервисы очевидно не поддерживают)
// проверка на соответствие стандарту https://www.ietf.org/rfc/rfc0822.txt?number=822.
// Многие почтовые сервисы используют свои правила и потому единственный 100% способ провалидировать почту - это
// отправка письма на ящик. Данный метод проводит лишь первичную проверку, а потому адрес нельзя считать достоверным.
// Поэтому первичная валидация происходит на соответствие regexp:
//
// /^([a-z0-9!#$%&'*+\\/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+\\/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)$/
//
// Также проверяется длина почты: не более 256 символов
func ValidateEmail(email string) error {
	re := regexp.MustCompile("^([a-z0-9!#$%&'*+\\\\/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+\\\\/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)$")
	// т.к. почта состоит из ascii символов, то можно использовать len()
	if !re.MatchString(email) {
		return exc.New(exc.Service, exc.BadRequest, "невалидная почта")
	}
	if len(email) > 256 {
		return exc.New(exc.Service, exc.BadRequest, "почта не может быть длиннее 256 символов")
	}
	return nil
}

// HashPassword генерирует соль длинной в 8 символов и хэширует пароль с этой солью
func HashPassword(password string) (string, string) {
	salt := random.RandomString(8)
	hash := pbkdf2.Key([]byte(password), []byte(salt), PasswordHashIterations, PasswordHashIterations, PasswordHashFunction)
	return salt, string(hash)
}

// CheckPassword сравнивает пароль с хешем и возвращает true, если они совпадают.
func (u *User) CheckPassword(password string) bool {
	return string(pbkdf2.Key([]byte(password), []byte(u.PasswordSalt), PasswordHashIterations, PasswordHashKeyLength, PasswordHashFunction)) == u.PasswordHash
}

func NewUserEmpty() *User {
	return new(User)
}
