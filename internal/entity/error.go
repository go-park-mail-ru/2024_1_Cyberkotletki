package entity

import (
	"errors"
	"fmt"
)

// ClientError - это ошибка, которая доставляется конечному клиенту
type ClientError struct {
	Msg        string
	Additional error
}

func (err ClientError) Error() string {
	// Отображаем только сообщение, которое можно видеть клиенту, Additional исключительно для внутренних нужд!
	return fmt.Sprintf("%s", err.Msg)
}

// NewClientError генерирует ошибку, содержащую сообщение для клиента и вспомогательные ошибки
func NewClientError(msg string, errs ...error) error {
	return ClientError{
		Msg:        msg,
		Additional: errors.Join(errs...),
	}
}

// Contains позволяет проверить, содержит ли ошибка другую ошибку. В отличие от errors.Is, корректно работает с
// ClientError и проверяет, содержится ли ошибка в поле ClientError.Additional
func Contains(err error, target error) bool {
	var er ClientError
	if errors.As(err, &er) {
		return errors.Is(er.Additional, target)
	}
	return errors.Is(err, target)
}

var (
	ErrNotFound      = errors.New("not found")
	ErrBadRequest    = errors.New("bad request")
	ErrForbidden     = errors.New("forbidden")
	ErrAlreadyExists = errors.New("already exists")
	ErrRedis         = errors.New("redis error")
)
