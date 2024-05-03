package service

import (
	"bytes"
	"fmt"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
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
			ExpectedErr: fmt.Errorf("не удалось получить аватар"),
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
				repo.EXPECT().UploadStatic(gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
			},
		},
		{
			Name: "Невалидный файл",
			Input: func() io.Reader {
				// возвращаем невалидные данные
				return bytes.NewReader([]byte("invalid data"))
			},
			ExpectedErr:         fmt.Errorf(""),
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {},
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
			ExpectedErr:         fmt.Errorf(""),
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			staticService := StaticService{
				staticRepo: mockStaticRepo,
			}
			tc.SetupStaticRepoMock(mockStaticRepo)

			_, err := staticService.UploadAvatar(tc.Input())
			require.EqualValues(t, tc.ExpectedErr, err)
		})
	}
}
