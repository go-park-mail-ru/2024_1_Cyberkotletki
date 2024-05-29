package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	//"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"

	"io"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFavouriteEndpoints_CreateFavourite1(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                      string
		Input                     func() io.Reader
		Cookies                   *http.Cookie
		ExpectedErr               error
		SetupFavouriteUsecaseMock func(usecase *mockusecase.MockFavourite)
		SetupAuthUsecaseMock      func(usecase *mockusecase.MockAuth)
	}{
		{
			Name: "1111111111111Успешное создание",
			Input: func() io.Reader {
				return strings.NewReader(`{"contentID":1,"category":"favourite"}`)
			},

			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {
				usecase.EXPECT().CreateFavourite(1, 1, "favourite").Return(nil).AnyTimes()
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil).AnyTimes()
			},
		},
		{
			Name: "6666666666666Ошибка авторизации",
			Input: func() io.Reader {
				return strings.NewReader(`{"contentID":1,"category":"favourite"}`)
			},
			Cookies:                   &http.Cookie{Name: "session_id", Value: "xxx"},
			ExpectedErr:               &echo.HTTPError{Code: 401, Message: "Не авторизован"},
			SetupFavouriteUsecaseMock: func(mock *mockusecase.MockFavourite) {},
			SetupAuthUsecaseMock: func(mock *mockusecase.MockAuth) {
				mock.EXPECT().GetUserIDBySession(gomock.Any()).Return(-1, utils.ErrUnauthorized).AnyTimes()
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFavouriteUsecase := mockusecase.NewMockFavourite(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			tc.SetupFavouriteUsecaseMock(mockFavouriteUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)

			favouriteHandler := NewFavouriteEndpoints(mockFavouriteUsecase, mockAuthUsecase)
			req := httptest.NewRequest(http.MethodPut, "/favourite", tc.Input())

			req.AddCookie(tc.Cookies)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/favourite")
			c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			err := favouriteHandler.CreateFavourite(c)
			require.Equal(t, tc.ExpectedErr, err)

		})
	}
}

func TestFavouriteEndpoints_DeleteFavourite(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                      string
		FavouriteID               string
		ExpectedErr               error
		Cookies                   *http.Cookie
		SetupFavouriteUsecaseMock func(usecase *mockusecase.MockFavourite)
		SetupAuthUsecaseMock      func(usecase *mockusecase.MockAuth)
	}{
		{
			Name:        "Успешное удаление",
			FavouriteID: "1",
			ExpectedErr: nil,
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {
				usecase.EXPECT().DeleteFavourite(1, 1).Return(nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:        "Не найден",
			FavouriteID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "Избранное не найдено",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupFavouriteUsecaseMock: func(uc *mockusecase.MockFavourite) {
				uc.EXPECT().DeleteFavourite(1, 1).Return(usecase.ErrFavouriteNotFound)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:        "Пользователь не авторизован",
			FavouriteID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    401,
				Message: "Не авторизован",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession(gomock.Any()).Return(-1, utils.ErrUnauthorized)
			},
		},
		{
			Name:        "Неожиданная ошибка",
			FavouriteID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:     500,
				Message:  "Внутренняя ошибка сервера",
				Internal: errors.New("123"),
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {
				usecase.EXPECT().DeleteFavourite(1, 1).Return(errors.New("123"))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFavouriteUsecase := mockusecase.NewMockFavourite(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			tc.SetupFavouriteUsecaseMock(mockFavouriteUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)

			favouriteHandler := NewFavouriteEndpoints(mockFavouriteUsecase, mockAuthUsecase)
			req := httptest.NewRequest(http.MethodDelete, "/favourite/", nil)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/favourite/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.FavouriteID)
			err := favouriteHandler.DeleteFavourite(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewEndpoints_GetFavourites(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                      string
		UserID                    string
		Page                      string
		ExpectedErr               error
		ExpectedOutput            *dto.FavouritesResponse
		SetupFavouriteUsecaseMock func(usecase *mockusecase.MockFavourite)
	}{
		{
			Name:        "Успешное получение",
			UserID:      "1",
			Page:        "1",
			ExpectedErr: nil,
			ExpectedOutput: &dto.FavouritesResponse{
				Favourites: []dto.Favourite{
					{
						Content: dto.PreviewContent{
							ID: 1,
						},
						Category: "favourite",
					},
				},
			},

			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {
				usecase.EXPECT().GetFavourites(1).Return(&dto.FavouritesResponse{
					Favourites: []dto.Favourite{
						{
							Content: dto.PreviewContent{
								ID: 1,
							},
							Category: "favourite",
						},
					},
				}, nil)
			},
		},
		{
			Name:   "Неожиданная ошибка",
			UserID: "1",
			Page:   "1",
			ExpectedErr: &echo.HTTPError{
				Code:     500,
				Message:  "Внутренняя ошибка сервера",
				Internal: errors.New("123"),
			},
			ExpectedOutput: nil,
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {
				usecase.EXPECT().GetFavourites(1).Return(nil, errors.New("123"))
			},
		},
		{
			Name:   "Невалидный айди",
			UserID: "ogo!",
			Page:   "1",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "Невалидный ID",
			},
			ExpectedOutput:            nil,
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {},
		},
		{
			Name:   "Внутренняя ошибка сервера",
			UserID: "1",
			Page:   "1",
			ExpectedErr: &echo.HTTPError{
				Code:     500,
				Message:  "Внутренняя ошибка сервера",
				Internal: errors.New("123"),
			},
			ExpectedOutput: nil,
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {
				usecase.EXPECT().GetFavourites(1).Return(nil, errors.New("123"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockFavouriteUsecase := mockusecase.NewMockFavourite(ctrl)
			tc.SetupFavouriteUsecaseMock(mockFavouriteUsecase)

			favouriteHandler := NewFavouriteEndpoints(mockFavouriteUsecase, nil)
			req := httptest.NewRequest(http.MethodDelete, "/favourite/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/favourite/:id")
			c.SetParamNames("id", "page")
			c.SetParamValues(tc.UserID, tc.Page)
			err := favouriteHandler.GetFavouritesByUser(c)
			require.Equal(t, tc.ExpectedErr, err)
			if tc.ExpectedErr == nil {
				var reviews dto.FavouritesResponse
				err = json.NewDecoder(rec.Body).Decode(&reviews)
				require.NoError(t, err)
				require.Equal(t, *tc.ExpectedOutput, reviews)
			}
		})
	}
}

func TestReviewEndpoints_GetMyFavourites(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                      string
		Cookies                   *http.Cookie
		ExpectedErr               error
		ExpectedOutput            *dto.FavouritesResponse
		SetupFavouriteUsecaseMock func(usecase *mockusecase.MockFavourite)
		SetupAuthUsecaseMock      func(usecase *mockusecase.MockAuth)
	}{
		{
			Name: "Успешное получение",
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			ExpectedErr: nil,
			ExpectedOutput: &dto.FavouritesResponse{
				Favourites: []dto.Favourite{
					{
						Content: dto.PreviewContent{
							ID: 1,
						},
						Category: "favourite",
					},
				},
			},
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {
				usecase.EXPECT().GetFavourites(1).Return(&dto.FavouritesResponse{
					Favourites: []dto.Favourite{
						{
							Content: dto.PreviewContent{
								ID: 1,
							},
							Category: "favourite",
						},
					},
				}, nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name: "Неожиданная ошибка",
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			ExpectedErr: &echo.HTTPError{
				Code:     500,
				Message:  "Внутренняя ошибка сервера",
				Internal: errors.New("123"),
			},
			ExpectedOutput: nil,
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {
				usecase.EXPECT().GetFavourites(1).Return(nil, errors.New("123"))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name: "Пользователь не авторизован",
			ExpectedErr: &echo.HTTPError{
				Code:    401,
				Message: "Не авторизован",
			},
			ExpectedOutput:            nil,
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {},
			SetupAuthUsecaseMock:      func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name: "Пользователь не найден",
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "Пользователь не найден",
			},
			ExpectedOutput: nil,
			SetupFavouriteUsecaseMock: func(uc *mockusecase.MockFavourite) {
				uc.EXPECT().GetFavourites(1).Return(nil, usecase.ErrFavouriteUserNotFound)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockFavouriteUsecase := mockusecase.NewMockFavourite(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			tc.SetupFavouriteUsecaseMock(mockFavouriteUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)

			favouriteHandler := NewFavouriteEndpoints(mockFavouriteUsecase, mockAuthUsecase)
			req := httptest.NewRequest(http.MethodDelete, "/favourite/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/favourite/my")
			c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			err := favouriteHandler.GetMyFavourites(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewEndpoints_GetStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                      string
		Cookies                   *http.Cookie
		ContentID                 string
		ExpectedErr               error
		ExpectedOutput            *dto.FavouriteStatusResponse
		SetupFavouriteUsecaseMock func(usecase *mockusecase.MockFavourite)
		SetupAuthUsecaseMock      func(usecase *mockusecase.MockAuth)
	}{
		{
			Name: "Успешное получение",
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			ContentID:   "1",
			ExpectedErr: nil,
			ExpectedOutput: &dto.FavouriteStatusResponse{
				Status: "favourite",
			},
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {
				usecase.EXPECT().GetStatus(1, 1).Return(&dto.FavouriteStatusResponse{Status: "favourite"}, nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name: "Не найден",
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			ContentID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "Не добавлено в избранное",
			},
			ExpectedOutput: nil,
			SetupFavouriteUsecaseMock: func(uc *mockusecase.MockFavourite) {
				uc.EXPECT().GetStatus(1, 1).Return(nil, usecase.ErrFavouriteNotFound)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:      "Пользователь не авторизован",
			ContentID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    401,
				Message: "Не авторизован",
			},
			ExpectedOutput:            nil,
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {},
			SetupAuthUsecaseMock:      func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name: "Внутренняя ошибка сервера",
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			ContentID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:     500,
				Message:  "Внутренняя ошибка сервера",
				Internal: errors.New("123"),
			},
			ExpectedOutput: nil,
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {
				usecase.EXPECT().GetStatus(1, 1).Return(nil, errors.New("123"))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name: "Невалидный айди",
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			ContentID: "ogo!",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "Невалидный ID",
			},
			ExpectedOutput:            nil,
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name: "Неожиданная ошибка",
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			ContentID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:     500,
				Message:  "Внутренняя ошибка сервера",
				Internal: errors.New("123"),
			},
			ExpectedOutput: nil,
			SetupFavouriteUsecaseMock: func(usecase *mockusecase.MockFavourite) {
				usecase.EXPECT().GetStatus(1, 1).Return(nil, errors.New("123"))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockFavouriteUsecase := mockusecase.NewMockFavourite(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			tc.SetupFavouriteUsecaseMock(mockFavouriteUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			favouriteHandler := NewFavouriteEndpoints(mockFavouriteUsecase, mockAuthUsecase)
			req := httptest.NewRequest(http.MethodDelete, "/favourite/", nil)
			rec := httptest.NewRecorder()
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			c := e.NewContext(req, rec)
			c.SetPath("/favourite/status/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.ContentID)
			err := favouriteHandler.GetStatus(c)
			require.Equal(t, tc.ExpectedErr, err)
			if tc.ExpectedErr == nil {
				var reviews dto.FavouriteStatusResponse
				err = json.NewDecoder(rec.Body).Decode(&reviews)
				require.NoError(t, err)
				require.Equal(t, *tc.ExpectedOutput, reviews)
			}
		})
	}
}
