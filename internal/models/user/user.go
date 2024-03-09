package user

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

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

func (u *User) ValidatePassword(password string) error {
	if len(password) < 8 {
		return exc.New(exc.Server, exc.BadRequest, "пароль должен содержать не менее 8 символов")
	}
	if len(password) > 32 {
		return exc.New(exc.Server, exc.BadRequest, "пароль должен содержать не более 32 символов")
	}
	return nil
}

// ValidateEmail проверяет валидность Email.
// Используется приближенная (т.е. без учёта исключительных случаев, которые почтовые сервисы очевидно не поддерживают)
// проверка на соответствие стандарту https://www.ietf.org/rfc/rfc0822.txt?number=822.
// Многие почтовые сервисы используют свои правила и потому единственный 100% способ провалидировать почту - это
// отправка письма на ящик. Данный метод проводит лишь первичную проверку, а потому адрес нельзя считать достоверным.
// Проверяются следующие факторы:
// 0) Адрес содержит не более 256 символов;
// 1) Адрес содержит @, который разбивает адрес на имя и хост;
// 2) Хост содержит хотя бы одну точку, при этом до и после точки должна находиться подстрока, состоящая только из
// буквенных и циферных символов;
// 3) Имя содержит хотя бы один буквенный или циферный символ.
func (u *User) ValidateEmail(email string) error {
	invalidEmail := exc.New(exc.Service, exc.BadRequest, "невалидная почта")
	if utf8.RuneCountInString(email) > 256 || strings.Count(email, "@") != 1 {
		return invalidEmail
	}
	parts := strings.Split(email, "@")
	name := parts[0]
	host := parts[1]

	// regexp не позволяет сделать проверку на наличие буквы из любого алфавита, поэтому лучше обойтись без регулярок
	flag := false
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			flag = true
			break
		}
	}
	if !flag {
		return invalidEmail
	}
	// ----------------
	if !strings.Contains(host, ".") {
		return invalidEmail
	}
	for _, part := range strings.Split(host, ".") {
		if part == "" {
			return invalidEmail
		}
		for _, r := range part {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
				return invalidEmail
			}
		}
	}

	return nil
}

func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

// CheckPassword сравнивает пароль с хешем и возвращает true, если они совпадают.
func (u *User) CheckPassword(password string) bool {
	hash := HashPassword(password)
	return hash == u.PasswordHash
}

func NewUserEmpty() *User {
	return new(User)
}
