package usecase

import (
	"errors"
	"io"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_static.go
type Static interface {
	// GetStatic возвращает путь к статике по ID
	// Возможные ошибки:
	// ErrStaticNotFound - статика с таким id не найдена
	GetStatic(staticID int) (string, error)
	// UploadAvatar загружает аватар на сервер
	// Возможные ошибки:
	// ErrStaticTooBigFile - файл слишком большой
	// ErrStaticNotImage - файл не является валидным изображением
	// ErrStaticImageDimensions - изображение имеет недопустимые размеры
	UploadAvatar(reader io.Reader) (int, error)
}

var (
	ErrStaticNotFound        = errors.New("статика с таким id не найдена")
	ErrStaticTooBigFile      = errors.New("слишком большой файл")
	ErrStaticNotImage        = errors.New("файл не является валидным изображением")
	ErrStaticImageDimensions = errors.New("изображение имеет недопустимые размеры")
)
