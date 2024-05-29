package service

import (
	"database/sql"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"testing"
	"time"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	mock_usecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
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
		SetupStaticRepoMock  func(repo *mock_usecase.MockStatic)
	}{
		{
			Name:      "Существующий контент",
			ContentID: 1,
			ExpectedOutput: &dto.Content{
				ID:             1,
				Title:          "Бэтмен",
				OriginalTitle:  "Batman",
				Slogan:         "I'm Batman",
				Budget:         "1000000",
				AgeRestriction: 18,
				Rating:         9.1,
				IMDBRating:     9.1,
				PosterURL:      "http://localhost:8080/static/1",
				BackdropURL:    "http://localhost:8080/static/1",
				PicturesURL:    []string{"http://localhost:8080/static/1"},
				Countries:      []string{"USA"},
				Genres:         []string{"Боевик"},
				Actors:         []dto.PersonPreview{{ID: 1, Name: "Кристиан Бэйл"}},
				Directors:      []dto.PersonPreview{},
				Writers:        []dto.PersonPreview{},
				Producers:      []dto.PersonPreview{},
				Operators:      []dto.PersonPreview{},
				Composers:      []dto.PersonPreview{},
				Editors:        []dto.PersonPreview{},
				SimilarContent: []dto.PreviewContentCardVertical{
					{
						ID:     2,
						Poster: "http://localhost:8080/static/1",
						Genres: []string{"Боевик"},
					},
				},
				Series: dto.SeriesContent{Seasons: []dto.Season{{ID: 1, Episodes: []dto.Episode{{ID: 1, Title: "Эпизод 1"}}}}},
				Type:   entity.ContentTypeSeries,
			},
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
					Description:      "",
					Country:          []entity.Country{{Name: "USA"}},
					Genres:           []entity.Genre{{Name: "Боевик"}},
					Actors:           []entity.Person{{ID: 1, Name: "Кристиан Бэйл"}},
					Type:             entity.ContentTypeSeries,
					Series: &entity.Series{
						Seasons: []entity.Season{
							{
								ID: 1,
								Episodes: []entity.Episode{
									{
										ID:    1,
										Title: "Эпизод 1",
									},
								},
							},
						},
					},
				}, nil).AnyTimes()
				repo.EXPECT().GetSimilarContent(1).Return([]entity.Content{
					{
						ID:             2,
						PosterStaticID: 1,
						Genres:         []entity.Genre{{Name: "Боевик"}},
					},
				}, nil).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mock_usecase.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("http://localhost:8080/static/1", nil).AnyTimes()
			},
		},
		{
			Name:        "Не существующий контент",
			ContentID:   2,
			ExpectedErr: entity.UsecaseWrap(fmt.Errorf("ошибка при получении контента"), fmt.Errorf("ошибка при получении контента")),
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetContent(2).Return(nil, fmt.Errorf("ошибка при получении контента")).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mock_usecase.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("", usecase.ErrStaticNotFound).AnyTimes()
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
				repo.EXPECT().GetSimilarContent(1).Return([]entity.Content{{
					ID:             2,
					PosterStaticID: 1,
				}}, nil).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mock_usecase.MockStatic) {
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
			mockStaticRepo := mock_usecase.NewMockStatic(ctrl)
			contentService := NewContentService(mockContentRepo, mockStaticRepo, "")
			tc.SetupContentRepoMock(mockContentRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			content, err := contentService.GetContentByID(tc.ContentID)
			require.EqualValues(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOutput, content)
		})
	}
}

func TestContentService_GetPersonByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		PersonID             int
		ExpectedOutput       *dto.Person
		ExpectedErr          error
		SetupContentRepoMock func(repo *mockrepo.MockContent)
		SetupStaticRepoMock  func(repo *mock_usecase.MockStatic)
	}{
		{
			Name:     "Существующий человек",
			PersonID: 1,
			ExpectedOutput: &dto.Person{
				ID:        1,
				Name:      "Кристиан Бэйл",
				EnName:    "Christian Bale",
				BirthDate: &time.Time{},
				DeathDate: &time.Time{},
				Sex:       "M",
				Height:    183,
				PhotoURL:  "http://localhost:8080/static/1",
				Roles:     map[string][]dto.PreviewContentCardVertical{"actor": {{ID: 1, Poster: "http://localhost:8080/static/1", Genres: []string{}}}},
			},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetPerson(1).Return(&entity.Person{
					ID:            1,
					Name:          "Кристиан Бэйл",
					EnName:        "Christian Bale",
					BirthDate:     sql.NullTime{Valid: true},
					DeathDate:     sql.NullTime{Valid: true},
					Sex:           "M",
					Height:        sql.NullInt64{Int64: 183, Valid: true},
					PhotoStaticID: sql.NullInt64{Int64: 1, Valid: true},
				}, nil).AnyTimes()
				repo.EXPECT().GetPersonRoles(1).Return([]entity.PersonRole{
					{PersonID: 1, Role: entity.Role{Name: "actor"}, ContentID: 1},
				}, nil).AnyTimes()
				repo.EXPECT().GetPreviewContent(1).Return(&entity.Content{
					ID:             1,
					PosterStaticID: 1,
				}, nil).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mock_usecase.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("http://localhost:8080/static/1", nil).AnyTimes()
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
			mockStaticRepo := mock_usecase.NewMockStatic(ctrl)
			contentService := NewContentService(mockContentRepo, mockStaticRepo, "")
			tc.SetupContentRepoMock(mockContentRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			person, err := contentService.GetPersonByID(tc.PersonID)
			require.EqualValues(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOutput, person)
		})
	}
}

func TestContentService_GetPreviewPersonByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		PersonID             int
		ExpectedOutput       *dto.PersonPreviewWithPhoto
		ExpectedErr          error
		SetupContentRepoMock func(repo *mockrepo.MockContent)
		SetupStaticRepoMock  func(repo *mock_usecase.MockStatic)
	}{
		{
			Name:     "Существующий человек",
			PersonID: 1,
			ExpectedOutput: &dto.PersonPreviewWithPhoto{
				ID:       1,
				Name:     "Кристиан Бэйл",
				EnName:   "Christian Bale",
				PhotoURL: "http://localhost:8080/static/1",
			},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetPerson(1).Return(&entity.Person{
					ID:            1,
					Name:          "Кристиан Бэйл",
					EnName:        "Christian Bale",
					PhotoStaticID: sql.NullInt64{Int64: 1, Valid: true},
				}, nil).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mock_usecase.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("http://localhost:8080/static/1", nil).AnyTimes()
			},
		},
		{
			Name:        "Ошибка при получении персоны",
			PersonID:    2,
			ExpectedErr: entity.UsecaseWrap(fmt.Errorf("ошибка при получении персоны"), fmt.Errorf("ошибка при получении человека")),
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetPerson(2).Return(nil, fmt.Errorf("ошибка при получении человека")).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mock_usecase.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("", usecase.ErrStaticNotFound).AnyTimes()
			},
		},
		{
			Name:        "Ошибка при получении фото",
			PersonID:    1,
			ExpectedErr: entity.UsecaseWrap(fmt.Errorf("ошибка при получении фото персоны"), fmt.Errorf("ошибка при получении фото")),
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetPerson(1).Return(&entity.Person{
					ID:            1,
					Name:          "Кристиан Бэйл",
					EnName:        "Christian Bale",
					PhotoStaticID: sql.NullInt64{Int64: 1, Valid: true},
				}, nil).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mock_usecase.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("", fmt.Errorf("ошибка при получении фото")).AnyTimes()
			},
		},
		{
			Name:        "Несуществующий человек",
			PersonID:    2,
			ExpectedErr: usecase.ErrPersonNotFound,
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetPerson(2).Return(nil, repository.ErrPersonNotFound).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mock_usecase.MockStatic) {},
		},
		{
			Name:        "Без фото",
			PersonID:    1,
			ExpectedErr: nil,
			ExpectedOutput: &dto.PersonPreviewWithPhoto{
				ID:       1,
				Name:     "Кристиан Бэйл",
				EnName:   "Christian Bale",
				PhotoURL: "",
			},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetPerson(1).Return(&entity.Person{
					ID:            1,
					Name:          "Кристиан Бэйл",
					EnName:        "Christian Bale",
					PhotoStaticID: sql.NullInt64{Int64: 1, Valid: true},
				}, nil).AnyTimes()
			},
			SetupStaticRepoMock: func(uc *mock_usecase.MockStatic) {
				uc.EXPECT().GetStatic(1).Return("", usecase.ErrStaticNotFound).AnyTimes()
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
			mockStaticRepo := mock_usecase.NewMockStatic(ctrl)
			contentService := NewContentService(mockContentRepo, mockStaticRepo, "")
			tc.SetupContentRepoMock(mockContentRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			person, err := contentService.GetPreviewPersonByID(tc.PersonID)
			require.EqualValues(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOutput, person)
		})
	}
}

func TestContentService_GetPreviewContentByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		ContentID            int
		ExpectedOutput       *dto.PreviewContent
		ExpectedErr          error
		SetupContentRepoMock func(repo *mockrepo.MockContent)
		SetupStaticRepoMock  func(repo *mock_usecase.MockStatic)
	}{
		{
			Name:      "Существующий контент",
			ContentID: 1,
			ExpectedOutput: &dto.PreviewContent{
				ID:      1,
				Poster:  "http://localhost:8080/static/1",
				Genre:   "Боевик",
				Country: "USA",
				Actors:  make([]string, 0),
			},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetPreviewContent(1).Return(&entity.Content{
					ID:             1,
					PosterStaticID: 1,
					Genres:         []entity.Genre{{Name: "Боевик"}},
					Country:        []entity.Country{{Name: "USA"}},
				}, nil).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mock_usecase.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("http://localhost:8080/static/1", nil).AnyTimes()
			},
		},
		{
			Name:        "Ошибка при получении контента",
			ContentID:   2,
			ExpectedErr: entity.UsecaseWrap(fmt.Errorf("ошибка при получении контента"), fmt.Errorf("ошибка при получении контента")),
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetPreviewContent(2).Return(nil, fmt.Errorf("ошибка при получении контента")).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mock_usecase.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("", usecase.ErrStaticNotFound).AnyTimes()
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
			mockStaticRepo := mock_usecase.NewMockStatic(ctrl)
			contentService := NewContentService(mockContentRepo, mockStaticRepo, "")
			tc.SetupContentRepoMock(mockContentRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			content, err := contentService.GetPreviewContentByID(tc.ContentID)
			require.EqualValues(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOutput, content)
		})
	}
}

func TestContentService_GetNearestOngoings(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		ExpectedOutput       *dto.PreviewOngoingContentList
		ExpectedErr          error
		SetupContentRepoMock func(repo *mockrepo.MockContent)
		SetupStaticRepoMock  func(repo *mock_usecase.MockStatic)
	}{
		{
			Name: "Существующие контенты",
			ExpectedOutput: &dto.PreviewOngoingContentList{
				OnGoingContentList: []*dto.PreviewContent{
					{
						ID:     1,
						Poster: "http://localhost:8080/static/1",
						Genre:  "Боевик",
						Actors: make([]string, 0),
					},
				},
			},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetNearestOngoings(10).Return([]int{1}, nil).AnyTimes()
				repo.EXPECT().GetPreviewContent(1).Return(&entity.Content{
					ID:             1,
					PosterStaticID: 1,
					Genres:         []entity.Genre{{Name: "Боевик"}},
				}, nil).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mock_usecase.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("http://localhost:8080/static/1", nil).AnyTimes()
			},
		},
		{
			Name:        "Ошибка при получении контента",
			ExpectedErr: entity.UsecaseWrap(fmt.Errorf("ошибка при получении ближайших релизов"), fmt.Errorf("ошибка при получении контента")),
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetNearestOngoings(10).Return(nil, fmt.Errorf("ошибка при получении контента")).AnyTimes()
			},
			SetupStaticRepoMock: func(repo *mock_usecase.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("", usecase.ErrStaticNotFound).AnyTimes()
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
			mockStaticRepo := mock_usecase.NewMockStatic(ctrl)
			contentService := NewContentService(mockContentRepo, mockStaticRepo, "")
			tc.SetupContentRepoMock(mockContentRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			content, err := contentService.GetNearestOngoings()
			require.EqualValues(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOutput, content)
		})
	}
}

func TestContentService_GetAllOngoingsYears(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		ExpectedOutput       *dto.ReleaseYearsResponse
		ExpectedErr          error
		SetupContentRepoMock func(repo *mockrepo.MockContent)
	}{
		{
			Name:           "Существующие года",
			ExpectedOutput: &dto.ReleaseYearsResponse{Years: []int{2021}},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetAllOngoingsYears().Return([]int{2021}, nil).AnyTimes()
			},
		},
		{
			Name:        "Ошибка при получении годов",
			ExpectedErr: entity.UsecaseWrap(fmt.Errorf("ошибка при получении всех годов релизов"), fmt.Errorf("ошибка при получении годов")),
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetAllOngoingsYears().Return(nil, fmt.Errorf("ошибка при получении годов")).AnyTimes()
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
			contentService := NewContentService(mockContentRepo, nil, "")
			tc.SetupContentRepoMock(mockContentRepo)
			content, err := contentService.GetAllOngoingsYears()
			require.EqualValues(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOutput, content)
		})
	}
}

func TestContentService_SubscribeOnContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		UserID               int
		ContentID            int
		ExpectedErr          error
		SetupContentRepoMock func(repo *mockrepo.MockContent)
	}{
		{
			Name:      "Успешная подписка",
			UserID:    1,
			ContentID: 1,
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().SubscribeOnContent(1, 1).Return(nil).AnyTimes()
			},
		},
		{
			Name:        "Ошибка при подписке",
			UserID:      1,
			ContentID:   1,
			ExpectedErr: entity.UsecaseWrap(fmt.Errorf("ошибка при подписке на контент"), fmt.Errorf("ошибка при подписке на контент")),
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().SubscribeOnContent(1, 1).Return(fmt.Errorf("ошибка при подписке на контент")).AnyTimes()
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
			contentService := NewContentService(mockContentRepo, nil, "")
			tc.SetupContentRepoMock(mockContentRepo)
			err := contentService.SubscribeOnContent(tc.UserID, tc.ContentID)
			require.EqualValues(t, tc.ExpectedErr, err)
		})
	}
}

func TestContentService_UnsubscribeFromContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		UserID               int
		ContentID            int
		ExpectedErr          error
		SetupContentRepoMock func(repo *mockrepo.MockContent)
	}{
		{
			Name:      "Успешная отписка",
			UserID:    1,
			ContentID: 1,
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().UnsubscribeFromContent(1, 1).Return(nil).AnyTimes()
			},
		},
		{
			Name:        "Ошибка при отписке",
			UserID:      1,
			ContentID:   1,
			ExpectedErr: entity.UsecaseWrap(fmt.Errorf("ошибка при отписке от контента"), fmt.Errorf("ошибка при отписке от контента")),
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().UnsubscribeFromContent(1, 1).Return(fmt.Errorf("ошибка при отписке от контента")).AnyTimes()
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
			contentService := NewContentService(mockContentRepo, nil, "")
			tc.SetupContentRepoMock(mockContentRepo)
			err := contentService.UnsubscribeFromContent(tc.UserID, tc.ContentID)
			require.EqualValues(t, tc.ExpectedErr, err)
		})
	}
}

func TestContentService_GetSubscribedContentIDs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		UserID               int
		ExpectedOutput       *dto.SubscriptionsResponse
		ExpectedErr          error
		SetupContentRepoMock func(repo *mockrepo.MockContent)
	}{
		{
			Name:           "Существующие контенты",
			UserID:         1,
			ExpectedOutput: &dto.SubscriptionsResponse{Subscriptions: []int{1}},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetSubscribedContentIDs(1).Return([]int{1}, nil).AnyTimes()
			},
		},
		{
			Name:        "Ошибка при получении контентов",
			UserID:      1,
			ExpectedErr: entity.UsecaseWrap(fmt.Errorf("ошибка при получении подписок пользователя"), fmt.Errorf("ошибка при получении подписок")),
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetSubscribedContentIDs(1).Return(nil, fmt.Errorf("ошибка при получении подписок")).AnyTimes()
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
			contentService := NewContentService(mockContentRepo, nil, "")
			tc.SetupContentRepoMock(mockContentRepo)
			content, err := contentService.GetSubscribedContentIDs(tc.UserID)
			require.EqualValues(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOutput, content)
		})
	}
}
