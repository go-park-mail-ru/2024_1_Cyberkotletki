package http

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearchEndpoints_Search(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		Input                  func() string
		ExpectedErr            error
		SetupSearchUsecaseMock func(usecase *mockusecase.MockSearch)
	}{
		{
			Name: "Пустой запрос",
			Input: func() string {
				return ""
			},
			ExpectedErr:            &echo.HTTPError{Code: 400, Message: "Пустой запрос"},
			SetupSearchUsecaseMock: func(usecase *mockusecase.MockSearch) {},
		},
		{
			Name: "Слишком длинный запрос",
			Input: func() string {
				return strings.Repeat("a", 101)
			},
			ExpectedErr:            &echo.HTTPError{Code: 400, Message: "Слишком длинный запрос"},
			SetupSearchUsecaseMock: func(usecase *mockusecase.MockSearch) {},
		},
		{
			Name: "Успешный поиск",
			Input: func() string {
				return "hello"
			},
			ExpectedErr: nil,
			SetupSearchUsecaseMock: func(usecase *mockusecase.MockSearch) {
				usecase.EXPECT().Search("hello").Return(&dto.SearchResult{}, nil)
			},
		},
		{
			Name: "Внутренняя ошибка сервера",
			Input: func() string {
				return "hello"
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "Внутренняя ошибка сервера", Internal: errors.New("123")},
			SetupSearchUsecaseMock: func(usecase *mockusecase.MockSearch) {
				usecase.EXPECT().Search("hello").Return(nil, errors.New("123"))
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
			mockSearchUsecase := mockusecase.NewMockSearch(ctrl)
			searchEndpoints := NewSearchEndpoints(mockSearchUsecase)
			tc.SetupSearchUsecaseMock(mockSearchUsecase)
			req := httptest.NewRequest(http.MethodGet, "/search", nil)
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)
			ctx.QueryParams().Set("query", tc.Input())
			ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			err := searchEndpoints.Search(ctx)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
