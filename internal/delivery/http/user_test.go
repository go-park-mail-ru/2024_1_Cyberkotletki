package http

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserEndpoints_Register(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		Input                  func() io.Reader
		ExpectedErr            error
		SetupUserUsecaseMock   func(usecase *mockusecase.MockUser)
		SetupAuthUsecaseMock   func(usecase *mockusecase.MockAuth)
		SetupStaticUsecaseMock func(usecase *mockusecase.MockStatic)
	}{
		{
			Name: "Невалидный JSON",
			Input: func() io.Reader {
				return strings.NewReader("invalid")
			},
			ExpectedErr:            &echo.HTTPError{Code: 400, Message: "Bad Request"},
			SetupUserUsecaseMock:   func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock:   func(usecase *mockusecase.MockAuth) {},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {},
		},
		{
			Name: "Успешная регистрация",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.Register{
					Email:    "email@email.com",
					Password: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: nil,
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().Register(gomock.Any()).Return(1, nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().CreateSession(1).Return("session", nil)
			},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {},
		},
		{
			Name: "Ошибка при регистрации",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.Register{
					Email:    "email",
					Password: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "ошибка при регистрации"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().Register(gomock.Any()).Return(0, entity.NewClientError("ошибка при регистрации", entity.ErrInternal))
			},
			SetupAuthUsecaseMock:   func(usecase *mockusecase.MockAuth) {},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {},
		},
		{
			Name: "Пользователь с таким email уже существует",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.Register{
					Email:    "email",
					Password: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 409, Message: "пользователь с таким email уже существует"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().Register(gomock.Any()).Return(0, entity.NewClientError("пользователь с таким email уже существует", entity.ErrAlreadyExists))
			},
			SetupAuthUsecaseMock:   func(usecase *mockusecase.MockAuth) {},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {},
		},
		{
			Name: "Невалидные данные",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.Register{
					Email:    "email",
					Password: "purepass",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 400, Message: "невалидные данные"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().Register(gomock.Any()).Return(0, entity.NewClientError("невалидные данные", entity.ErrBadRequest))
			},
			SetupAuthUsecaseMock:   func(usecase *mockusecase.MockAuth) {},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {},
		},
		{
			Name: "Ошибка при создании сессии",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.Register{
					Email:    "email",
					Password: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "ошибка при создании сессии"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().Register(gomock.Any()).Return(1, nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().CreateSession(1).Return("", entity.NewClientError("ошибка при создании сессии", entity.ErrInternal))
			},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserUsecase := mockusecase.NewMockUser(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			mockStaticUseCase := mockusecase.NewMockStatic(ctrl)
			userEndpoints := NewUserEndpoints(mockUserUsecase, mockAuthUsecase, mockStaticUseCase)
			tc.SetupUserUsecaseMock(mockUserUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			tc.SetupStaticUsecaseMock(mockStaticUseCase)
			req := httptest.NewRequest(http.MethodPost, "/user/register", tc.Input())
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			cfg := config.Config{}
			cfg.Auth.SessionAliveTime = 1
			ctx.Set("params", cfg)
			err := userEndpoints.Register(ctx)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestUserEndpoints_Login(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		Input                func() io.Reader
		ExpectedErr          error
		SetupUserUsecaseMock func(usecase *mockusecase.MockUser)
		SetupAuthUsecaseMock func(usecase *mockusecase.MockAuth)
	}{
		{
			Name: "Успешная авторизация",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.Login{
					Login:    "email",
					Password: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: nil,
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().Login(gomock.Any()).Return(1, nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().CreateSession(1).Return("session", nil)
			},
		},
		{
			Name: "Невалидный JSON",
			Input: func() io.Reader {
				return strings.NewReader("invalid")
			},
			ExpectedErr:          &echo.HTTPError{Code: 400, Message: "Bad Request"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name: "Ошибка при авторизации",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.Login{
					Login:    "email",
					Password: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "ошибка при авторизации"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().Login(gomock.Any()).Return(0, entity.NewClientError("ошибка при авторизации", entity.ErrInternal))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name: "Пользователь не найден",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.Login{
					Login:    "email",
					Password: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 404, Message: "пользователь не найден"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().Login(gomock.Any()).Return(0, entity.NewClientError("пользователь не найден", entity.ErrNotFound))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name: "Неверный пароль",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.Login{
					Login:    "email",
					Password: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 403, Message: "неверный пароль"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().Login(gomock.Any()).Return(0, entity.NewClientError("неверный пароль", entity.ErrForbidden))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {},
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
			mockUserUsecase := mockusecase.NewMockUser(ctrl)
			mockStaticUseCase := mockusecase.NewMockStatic(ctrl)
			userEndpoints := NewUserEndpoints(mockUserUsecase, mockAuthUsecase, mockStaticUseCase)
			tc.SetupUserUsecaseMock(mockUserUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			req := httptest.NewRequest(http.MethodPost, "/user/login", tc.Input())
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			cfg := config.Config{}
			cfg.Auth.SessionAliveTime = 1
			ctx.Set("params", cfg)
			err := userEndpoints.Login(ctx)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
