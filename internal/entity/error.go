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
	return fmt.Sprint(err.Msg)
}

func PSQLWrap(errs ...error) error {
	return errors.Join(ErrPSQL, errors.Join(errs...))
}

func PSQLQueryErr(queryName string, err error) error {
	return PSQLWrap(fmt.Errorf("ошибка при выполнении запроса %s", queryName), err)
}

func RedisWrap(errs ...error) error {
	return errors.Join(ErrRedis, errors.Join(errs...))
}

func UsecaseWrap(errs ...error) error {
	return errors.Join(ErrInternal, errors.Join(errs...))
}

var (
	ErrRedis    = errors.New("redis error")
	ErrPSQL     = errors.New("postgres error")
	ErrInternal = errors.New("internal server error")
)

const (
	PSQLUniqueViolation     = "23505"
	PSQLCheckViolation      = "23514"
	PSQLForeignKeyViolation = "23503"
)
