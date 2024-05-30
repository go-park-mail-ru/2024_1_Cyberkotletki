package service

import (
	"database/sql"
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
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
			PhotoStaticID: sql.NullInt64{Int64: 2, Valid: true},
		},
		{
			ID:            2,
			Name:          "Петр",
			EnName:        "Petr",
			PhotoStaticID: sql.NullInt64{Int64: 3, Valid: true},
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
		SetupContentUCMock  func(uc *mockusecase.MockContent)
	}{
		{
			Name:        "Успешный поиск фильма",
			Input:       "query",
			ExpectedErr: nil,
			SetupSearchRepoMock: func(repo *mockrepo.MockSearch) {
				repo.EXPECT().SearchContent(gomock.Any()).Return([]int{1}, nil)
				repo.EXPECT().SearchPerson(gomock.Any()).Return([]int{1}, nil)
			},
			SetupContentUCMock: func(uc *mockusecase.MockContent) {
				uc.EXPECT().GetPreviewContentByID(1).Return(&dto.PreviewContent{}, nil)
				uc.EXPECT().GetPreviewPersonByID(1).Return(&dto.PersonPreviewWithPhoto{}, nil)
			},
		},
		{
			Name:        "Успешный поиск сериала",
			Input:       "query",
			ExpectedErr: nil,
			SetupSearchRepoMock: func(repo *mockrepo.MockSearch) {
				repo.EXPECT().SearchContent(gomock.Any()).Return([]int{1}, nil)
				repo.EXPECT().SearchPerson(gomock.Any()).Return([]int{1}, nil)
			},
			SetupContentUCMock: func(uc *mockusecase.MockContent) {
				uc.EXPECT().GetPreviewContentByID(1).Return(&dto.PreviewContent{}, nil)
				uc.EXPECT().GetPreviewPersonByID(1).Return(&dto.PersonPreviewWithPhoto{}, nil)
			},
		},
		{
			Name:        "Ошибка поиска контента",
			Input:       "query",
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при поиске контента в SearchService"), errors.New("ошибка при поиске контента в SearchService")),
			SetupSearchRepoMock: func(repo *mockrepo.MockSearch) {
				repo.EXPECT().SearchContent(gomock.Any()).Return(nil, errors.New("ошибка при поиске контента в SearchService"))
			},
			SetupContentUCMock: func(uc *mockusecase.MockContent) {

			},
		},
		{
			Name:        "Ошибка поиска персон",
			Input:       "query",
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при поиске персон в SearchService"), errors.New("ошибка при поиске персон в SearchService")),
			SetupSearchRepoMock: func(repo *mockrepo.MockSearch) {
				repo.EXPECT().SearchContent(gomock.Any()).Return([]int{1}, nil)
				repo.EXPECT().SearchPerson(gomock.Any()).Return(nil, errors.New("ошибка при поиске персон в SearchService"))
			},
			SetupContentUCMock: func(uc *mockusecase.MockContent) {

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
			mockContentUC := mockusecase.NewMockContent(ctrl)
			searchService := NewSearchService(mockSearchRepo, mockContentUC)
			tc.SetupSearchRepoMock(mockSearchRepo)
			tc.SetupContentUCMock(mockContentUC)
			_, err := searchService.Search(tc.Input)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
