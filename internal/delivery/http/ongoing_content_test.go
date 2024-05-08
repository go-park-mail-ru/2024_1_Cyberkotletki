package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"time"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestOngoingContentEndpoints_GetOngoingContentByContentID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                           string
		OngoingID                      string
		ExpectedErr                    error
		SetupOngoingContentUsecaseMock func(mock *mockusecase.MockOngoingContent)
	}{
		{
			Name:        "Успех",
			OngoingID:   "1",
			ExpectedErr: nil,
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {
				mock.EXPECT().GetOngoingContentByContentID(1).Return(nil, nil)
			},
		},
		{
			Name:                           "Ошибка валидации",
			OngoingID:                      "abc",
			ExpectedErr:                    &echo.HTTPError{Code: 400, Message: "невалидный id контента календаря релизов"},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {},
		},
		{
			Name:        "Контент не найден",
			OngoingID:   "1",
			ExpectedErr: &echo.HTTPError{Code: 404, Message: "контент календаря релизов не найден"},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {
				mock.EXPECT().GetOngoingContentByContentID(1).Return(nil, usecase.ErrOngoingContentNotFound)
			},
		},
		{
			Name:        "Неожиданная ошибка",
			OngoingID:   "1",
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "ошибка при получении контента календаря релизов", Internal: errors.New("123")},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {
				mock.EXPECT().GetOngoingContentByContentID(1).Return(nil, errors.New("123"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentUsecase := mockusecase.NewMockOngoingContent(ctrl)
			ongoingContentEndpoints := NewOngoingContentEndpoints(mockOngoingContentUsecase)
			tc.SetupOngoingContentUsecaseMock(mockOngoingContentUsecase)
			req := httptest.NewRequest(http.MethodGet, "/ongoing/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/ongoing/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.OngoingID)
			err := ongoingContentEndpoints.GetOngoingContentByContentID(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestOngoingContentEndpoints_GetNearestOngoings(t *testing.T) {
	t.Parallel()
	releaseTime := time.Now().Add(time.Hour)

	testCases := []struct {
		Name                           string
		Limit                          string
		ExpectedErr                    error
		ExpectedOutput                 []*dto.PreviewOngoingContentCardVertical
		SetupOngoingContentUsecaseMock func(mock *mockusecase.MockOngoingContent)
	}{
		{
			Name:        "Успех",
			Limit:       "1",
			ExpectedErr: nil,
			// несколько примеров, из которых будет 1 ближайший
			ExpectedOutput: []*dto.PreviewOngoingContentCardVertical{
				{
					ID:          1,
					Title:       "Бэтмен",
					Genres:      []string{"Боевик"},
					Poster:      "/static/poster.jpg",
					ReleaseDate: releaseTime,
					Type:        "movie",
				},
				{
					ID:          2,
					Title:       "Супермен",
					Genres:      []string{"Боевик"},
					Poster:      "/static/poster.jpg",
					ReleaseDate: releaseTime.Add(time.Hour),
					Type:        "movie",
				},
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {
				mock.EXPECT().GetNearestOngoings(1).Return([]*dto.PreviewOngoingContentCardVertical{
					{
						ID:          1,
						Title:       "Бэтмен",
						Genres:      []string{"Боевик"},
						Poster:      "/static/poster.jpg",
						ReleaseDate: releaseTime,
						Type:        "movie",
					},
				}, nil)
			},
		},
		{
			Name:  "Ошибка валидации",
			Limit: "abc",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидное количество релизов",
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {},
		},
		{
			Name:        "Контент не найден",
			Limit:       "1",
			ExpectedErr: &echo.HTTPError{Code: 404, Message: "контент календаря релизов не найден"},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {
				mock.EXPECT().GetNearestOngoings(1).Return(nil, usecase.ErrOngoingContentNotFound)
			},
		},
		{
			Name:        "Неожиданная ошибка",
			Limit:       "1",
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "ошибка при получении ближайших релизов", Internal: errors.New("123")},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {
				mock.EXPECT().GetNearestOngoings(1).Return(nil, errors.New("123"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentUsecase := mockusecase.NewMockOngoingContent(ctrl)
			ongoingContentEndpoints := NewOngoingContentEndpoints(mockOngoingContentUsecase)
			tc.SetupOngoingContentUsecaseMock(mockOngoingContentUsecase)
			req := httptest.NewRequest(http.MethodGet, "/ongoing/nearest", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/ongoing/nearest")
			c.QueryParams().Add("limit", tc.Limit)
			err := ongoingContentEndpoints.GetNearestOngoings(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestOngoingContentEndpoints_GetOngoingContentByMonthAndYear(t *testing.T) {
	t.Parallel()
	releaseMounth := 5
	releaseYear := 2025
	releaseTime := time.Date(releaseYear, time.Month(releaseMounth), 1, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		Name                           string
		Month                          string
		Year                           string
		ExpectedErr                    error
		ExpectedOutput                 []*dto.PreviewOngoingContentCardVertical
		SetupOngoingContentUsecaseMock func(mock *mockusecase.MockOngoingContent)
	}{
		{
			Name:        "Успех",
			Month:       "5",
			Year:        "2025",
			ExpectedErr: nil,
			// несколько примеров, из которых будет 1 ближайший
			ExpectedOutput: []*dto.PreviewOngoingContentCardVertical{
				{
					ID:          1,
					Title:       "Бэтмен",
					Genres:      []string{"Боевик"},
					Poster:      "/static/poster.jpg",
					ReleaseDate: releaseTime,
					Type:        "movie",
				},
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {
				mock.EXPECT().GetOngoingContentByMonthAndYear(releaseMounth, releaseYear).Return([]*dto.PreviewOngoingContentCardVertical{
					{
						ID:          1,
						Title:       "Бэтмен",
						Genres:      []string{"Боевик"},
						Poster:      "/static/poster.jpg",
						ReleaseDate: releaseTime,
						Type:        "movie",
					},
				}, nil)
			},
		},
		{
			Name:  "Ошибка валидации",
			Month: "abc",
			Year:  "2022",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный месяц",
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {},
		},
		{
			Name:  "Ошибка валидации",
			Month: "1",
			Year:  "abc",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный год",
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {},
		},
		{
			Name:  "Контент не найден",
			Month: "1",
			Year:  "2025",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "контент календаря релизов не найден",
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {
				mock.EXPECT().GetOngoingContentByMonthAndYear(1, 2025).Return(nil, usecase.ErrOngoingContentNotFound)
			},
		},
		{
			Name:  "Неожиданная ошибка",
			Month: "1",
			Year:  "2025",
			ExpectedErr: &echo.HTTPError{
				Code:     500,
				Message:  "ошибка при получении релизов по месяцу и году",
				Internal: errors.New("123"),
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {
				mock.EXPECT().GetOngoingContentByMonthAndYear(1, 2025).Return(nil, errors.New("123"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentUsecase := mockusecase.NewMockOngoingContent(ctrl)
			ongoingContentEndpoints := NewOngoingContentEndpoints(mockOngoingContentUsecase)
			tc.SetupOngoingContentUsecaseMock(mockOngoingContentUsecase)
			req := httptest.NewRequest(http.MethodGet, "/ongoing/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/ongoing/:month/:year")
			c.SetParamNames("month", "year")
			c.SetParamValues(tc.Month, tc.Year)
			err := ongoingContentEndpoints.GetOngoingContentByMonthAndYear(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestOngoingContentEndpoints_GetAllReleaseYears(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                           string
		ExpectedErr                    error
		ExpectedOutput                 []int
		SetupOngoingContentUsecaseMock func(mock *mockusecase.MockOngoingContent)
	}{
		{
			Name:        "Успех",
			ExpectedErr: nil,
			ExpectedOutput: []int{
				2022,
				2023,
				2024,
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {
				mock.EXPECT().GetAllReleaseYears().Return([]int{2022, 2023, 2024}, nil)
			},
		},
		{
			Name: "Контент не найден",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "года релизов не найдены",
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {
				mock.EXPECT().GetAllReleaseYears().Return(nil, usecase.ErrOngoingContentYearsNotFound)
			},
		},
		{
			Name: "Неожиданная ошибка",
			ExpectedErr: &echo.HTTPError{
				Code:     500,
				Message:  "ошибка при получении годов релизов",
				Internal: errors.New("123"),
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockOngoingContent) {
				mock.EXPECT().GetAllReleaseYears().Return(nil, errors.New("123"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentUsecase := mockusecase.NewMockOngoingContent(ctrl)
			ongoingContentEndpoints := NewOngoingContentEndpoints(mockOngoingContentUsecase)
			tc.SetupOngoingContentUsecaseMock(mockOngoingContentUsecase)
			req := httptest.NewRequest(http.MethodGet, "/ongoing/years", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/ongoing/years")
			err := ongoingContentEndpoints.GetAllReleaseYears(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
