package service

import (
	"bytes"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
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

func (s StaticService) GetAvatar(staticID int) (string, error) {
	path, err := s.staticRepo.GetStatic(staticID)
	if err != nil {
		return "", err
	}
	return path, nil
}

func (s StaticService) UploadAvatar(data []byte) (int, error) {
	// Определение типа файла
	contentType := http.DetectContentType(data)

	// Проверка, является ли файл изображением
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" {
		return -1, entity.NewClientError("файл не является изображением", entity.ErrBadRequest)
	}

	// Чтение данных файла в структуру изображения
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return -1, err
	}

	// Конвертация изображения в формат JPEG и запись в переменную
	var out bytes.Buffer
	var opts jpeg.Options
	opts.Quality = 60
	err = jpeg.Encode(&out, img, &opts)

	// Загрузка обработанного изображения на сервер
	id, err := s.staticRepo.UploadStatic("avatars", uuid.New().String()+".jpg", out.Bytes())
	if err != nil {
		return -1, err
	}
	return id, nil
}
