package http

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/http/utils"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthEndpoints_NewAuthEndpoints(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAuth := mockusecase.NewMockAuth(ctrl)
	h := NewAuthEndpoints(mockAuth)

	if h.authUC != mockAuth {
		t.Errorf("NewAuthEndpoints() = %v, want %v", h.authUC, mockAuth)
	}
}

func TestAuthEndpoints_Register(t *testing.T) {
	t.Parallel()
	e := echo.New()

	testCases := []struct {
		Name        string
		Input       dto.Register
		ExpectedErr func(echo.Context) error
		SetupMock   func(*mockusecase.MockAuth)
	}{
		{
			Name: "Успешная регистрация",
			Input: dto.Register{
				Email:    "test@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: nil,
			SetupMock: func(mockAuth *mockusecase.MockAuth) {
				mockAuth.EXPECT().Register(gomock.Eq(&dto.Register{
					Email:    "test@example.com",
					Password: "AmazingPassword1!",
				})).Return("session", nil)
			},
		},
		{
			Name: "Пользователь уже существует",
			Input: dto.Register{
				Email:    "111@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: func(ctx echo.Context) error {
				return utils.NewError(ctx, http.StatusConflict, entity.ErrAlreadyExists)
			},
			SetupMock: func(mockAuth *mockusecase.MockAuth) {
				mockAuth.EXPECT().Register(gomock.Eq(&dto.Register{
					Email:    "111@example.com",
					Password: "AmazingPassword1!",
				})).Return("", entity.ErrAlreadyExists)
			},
		},
		{
			Name: "Какая-то внутренняя ошибка сервера",
			Input: dto.Register{
				Email:    "111@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: func(ctx echo.Context) error {
				return utils.NewError(ctx, http.StatusInternalServerError, entity.ErrRedis)
			},
			SetupMock: func(mockAuth *mockusecase.MockAuth) {
				mockAuth.EXPECT().Register(gomock.Eq(&dto.Register{
					Email:    "111@example.com",
					Password: "AmazingPassword1!",
				})).Return("", entity.ErrRedis)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockAuth := mockusecase.NewMockAuth(ctrl)
			h := NewAuthEndpoints(mockAuth)
			tc.SetupMock(mockAuth)

			// запрос
			reqBody, err := json.Marshal(tc.Input)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// конфиг
			c.Set("params", config.Config{})
			// сам метод
			err = h.Register(c)
			// проверка
			if tc.ExpectedErr != nil {
				require.ErrorContains(t, err, tc.ExpectedErr(c).Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAuthEndpoints_Login(t *testing.T) {
	t.Parallel()
	e := echo.New()

	testCases := []struct {
		Name        string
		Input       dto.Login
		ExpectedErr func(echo.Context) error
		SetupMock   func(*mockusecase.MockAuth)
	}{
		{
			Name: "Успешный вход",
			Input: dto.Login{
				Login:    "test@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: nil,
			SetupMock: func(mockAuth *mockusecase.MockAuth) {
				mockAuth.EXPECT().Login(gomock.Eq(&dto.Login{
					Login:    "test@example.com",
					Password: "AmazingPassword1!",
				})).Return("session", nil)
			},
		},
		{
			Name: "Неверный пароль",
			Input: dto.Login{
				Login:    "test@example.com",
				Password: "WrongPassword!",
			},
			ExpectedErr: func(ctx echo.Context) error {
				return utils.NewError(ctx, http.StatusForbidden, entity.ErrForbidden)
			},
			SetupMock: func(mockAuth *mockusecase.MockAuth) {
				mockAuth.EXPECT().Login(gomock.Eq(&dto.Login{
					Login:    "test@example.com",
					Password: "WrongPassword!",
				})).Return("", entity.ErrForbidden)
			},
		},
		{
			Name: "Пользователь не найден",
			Input: dto.Login{
				Login:    "notfound@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: func(ctx echo.Context) error {
				return utils.NewError(ctx, http.StatusNotFound, entity.ErrNotFound)
			},
			SetupMock: func(mockAuth *mockusecase.MockAuth) {
				mockAuth.EXPECT().Login(gomock.Eq(&dto.Login{
					Login:    "notfound@example.com",
					Password: "AmazingPassword1!",
				})).Return("", entity.ErrNotFound)
			},
		},
		{
			Name: "Какая-то внутренняя ошибка сервера",
			Input: dto.Login{
				Login:    "notfound@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: func(ctx echo.Context) error {
				return utils.NewError(ctx, http.StatusInternalServerError, entity.ErrPSQL)
			},
			SetupMock: func(mockAuth *mockusecase.MockAuth) {
				mockAuth.EXPECT().Login(gomock.Eq(&dto.Login{
					Login:    "notfound@example.com",
					Password: "AmazingPassword1!",
				})).Return("", entity.ErrPSQL)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockAuth := mockusecase.NewMockAuth(ctrl)
			h := NewAuthEndpoints(mockAuth)
			tc.SetupMock(mockAuth)
			// запрос
			reqBody, err := json.Marshal(tc.Input)
			require.NoError(t, err)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// конфиг
			c.Set("params", config.Config{})
			// сам метод
			err = h.Login(c)
			// проверка
			if tc.ExpectedErr != nil {
				require.ErrorContains(t, err, tc.ExpectedErr(c).Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAuthEndpoints_IsAuth(t *testing.T) {
	t.Parallel()
	e := echo.New()

	testCases := []struct {
		Name        string
		Cookie      *http.Cookie
		ExpectedErr func(echo.Context) error
		SetupMock   func(*mockusecase.MockAuth)
	}{
		{
			Name:   "Не авторизован (нет cookie)",
			Cookie: nil,
			ExpectedErr: func(ctx echo.Context) error {
				return utils.NewError(ctx, http.StatusUnauthorized, entity.NewClientError("не авторизован"))
			},
			SetupMock: func(*mockusecase.MockAuth) {}, //ошибка до мока - мок не нужен
		},
		{
			Name:   "Не авторизован (неверное значение cookie)",
			Cookie: &http.Cookie{Name: "session", Value: "invalid"},
			ExpectedErr: func(ctx echo.Context) error {
				return utils.NewError(ctx, http.StatusUnauthorized, entity.NewClientError("не авторизован"))
			},
			SetupMock: func(mockAuth *mockusecase.MockAuth) {
				mockAuth.EXPECT().IsAuth("invalid").Return(false, nil)
			},
		},
		{
			Name:        "Успех",
			Cookie:      &http.Cookie{Name: "session", Value: "valid"},
			ExpectedErr: nil,
			SetupMock: func(mockAuth *mockusecase.MockAuth) {
				mockAuth.EXPECT().IsAuth("valid").Return(true, nil)
			},
		},
		{
			Name:   "Какая-то внутренняя ошибка сервера",
			Cookie: &http.Cookie{Name: "session", Value: "valid"},
			ExpectedErr: func(ctx echo.Context) error {
				return utils.NewError(ctx, http.StatusInternalServerError, entity.ErrRedis)
			},
			SetupMock: func(mockAuth *mockusecase.MockAuth) {
				mockAuth.EXPECT().IsAuth("valid").Return(false, entity.ErrRedis)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockAuth := mockusecase.NewMockAuth(ctrl)
			h := NewAuthEndpoints(mockAuth)
			tc.SetupMock(mockAuth)
			// запрос
			req := httptest.NewRequest(http.MethodGet, "/isAuth", nil)
			if tc.Cookie != nil {
				req.AddCookie(tc.Cookie)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// конфиг
			c.Set("params", config.Config{})
			// сам метод
			err := h.IsAuth(c)
			// проверка
			if tc.ExpectedErr != nil {
				require.ErrorContains(t, err, tc.ExpectedErr(c).Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAuthEndpoints_Logout(t *testing.T) {
	t.Parallel()
	e := echo.New()

	testCases := []struct {
		Name        string
		Cookie      *http.Cookie
		ExpectedErr error
		SetupMock   func(*mockusecase.MockAuth)
	}{
		{
			Name:        "Не авторизован (нет cookie)",
			Cookie:      nil, // нет кук
			ExpectedErr: nil,
			SetupMock:   func(*mockusecase.MockAuth) {},
		},
		{
			Name:        "Успешный выход",
			Cookie:      &http.Cookie{Name: "session", Value: "valid"},
			ExpectedErr: nil,
			SetupMock: func(mockAuth *mockusecase.MockAuth) {
				mockAuth.EXPECT().Logout("valid")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockAuth := mockusecase.NewMockAuth(ctrl)
			h := NewAuthEndpoints(mockAuth)
			tc.SetupMock(mockAuth)
			// запрос
			req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
			if tc.Cookie != nil {
				req.AddCookie(tc.Cookie)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// конфиг
			c.Set("params", config.Config{})
			// сам метод
			err := h.Logout(c)
			// проверка
			if tc.ExpectedErr != nil {
				require.ErrorContains(t, err, tc.ExpectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
