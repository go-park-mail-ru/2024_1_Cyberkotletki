package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/echoutil"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestContentEndpoints_NewContentEndpoints(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockContent := mockusecase.NewMockContent(ctrl)
	h := NewContentEndpoints(mockContent)

	if h.useCase != mockContent {
		t.Errorf("NewContentEndpoints() = %v, want %v", h.useCase, mockContent)
	}
}

func TestContentEndpoints_GetContentPreview(t *testing.T) {
	t.Parallel()
	e := echo.New()

	testCases := []struct {
		Name        string
		ID          string
		ExpectedErr error
		SetupMock   func(*mockusecase.MockContent)
	}{
		{
			Name:        "Успех",
			ID:          "1",
			ExpectedErr: nil,
			SetupMock: func(mockContent *mockusecase.MockContent) {
				mockContent.EXPECT().GetContentPreviewCard(1).Return(&dto.PreviewContentCard{}, nil)
			},
		},
		{
			Name:        "Ошибка преобразования ID",
			ID:          "-",
			ExpectedErr: echoutil.NewError(nil, http.StatusBadRequest, strconv.ErrSyntax),
			SetupMock:   func(*mockusecase.MockContent) {},
		},
		{
			Name:        "Контент не найден",
			ID:          "2",
			ExpectedErr: echoutil.NewError(nil, http.StatusNotFound, entity.ErrNotFound),
			SetupMock: func(mockContent *mockusecase.MockContent) {
				mockContent.EXPECT().GetContentPreviewCard(2).Return(nil, entity.ErrNotFound)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockContent := mockusecase.NewMockContent(ctrl)
			h := NewContentEndpoints(mockContent)
			tc.SetupMock(mockContent)

			req := httptest.NewRequest(http.MethodGet, "/content?id="+tc.ID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := h.GetContentPreview(c)
			if tc.ExpectedErr != nil {
				require.ErrorContains(t, err, tc.ExpectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
