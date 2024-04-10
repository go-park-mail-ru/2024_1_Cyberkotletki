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

func TestCompilationEndpoints_GetCompilationTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                        string
		ExpectedErr                 error
		SetupCompilationUsecaseMock func(usecase *mockusecase.MockCompilation)
	}{
		{
			Name:        "Успех",
			ExpectedErr: nil,
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationTypes().Return(nil, nil)
			},
		},
		{
			Name:        "Не найдено",
			ExpectedErr: &echo.HTTPError{Code: 404, Message: "категория не найдена"},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationTypes().Return(nil,
					entity.NewClientError("категория не найдена", entity.ErrNotFound))
			},
		},
		{
			Name:        "Неизвестная ошибка",
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "Неизвестно"},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationTypes().Return(nil,
					entity.NewClientError("Неизвестно", entity.ErrInternal))
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
			mockCompilationUsecase := mockusecase.NewMockCompilation(ctrl)
			tc.SetupCompilationUsecaseMock(mockCompilationUsecase)
			compilationHandler := NewCompilationEndpoints(mockCompilationUsecase)
			req := httptest.NewRequest(http.MethodGet, "/compilation/types", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := compilationHandler.GetCompilationTypes(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestCompilationEndpoints_GetCompilationsByCompilationType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                        string
		CompilationTypeID           string
		ExpectedErr                 error
		SetupCompilationUsecaseMock func(usecase *mockusecase.MockCompilation)
	}{
		{
			Name:              "Успех",
			CompilationTypeID: "1",
			ExpectedErr:       nil,
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationsByCompilationType(1).Return(nil, nil)
			},
		},
		{
			Name:              "Не найдено",
			CompilationTypeID: "2",
			ExpectedErr:       &echo.HTTPError{Code: 404, Message: "Подборки не найдены"},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationsByCompilationType(2).Return(nil,
					entity.NewClientError("Подборки не найдены", entity.ErrNotFound))
			},
		},
		{
			Name:              "Неизвестная ошибка",
			CompilationTypeID: "3",
			ExpectedErr:       &echo.HTTPError{Code: 500, Message: "Неизвестная ошибка"},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationsByCompilationType(3).Return(nil,
					entity.NewClientError("Неизвестная ошибка", entity.ErrInternal))
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
			mockCompilationUsecase := mockusecase.NewMockCompilation(ctrl)
			tc.SetupCompilationUsecaseMock(mockCompilationUsecase)
			compilationHandler := NewCompilationEndpoints(mockCompilationUsecase)
			req := httptest.NewRequest(http.MethodGet, "/compilation/type/"+tc.CompilationTypeID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/compilation/type/:compilationType")
			c.SetParamNames("compilationType")
			c.SetParamValues(tc.CompilationTypeID)
			err := compilationHandler.GetCompilationsByCompilationType(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestCompilationEndpoints_GetCompilationContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                        string
		CompilationID               string
		Page                        string
		ExpectedErr                 error
		SetupCompilationUsecaseMock func(usecase *mockusecase.MockCompilation)
	}{
		{
			Name:          "Успешно",
			CompilationID: "1",
			Page:          "1",
			ExpectedErr:   nil,
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationContent(1, 1, 10).Return(nil, nil)
			},
		},
		{
			Name:          "Не найдено",
			CompilationID: "2",
			Page:          "1",
			ExpectedErr:   &echo.HTTPError{Code: 404, Message: "Контент не найден"},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationContent(2, 1, 10).Return(nil,
					entity.NewClientError("Контент не найден", entity.ErrNotFound))
			},
		},
		{
			Name:          "Неизвестная ошибка",
			CompilationID: "3",
			Page:          "1",
			ExpectedErr:   &echo.HTTPError{Code: 500, Message: "Неизвестная ошибка"},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationContent(3, 1, 10).Return(nil,
					entity.NewClientError("Неизвестная ошибка", entity.ErrInternal))
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
			mockCompilationUsecase := mockusecase.NewMockCompilation(ctrl)
			tc.SetupCompilationUsecaseMock(mockCompilationUsecase)
			compilationHandler := NewCompilationEndpoints(mockCompilationUsecase)
			req := httptest.NewRequest(http.MethodGet, "/compilation/"+tc.CompilationID+"/"+tc.Page, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/compilation/:id/:page")
			c.SetParamNames("id", "page")
			c.SetParamValues(tc.CompilationID, tc.Page)
			err := compilationHandler.GetCompilationContent(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
