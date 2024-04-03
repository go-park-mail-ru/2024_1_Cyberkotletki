package service

import (
	"bytes"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/google/uuid"
	"image"
	"image/draw"
	_ "image/gif" // для поддержки формата gif
	"image/jpeg"
	_ "image/png" // для поддержки формата png
	"net/http"
)

type StaticService struct {
	staticRepo repository.Static
}

func NewStaticService(staticRepo repository.Static) usecase.Static {
	return &StaticService{
		staticRepo: staticRepo,
	}
}

func (s *StaticService) GetAvatar(staticID int) (string, error) {
	path, err := s.staticRepo.GetStatic(staticID)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (s *StaticService) UploadAvatar(data []byte) (int, error) {
	// Определение типа файла
	contentType := http.DetectContentType(data)

	// Проверка, является ли файл изображением
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		return -1, entity.NewClientError("файл не является изображением", entity.ErrBadRequest)
	}

	// Чтение данных файла в структуру изображения
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return -1, entity.NewClientError("файл не является изображением", entity.ErrBadRequest, err)
	}

	// Проверка размеров изображения
	const minImageWidth, minImageHeight = 100, 100
	if img.Bounds().Dx() < minImageWidth || img.Bounds().Dy() < minImageHeight {
		return -1, entity.NewClientError(
			fmt.Sprintf("размеры изображения меньше %d x %d", minImageWidth, minImageHeight),
			entity.ErrBadRequest,
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

	// Конвертация изображения в формат JPEG и запись в переменную
	var out bytes.Buffer
	var opts jpeg.Options
	opts.Quality = 60
	err = jpeg.Encode(&out, squareImage, &opts)
	if err != nil {
		return -1, entity.NewClientError("ошибка при обработке изображения", entity.ErrBadRequest, err)
	}

	// Загрузка обработанного изображения на сервер
	id, err := s.staticRepo.UploadStatic("avatars", uuid.New().String()+".jpg", out.Bytes())
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s *StaticService) GetStaticURL(id int) (string, error) {
	return s.staticRepo.GetStatic(id)
}
