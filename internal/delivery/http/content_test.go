package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestContentEndpoints_GetContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                    string
		ContentID               string
		ExpectedErr             error
		ExpectedOutput          *dto.Content
		SetupContentUsecaseMock func(mock *mockusecase.MockContent)
	}{
		{
			Name:        "Успех",
			ContentID:   "1",
			ExpectedErr: nil,
			ExpectedOutput: &dto.Content{
				ID:             1,
				Title:          "Бэтмен",
				OriginalTitle:  "Batman",
				Slogan:         "I'm Batman",
				Budget:         "1000000",
				AgeRestriction: 18,
				Rating:         9.1,
				IMDBRating:     9.1,
				Description:    "Описание фильма или сериала",
				Facts:          []string{"Факт1", "Факт2"},
				TrailerLink:    "https://www.youtube.com/watch?v=123456",
				BackdropURL:    "/static/backdrop.jpg",
				PicturesURL:    []string{"/static/picture1.jpg", "/static/picture2.jpg"},
				PosterURL:      "/static/poster.jpg",
				Countries:      []string{"Россия", "США"},
				Genres:         []string{"Боевик", "Драма"},
				Type:           "movie",
			},
			SetupContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetContentByID(1).Return(&dto.Content{
					ID:             1,
					Title:          "Бэтмен",
					OriginalTitle:  "Batman",
					Slogan:         "I'm Batman",
					Budget:         "1000000",
					AgeRestriction: 18,
					Rating:         9.1,
					IMDBRating:     9.1,
					Description:    "Описание фильма или сериала",
					Facts:          []string{"Факт1", "Факт2"},
					TrailerLink:    "https://www.youtube.com/watch?v=123456",
					BackdropURL:    "/static/backdrop.jpg",
					PicturesURL:    []string{"/static/picture1.jpg", "/static/picture2.jpg"},
					PosterURL:      "/static/poster.jpg",
					Countries:      []string{"Россия", "США"},
					Genres:         []string{"Боевик", "Драма"},
					Type:           "movie",
				}, nil)
			},
		},
		{
			Name:                    "Ошибка валидации",
			ContentID:               "invalid",
			ExpectedErr:             &echo.HTTPError{Code: 400, Message: "Невалидный id контента"},
			ExpectedOutput:          nil,
			SetupContentUsecaseMock: func(mock *mockusecase.MockContent) {},
		},
		{
			Name:           "Контент не найден",
			ContentID:      "1",
			ExpectedErr:    &echo.HTTPError{Code: 404, Message: "Контент с таким id не найден"},
			ExpectedOutput: nil,
			SetupContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetContentByID(1).Return(nil, usecase.ErrContentNotFound)
			},
		},
		{
			Name:           "Внутренняя ошибка сервера",
			ContentID:      "1",
			ExpectedErr:    &echo.HTTPError{Code: 500, Message: "Внутренняя ошибка сервера", Internal: errors.New("123")},
			ExpectedOutput: nil,
			SetupContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetContentByID(1).Return(nil, errors.New("123"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockContentUsecase := mockusecase.NewMockContent(ctrl)
			contentEndpoints := NewContentEndpoints(mockContentUsecase)
			tc.SetupContentUsecaseMock(mockContentUsecase)
			req := httptest.NewRequest(http.MethodGet, "/content/"+tc.ContentID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/content/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.ContentID)
			err := contentEndpoints.GetContent(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestContentEndpoints_GetPerson(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                    string
		PersonID                string
		ExpectedErr             error
		ExpectedOutput          *dto.Person
		SetupContentUsecaseMock func(mock *mockusecase.MockContent)
	}{
		{
			Name:        "Успех",
			PersonID:    "1",
			ExpectedErr: nil,
			ExpectedOutput: &dto.Person{
				ID:        1,
				Name:      "Киану Ривз",
				EnName:    "Keanu Reeves",
				BirthDate: nil,
				DeathDate: nil,
				Sex:       "M",
				PhotoURL:  "/static/photo.jpg",
				Height:    185,
			},
			SetupContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetPersonByID(1).Return(&dto.Person{
					ID:        1,
					Name:      "Киану Ривз",
					EnName:    "Keanu Reeves",
					BirthDate: nil,
					DeathDate: nil,
					Sex:       "M",
					PhotoURL:  "/static/photo.jpg",
					Height:    185,
				}, nil)
			},
		},
		{
			Name:     "Ошибка валидации",
			PersonID: "invalid",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "Невалидный id персоны",
			},
			ExpectedOutput:          nil,
			SetupContentUsecaseMock: func(mock *mockusecase.MockContent) {},
		},
		{
			Name:           "Персона не найдена",
			PersonID:       "1",
			ExpectedErr:    &echo.HTTPError{Code: 404, Message: "Персона с таким id не найдена"},
			ExpectedOutput: nil,
			SetupContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetPersonByID(1).Return(nil, usecase.ErrPersonNotFound)
			},
		},
		{
			Name:           "Внутренняя ошибка сервера",
			PersonID:       "1",
			ExpectedErr:    &echo.HTTPError{Code: 500, Message: "Внутренняя ошибка сервера", Internal: errors.New("123")},
			ExpectedOutput: nil,
			SetupContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetPersonByID(1).Return(nil, errors.New("123"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockContentUsecase := mockusecase.NewMockContent(ctrl)
			contentEndpoints := NewContentEndpoints(mockContentUsecase)
			tc.SetupContentUsecaseMock(mockContentUsecase)
			req := httptest.NewRequest(http.MethodGet, "/content/person/"+tc.PersonID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/content/person/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.PersonID)
			err := contentEndpoints.GetPerson(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}

}
