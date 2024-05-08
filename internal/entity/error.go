package entity

import (
	"errors"
	"fmt"
)

func PSQLWrap(errs ...error) error {
	return errors.Join(ErrPSQL, errors.Join(errs...))
}

func PSQLQueryErr(queryName string, err error) error {
	return PSQLWrap(fmt.Errorf("ошибка при выполнении запроса %s", queryName), err)
}

func RedisWrap(errs ...error) error {
	return errors.Join(ErrRedis, errors.Join(errs...))
}

func S3Wrap(errs ...error) error {
	return errors.Join(ErrS3, errors.Join(errs...))
}

func UsecaseWrap(errs ...error) error {
	return errors.Join(ErrInternal, errors.Join(errs...))
}

var (
	ErrRedis    = errors.New("redis error")
	ErrPSQL     = errors.New("postgres error")
	ErrS3       = errors.New("s3 error")
	ErrInternal = errors.New("internal server error")
)

const (
	PSQLUniqueViolation     = "23505"
	PSQLCheckViolation      = "23514"
	PSQLForeignKeyViolation = "23503"
)
