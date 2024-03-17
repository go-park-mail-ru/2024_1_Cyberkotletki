package entity

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/random"
	"golang.org/x/crypto/argon2"
	"regexp"
	"time"
)

// Константы подобраны в соответствии с рекомендациями по безопасности. В рамках оптимизации можно менять кол-во потоков
const (
	PasswordHashTime      = 1
	PasswordHashKibMemory = 64 * 1024
	PasswordHashThreads   = 4
)

type User struct {
	Id               int       `json:"id"`    // Уникальный идентификатор
	Name             string    `json:"name"`  // Имя пользователя
	Email            string    `json:"email"` // Электронная почта
	PasswordHash     string    // Хэш пароля пользователя
	PasswordSalt     string    // Соль для генерации хэша пароля
	BirthDate        time.Time `json:"birth_date"`        // День рождения
	SavedFilms       []Film    `json:"saved_films"`       // Сохраненные фильмы
	SavedSeries      []Series  `json:"saved_series"`      // Сохраненные сериалы
	SavedPersons     []Person  `json:"saved_persons"`     // Сохраненные персоны
	Friends          []User    `json:"friends"`           // Друзья
	ExpectedFilms    []Film    `json:"expected_films"`    // Ожидаемые фильмы
	RegistrationDate time.Time `json:"registration_date"` // Дата регистрации пользователя
}

func NewUserEmpty() *User {
	return new(User)
}

// ValidatePassword проверяет валидность пароля.
//
// Критерии валидности пароля:
//
// 1) пароль содержит от 8 до 32 символов включительно
//
// 2) пароль не содержит ничего, кроме латинских букв, цифр и символов !@#$%^&*
//
// 3) пароль содержит как минимум одну заглавную, одну строчную букву, одну цифру и один из символов !@#$%^&*
func ValidatePassword(password string) error {
	switch {
	case len(password) < 8:
		return NewClientError("пароль должен содержать не менее 8 символов", ErrBadRequest)
	case len(password) > 32:
		return NewClientError("пароль должен содержать не более 32 символов", ErrBadRequest)
	case !regexp.MustCompile("^[!@#$%^&*\\w]+$").MatchString(password):
		return NewClientError("пароль должен состоять из латинских букв, цифр и специальных символов !@#$%^&*", ErrBadRequest)
	case !regexp.MustCompile("[A-Z]").MatchString(password):
		return NewClientError("пароль должен содержать как минимум одну заглавную букву", ErrBadRequest)
	case !regexp.MustCompile("[a-z]").MatchString(password):
		return NewClientError("пароль должен содержать как минимум одну строчную букву", ErrBadRequest)
	case !regexp.MustCompile("\\d").MatchString(password):
		return NewClientError("пароль должен содержать как минимум одну цифру", ErrBadRequest)
	case !regexp.MustCompile("[!@#$%^&*]").MatchString(password):
		return NewClientError("пароль должен содержать как минимум один из специальных символов !@#$%^&*", ErrBadRequest)
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
		return NewClientError("невалидная почта", ErrBadRequest)
	}
	if len(email) > 256 {
		return NewClientError("почта не может быть длиннее 256 символов", ErrBadRequest)
	}
	return nil
}

// HashPassword генерирует соль длинной в 8 символов и хэширует пароль с этой солью.
// Первым параметром возвращает полученную соль, вторым параметром возвращает хеш
func HashPassword(password string) (salt string, hash string) {
	salt = random.String(8)
	hash = string(argon2.IDKey([]byte(password), []byte(salt), PasswordHashTime, PasswordHashKibMemory, PasswordHashThreads, 32))
	return salt, hash
}

// CheckPassword принимает пароль и хэширует его с подмешиванием соли пользователя. При совпадении хешей возвращает true
func (u *User) CheckPassword(password string) bool {
	return string(argon2.IDKey([]byte(password), []byte(u.PasswordSalt), PasswordHashTime, PasswordHashKibMemory, PasswordHashThreads, 32)) == u.PasswordHash
}
