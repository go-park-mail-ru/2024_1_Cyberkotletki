package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"os"
	"testing"
)

func TestStaticService_GetStatic(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		Input               int
		ExpectedErr         error
		SetupStaticRepoMock func(repo *mockrepo.MockStatic)
	}{
		{
			Name:        "Существующий аватар",
			Input:       1,
			ExpectedErr: nil,
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("path", nil)
			},
		},
		{
			Name:        "Не существующий аватар",
			Input:       2,
			ExpectedErr: entity.UsecaseWrap(fmt.Errorf("не удалось получить аватар"), fmt.Errorf("ошибка при получении статики")),
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(2).Return("", fmt.Errorf("не удалось получить аватар"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			staticService := NewStaticService(mockStaticRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			_, err := staticService.GetStatic(tc.Input)
			require.EqualValues(t, tc.ExpectedErr, err)
		})
	}
}

func TestStaticService_UploadAvatar(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		Input               func() io.Reader
		ExpectedErr         error
		SetupStaticRepoMock func(repo *mockrepo.MockStatic)
	}{
		{
			Name: "Валидный файл",
			Input: func() io.Reader {
				// открываем изображение valid_picture.png и читаем его в байты
				data, err := os.ReadFile("../../../assets/tests/valid_picture.png")
				if err != nil {
					t.Fatal(err)
				}
				// преобразуем data в io.Reader
				return bytes.NewReader(data)
			},
			ExpectedErr: nil,
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetMaxSize().Return(1000000).AnyTimes()
				repo.EXPECT().UploadStatic(gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
			},
		},
		{
			Name: "Невалидный файл",
			Input: func() io.Reader {
				// возвращаем невалидные данные
				return bytes.NewReader([]byte("invalid data"))
			},
			ExpectedErr: usecase.ErrStaticNotImage,
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetMaxSize().Return(1000000).AnyTimes()
			},
		},
		{
			Name: "Слишком большой файл",
			Input: func() io.Reader {
				// возвращаем слишком большой файл
				return bytes.NewReader(bytes.Repeat([]byte("a"), 1000001))
			},
			ExpectedErr: usecase.ErrStaticTooBigFile,
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetMaxSize().Return(1000000).AnyTimes()
			},
		},
		{
			Name: "Неподходящий размер изображения",
			Input: func() io.Reader {
				// открываем изображение small_picture.png и читаем его в байты
				data, err := os.ReadFile("../../../assets/tests/invalid_size.png")
				if err != nil {
					t.Fatal(err)
				}
				// преобразуем data в io.Reader
				return bytes.NewReader(data)
			},
			ExpectedErr: errors.Join(
				usecase.ErrStaticImageDimensions,
				fmt.Errorf("изображение имеет размеры 62x55, а должно быть как минимум 100x100"),
			),
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetMaxSize().Return(1000000).AnyTimes()
			},
		},
		{
			Name: "Ошибка при создании записи в БД",
			Input: func() io.Reader {
				// открываем изображение valid_picture.png и читаем его в байты
				data, err := os.ReadFile("../../../assets/tests/valid_picture.png")
				if err != nil {
					t.Fatal(err)
				}
				// преобразуем data в io.Reader
				return bytes.NewReader(data)
			},
			ExpectedErr: fmt.Errorf("ошибка при создании записи в БД"),
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetMaxSize().Return(1000000).AnyTimes()
				repo.EXPECT().UploadStatic(gomock.Any(), gomock.Any(), gomock.Any()).Return(-1, fmt.Errorf("ошибка при создании записи в БД"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			staticService := NewStaticService(mockStaticRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			_, err := staticService.UploadAvatar(tc.Input())
			require.EqualValues(t, tc.ExpectedErr, err)
		})
	}
}
