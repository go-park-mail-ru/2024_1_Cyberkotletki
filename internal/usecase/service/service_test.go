package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

var (
	defaultEntityContentMovie = []entity.Content{
		{
			ID: 1,
			Country: []entity.Country{
				{
					Name: "Россия",
				},
			},
			Genres: []entity.Genre{
				{
					Name: "Комедия",
				},
			},
			PosterStaticID: 1,
			Type:           "movie",
			Movie: &entity.Movie{
				Duration: 1,
			},
			Actors: []entity.Person{
				{
					Name: "Иван",
				},
			},
		},
	}
	defaultEntityContentSeries = []entity.Content{
		{
			ID:             1,
			Country:        []entity.Country{},
			Genres:         []entity.Genre{},
			PosterStaticID: 1,
			Type:           "series",
			Series: &entity.Series{
				Seasons: []entity.Season{},
			},
			Actors: []entity.Person{
				{
					Name: "Иван",
				},
			},
			Directors: []entity.Person{
				{
					Name: "Петр",
				},
			},
		},
	}
	defaultEntityPerson = []entity.Person{
		{
			ID:            1,
			Name:          "Иван",
			EnName:        "Ivan",
			PhotoStaticID: 2,
		},
		{
			ID:            2,
			Name:          "Петр",
			EnName:        "Petr",
			PhotoStaticID: 3,
		},
	}
)

func TestSearchService_Search(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		Input               string
		ExpectedErr         error
		SetupSearchRepoMock func(repo *mockrepo.MockSearch)
		SetupStaticUCMock   func(uc *mockusecase.MockStatic)
	}{
		{
			Name:        "Успешный поиск фильма",
			Input:       "query",
			ExpectedErr: nil,
			SetupSearchRepoMock: func(repo *mockrepo.MockSearch) {
				repo.EXPECT().SearchContent(gomock.Any()).Return(defaultEntityContentMovie, nil)
				repo.EXPECT().SearchPerson(gomock.Any()).Return(defaultEntityPerson, nil)
			},
			SetupStaticUCMock: func(uc *mockusecase.MockStatic) {
				uc.EXPECT().GetStatic(gomock.Any()).Return("posterURL", nil).AnyTimes()
			},
		},
		{
			Name:        "Успешный поиск сериала",
			Input:       "query",
			ExpectedErr: nil,
			SetupSearchRepoMock: func(repo *mockrepo.MockSearch) {
				repo.EXPECT().SearchContent(gomock.Any()).Return(defaultEntityContentSeries, nil)
				repo.EXPECT().SearchPerson(gomock.Any()).Return(defaultEntityPerson, nil)
			},
			SetupStaticUCMock: func(uc *mockusecase.MockStatic) {
				uc.EXPECT().GetStatic(gomock.Any()).Return("posterURL", nil).AnyTimes()
			},
		},
		{
			Name:        "Ошибка поиска контента",
			Input:       "query",
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при поиске контента в SearchService"), errors.New("ошибка при поиске контента в SearchService")),
			SetupSearchRepoMock: func(repo *mockrepo.MockSearch) {
				repo.EXPECT().SearchContent(gomock.Any()).Return(nil, errors.New("ошибка при поиске контента в SearchService"))
			},
			SetupStaticUCMock: func(uc *mockusecase.MockStatic) {

			},
		},
		{
			Name:        "Ошибка поиска персон",
			Input:       "query",
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при поиске персон в SearchService"), errors.New("ошибка при поиске персон в SearchService")),
			SetupSearchRepoMock: func(repo *mockrepo.MockSearch) {
				repo.EXPECT().SearchContent(gomock.Any()).Return([]entity.Content{}, nil)
				repo.EXPECT().SearchPerson(gomock.Any()).Return(nil, errors.New("ошибка при поиске персон в SearchService"))
			},
			SetupStaticUCMock: func(uc *mockusecase.MockStatic) {

			},
		},
		{
			Name:        "Ошибка поиска статики",
			Input:       "query",
			ExpectedErr: entity.UsecaseWrap(errors.New("123"), errors.New("ошибка при получении статики контента из Search")),
			SetupSearchRepoMock: func(repo *mockrepo.MockSearch) {
				repo.EXPECT().SearchContent(gomock.Any()).Return(defaultEntityContentMovie, nil)
				repo.EXPECT().SearchPerson(gomock.Any()).Return(defaultEntityPerson, nil)
			},
			SetupStaticUCMock: func(uc *mockusecase.MockStatic) {
				uc.EXPECT().GetStatic(gomock.Any()).Return("", errors.New("123")).AnyTimes()
			},
		},
		{
			Name:        "Ошибка поиска статики контента",
			Input:       "query",
			ExpectedErr: nil,
			SetupSearchRepoMock: func(repo *mockrepo.MockSearch) {
				repo.EXPECT().SearchContent(gomock.Any()).Return(defaultEntityContentMovie, nil)
				repo.EXPECT().SearchPerson(gomock.Any()).Return(defaultEntityPerson, nil)
			},
			SetupStaticUCMock: func(uc *mockusecase.MockStatic) {
				uc.EXPECT().GetStatic(gomock.Any()).Return("", usecase.ErrStaticNotFound).AnyTimes()
			},
		},
		{
			Name:        "Ошибка поиска статики персоны",
			Input:       "query",
			ExpectedErr: entity.UsecaseWrap(errors.New("123"), errors.New("ошибка при получении статики персоны из Search")),
			SetupSearchRepoMock: func(repo *mockrepo.MockSearch) {
				repo.EXPECT().SearchContent(gomock.Any()).Return(defaultEntityContentMovie, nil)
				repo.EXPECT().SearchPerson(gomock.Any()).Return(defaultEntityPerson, nil)
			},
			SetupStaticUCMock: func(uc *mockusecase.MockStatic) {
				uc.EXPECT().GetStatic(1).Return("", nil)
				uc.EXPECT().GetStatic(gomock.Any()).Return("", errors.New("123")).AnyTimes()
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSearchRepo := mockrepo.NewMockSearch(ctrl)
			mockStaticUC := mockusecase.NewMockStatic(ctrl)
			searchService := NewSearchService(mockSearchRepo, mockStaticUC)
			tc.SetupSearchRepoMock(mockSearchRepo)
			tc.SetupStaticUCMock(mockStaticUC)
			_, err := searchService.Search(tc.Input)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
