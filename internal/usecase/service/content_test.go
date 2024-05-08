package service

import (
	"fmt"
	"testing"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestContentService_GetContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		ContentID            int
		ExpectedOutput       *dto.Content
		ExpectedErr          error
		SetupContentRepoMock func(repo *mockrepo.MockContent)
		SetupStaticRepoMock  func(repo *mockrepo.MockStatic)
	}{
		{
			Name:        "Не существующий контент",
			ContentID:   2,
			ExpectedErr: entity.UsecaseWrap(fmt.Errorf("ошибка при получении контента"), fmt.Errorf("ошибка при получении контента")),
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetContent(2).Return(nil, fmt.Errorf("ошибка при получении контента")).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("", repository.ErrStaticNotFound).AnyTimes()
			},
		},
		{
			Name:        "ошибка при получении постера",
			ContentID:   1,
			ExpectedErr: entity.UsecaseWrap(fmt.Errorf("ошибка при получении постера"), fmt.Errorf("ошибка при получении постера")),
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetContent(1).Return(&entity.Content{
					ID:               1,
					Title:            "Бэтмен",
					OriginalTitle:    "Batman",
					Slogan:           "I'm Batman",
					Budget:           "1000000",
					AgeRestriction:   18,
					Rating:           9.1,
					IMDBRating:       9.1,
					PosterStaticID:   1,
					BackdropStaticID: 1,
					PicturesStaticID: []int{1},
					Description:      "Описание фильма или сериала",
				}, nil).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("", fmt.Errorf("ошибка при получении постера")).AnyTimes()
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			contentService := NewContentService(mockContentRepo, mockStaticRepo)
			tc.SetupContentRepoMock(mockContentRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			content, err := contentService.GetContentByID(tc.ContentID)
			require.EqualValues(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOutput, content)
		})
	}
}
