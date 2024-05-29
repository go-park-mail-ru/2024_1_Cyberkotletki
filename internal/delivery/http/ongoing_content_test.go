package http

import (
	"errors"
	"fmt"
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

func TestOngoingContentEndpoints_GetNearestOngoings(t *testing.T) {
	t.Parallel()
	releaseYear := 2022

	testCases := []struct {
		Name                           string
		ExpectedErr                    error
		ExpectedOutput                 *dto.PreviewOngoingContentList
		SetupOngoingContentUsecaseMock func(mock *mockusecase.MockContent)
	}{
		{
			Name:        "Успех",
			ExpectedErr: nil,
			// несколько примеров, из которых будет 1 ближайший
			ExpectedOutput: &dto.PreviewOngoingContentList{
				OnGoingContentList: []*dto.PreviewContent{
					{
						ID:          1,
						Title:       "Бэтмен",
						Genre:       "Боевик",
						Poster:      "/static/poster.jpg",
						ReleaseYear: releaseYear,
						Type:        "movie",
					},
					{
						ID:          2,
						Title:       "Супермен",
						Genre:       "Боевик",
						Poster:      "/static/poster.jpg",
						ReleaseYear: releaseYear,
						Type:        "movie",
					},
				},
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetNearestOngoings().Return(&dto.PreviewOngoingContentList{
					OnGoingContentList: []*dto.PreviewContent{
						{
							ID:          1,
							Title:       "Бэтмен",
							Genre:       "Боевик",
							Poster:      "/static/poster.jpg",
							ReleaseYear: releaseYear,
							Type:        "movie",
						},
					},
				}, nil)
			},
		},
		{
			Name:        "Контент не найден",
			ExpectedErr: &echo.HTTPError{Code: 404, Message: "контент календаря релизов не найден"},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetNearestOngoings().Return(nil, usecase.ErrContentNotFound)
			},
		},
		{
			Name:        "Неожиданная ошибка",
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "ошибка при получении ближайших релизов", Internal: errors.New("123")},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetNearestOngoings().Return(nil, errors.New("123"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentUsecase := mockusecase.NewMockContent(ctrl)
			mockAuthUseCase := mockusecase.NewMockAuth(ctrl)
			ongoingContentEndpoints := NewOngoingContentEndpoints(mockOngoingContentUsecase, mockAuthUseCase)
			tc.SetupOngoingContentUsecaseMock(mockOngoingContentUsecase)
			req := httptest.NewRequest(http.MethodGet, "/ongoing/nearest", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/ongoing/nearest")
			err := ongoingContentEndpoints.GetNearestOngoings(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestOngoingContentEndpoints_GetOngoingContentByMonthAndYear(t *testing.T) {
	t.Parallel()
	releaseMonth := 5
	releaseYear := 2025

	testCases := []struct {
		Name                           string
		Month                          string
		Year                           string
		ExpectedErr                    error
		ExpectedOutput                 *dto.PreviewOngoingContentList
		SetupOngoingContentUsecaseMock func(mock *mockusecase.MockContent)
	}{
		{
			Name:        "Успех",
			Month:       "5",
			Year:        "2025",
			ExpectedErr: nil,
			// несколько примеров, из которых будет 1 ближайший
			ExpectedOutput: &dto.PreviewOngoingContentList{
				OnGoingContentList: []*dto.PreviewContent{
					{
						ID:          1,
						Title:       "Бэтмен",
						Genre:       "Боевик",
						Poster:      "/static/poster.jpg",
						ReleaseYear: releaseYear,
						Type:        "movie",
					},
				},
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetOngoingContentByMonthAndYear(releaseMonth, releaseYear).Return(&dto.PreviewOngoingContentList{
					OnGoingContentList: []*dto.PreviewContent{
						{
							ID:          1,
							Title:       "Бэтмен",
							Genre:       "Боевик",
							Poster:      "/static/poster.jpg",
							ReleaseYear: releaseYear,
							Type:        "movie",
						},
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
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {},
		},
		{
			Name:  "Ошибка валидации",
			Month: "1",
			Year:  "abc",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный год",
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {},
		},
		{
			Name:  "Контент не найден",
			Month: "1",
			Year:  "2025",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "контент календаря релизов не найден",
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetOngoingContentByMonthAndYear(1, 2025).Return(nil, usecase.ErrContentNotFound)
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
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
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
			mockOngoingContentUsecase := mockusecase.NewMockContent(ctrl)
			mockAuthUseCase := mockusecase.NewMockAuth(ctrl)
			ongoingContentEndpoints := NewOngoingContentEndpoints(mockOngoingContentUsecase, mockAuthUseCase)
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
		ExpectedOutput                 *dto.ReleaseYearsResponse
		SetupOngoingContentUsecaseMock func(mock *mockusecase.MockContent)
	}{
		{
			Name:        "Успех",
			ExpectedErr: nil,
			ExpectedOutput: &dto.ReleaseYearsResponse{
				Years: []int{2022, 2023, 2024},
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetAllOngoingsYears().Return(&dto.ReleaseYearsResponse{
					Years: []int{2022, 2023, 2024},
				}, nil)
			},
		},
		{
			Name: "Контент не найден",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "года релизов не найдены",
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetAllOngoingsYears().Return(nil, usecase.ErrContentNotFound)
			},
		},
		{
			Name: "Неожиданная ошибка",
			ExpectedErr: &echo.HTTPError{
				Code:     500,
				Message:  "ошибка при получении годов релизов",
				Internal: errors.New("123"),
			},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetAllOngoingsYears().Return(nil, errors.New("123"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentUsecase := mockusecase.NewMockContent(ctrl)
			mockAuthUseCase := mockusecase.NewMockAuth(ctrl)
			ongoingContentEndpoints := NewOngoingContentEndpoints(mockOngoingContentUsecase, mockAuthUseCase)
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

func TestOngoingContentEndpoints_SetReleasedState(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                           string
		ExpectedErr                    error
		ID                             string
		SecretKey                      string
		IsReleased                     string
		SetupOngoingContentUsecaseMock func(mock *mockusecase.MockContent)
	}{
		{
			Name:        "Успех",
			ExpectedErr: nil,
			ID:          "1",
			SecretKey:   "123",
			IsReleased:  "true",
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().SetReleasedState("123", 1, true).Return(nil)
			},
		},
		{
			Name: "Неожиданная ошибка",
			ExpectedErr: &echo.HTTPError{
				Code:     500,
				Message:  "ошибка при установке состояния релиза",
				Internal: errors.New("123"),
			},
			ID:         "1",
			SecretKey:  "123",
			IsReleased: "true",
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().SetReleasedState("123", 1, true).Return(errors.New("123"))
			},
		},
		{
			Name: "Ошибка валидации id",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный ID",
			},
			ID:                             "dd0",
			SecretKey:                      "123",
			IsReleased:                     "true",
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {},
		},
		{
			Name: "Ошибка валидации is_released",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидное значение is_released",
			},
			ID:                             "1",
			SecretKey:                      "123",
			IsReleased:                     "abc",
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {},
		},
		{
			Name: "неверный секретный ключ",
			ExpectedErr: &echo.HTTPError{
				Code:    403,
				Message: "неверный секретный ключ",
			},
			ID:         "1",
			SecretKey:  "1234",
			IsReleased: "true",
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().SetReleasedState("1234", 1, true).Return(usecase.ErrContentInvalidSecretKey)
			},
		},
		{
			Name: "контент не найден",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "контент календаря релизов не найден",
			},
			ID:         "1",
			SecretKey:  "123",
			IsReleased: "true",
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().SetReleasedState("123", 1, true).Return(usecase.ErrContentNotFound)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentUsecase := mockusecase.NewMockContent(ctrl)
			mockAuthUseCase := mockusecase.NewMockAuth(ctrl)
			ongoingContentEndpoints := NewOngoingContentEndpoints(mockOngoingContentUsecase, mockAuthUseCase)
			tc.SetupOngoingContentUsecaseMock(mockOngoingContentUsecase)
			req := httptest.NewRequest(
				http.MethodPut,
				fmt.Sprintf("/ongoing/1/is_released?is_released=%v&secret_key=%v", tc.IsReleased, tc.SecretKey),
				nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/ongoing/:id/is_released")
			c.SetParamNames("id")
			c.SetParamValues(tc.ID)
			err := ongoingContentEndpoints.SetReleasedState(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestOngoingContentEndpoints_SubscribeOnContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                           string
		ID                             string
		ExpectedErr                    error
		Cookies                        *http.Cookie
		SetupOngoingContentUsecaseMock func(mock *mockusecase.MockContent)
	}{
		{
			Name:        "Успех",
			ExpectedErr: nil,
			ID:          "1",
			Cookies:     &http.Cookie{Name: "session", Value: "123"},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().SubscribeOnContent(1, 1).Return(nil)
			},
		},
		{
			Name: "Неожиданная ошибка",
			ExpectedErr: &echo.HTTPError{
				Code:     500,
				Message:  "ошибка при подписке на контент",
				Internal: errors.New("123"),
			},
			ID:      "1",
			Cookies: &http.Cookie{Name: "session", Value: "123"},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().SubscribeOnContent(1, 1).Return(errors.New("123"))
			},
		},
		{
			Name: "контент не найден",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "контент не найден",
			},
			ID:      "1",
			Cookies: &http.Cookie{Name: "session", Value: "123"},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().SubscribeOnContent(1, 1).Return(usecase.ErrContentNotFound)
			},
		},
		{
			Name: "Пользователь не найден",
			ExpectedErr: &echo.HTTPError{
				Code:    401,
				Message: "не авторизован",
			},
			ID:      "1",
			Cookies: &http.Cookie{Name: "session", Value: "123"},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().SubscribeOnContent(1, 1).Return(usecase.ErrUserNotFound)
			},
		},
		{
			Name:        "Не авторизован",
			ExpectedErr: &echo.HTTPError{Code: 401, Message: "Не авторизован"},
			ID:          "1",
			Cookies:     nil,
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().SubscribeOnContent(1, 1).Times(0)
			},
		},
		{
			Name:        "Невалидный ID",
			ExpectedErr: &echo.HTTPError{Code: 400, Message: "невалидный ID"},
			ID:          "abc",
			Cookies:     &http.Cookie{Name: "session", Value: "123"},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().SubscribeOnContent(1, 1).Times(0)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentUsecase := mockusecase.NewMockContent(ctrl)
			mockAuthUseCase := mockusecase.NewMockAuth(ctrl)
			ongoingContentEndpoints := NewOngoingContentEndpoints(mockOngoingContentUsecase, mockAuthUseCase)
			tc.SetupOngoingContentUsecaseMock(mockOngoingContentUsecase)
			req := httptest.NewRequest(http.MethodPost, "/ongoing/1/subscribe", nil)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/ongoing/:id/subscribe")
			c.SetParamNames("id")
			c.SetParamValues(tc.ID)
			mockAuthUseCase.EXPECT().GetUserIDBySession("123").Return(1, nil).AnyTimes()
			err := ongoingContentEndpoints.SubscribeOnContent(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestOngoingContentEndpoints_UnsubscribeFromContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                           string
		ID                             string
		ExpectedErr                    error
		Cookies                        *http.Cookie
		SetupOngoingContentUsecaseMock func(mock *mockusecase.MockContent)
	}{
		{
			Name:        "Успех",
			ExpectedErr: nil,
			ID:          "1",
			Cookies:     &http.Cookie{Name: "session", Value: "123"},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().UnsubscribeFromContent(1, 1).Return(nil)
			},
		},
		{
			Name: "Невалидный ID",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный ID",
			},
			ID:                             "abc",
			Cookies:                        &http.Cookie{Name: "session", Value: "123"},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {},
		},
		{
			Name: "Не авторизован",
			ExpectedErr: &echo.HTTPError{
				Code:    401,
				Message: "Не авторизован",
			},
			ID:                             "1",
			Cookies:                        nil,
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentUsecase := mockusecase.NewMockContent(ctrl)
			mockAuthUseCase := mockusecase.NewMockAuth(ctrl)
			ongoingContentEndpoints := NewOngoingContentEndpoints(mockOngoingContentUsecase, mockAuthUseCase)
			tc.SetupOngoingContentUsecaseMock(mockOngoingContentUsecase)
			req := httptest.NewRequest(http.MethodDelete, "/ongoing/1/subscribe", nil)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/ongoing/:id/subscribe")
			c.SetParamNames("id")
			c.SetParamValues(tc.ID)
			mockAuthUseCase.EXPECT().GetUserIDBySession("123").Return(1, nil).AnyTimes()
			err := ongoingContentEndpoints.UnsubscribeFromContent(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestOngoingContentEndpoints_GetSubscribedContentIDs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                           string
		ExpectedErr                    error
		Cookies                        *http.Cookie
		SetupOngoingContentUsecaseMock func(mock *mockusecase.MockContent)
	}{
		{
			Name:        "Успех",
			ExpectedErr: nil,
			Cookies:     &http.Cookie{Name: "session", Value: "123"},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetSubscribedContentIDs(1).Return(&dto.SubscriptionsResponse{Subscriptions: []int{1}}, nil)
			},
		},
		{
			Name: "Не авторизован",
			ExpectedErr: &echo.HTTPError{
				Code:    401,
				Message: "Не авторизован",
			},
			Cookies:                        nil,
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {},
		},
		{
			Name: "Неожиданная ошибка",
			ExpectedErr: &echo.HTTPError{
				Code:     500,
				Message:  "ошибка при получении подписок",
				Internal: errors.New("123"),
			},
			Cookies: &http.Cookie{Name: "session", Value: "123"},
			SetupOngoingContentUsecaseMock: func(mock *mockusecase.MockContent) {
				mock.EXPECT().GetSubscribedContentIDs(1).Return(nil, errors.New("123"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockOngoingContentUsecase := mockusecase.NewMockContent(ctrl)
			mockAuthUseCase := mockusecase.NewMockAuth(ctrl)
			ongoingContentEndpoints := NewOngoingContentEndpoints(mockOngoingContentUsecase, mockAuthUseCase)
			tc.SetupOngoingContentUsecaseMock(mockOngoingContentUsecase)
			req := httptest.NewRequest(http.MethodGet, "/ongoing/subscribed", nil)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/ongoing/subscribed")
			mockAuthUseCase.EXPECT().GetUserIDBySession("123").Return(1, nil).AnyTimes()
			err := ongoingContentEndpoints.GetSubscribedContentIDs(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
