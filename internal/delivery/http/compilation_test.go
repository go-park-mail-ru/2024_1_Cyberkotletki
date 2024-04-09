package http

/*import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompilationEndpoints_GetCompilation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                        string
		CompilationID               string
		ExpectedErr                 error
		SetupCompilationUsecaseMock func(usecase *mockusecase.MockCompilation)
	}{
		{
			Name:          "Успешное получение подборки",
			CompilationID: "1",
			ExpectedErr:   nil,
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilation(1).Return(nil, nil)
			},
		},
		{
			Name:          "Подборка не найдена",
			CompilationID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "подборка не найдена",
			},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilation(1).Return(nil, entity.NewClientError("подборка не найдена", entity.ErrNotFound))
			},
		},
		{
			Name:          "Неожиданная ошибка",
			CompilationID: "1",
			ExpectedErr:   &echo.HTTPError{Code: 500, Message: "internal"},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilation(1).Return(nil, entity.NewClientError("internal", entity.ErrInternal))
			},
		},
		{
			Name:                        "Невалидный айди",
			CompilationID:               "ogo!",
			ExpectedErr:                 &echo.HTTPError{Code: 400, Message: "невалидный id подборки"},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {},
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
			req := httptest.NewRequest(http.MethodGet, "/compilation/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/compilation/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.CompilationID)
			err := compilationHandler.GetCompilation(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestCompilationEndpoints_GetCompilationTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                        string
		ExpectedErr                 error
		SetupCompilationUsecaseMock func(usecase *mockusecase.MockCompilation)
	}{
		{
			Name:        "Успешное получение списка подборок",
			ExpectedErr: nil,
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationTypes().Return(nil, nil)
			},
		},
		{
			Name: "Подборки не найдены",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "подборки не найдены",
			},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationTypes().Return(nil, entity.NewClientError("подборки не найдены", entity.ErrNotFound))
			},
		},
		{
			Name:        "Неожиданная ошибка",
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "internal"},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationTypes().Return(nil, entity.NewClientError("internal", entity.ErrInternal))
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
			Name:              "Успешное получение списка подборок по типу",
			CompilationTypeID: "1",
			ExpectedErr:       nil,
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationsByCompilationType(1).Return(nil, nil)
			},
		},
		{
			Name:              "Подборки по данному типу не найдены",
			CompilationTypeID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "подборки по данному типу не найдены",
			},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationsByCompilationType(1).Return(nil, entity.NewClientError("подборки по данному типу не найдены", entity.ErrNotFound))
			},
		},
		{
			Name:              "Неожиданная ошибка",
			CompilationTypeID: "1",
			ExpectedErr:       &echo.HTTPError{Code: 500, Message: "internal"},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {
				usecase.EXPECT().GetCompilationsByCompilationType(1).Return(nil, entity.NewClientError("internal", entity.ErrInternal))
			},
		},
		{
			Name:                        "Невалидный айди типа подборки",
			CompilationTypeID:           "ogo!",
			ExpectedErr:                 &echo.HTTPError{Code: 400, Message: "невалидный id типа подборки"},
			SetupCompilationUsecaseMock: func(usecase *mockusecase.MockCompilation) {},
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
			req := httptest.NewRequest(http.MethodGet, "/compilation/type/", nil)
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
*/
