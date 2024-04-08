package http

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
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
			ExpectedErr:          &echo.HTTPError{Code: 400, Message: "Bad Request"},
			Cookies:              &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
		},
		{
			Name: "Ошибка при обновлении пароля",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.UpdatePassword{
					OldPassword: "AmaziNgPassw0rd!",
					NewPassword: "AmaziNgPassw0rd!",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "ошибка при обновлении пароля"},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().UpdatePassword(1, gomock.Any()).Return(entity.NewClientError("ошибка при обновлении пароля", entity.ErrInternal))
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
			ExpectedErr: &echo.HTTPError{Code: 400, Message: "неверный пароль"},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().UpdatePassword(1, gomock.Any()).Return(entity.NewClientError("неверный пароль", entity.ErrBadRequest))
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
			ExpectedErr:          &echo.HTTPError{Code: 401, Message: "отсутствует cookies с сессией"},
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
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "ошибка при создании сессии"},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().UpdatePassword(1, gomock.Any()).Return(nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
				usecase.EXPECT().CreateSession(1).Return("", entity.NewClientError("ошибка при создании сессии", entity.ErrInternal))
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
			userEndpoints := NewUserEndpoints(mockUserUsecase, mockAuthUsecase, mockStaticUseCase)
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
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {
				usecase.EXPECT().UploadAvatar(gomock.Any()).Return(1, nil)
			},
		},
		{
			Name: "Файл не прикреплён",
			Input: func() (io.Reader, string) {
				var buffer bytes.Buffer
				writer := multipart.NewWriter(&buffer)
				writer.Close()
				return &buffer, writer.FormDataContentType()
			},
			ExpectedErr:          &echo.HTTPError{Code: 400, Message: "файл не прикреплен"},
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
			ExpectedErr:          &echo.HTTPError{Code: 400, Message: "невалидное изображение"},
			Cookies:              &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {
				usecase.EXPECT().UploadAvatar(gomock.Any()).Return(0, entity.NewClientError("невалидное изображение", entity.ErrBadRequest))
			},
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
			ExpectedErr:            &echo.HTTPError{Code: 401, Message: "отсутствует cookies с сессией"},
			Cookies:                nil,
			SetupUserUsecaseMock:   func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock:   func(usecase *mockusecase.MockAuth) {},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {},
		},
		{
			Name: "Не удалось загрузить аватар",
			Input: func() (io.Reader, string) {
				var buffer bytes.Buffer
				writer := multipart.NewWriter(&buffer)
				part, _ := writer.CreateFormFile("avatar", "test.jpg")
				part.Write([]byte("test image data")) // здесь можно использовать реальные данные изображения
				writer.Close()
				return &buffer, writer.FormDataContentType()
			},
			ExpectedErr:          &echo.HTTPError{Code: 500, Message: "ошибка при загрузке аватара"},
			Cookies:              &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {
				usecase.EXPECT().UploadAvatar(gomock.Any()).Return(0, entity.NewClientError("ошибка при загрузке аватара", entity.ErrInternal))
			},
		},
		{
			Name: "Ошибка при обновлении аватара",
			Input: func() (io.Reader, string) {
				var buffer bytes.Buffer
				writer := multipart.NewWriter(&buffer)
				part, _ := writer.CreateFormFile("avatar", "test.jpg")
				part.Write([]byte("test image data")) // здесь можно использовать реальные данные изображения
				writer.Close()
				return &buffer, writer.FormDataContentType()
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "ошибка при обновлении аватара"},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().UpdateAvatar(1, 1).Return(entity.NewClientError("ошибка при обновлении аватара", entity.ErrInternal))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {
				usecase.EXPECT().UploadAvatar(gomock.Any()).Return(1, nil)
			},
		},
		{
			Name: "Ошибка при получении ID пользователя",
			Input: func() (io.Reader, string) {
				var buffer bytes.Buffer
				writer := multipart.NewWriter(&buffer)
				part, _ := writer.CreateFormFile("avatar", "test.jpg")
				part.Write([]byte("test image data")) // здесь можно использовать реальные данные изображения
				writer.Close()
				return &buffer, writer.FormDataContentType()
			},
			ExpectedErr:          &echo.HTTPError{Code: 401, Message: "необходима авторизация"},
			Cookies:              &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(0, entity.NewClientError("ошибка при получении ID пользователя", entity.ErrInternal))
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
			ExpectedErr:          &echo.HTTPError{Code: 400, Message: "Bad Request"},
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
			ExpectedErr: &echo.HTTPError{Code: 400, Message: "невалидные данные"},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().UpdateInfo(1, gomock.Any()).Return(entity.NewClientError("невалидные данные", entity.ErrBadRequest))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(1, nil)
			},
		},
		{
			Name: "Ошибка при обновлении информации",
			Input: func() io.Reader {
				body, _ := json.Marshal(dto.UserUpdate{
					Name:  "name",
					Email: "email",
				})
				return strings.NewReader(string(body))
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "ошибка при обновлении информации"},
			Cookies:     &http.Cookie{Name: "session", Value: "session"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().UpdateInfo(1, gomock.Any()).Return(entity.NewClientError("ошибка при обновлении информации", entity.ErrInternal))
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
			ExpectedErr:          &echo.HTTPError{Code: 401, Message: "отсутствует cookies с сессией"},
			Cookies:              nil,
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {},
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
			mockUserUsecase := mockusecase.NewMockUser(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			mockStaticUseCase := mockusecase.NewMockStatic(ctrl)
			userEndpoints := NewUserEndpoints(mockUserUsecase, mockAuthUsecase, mockStaticUseCase)
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
			Name:        "Ошибка при получении профиля",
			RequestID:   "1",
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "ошибка при получении профиля"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().GetUser(1).Return(nil, entity.NewClientError("ошибка при получении профиля", entity.ErrInternal))
			},
		},
		{
			Name:        "Пользователь не найден",
			RequestID:   "1",
			ExpectedErr: &echo.HTTPError{Code: 404, Message: "пользователь не найден"},
			SetupUserUsecaseMock: func(usecase *mockusecase.MockUser) {
				usecase.EXPECT().GetUser(1).Return(nil, entity.NewClientError("пользователь не найден", entity.ErrNotFound))
			},
		},
		{
			Name:                 "Неверный ID",
			RequestID:            "invalid",
			ExpectedErr:          &echo.HTTPError{Code: 400, Message: "неверный id"},
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
			userEndpoints := NewUserEndpoints(mockUserUsecase, mockAuthUsecase, mockStaticUseCase)
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
			Name:        "Нет cookies",
			ExpectedErr: &echo.HTTPError{Code: 401, Message: "необходима авторизация"},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("session").Return(0, entity.NewClientError("ошибка при получении ID", entity.ErrInternal))
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
			userEndpoints := NewUserEndpoints(mockUserUsecase, mockAuthUsecase, mockStaticUseCase)
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
