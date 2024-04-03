package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthEndpoints_IsAuth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		ExpectedErr          error
		Cookies              *http.Cookie
		SetupAuthUsecaseMock func(usecase *mockusecase.MockAuth)
	}{
		{
			Name:        "Пользователь авторизован",
			ExpectedErr: nil,
			Cookies:     &http.Cookie{Name: "session", Value: "xxx"},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession(gomock.Any()).Return(1, nil)
			},
		},
		{
			Name:                 "Нет cookies",
			ExpectedErr:          &echo.HTTPError{Code: 401, Message: "отсутствует cookies с сессией"},
			Cookies:              nil,
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name:        "Пользователь не авторизован",
			ExpectedErr: &echo.HTTPError{Code: 401, Message: "необходима авторизация"},
			Cookies:     &http.Cookie{Name: "session", Value: "xxx"},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession(gomock.Any()).Return(-1, entity.NewClientError("необходима авторизация", entity.ErrUnauthorized))
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
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			authEndpoints := NewAuthEndpoints(mockAuthUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			req := httptest.NewRequest(http.MethodGet, "/auth/isAuth", nil)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := authEndpoints.IsAuth(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestAuthEndpoints_Logout(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		ExpectedErr          error
		Cookies              *http.Cookie
		SetupAuthUsecaseMock func(usecase *mockusecase.MockAuth)
	}{
		{
			Name:                 "Пользователь вышел, уже не имея сессии",
			ExpectedErr:          nil,
			Cookies:              nil,
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name:        "Пользователь вышел, изначально имея сессию",
			ExpectedErr: nil,
			Cookies:     &http.Cookie{Name: "session", Value: "xxx"},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().Logout("xxx").Return(nil)
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
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			authEndpoints := NewAuthEndpoints(mockAuthUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := authEndpoints.Logout(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestAuthEndpoints_LogoutAll(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		ExpectedErr          error
		Cookies              *http.Cookie
		SetupAuthUsecaseMock func(usecase *mockusecase.MockAuth)
	}{
		{
			Name:                 "Пользователь вышел, уже не имея сессии",
			ExpectedErr:          nil,
			Cookies:              nil,
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name:        "Пользователь вышел, изначально имея сессию",
			ExpectedErr: nil,
			Cookies:     &http.Cookie{Name: "session", Value: "xxx"},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
				usecase.EXPECT().LogoutAll(1).Return(nil)
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
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			authEndpoints := NewAuthEndpoints(mockAuthUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			req := httptest.NewRequest(http.MethodPost, "/auth/logoutAll", nil)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := authEndpoints.LogoutAll(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
