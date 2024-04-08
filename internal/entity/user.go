package entity

import (
	"bytes"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/random"
	"golang.org/x/crypto/argon2"
	"regexp"
	"unicode/utf8"
)

// Константы подобраны в соответствии с рекомендациями по безопасности. В рамках оптимизации можно менять кол-во потоков
const (
	PasswordHashTime      = 1
	PasswordHashKibMemory = 64 * 1024
	PasswordHashThreads   = 4
)

type User struct {
	ID             int    `json:"id"`    // Уникальный идентификатор
	Name           string `json:"name"`  // Имя пользователя
	Email          string `json:"email"` // Электронная почта
	PasswordHash   []byte // Хэш пароля пользователя
	PasswordSalt   []byte // Соль для генерации хэша пароля
	AvatarUploadID int    `json:"avatar_upload_id"` // Ссылка на аватар
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
	case !regexp.MustCompile(`^[!@#$%^&*\w]+$`).MatchString(password):
		return NewClientError("пароль должен состоять из латинских букв, цифр и специальных символов !@#$%^&*", ErrBadRequest)
	case !regexp.MustCompile(`[A-Z]`).MatchString(password):
		return NewClientError("пароль должен содержать как минимум одну заглавную букву", ErrBadRequest)
	case !regexp.MustCompile(`[a-z]`).MatchString(password):
		return NewClientError("пароль должен содержать как минимум одну строчную букву", ErrBadRequest)
	case !regexp.MustCompile(`\d`).MatchString(password):
		return NewClientError("пароль должен содержать как минимум одну цифру", ErrBadRequest)
	case !regexp.MustCompile(`[!@#$%^&*]`).MatchString(password):
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
// `/^([a-z0-9!#$%&'*+\\/=?^_`{|}~-]+(?:\.[a-z0-9!#$%&'*+\\/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+
// [a-z0-9](?:[a-z0-9-]*[a-z0-9])?)$/`
//
// Также проверяется длина почты: не более 256 символов
func ValidateEmail(email string) error {
	re := regexp.MustCompile("^([a-z0-9!#$%&'*+\\\\/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+\\\\/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)$") // nolint: lll
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
func HashPassword(password string) (salt []byte, hash []byte, err error) {
	salt, err = random.Bytes(8)
	if err != nil {
		return nil, nil, NewClientError("произошла непредвиденная ошибка", ErrInternal)
	}
	hash = argon2.IDKey(
		[]byte(password),
		salt,
		PasswordHashTime,
		PasswordHashKibMemory,
		PasswordHashThreads,
		32,
	)
	return salt, hash, nil
}

// CheckPassword принимает пароль и хэширует его с подмешиванием соли пользователя. При совпадении хешей возвращает true
func (u *User) CheckPassword(password string) bool {
	return bytes.Equal(
		argon2.IDKey(
			[]byte(password),
			u.PasswordSalt,
			PasswordHashTime,
			PasswordHashKibMemory,
			PasswordHashThreads,
			32,
		),
		u.PasswordHash,
	)
}

// ValidateName проверяет валидность имени.
//
// Имя необязательно должно быть настоящим, оно служит в роли никнейма пользователя, поэтому допустимы любые символы,
// но не более 30 символов
func ValidateName(name string) error {
	if utf8.RuneCountInString(name) > 30 {
		return NewClientError("имя не может быть длиннее 30 символов", ErrBadRequest)
	}
	return nil
}
