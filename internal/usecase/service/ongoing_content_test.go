package service

import (
	"errors"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
)

func TestGetOngoingContentByContentID(t *testing.T) {
	t.Parallel()
	releaseDate := time.Now()

	testCases := []struct {
		Name                        string
		ContentID                   int
		ExpectedErr                 error
		SetupOngoingContentRepoMock func(repo *mockrepo.MockOngoingContent)
		SetupStaticRepoMock         func(repo *mockrepo.MockStatic)
	}{
		{
			Name:        "Успех",
			ContentID:   1,
			ExpectedErr: nil,
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().GetOngoingContentByID(1).Return(&entity.OngoingContent{
					ID:             1,
					Title:          "title",
					PosterStaticID: 1,
					ReleaseDate:    releaseDate,
					Type:           "movie",
				}, nil)
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(gomock.Any()).Return("posterURL", nil)
			},
		},
		{
			Name:        "Ошибка получения контента калндаря",
			ContentID:   1,
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при получении контента календаря релизов"), errors.New("database error")),
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().GetOngoingContentByID(1).Return(nil, errors.New("database error"))
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {},
		},
		{
			Name:        "Не найден",
			ContentID:   1,
			ExpectedErr: usecase.ErrOngoingContentNotFound,
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().GetOngoingContentByID(1).Return(nil, repository.ErrOngoingContentNotFound)
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {},
		},
		{
			Name:        "Статика пустая",
			ContentID:   1,
			ExpectedErr: nil,
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().GetOngoingContentByID(1).Return(&entity.OngoingContent{
					ID:             1,
					Title:          "title",
					PosterStaticID: 1,
					ReleaseDate:    releaseDate,
					Type:           "movie",
				}, nil)
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(gomock.Any()).Return("", repository.ErrStaticNotFound)
			},
		},
		{
			Name:        "Ошибка получения статики",
			ContentID:   1,
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при получении постера"), errors.New("database error")),
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().GetOngoingContentByID(1).Return(&entity.OngoingContent{
					ID:             1,
					Title:          "title",
					PosterStaticID: 1,
					ReleaseDate:    releaseDate,
					Type:           "movie",
				}, nil)
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(gomock.Any()).Return("", errors.New("database error"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentRepo := mockrepo.NewMockOngoingContent(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			service := NewOngoingContentService(mockOngoingContentRepo, mockStaticRepo)
			tc.SetupOngoingContentRepoMock(mockOngoingContentRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			_, err := service.GetOngoingContentByContentID(tc.ContentID)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestGetNearestOngoings(t *testing.T) {
	t.Parallel()
	releaseDate := time.Now()

	testCases := []struct {
		Name                        string
		Limit                       int
		ExpectedErr                 error
		SetupOngoingContentRepoMock func(repo *mockrepo.MockOngoingContent)
		SetupStaticRepoMock         func(repo *mockrepo.MockStatic)
	}{
		{
			Name:        "Успешно",
			Limit:       5,
			ExpectedErr: nil,
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().GetNearestOngoings(5).Return([]*entity.OngoingContent{
					{
						ID:             1,
						Title:          "title",
						PosterStaticID: 1,
						ReleaseDate:    releaseDate,
						Type:           "movie",
					},
				}, nil)
				repo.EXPECT().GetOngoingContentByID(1).Return(&entity.OngoingContent{
					ID:             1,
					Title:          "title",
					PosterStaticID: 1,
					ReleaseDate:    releaseDate,
					Type:           "movie",
				}, nil).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("path/to/static", nil).AnyTimes()
			},
		},
		{
			Name:        "Ошибка при получении ближайших релизов",
			Limit:       5,
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при получении ближайших релизов"), errors.New("database error")),
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().GetNearestOngoings(5).Return(nil, errors.New("database error"))
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {},
		},
		{
			Name:        "Ошибка в преобразовании в dto",
			Limit:       5,
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при получении контента календаря релизов"), errors.New("conversion error")),
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().GetNearestOngoings(5).Return([]*entity.OngoingContent{
					{
						ID:             1,
						Title:          "title",
						PosterStaticID: 1,
						ReleaseDate:    releaseDate,
						Type:           "movie",
					},
				}, nil)
				repo.EXPECT().GetOngoingContentByID(1).Return(nil, errors.New("conversion error")).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentRepo := mockrepo.NewMockOngoingContent(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			service := NewOngoingContentService(mockOngoingContentRepo, mockStaticRepo)
			tc.SetupOngoingContentRepoMock(mockOngoingContentRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			_, err := service.GetNearestOngoings(tc.Limit)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestGetOngoingContentByMonthAndYear(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                        string
		Month, Year                 int
		ExpectedErr                 error
		SetupOngoingContentRepoMock func(repo *mockrepo.MockOngoingContent)
	}{
		{
			Name:  "Database error",
			Month: 1,
			Year:  2022,
			ExpectedErr: entity.UsecaseWrap(
				errors.New("ошибка при получении релизов по месяцу и году"),
				errors.New("database error"),
			),
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().GetOngoingContentByMonthAndYear(1, 2022).Return(nil, errors.New("database error"))
			},
		},
		{
			Name:  "No content found",
			Month: 1,
			Year:  2022,
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().GetOngoingContentByMonthAndYear(1, 2022).Return(nil, nil)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentRepo := mockrepo.NewMockOngoingContent(ctrl)
			service := NewOngoingContentService(mockOngoingContentRepo, nil)
			tc.SetupOngoingContentRepoMock(mockOngoingContentRepo)
			_, err := service.GetOngoingContentByMonthAndYear(tc.Month, tc.Year)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestGetAllReleaseYears(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                        string
		ExpectedErr                 error
		ExpectedResult              []int
		SetupOngoingContentRepoMock func(repo *mockrepo.MockOngoingContent)
	}{
		{
			Name:           "Success",
			ExpectedErr:    nil,
			ExpectedResult: []int{2020, 2021, 2022},
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().GetAllReleaseYears().Return([]int{2020, 2021, 2022}, nil)
			},
		},
		{
			Name:           "Error retrieving release years",
			ExpectedErr:    errors.New("database error"),
			ExpectedResult: nil,
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().GetAllReleaseYears().Return(nil, errors.New("database error"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentRepo := mockrepo.NewMockOngoingContent(ctrl)
			service := NewOngoingContentService(mockOngoingContentRepo, nil)
			tc.SetupOngoingContentRepoMock(mockOngoingContentRepo)
			result, err := service.GetAllReleaseYears()
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedResult, result)
		})
	}
}

func TestIsOngoingContentFinished(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                        string
		ContentID                   int
		ExpectedErr                 error
		ExpectedResult              bool
		SetupOngoingContentRepoMock func(repo *mockrepo.MockOngoingContent)
	}{
		{
			Name:           "Success",
			ContentID:      1,
			ExpectedErr:    nil,
			ExpectedResult: true,
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().IsOngoingContentFinished(1).Return(true, nil)
			},
		},
		{
			Name:           "Error retrieving content status",
			ContentID:      1,
			ExpectedErr:    errors.New("database error"),
			ExpectedResult: false,
			SetupOngoingContentRepoMock: func(repo *mockrepo.MockOngoingContent) {
				repo.EXPECT().IsOngoingContentFinished(1).Return(false, errors.New("database error"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentRepo := mockrepo.NewMockOngoingContent(ctrl)
			service := NewOngoingContentService(mockOngoingContentRepo, nil)
			tc.SetupOngoingContentRepoMock(mockOngoingContentRepo)
			result, err := service.IsOngoingContentFinished(tc.ContentID)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedResult, result)
		})
	}
}
