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
	staticURL, err := s.staticRepo.GetStatic(staticID)
	switch {
	case errors.Is(err, repository.ErrStaticNotFound):
		return "", usecase.ErrStaticNotFound
	case err != nil:
		return "", entity.UsecaseWrap(err, errors.New("ошибка при получении статики"))
	default:
		return staticURL, nil
	}
}

func (s *StaticService) UploadAvatar(reader io.Reader) (int, error) {
	data := make([]byte, s.staticRepo.GetMaxSize())
	bytesCount, err := reader.Read(data)
	if err != nil {
		return -1, errors.New("UploadAvatar: ошибка при чтении файла")
	}
	if bytesCount >= s.staticRepo.GetMaxSize() {
		return -1, usecase.ErrStaticTooBigFile
	}
	data = data[:bytesCount]

	// Определение типа файла
	contentType := http.DetectContentType(data)

	// Проверка, является ли файл изображением
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		return -1, usecase.ErrStaticNotImage
	}

	// Чтение данных файла в структуру изображения
	img, _, err := image.Decode(bytes.NewReader(data))
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
	id, err := s.staticRepo.UploadStatic("avatars", uuid.New().String()+".webp", out)
	if err != nil {
		return -1, err
	}
	return id, nil
}
