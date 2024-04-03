package service

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestContent_GetContentPreviewCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockContentRepo := mockrepo.NewMockContent(ctrl)
	contentService := ContentService{
		contentRepo: mockContentRepo,
	}

	testCases := []struct {
		Name                 string
		Input                int
		ExpectedErr          error
		SetupContentRepoMock func(repo *mockrepo.MockContent)
	}{
		{
			Name:        "Успешное получение карточки фильма",
			Input:       1,
			ExpectedErr: nil,
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetFilm(1).Return(&entity.Film{}, nil)
			},
		},
		{
			Name:        "Фильм не найден",
			Input:       100,
			ExpectedErr: entity.ErrNotFound,
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetFilm(100).Return(nil, entity.ErrNotFound)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			tc.SetupContentRepoMock(mockContentRepo)
			_, err := contentService.GetContentPreviewCard(tc.Input)
			if tc.ExpectedErr != nil {
				require.EqualError(t, err, tc.ExpectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
