package service

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"os"
	"testing"
)

func TestStaticService_GetAvatar(t *testing.T) {
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
			staticService := StaticService{
				staticRepo: mockStaticRepo,
			}
			tc.SetupStaticRepoMock(mockStaticRepo)
			_, err := staticService.GetAvatar(tc.Input)
			require.EqualValues(t, tc.ExpectedErr, err)
		})
	}
}

func TestStaticService_GetStaticUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		Input               int
		ExpectedErr         error
		SetupStaticRepoMock func(repo *mockrepo.MockStatic)
	}{
		{
			Name:        "Существующий файл",
			Input:       1,
			ExpectedErr: nil,
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("path", nil)
			},
		},
		{
			Name:        "Не существующий файл",
			Input:       2,
			ExpectedErr: fmt.Errorf("не удалось получить файл"),
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(2).Return("", fmt.Errorf("не удалось получить файл"))
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
			staticService := StaticService{
				staticRepo: mockStaticRepo,
			}
			tc.SetupStaticRepoMock(mockStaticRepo)
			_, err := staticService.GetStaticURL(tc.Input)
			require.EqualValues(t, tc.ExpectedErr, err)
		})
	}
}

func TestStaticService_UploadAvatar(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		Input               func() []byte
		ExpectedErr         error
		SetupStaticRepoMock func(repo *mockrepo.MockStatic)
	}{
		{
			Name: "Валидный файл",
			Input: func() []byte {
				// открываем изображение valid_picture.png и читаем его в байты
				data, err := os.ReadFile("../../../assets/tests/valid_picture.png")
				if err != nil {
					t.Fatal(err)
				}
				return data
			},
			ExpectedErr: nil,
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().UploadStatic(gomock.Any(), gomock.Any(), gomock.Any()).Return(1, nil)
			},
		},
		{
			Name: "Невалидный файл",
			Input: func() []byte {
				return []byte("data")
			},
			ExpectedErr:         entity.NewClientError("файл не является изображением", entity.ErrBadRequest),
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {},
		},
		{
			Name: "Неподходящий размер изображения",
			Input: func() []byte {
				// открываем изображение small_picture.png и читаем его в байты
				data, err := os.ReadFile("../../../assets/tests/invalid_size.png")
				if err != nil {
					t.Fatal(err)
				}
				return data
			},
			ExpectedErr:         entity.NewClientError("размеры изображения меньше 100 x 100", entity.ErrBadRequest),
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
