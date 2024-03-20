package http

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/DTO"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/echoutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthEndpoints_NewAuthEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuth := mockusecase.NewMockAuth(ctrl)

	h := NewAuthEndpoints(mockAuth)

	if h.useCase != mockAuth {
		t.Errorf("NewAuthEndpoints() = %v, want %v", h.useCase, mockAuth)
	}
}

func TestAuthEndpoints_Register(t *testing.T) {
	// управляет ж циклом мока
	ctrl := gomock.NewController(t)
	// чистка мока
	defer ctrl.Finish()
	// мок
	mockAuth := mockusecase.NewMockAuth(ctrl)
	// экземпляр для тестов
	h := NewAuthEndpoints(mockAuth)

	e := echo.New()

	testCases := []struct {
		Name        string
		Input       DTO.Register
		ExpectedErr error
		SetupMock   func()
	}{
		{
			Name: "Успешная регистрация",
			Input: DTO.Register{
				Email:    "test@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: nil,
			SetupMock: func() {
				mockAuth.EXPECT().Register(gomock.Eq(DTO.Register{
					Email:    "test@example.com",
					Password: "AmazingPassword1!",
				})).Return("session", nil)
			},
		},
		{
			Name: "Пароль уже существует",
			Input: DTO.Register{
				Email:    "111@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: echoutil.NewError(nil, http.StatusBadRequest, entity.ErrAlreadyExists),
			SetupMock: func() {
				mockAuth.EXPECT().Register(gomock.Eq(DTO.Register{
					Email:    "111@example.com",
					Password: "AmazingPassword1!",
				})).Return("", entity.ErrAlreadyExists)
			},
		},
		{
			Name: "Почта уже существует",
			Input: DTO.Register{
				Email:    "test@example.com",
				Password: "UniqPassword1!",
			},
			ExpectedErr: echoutil.NewError(nil, http.StatusBadRequest, entity.ErrAlreadyExists),
			SetupMock: func() {
				mockAuth.EXPECT().Register(gomock.Eq(DTO.Register{
					Email:    "test@example.com",
					Password: "UniqPassword1!",
				})).Return("", entity.ErrAlreadyExists)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.SetupMock()
			// запрос
			reqBody, _ := json.Marshal(tc.Input)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// конфиг
			c.Set("params", config.Config{})
			// сам метод
			err := h.Register(c)
			// проверка
			if tc.ExpectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.ExpectedErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, rec.Code)
			}
		})
	}
}

func TestAuthEndpoints_Login(t *testing.T) {
	// управляет ж циклом мока
	ctrl := gomock.NewController(t)
	// чистка мока
	defer ctrl.Finish()
	// мок
	mockAuth := mockusecase.NewMockAuth(ctrl)
	// экземпляр для тестов
	h := NewAuthEndpoints(mockAuth)

	e := echo.New()

	testCases := []struct {
		Name        string
		Input       DTO.Login
		ExpectedErr error
		SetupMock   func()
	}{
		{
			Name: "Успешный вход",
			Input: DTO.Login{
				Login:    "test@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: nil,
			SetupMock: func() {
				mockAuth.EXPECT().Login(gomock.Eq(DTO.Login{
					Login:    "test@example.com",
					Password: "AmazingPassword1!",
				})).Return("session", nil)
			},
		},
		{
			Name: "Неверный пароль",
			Input: DTO.Login{
				Login:    "test@example.com",
				Password: "WrongPassword!",
			},
			ExpectedErr: echoutil.NewError(nil, http.StatusBadRequest, entity.ErrBadRequest),
			SetupMock: func() {
				mockAuth.EXPECT().Login(gomock.Eq(DTO.Login{
					Login:    "test@example.com",
					Password: "WrongPassword!",
				})).Return("", entity.ErrBadRequest)
			},
		},
		{
			Name: "Пользователь не найден",
			Input: DTO.Login{
				Login:    "notfound@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: echoutil.NewError(nil, http.StatusNotFound, entity.ErrNotFound),
			SetupMock: func() {
				mockAuth.EXPECT().Login(gomock.Eq(DTO.Login{
					Login:    "notfound@example.com",
					Password: "AmazingPassword1!",
				})).Return("", entity.ErrNotFound)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.SetupMock()
			// запрос
			reqBody, _ := json.Marshal(tc.Input)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			// конфиг
			c.Set("params", config.Config{})
			// сам метод
			err := h.Login(c)
			// проверка
			if tc.ExpectedErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.ExpectedErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, rec.Code)
			}
		})
	}
}

func TestAuthEndpoints_IsAuth(t *testing.T) {
	// управляет ж циклом мока
	ctrl := gomock.NewController(t)
	// чистка мока
	defer ctrl.Finish()
	// мок
	mockAuth := mockusecase.NewMockAuth(ctrl)
	// экземпляр для тестов
	h := NewAuthEndpoints(mockAuth)

	e := echo.New()

	testCases := []struct {
		Name        string
		Cookie      *http.Cookie
		ExpectedErr error
		SetupMock   func()
	}{
		{
			Name:        "Не авторизован (нет cookie)",
			Cookie:      nil,
			ExpectedErr: echoutil.NewError(nil, http.StatusUnauthorized, entity.NewClientError("не авторизован")),
			SetupMock:   func() {}, //ошибка до мока - мок не нужен
		},
		{
			Name:        "Не авторизован (неверное значение cookie)",
			Cookie:      &http.Cookie{Name: "session", Value: "invalid"},
			ExpectedErr: echoutil.NewError(nil, http.StatusUnauthorized, entity.NewClientError("не авторизован")),
			SetupMock: func() {
				mockAuth.EXPECT().IsAuth("invalid").Return(false)
			},
		},
		{
			Name:        "Успех",
			Cookie:      &http.Cookie{Name: "session", Value: "valid"},
			ExpectedErr: nil,
			SetupMock: func() {
				mockAuth.EXPECT().IsAuth("valid").Return(true)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.SetupMock()
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
				require.Error(t, err)
				require.Equal(t, tc.ExpectedErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, rec.Code)
			}
		})
	}
}

func TestAuthEndpoints_Logout(t *testing.T) {
	// управляет ж циклом мока
	ctrl := gomock.NewController(t)
	// чистка мока
	defer ctrl.Finish()
	// мок
	mockAuth := mockusecase.NewMockAuth(ctrl)
	// экземпляр для тестов
	h := NewAuthEndpoints(mockAuth)

	e := echo.New()

	testCases := []struct {
		Name        string
		Cookie      *http.Cookie
		ExpectedErr error
		SetupMock   func()
	}{
		{
			Name:        "Не авторизован (нет cookie)",
			Cookie:      nil, // нет кук
			ExpectedErr: nil,
			SetupMock:   func() {},
		},
		{
			Name:        "Успешный выход",
			Cookie:      &http.Cookie{Name: "session", Value: "valid"},
			ExpectedErr: nil,
			SetupMock: func() {
				mockAuth.EXPECT().Logout("valid")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.SetupMock()
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
				require.Error(t, err)
				require.Equal(t, tc.ExpectedErr, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, http.StatusOK, rec.Code)
			}
		})
	}
}
