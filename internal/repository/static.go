package repository

import (
	"bytes"
	"errors"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_static.go
type Static interface {
	// GetStatic возвращает путь к статике по его ID
	// Возможные ошибки:
	// ErrStaticNotFound - статика с таким id не найдена
	GetStatic(staticID int) (string, error)
	// UploadStatic загружает статику на сервер
	// Возможные ошибки:
	// ErrStaticTooBigFile - файл слишком большой
	UploadStatic(path, filename string, buf bytes.Buffer) (int, error)
	// GetBasicPath возвращает базовый путь для статики
	GetBasicPath() string
	// GetMaxSize возвращает максимальный размер файла
	GetMaxSize() int
}

var (
	ErrStaticNotFound   = errors.New("статика с таким id не найдена")
	ErrStaticTooBigFile = errors.New("слишком большой файл")
)
