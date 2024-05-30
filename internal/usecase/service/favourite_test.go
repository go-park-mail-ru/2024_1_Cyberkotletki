package service

import (
	//"database/sql"
	"errors"
	mock_usecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"testing"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"

	//"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestFavouriteService_DeleteFavourite(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		UserID                 int
		ContentID              int
		ExpectedErr            error
		SetupFavouriteRepoMock func(repo *mockrepo.MockFavourite)
	}{
		{
			Name:        "Успех",
			UserID:      1,
			ContentID:   1,
			ExpectedErr: nil,
			SetupFavouriteRepoMock: func(repo *mockrepo.MockFavourite) {
				repo.EXPECT().DeleteFavourite(1, 1).Return(nil).AnyTimes()
			},
		},
		{
			Name:        "Ошибка при удалении",
			UserID:      1,
			ContentID:   1,
			ExpectedErr: entity.UsecaseWrap(errors.New("database error"), errors.New("ошибка при удалении из избранного в FavouriteService")),
			SetupFavouriteRepoMock: func(repo *mockrepo.MockFavourite) {
				repo.EXPECT().DeleteFavourite(1, 1).Return(errors.New("database error"))
			},
		},
		{
			Name:        "Избранное не найдено",
			UserID:      1,
			ContentID:   1,
			ExpectedErr: repository.ErrFavouriteNotFound,
			SetupFavouriteRepoMock: func(repo *mockrepo.MockFavourite) {
				repo.EXPECT().DeleteFavourite(1, 1).Return(repository.ErrFavouriteNotFound)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockFavouriteRepo := mockrepo.NewMockFavourite(ctrl)
			mockContentUC := mock_usecase.NewMockContent(ctrl)
			favService := NewFavouriteService(mockFavouriteRepo, mockContentUC)
			tc.SetupFavouriteRepoMock(mockFavouriteRepo)
			err := favService.DeleteFavourite(tc.UserID, tc.ContentID)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestFavouriteService_GetFavourites(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		UserID                 int
		ExpectedErr            error
		ExpectedFavourites     *dto.FavouritesResponse
		SetupFavouriteRepoMock func(repo *mockrepo.MockFavourite)
		SetupContentUCMock     func(uc *mock_usecase.MockContent)
	}{
		{
			Name:        "Успех",
			UserID:      1,
			ExpectedErr: nil,
			ExpectedFavourites: &dto.FavouritesResponse{
				Favourites: []dto.Favourite{
					{Content: dto.PreviewContent{ID: 1}, Category: "favourite"},
				},
			},
			SetupFavouriteRepoMock: func(repo *mockrepo.MockFavourite) {
				repo.EXPECT().GetFavourites(1).Return([]*entity.Favourite{
					{ContentID: 1, Category: "favourite"},
				}, nil).AnyTimes()
			},
			SetupContentUCMock: func(uc *mock_usecase.MockContent) {
				uc.EXPECT().GetPreviewContentByID(1).Return(&dto.PreviewContent{ID: 1}, nil).AnyTimes()
			},
		},
		{
			Name:        "Ошибка при получении избранного",
			UserID:      1,
			ExpectedErr: entity.UsecaseWrap(errors.New("database error"), errors.New("ошибка при получении избранного контента в FavouriteService")),
			SetupFavouriteRepoMock: func(repo *mockrepo.MockFavourite) {
				repo.EXPECT().GetFavourites(1).Return(nil, errors.New("database error"))
			},
			SetupContentUCMock: func(uc *mock_usecase.MockContent) {

			},
		},
		{
			Name:        "Ошибка при получении контента в избранном",
			UserID:      1,
			ExpectedErr: entity.UsecaseWrap(errors.New("database error"), errors.New("ошибка при получении контента из избранного в FavouriteService")),
			SetupFavouriteRepoMock: func(repo *mockrepo.MockFavourite) {
				repo.EXPECT().GetFavourites(1).Return([]*entity.Favourite{
					{ContentID: 1, Category: "favourite"},
				}, nil).AnyTimes()
			},
			SetupContentUCMock: func(uc *mock_usecase.MockContent) {
				uc.EXPECT().GetPreviewContentByID(1).Return(nil, errors.New("database error"))
			},
		},
		{
			Name:        "Пользователь не найден",
			UserID:      1,
			ExpectedErr: repository.ErrFavouriteUserNotFound,
			SetupFavouriteRepoMock: func(repo *mockrepo.MockFavourite) {
				repo.EXPECT().GetFavourites(1).Return(nil, repository.ErrFavouriteUserNotFound)
			},
			SetupContentUCMock: func(uc *mock_usecase.MockContent) {

			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockFavouriteRepo := mockrepo.NewMockFavourite(ctrl)
			mockContentUC := mock_usecase.NewMockContent(ctrl)
			favService := NewFavouriteService(mockFavouriteRepo, mockContentUC)
			tc.SetupFavouriteRepoMock(mockFavouriteRepo)
			tc.SetupContentUCMock(mockContentUC)
			response, err := favService.GetFavourites(tc.UserID)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedFavourites, response)
		})
	}
}

func TestFavouriteService_GetStatus(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		UserID                 int
		ContentID              int
		ExpectedErr            error
		ExpectedStatus         *dto.FavouriteStatusResponse
		SetupFavouriteRepoMock func(repo *mockrepo.MockFavourite)
	}{
		{
			Name:        "Успех",
			UserID:      1,
			ContentID:   1,
			ExpectedErr: nil,
			ExpectedStatus: &dto.FavouriteStatusResponse{
				Status: "favourite",
			},
			SetupFavouriteRepoMock: func(repo *mockrepo.MockFavourite) {
				repo.EXPECT().GetFavourite(1, 1).Return(&entity.Favourite{
					Category: "favourite",
				}, nil).AnyTimes()
			},
		},
		{
			Name:        "Ошибка при получении статуса",
			UserID:      1,
			ContentID:   1,
			ExpectedErr: entity.UsecaseWrap(errors.New("database error"), errors.New("ошибка при получении статуса контента в избранном в FavouriteService")),
			SetupFavouriteRepoMock: func(repo *mockrepo.MockFavourite) {
				repo.EXPECT().GetFavourite(1, 1).Return(nil, errors.New("database error"))
			},
		},
		{
			Name:        "Избранное не найдено",
			UserID:      1,
			ContentID:   1,
			ExpectedErr: repository.ErrFavouriteNotFound,
			SetupFavouriteRepoMock: func(repo *mockrepo.MockFavourite) {
				repo.EXPECT().GetFavourite(1, 1).Return(nil, repository.ErrFavouriteNotFound)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockFavouriteRepo := mockrepo.NewMockFavourite(ctrl)
			mockContentUC := mock_usecase.NewMockContent(ctrl)
			favService := NewFavouriteService(mockFavouriteRepo, mockContentUC)
			tc.SetupFavouriteRepoMock(mockFavouriteRepo)
			response, err := favService.GetStatus(tc.UserID, tc.ContentID)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedStatus, response)
		})
	}
}

func TestFavouriteService_CreateFavourite(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		UserID                 int
		ContentID              int
		Category               string
		ExpectedErr            error
		SetupFavouriteRepoMock func(repo *mockrepo.MockFavourite)
	}{
		{
			Name:        "Успех",
			UserID:      1,
			ContentID:   1,
			Category:    "newFavourite",
			ExpectedErr: nil,
			SetupFavouriteRepoMock: func(repo *mockrepo.MockFavourite) {
				gomock.InOrder(
					repo.EXPECT().GetFavourite(1, 1).Return(&entity.Favourite{Category: "oldFavourite"}, nil),
					repo.EXPECT().DeleteFavourite(1, 1).Return(nil),
					repo.EXPECT().CreateFavourite(1, 1, "newFavourite").Return(nil),
				)
			},
		},
		{
			Name:        "Избранное уже в нужной категории",
			UserID:      1,
			ContentID:   1,
			Category:    "favourite",
			ExpectedErr: nil,
			SetupFavouriteRepoMock: func(repo *mockrepo.MockFavourite) {
				repo.EXPECT().GetFavourite(1, 1).Return(&entity.Favourite{Category: "favourite"}, nil)
			},
		},
		{
			Name:        "Ошибка при создании избранного",
			UserID:      1,
			ContentID:   1,
			Category:    "newFavourite",
			ExpectedErr: entity.UsecaseWrap(errors.New("database error"), errors.New("ошибка при добавлении в избранное в FavouriteService")),
			SetupFavouriteRepoMock: func(repo *mockrepo.MockFavourite) {
				gomock.InOrder(
					repo.EXPECT().GetFavourite(1, 1).Return(&entity.Favourite{Category: "oldFavourite"}, nil),
					repo.EXPECT().DeleteFavourite(1, 1).Return(nil),
					repo.EXPECT().CreateFavourite(1, 1, "newFavourite").Return(errors.New("database error")),
				)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockFavouriteRepo := mockrepo.NewMockFavourite(ctrl)
			mockContentUC := mock_usecase.NewMockContent(ctrl)
			favService := NewFavouriteService(mockFavouriteRepo, mockContentUC)
			tc.SetupFavouriteRepoMock(mockFavouriteRepo)
			err := favService.CreateFavourite(tc.UserID, tc.ContentID, tc.Category)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
