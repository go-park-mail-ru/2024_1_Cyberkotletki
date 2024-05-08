package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/google/uuid"
	"image"
	"image/draw"
	"io"
	"net/http"

	"github.com/chai2010/webp"
	_ "image/gif"  // для поддержки формата gif
	_ "image/jpeg" // для поддержки формата jpeg
	_ "image/png"  // для поддержки формата png
)

type StaticService struct {
	staticRepo repository.Static
}

func NewStaticService(staticRepo repository.Static) usecase.Static {
	return &StaticService{
		staticRepo: staticRepo,
	}
}

func (s *StaticService) GetStatic(staticID int) (string, error) {
	staticURL, err := s.staticRepo.GetStaticURL(staticID)
	switch {
	case errors.Is(err, repository.ErrStaticNotFound):
		return "", usecase.ErrStaticNotFound
	case err != nil:
		return "", entity.UsecaseWrap(err, errors.New("ошибка при получении статики"))
	default:
		return staticURL, nil
	}
}

func (s *StaticService) UploadAvatar(reader io.ReadSeeker) (int, error) {
	// Проверка размера файла
	size, err := reader.Seek(0, io.SeekEnd)
	if err != nil {
		return -1, errors.Join(err, errors.New("ошибка при определении размера файла"))
	}
	if size > int64(s.staticRepo.GetMaxSize()) {
		return -1, usecase.ErrStaticTooBigFile
	}
	_, err = reader.Seek(0, io.SeekStart) // Возвращаемся в начало файла
	if err != nil {
		return -1, errors.Join(err, errors.New("ошибка при возвращении io.ReadSeeker в начало файла"))
	}

	// Определение типа файла
	headerBytes := make([]byte, 512)
	_, err = reader.Read(headerBytes)
	if err != nil {
		return -1, errors.Join(err, errors.New("ошибка при чтении заголовка файла"))
	}
	contentType := http.DetectContentType(headerBytes)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		return -1, usecase.ErrStaticNotImage
	}
	_, err = reader.Seek(0, io.SeekStart) // Возвращаемся в начало файла
	if err != nil {
		return -1, errors.Join(err, errors.New("ошибка при возвращении io.ReadSeeker в начало файла"))
	}

	// Чтение данных файла в структуру изображения
	img, _, err := image.Decode(reader)
	if err != nil {
		return -1, usecase.ErrStaticNotImage
	}

	// Проверка размеров изображения
	const minImageWidth, minImageHeight = 100, 100
	if img.Bounds().Dx() < minImageWidth || img.Bounds().Dy() < minImageHeight {
		return -1, errors.Join(
			usecase.ErrStaticImageDimensions,
			fmt.Errorf(
				"изображение имеет размеры %dx%d, а должно быть как минимум %dx%d",
				img.Bounds().Dx(), img.Bounds().Dy(), minImageWidth, minImageHeight,
			),
		)
	}

	// Проверка размеров изображения и обрезка до квадратной формы
	width, height := img.Bounds().Dx(), img.Bounds().Dy()
	var squareImage *image.RGBA
	var start image.Point
	var squareSize int
	if width > height {
		start.X = (width - height) / 2
		squareSize = height
	} else {
		start.Y = (height - width) / 2
		squareSize = width
	}
	squareImage = image.NewRGBA(image.Rect(0, 0, squareSize, squareSize))
	draw.Draw(squareImage, squareImage.Bounds(), img, start, draw.Src)

	// Конвертация изображения в формат WEBP и запись в переменную
	var out bytes.Buffer
	var opts webp.Options
	opts.Lossless = false
	opts.Quality = 60
	if err = webp.Encode(&out, squareImage, &opts); err != nil {
		return -1, errors.Join(errors.New("ошибка при конвертации изображения в формат WEBP"), err)
	}

	// Загрузка обработанного изображения на сервер
	id, err := s.staticRepo.UploadStatic("avatars", uuid.New().String()+".webp", bytes.NewReader(out.Bytes()))
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *StaticService) GetStaticFile(staticURI string) (io.ReadSeeker, error) {
	static, err := s.staticRepo.GetStaticFile(staticURI)
	switch {
	case errors.Is(err, repository.ErrStaticNotFound):
		return nil, usecase.ErrStaticNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(err, errors.New("ошибка при получении статики"))
	}
	return static, nil
}
