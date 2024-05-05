package http

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStaticEndpoints_GetStatic(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		StaticId               string
		ExpectedErr            error
		SetupStaticUsecaseMock func(usecase *mockusecase.MockStatic)
	}{
		{
			Name:        "Получение URL",
			StaticId:    "1",
			ExpectedErr: nil,
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {
				usecase.EXPECT().GetStatic(1).Return("static_url", nil)
			},
		},
		{
			Name:                   "Невалидный id",
			StaticId:               "invalid_id",
			ExpectedErr:            &echo.HTTPError{Code: 400, Message: "Невалидный id статики"},
			SetupStaticUsecaseMock: func(uc *mockusecase.MockStatic) {},
		},
		{
			Name:        "Статики с таким id нет",
			StaticId:    "2",
			ExpectedErr: &echo.HTTPError{Code: 404, Message: "Статика не найдена"},
			SetupStaticUsecaseMock: func(uc *mockusecase.MockStatic) {
				uc.EXPECT().GetStatic(2).Return("", usecase.ErrStaticNotFound)
			},
		},
		{
			Name:        "Внутренняя ошибка сервера",
			StaticId:    "3",
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "Внутренняя ошибка сервера", Internal: errors.New("123")},
			SetupStaticUsecaseMock: func(uc *mockusecase.MockStatic) {
				uc.EXPECT().GetStatic(3).Return("", errors.New("123"))
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
			mockStaticUsecase := mockusecase.NewMockStatic(ctrl)
			staticEndpoints := NewStaticEndpoints(mockStaticUsecase)
			tc.SetupStaticUsecaseMock(mockStaticUsecase)
			req := httptest.NewRequest(http.MethodGet, "/static/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/static/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.StaticId)
			err := staticEndpoints.GetStaticURL(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
