package http

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mock_usecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewCollectionsEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollections := mock_usecase.NewMockCollections(ctrl)

	h := NewCollectionsEndpoints(mockCollections)

	if h.useCase != mockCollections {
		t.Errorf("NewCollectionsEndpoints() = %v, want %v", h.useCase, mockCollections)
	}
}

func TestCollectionsEndpoints_GetGenres(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollections := mock_usecase.NewMockCollections(ctrl)
	h := NewCollectionsEndpoints(mockCollections)

	e := echo.New()

	testCases := []struct {
		Name        string
		ExpectedErr error
		SetupMock   func()
	}{
		{
			Name:        "Успех",
			ExpectedErr: nil,
			SetupMock: func() {
				mockCollections.EXPECT().GetGenres().Return(&dto.Genres{
					Genres: []string{"action", "drama", "comedian"},
				}, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.SetupMock()
			req := httptest.NewRequest(http.MethodGet, "/collections/genres", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := h.GetGenres(c)
			if tc.ExpectedErr != nil {
				require.ErrorContains(t, err, tc.ExpectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
