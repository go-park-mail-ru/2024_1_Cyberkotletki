package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestStaticEndpoints_GetStaticUrl(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		StaticId               int
		ExpectedErr            error
		SetupStaticUsecaseMock func(usecase *mockusecase.MockStatic)
	}{
		{
			Name:        "Получение URL",
			StaticId:    1,
			ExpectedErr: nil,
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {
				usecase.EXPECT().GetStaticURL(1).Return("static_url", nil)
			},
		},
		{
			Name:        "Статики с таким id нет",
			StaticId:    2,
			ExpectedErr: &echo.HTTPError{Code: 404, Message: "файл не найден"},
			SetupStaticUsecaseMock: func(usecase *mockusecase.MockStatic) {
				usecase.EXPECT().GetStaticURL(2).Return("", entity.NewClientError("файл не найден", entity.ErrNotFound))
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
			req := httptest.NewRequest(http.MethodGet, "/static/1", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/static/:id")
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(tc.StaticId))
			err := staticEndpoints.GetStaticURL(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
