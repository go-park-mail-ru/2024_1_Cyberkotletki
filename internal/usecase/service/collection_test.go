package service

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCollections_GetCompilation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		Input                string
		ExpectedErr          error
		SetupContentRepoMock func(repo *mockrepo.MockContent)
	}{
		{
			Name:        "Успешное получение подборки",
			Input:       "drama",
			ExpectedErr: nil,
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetFilmsByGenre(1).Return([]entity.Film{}, nil)
			},
		},
		{
			Name:                 "Жанр не найден",
			Input:                "смешной хоррор",
			ExpectedErr:          fmt.Errorf("такого жанра не существует"),
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			collectionsService := CollectionsService{
				contentRepo: mockContentRepo,
			}
			tc.SetupContentRepoMock(mockContentRepo)
			_, err := collectionsService.GetCompilation(tc.Input)
			if tc.ExpectedErr != nil {
				require.EqualError(t, err, tc.ExpectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCollections_GetGenres(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContentRepo := mockrepo.NewMockContent(ctrl)
	collectionsService := CollectionsService{
		contentRepo: mockContentRepo,
	}
	genres, err := collectionsService.GetGenres()
	require.NoError(t, err)
	require.NotEmpty(t, genres)
}
