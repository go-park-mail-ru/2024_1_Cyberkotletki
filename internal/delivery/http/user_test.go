package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"mime/multipart"
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
			ExpectedErr:            &echo.HTTPError{Code: 400, Message: "Невалидный JSON"},
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
			Name: "Внутренняя ошибка сервера",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.Register{
					Email:    "email",
					Password: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "Внутренняя ошибка сервера", Internal: errors.New("123")},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().Register(gomock.Any()).Return(0, errors.New("123"))
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
			ExpectedErr: &echo.HTTPError{Code: 409, Message: "Пользователь с такой почтой уже существует"},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().Register(gomock.Any()).Return(0, usecase.ErrUserAlreadyExists)
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
			ExpectedErr: &echo.HTTPError{Code: 400, Message: "123"},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().Register(gomock.Any()).Return(0, usecase.UserIncorrectDataError{Err: errors.New("123")})
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
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "Внутренняя ошибка сервера", Internal: errors.New("123")},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().Register(gomock.Any()).Return(1, nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().CreateSession(1).Return("", errors.New("123"))
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
			sessionManager := utils.NewSessionManager(mockAuthUsecase, 1, false)
			userEndpoints := NewUserEndpoints(mockUserUsecase, mockAuthUsecase, mockStaticUseCase, sessionManager)
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
			ExpectedErr:          &echo.HTTPError{Code: 400, Message: "Невалидный JSON"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name: "Внутренняя ошибка сервера",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.Login{
					Login:    "email",
					Password: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "Внутренняя ошибка сервера", Internal: errors.New("123")},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().Login(gomock.Any()).Return(0, errors.New("123"))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name: "Внутренняя ошибка сервера при создании сессии",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.Login{
					Login:    "email",
					Password: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "Внутренняя ошибка сервера", Internal: errors.New("123")},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().Login(gomock.Any()).Return(0, nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().CreateSession(0).Return("", errors.New("123"))
			},
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
			ExpectedErr: &echo.HTTPError{Code: 404, Message: "Пользователь не найден"},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().Login(gomock.Any()).Return(0, usecase.ErrUserNotFound)
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
			ExpectedErr: &echo.HTTPError{Code: 403, Message: "123"},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().Login(gomock.Any()).Return(0, usecase.UserIncorrectDataError{Err: errors.New("123")})
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
			sessionManager := utils.NewSessionManager(mockAuthUsecase, 1, false)
			userEndpoints := NewUserEndpoints(mockUserUsecase, mockAuthUsecase, mockStaticUseCase, sessionManager)
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

func TestUserEndpoints_UpdatePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		Input                func() io.Reader
		ExpectedErr          error
		Cookies              *http.Cookie
		SetupUserUsecaseMock func(usecase *mockusecase.MockUser)
		SetupAuthUsecaseMock func(usecase *mockusecase.MockAuth)
	}{
		{
			Name: "Успешное обновление пароля",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.UpdatePassword{
					OldPassword: "AmaziNgPassw0rd!",
					NewPassword: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: nil,
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().UpdatePassword(1, gomock.Any()).Return(nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
				usecase.EXPECT().CreateSession(1).Return("session", nil)
			},
		},
		{
			Name: "Невалидный JSON",
			Input: func() io.Reader {
				return strings.NewReader("invalid")
			},
			ExpectedErr:          &echo.HTTPError{Code: 400, Message: "Невалидный JSON"},
			Cookies:              &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
		},
		{
			Name: "Внутренняя ошибка сервера",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.UpdatePassword{
					OldPassword: "AmaziNgPassw0rd!",
					NewPassword: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "Внутренняя ошибка сервера", Internal: errors.New("123")},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().UpdatePassword(1, gomock.Any()).Return(errors.New("123"))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
		},
		{
			Name: "Неверный пароль",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.UpdatePassword{
					OldPassword: "AmaziNgPassw0rd!",
					NewPassword: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 400, Message: "123"},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().UpdatePassword(1, gomock.Any()).Return(usecase.UserIncorrectDataError{Err: errors.New("123")})
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
		},
		{
			Name: "Пользователь не авторизован",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.UpdatePassword{
					OldPassword: "AmaziNgPassw0rd!",
					NewPassword: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr:          &echo.HTTPError{Code: 401, Message: "Не авторизован"},
			Cookies:              nil,
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name: "Не удалось создать сессию",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.UpdatePassword{
					OldPassword: "AmaziNgPassw0rd!",
					NewPassword: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "Внутренняя ошибка сервера", Internal: errors.New("123")},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().UpdatePassword(1, gomock.Any()).Return(nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
				usecase.EXPECT().CreateSession(1).Return("", errors.New("123"))
			},
		},
		{
			Name: "Пользователь не найден",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.UpdatePassword{
					OldPassword: "AmaziNgPassw0rd!",
					NewPassword: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 404, Message: "Пользователь не найден"},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().UpdatePassword(1, gomock.Any()).Return(usecase.ErrUserNotFound)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
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
			mockUserUsecase := mockusecase.NewMockUser(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			mockStaticUseCase := mockusecase.NewMockStatic(ctrl)
			sessionManager := utils.NewSessionManager(mockAuthUsecase, 1, false)
			userEndpoints := NewUserEndpoints(mockUserUsecase, mockAuthUsecase, mockStaticUseCase, sessionManager)
			tc.SetupUserUsecaseMock(mockUserUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			req := httptest.NewRequest(http.MethodPut, "/user/password", tc.Input())
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			cfg := config.Config{}
			cfg.Auth.SessionAliveTime = 1
			ctx.Set("params", cfg)
			ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			err := userEndpoints.UpdatePassword(ctx)
			require.Equal(t, tc.ExpectedErr, err)
		})

	}
}

func TestUserEndpoints_UploadAvatar(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		Input                  func() (io.Reader, string)
		ExpectedErr            error
		Cookies                *http.Cookie
		SetupUserUsecaseMock   func(usecase *mockusecase.MockUser)
		SetupAuthUsecaseMock   func(usecase *mockusecase.MockAuth)
		SetupStaticUsecaseMock func(usecase *mockusecase.MockStatic)
	}{
		{
			Name: "Успешная загрузка аватара",
			Input: func() (io.Reader, string) {
				var buffer bytes.Buffer
				writer := multipart.NewWriter(&buffer)
				part, _ := writer.CreateFormFile("avatar", "test.jpg")
				part.Write([]byte("bytes"))
				writer.Close()
				return &buffer, writer.FormDataContentType()
			},
			ExpectedErr: nil,
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().UpdateAvatar(1, gomock.Any()).Return(nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {},
		},
		{
			Name: "Файл не прикреплён",
			Input: func() (io.Reader, string) {
				var buffer bytes.Buffer
				writer := multipart.NewWriter(&buffer)
				writer.Close()
				return &buffer, writer.FormDataContentType()
			},
			ExpectedErr:          &echo.HTTPError{Code: 400, Message: "Файл не прикреплён"},
			Cookies:              &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {},
		},
		{
			Name: "Невалидное изображение",
			Input: func() (io.Reader, string) {
				var buffer bytes.Buffer
				writer := multipart.NewWriter(&buffer)
				part, _ := writer.CreateFormFile("avatar", "test.jpg")
				part.Write([]byte("test image data")) // здесь можно использовать реальные данные изображения
				writer.Close()
				return &buffer, writer.FormDataContentType()
			},
			ExpectedErr: &echo.HTTPError{Code: 400, Message: "123"},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().UpdateAvatar(1, gomock.Any()).Return(usecase.UserIncorrectDataError{Err: errors.New("123")})
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {},
		},
		{
			Name: "Пользователь не авторизован",
			Input: func() (io.Reader, string) {
				var buffer bytes.Buffer
				writer := multipart.NewWriter(&buffer)
				part, _ := writer.CreateFormFile("avatar", "test.jpg")
				part.Write([]byte("test image data")) // здесь можно использовать реальные данные изображения
				writer.Close()
				return &buffer, writer.FormDataContentType()
			},
			ExpectedErr:            &echo.HTTPError{Code: 401, Message: "Не авторизован"},
			Cookies:                nil,
			SetupUserUsecaseMock:   func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock:   func(usecase *mockusecase.MockAuth) {},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {},
		},
		{
			Name: "Внутренняя ошибка сервера",
			Input: func() (io.Reader, string) {
				var buffer bytes.Buffer
				writer := multipart.NewWriter(&buffer)
				part, _ := writer.CreateFormFile("avatar", "test.jpg")
				part.Write([]byte("test image data")) // здесь можно использовать реальные данные изображения
				writer.Close()
				return &buffer, writer.FormDataContentType()
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "Внутренняя ошибка сервера", Internal: errors.New("123")},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().UpdateAvatar(1, gomock.Any()).Return(errors.New("123"))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {},
		},
		{
			Name: "Пользователь не найден",
			Input: func() (io.Reader, string) {
				var buffer bytes.Buffer
				writer := multipart.NewWriter(&buffer)
				part, _ := writer.CreateFormFile("avatar", "test.jpg")
				part.Write([]byte("test image data")) // здесь можно использовать реальные данные изображения
				writer.Close()
				return &buffer, writer.FormDataContentType()
			},
			ExpectedErr: &echo.HTTPError{Code: 404, Message: "Пользователь не найден"},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().UpdateAvatar(1, gomock.Any()).Return(usecase.ErrUserNotFound)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {},
		},
		{
			Name: "Невалидный файл",
			Input: func() (io.Reader, string) {
				var buffer bytes.Buffer
				writer := multipart.NewWriter(&buffer)
				part, _ := writer.CreateFormFile("avatar", "test.jpg")
				part.Write([]byte("test image data")) // здесь можно использовать реальные данные изображения
				writer.Close()
				return &buffer, writer.FormDataContentType()
			},
			ExpectedErr: &echo.HTTPError{Code: 400, Message: "123"},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().UpdateAvatar(1, gomock.Any()).Return(usecase.UserIncorrectDataError{Err: errors.New("123")})
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
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
			sessionManager := utils.NewSessionManager(mockAuthUsecase, 1, false)
			userEndpoints := NewUserEndpoints(mockUserUsecase, mockAuthUsecase, mockStaticUseCase, sessionManager)
			tc.SetupUserUsecaseMock(mockUserUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			tc.SetupStaticUsecaseMock(mockStaticUseCase)
			body, contentType := tc.Input()
			req := httptest.NewRequest(http.MethodPut, "/user/avatar", body)
			req.Header.Set("Content-Type", contentType)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			cfg := config.Config{}
			cfg.Auth.SessionAliveTime = 1
			ctx.Set("params", cfg)
			err := userEndpoints.UploadAvatar(ctx)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestUserEndpoints_UpdateInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		Input                func() io.Reader
		ExpectedErr          error
		Cookies              *http.Cookie
		SetupUserUsecaseMock func(usecase *mockusecase.MockUser)
		SetupAuthUsecaseMock func(usecase *mockusecase.MockAuth)
	}{
		{
			Name: "Успешное обновление информации",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.UserUpdate{
					Name:  "name",
					Email: "email",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: nil,
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().UpdateInfo(1, gomock.Any()).Return(nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
		},
		{
			Name: "Невалидный JSON",
			Input: func() io.Reader {
				return strings.NewReader("invalid")
			},
			ExpectedErr:          &echo.HTTPError{Code: 400, Message: "Невалидный JSON"},
			Cookies:              &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
		},
		{
			Name: "Невалидные данные",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.UserUpdate{
					Name:  "name",
					Email: "email",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 400, Message: "123"},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().UpdateInfo(1, gomock.Any()).Return(usecase.UserIncorrectDataError{Err: errors.New("123")})
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
		},
		{
			Name: "Внутренняя ошибка сервера",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.UserUpdate{
					Name:  "name",
					Email: "email",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "Внутренняя ошибка сервера", Internal: errors.New("123")},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().UpdateInfo(1, gomock.Any()).Return(errors.New("123"))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
		},
		{
			Name: "Пользователь не авторизован",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.UserUpdate{
					Name:  "name",
					Email: "email",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 401, Message: "Не авторизован"},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "session",
			},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(0, utils.ErrUnauthorized)
			},
		},
		{
			Name: "Пользователь не найден",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.UserUpdate{
					Name:  "name",
					Email: "email",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 404, Message: "Пользователь не найден"},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "session",
			},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().UpdateInfo(1, gomock.Any()).Return(usecase.ErrUserNotFound)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
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
			mockUserUsecase := mockusecase.NewMockUser(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			mockStaticUseCase := mockusecase.NewMockStatic(ctrl)
			sessionManager := utils.NewSessionManager(mockAuthUsecase, 1, false)
			userEndpoints := NewUserEndpoints(mockUserUsecase, mockAuthUsecase, mockStaticUseCase, sessionManager)
			tc.SetupUserUsecaseMock(mockUserUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			req := httptest.NewRequest(http.MethodPut, "/user/info", tc.Input())
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			cfg := config.Config{}
			cfg.Auth.SessionAliveTime = 1
			ctx.Set("params", cfg)
			ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			err := userEndpoints.UpdateInfo(ctx)
			require.Equal(t, tc.ExpectedErr, err)
		})

	}
}

func TestUserEndpoints_GetProfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		RequestID            string
		ExpectedErr          error
		ExpectedOutput       *dto.UserProfile
		SetupUserUsecaseMock func(usecase *mockusecase.MockUser)
	}{
		{
			Name:        "Успешное получение профиля",
			RequestID:   "1",
			ExpectedErr: nil,
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().GetUser(1).Return(&dto.UserProfile{
					ID:     1,
					Name:   "name",
					Email:  "email",
					Rating: 15,
					Avatar: "avatar",
				}, nil)
			},
		},
		{
			Name:        "Внутренняя ошибка сервера",
			RequestID:   "1",
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "Внутренняя ошибка сервера", Internal: errors.New("123")},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().GetUser(1).Return(nil, errors.New("123"))
			},
		},
		{
			Name:        "Пользователь не найден",
			RequestID:   "1",
			ExpectedErr: &echo.HTTPError{Code: 404, Message: "Пользователь не найден"},
			SetupUserUsecaseMock: func(uc *mockusecase.MockUser) {
				uc.EXPECT().GetUser(1).Return(nil, usecase.ErrUserNotFound)
			},
		},
		{
			Name:                 "Неверный ID",
			RequestID:            "invalid",
			ExpectedErr:          &echo.HTTPError{Code: 400, Message: "Неверный id"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {},
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
			sessionManager := utils.NewSessionManager(mockAuthUsecase, 1, false)
			userEndpoints := NewUserEndpoints(mockUserUsecase, mockAuthUsecase, mockStaticUseCase, sessionManager)
			tc.SetupUserUsecaseMock(mockUserUsecase)
			req := httptest.NewRequest(http.MethodGet, "/user/profile", nil)
			req.AddCookie(&http.Cookie{Name: "session", Value: "session"})
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			cfg := config.Config{}
			cfg.Auth.SessionAliveTime = 1
			ctx.Set("params", cfg)
			ctx.QueryParams().Set("id", tc.RequestID)
			err := userEndpoints.GetProfile(ctx)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestUserEndpoints_GetMyID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		ExpectedErr          error
		ExpectedOutput       struct{ ID int }
		SetupAuthUsecaseMock func(usecase *mockusecase.MockAuth)
	}{
		{
			Name:        "Успешное получение ID",
			ExpectedErr: nil,
			ExpectedOutput: struct{ ID int }{
				ID: 1,
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
		},
		{
			Name:        "Не авторизован",
			ExpectedErr: &echo.HTTPError{Code: 401, Message: "Не авторизован"},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(0, utils.ErrUnauthorized)
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
			mockUserUsecase := mockusecase.NewMockUser(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			mockStaticUseCase := mockusecase.NewMockStatic(ctrl)
			sessionManager := utils.NewSessionManager(mockAuthUsecase, 1, false)
			userEndpoints := NewUserEndpoints(mockUserUsecase, mockAuthUsecase, mockStaticUseCase, sessionManager)
			req := httptest.NewRequest(http.MethodGet, "/user/id", nil)
			req.AddCookie(&http.Cookie{Name: "session", Value: "session"})
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			cfg := config.Config{}
			cfg.Auth.SessionAliveTime = 1
			ctx.Set("params", cfg)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			err := userEndpoints.GetMyID(ctx)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
