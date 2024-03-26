package http

import (
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

func TestNewCollectionsEndpoints(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCollections := mockusecase.NewMockCollections(ctrl)
	h := NewCollectionsEndpoints(mockCollections)

	if h.useCase != mockCollections {
		t.Errorf("NewCollectionsEndpoints() = %v, want %v", h.useCase, mockCollections)
	}
}

func TestCollectionsEndpoints_GetGenres(t *testing.T) {
	t.Parallel()
	e := echo.New()

	testCases := []struct {
		Name        string
		ExpectedErr error
		SetupMock   func(*mockusecase.MockCollections)
	}{
		{
			Name:        "Успех",
			ExpectedErr: nil,
			SetupMock: func(mockCollections *mockusecase.MockCollections) {
				mockCollections.EXPECT().GetGenres().Return(&dto.Genres{
					Genres: []string{"action", "drama", "comedian"},
				}, nil)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockCollections := mockusecase.NewMockCollections(ctrl)
			h := NewCollectionsEndpoints(mockCollections)
			tc.SetupMock(mockCollections)

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

func TestCollectionsEndpoints_GetCompilationByGenre(t *testing.T) {
	t.Parallel()
	e := echo.New()

	testCases := []struct {
		Name        string
		Genre       string
		ExpectedErr error
		SetupMock   func(*mockusecase.MockCollections)
	}{
		{
			Name:        "Успех",
			Genre:       "action",
			ExpectedErr: nil,
			SetupMock: func(mockCollections *mockusecase.MockCollections) {
				mockCollections.EXPECT().GetCompilation(gomock.Eq("action")).Return(&dto.Compilation{
					Genre:              "action",
					ContentIdentifiers: []int{1, 2, 3},
				}, nil)
			},
		},
		{
			Name:        "Несуществующий жанр",
			Genre:       "lolkek",
			ExpectedErr: entity.NewClientError("Такого жанра не существует", entity.ErrNotFound),
			SetupMock: func(mockCollections *mockusecase.MockCollections) {
				mockCollections.EXPECT().GetCompilation(gomock.Eq("lolkek")).Return(
					nil,
					entity.NewClientError("Такого жанра не существует", entity.ErrNotFound),
				)
			},
		},
		{
			Name:        "Не указан жанр",
			Genre:       "",
			ExpectedErr: entity.NewClientError("Не указан жанр", entity.ErrBadRequest),
			SetupMock:   func(mockCollections *mockusecase.MockCollections) {},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockCollections := mockusecase.NewMockCollections(ctrl)
			h := NewCollectionsEndpoints(mockCollections)
			tc.SetupMock(mockCollections)

			req := httptest.NewRequest(http.MethodGet, "/collections/compilation?genre="+tc.Genre, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := h.GetCompilationByGenre(c)
			if tc.ExpectedErr != nil {
				require.ErrorContains(t, err, tc.ExpectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
