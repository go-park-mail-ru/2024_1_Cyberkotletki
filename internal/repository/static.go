package repository

import (
	"errors"
	"io"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_static.go
type Static interface {
	// GetStaticURL возвращает путь к статике по его ID
	// Возможные ошибки:
	// ErrStaticNotFound - статика с таким id не найдена
	GetStaticURL(staticID int) (string, error)
	// UploadStatic загружает статику на сервер
	// Возможные ошибки:
	// ErrStaticTooBigFile - файл слишком большой
	UploadStatic(path, filename string, reader io.ReadSeeker) (int, error)
	// GetStaticFile возвращает статику по ID
	// Возможные ошибки:
	// ErrStaticNotFound - статика с таким id не найдена
	GetStaticFile(staticURI string) (io.ReadSeeker, error)
	// GetMaxSize возвращает максимальный размер файла
	GetMaxSize() int
}

var (
	ErrStaticNotFound   = errors.New("статика с таким id не найдена")
	ErrStaticTooBigFile = errors.New("слишком большой файл")
)
