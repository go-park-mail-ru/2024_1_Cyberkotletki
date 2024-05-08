package service

import (
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
)

func TestGetLatestReviews(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		Limit               int
		ExpectedErr         error
		SetupReviewRepoMock func(repo *mockrepo.MockReview)
		SetupUserRepoMock   func(repo *mockrepo.MockUser)
	}{
		{
			Name:        "Успешно",
			Limit:       10,
			ExpectedErr: nil,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().AddReview(&entity.Review{
					AuthorID:      1,
					ContentID:     1,
					ContentRating: 5,
					Title:         "Great Content!",
					Text:          "I really enjoyed this content.",
				}).Return(&entity.Review{
					ID:            1,
					AuthorID:      1,
					ContentID:     1,
					ContentRating: 5,
					Title:         "Great Content!",
					Text:          "I really enjoyed this content.",
				}, nil).AnyTimes()

				// Add this line
				repo.EXPECT().GetLatestReviews(gomock.Any()).Return(nil, nil).AnyTimes()
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockReviewRepo := mockrepo.NewMockReview(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			service := NewReviewService(mockReviewRepo, mockUserRepo,
				mockContentRepo, mockStaticRepo)
			tc.SetupReviewRepoMock(mockReviewRepo)
			_, err := service.GetLatestReviews(tc.Limit)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewService_GetUserReviews(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		UserID              int
		Count               int
		Page                int
		ExpectedErr         error
		SetupReviewRepoMock func(repo *mockrepo.MockReview)
	}{
		{
			Name:        "Success",
			UserID:      1,
			Count:       5,
			Page:        1,
			ExpectedErr: nil,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetReviewsByAuthorID(1, 1, 5).Return([]*entity.Review{}, nil).AnyTimes()
				repo.EXPECT().GetReviewsCountByAuthorID(1).Return(10, nil).AnyTimes()
			},
		},
		{
			Name:        "Error when getting reviews",
			UserID:      1,
			Count:       5,
			Page:        1,
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при получении отзывов пользователя"), errors.New("database error")),
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetReviewsByAuthorID(1, 1, 5).Return(nil, errors.New("database error"))
			},
		},
		{
			Name:        "Error when getting reviews count",
			UserID:      1,
			Count:       5,
			Page:        1,
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при получении количества отзывов пользователя"), errors.New("database error")),
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetReviewsByAuthorID(1, 1, 5).Return([]*entity.Review{}, nil).AnyTimes()
				repo.EXPECT().GetReviewsCountByAuthorID(1).Return(0, errors.New("database error"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockReviewRepo := mockrepo.NewMockReview(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			service := NewReviewService(mockReviewRepo, mockUserRepo,
				mockContentRepo, mockStaticRepo)
			tc.SetupReviewRepoMock(mockReviewRepo)
			_, err := service.GetUserReviews(tc.UserID, tc.Count, tc.Page)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewService_GetContentReviews(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		ContentID           int
		Count               int
		Page                int
		ExpectedErr         error
		SetupReviewRepoMock func(repo *mockrepo.MockReview)
	}{
		{
			Name:        "Success",
			ContentID:   1,
			Count:       5,
			Page:        1,
			ExpectedErr: nil,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetReviewsByContentID(1, 1, 5).Return([]*entity.Review{}, nil).AnyTimes()
				repo.EXPECT().GetReviewsCountByContentID(1).Return(10, nil).AnyTimes()
			},
		},
		{
			Name:        "Error when getting reviews",
			ContentID:   1,
			Count:       5,
			Page:        1,
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при получении отзывов контента"), errors.New("database error")),
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetReviewsByContentID(1, 1, 5).Return(nil, errors.New("database error"))
			},
		},
		{
			Name:        "Error when getting reviews count",
			ContentID:   1,
			Count:       5,
			Page:        1,
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при получении количества отзывов"), errors.New("database error")),
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetReviewsByContentID(1, 1, 5).Return([]*entity.Review{}, nil).AnyTimes()
				repo.EXPECT().GetReviewsCountByContentID(1).Return(0, errors.New("database error"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockReviewRepo := mockrepo.NewMockReview(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			service := NewReviewService(mockReviewRepo, mockUserRepo,
				mockContentRepo, mockStaticRepo)
			tc.SetupReviewRepoMock(mockReviewRepo)
			_, err := service.GetContentReviews(tc.ContentID, tc.Count, tc.Page)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewService_GetReview(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		ReviewID            int
		ExpectedErr         error
		SetupReviewRepoMock func(repo *mockrepo.MockReview)
	}{
		{
			Name:        "Review not found",
			ReviewID:    1,
			ExpectedErr: repository.ErrReviewNotFound,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetReviewByID(1).Return(nil, repository.ErrReviewNotFound)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockReviewRepo := mockrepo.NewMockReview(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			service := NewReviewService(mockReviewRepo, mockUserRepo,
				mockContentRepo, mockStaticRepo)
			tc.SetupReviewRepoMock(mockReviewRepo)
			_, err := service.GetReview(tc.ReviewID)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewService_GetContentReviewByAuthor(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		AuthorID            int
		ContentID           int
		UserID              int
		ExpectedErr         error
		SetupReviewRepoMock func(repo *mockrepo.MockReview)
		SetupUserRepoMock   func(repo *mockrepo.MockUser)
	}{
		{
			Name:        "Error when getting review",
			AuthorID:    1,
			ContentID:   1,
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при получении отзыва"), errors.New("database error")),
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetContentReviewByAuthor(1, 1).Return(nil, errors.New("database error"))
			},
		},
		{
			Name:        "Review not found",
			AuthorID:    1,
			ContentID:   1,
			ExpectedErr: usecase.ErrReviewNotFound,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetContentReviewByAuthor(1, 1).Return(nil, repository.ErrReviewNotFound)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockReviewRepo := mockrepo.NewMockReview(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			service := NewReviewService(mockReviewRepo, mockUserRepo,
				mockContentRepo, mockStaticRepo)
			tc.SetupReviewRepoMock(mockReviewRepo)
			_, err := service.GetContentReviewByAuthor(tc.AuthorID, tc.ContentID)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewService_CreateReview(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                string
		Review              dto.ReviewCreate
		ExpectedErr         error
		SetupReviewRepoMock func(repo *mockrepo.MockReview)
		SetupUserRepoMock   func(repo *mockrepo.MockUser)
	}{
		{
			Name: "Error when adding review",
			Review: dto.ReviewCreate{
				ReviewCreateRequest: dto.ReviewCreateRequest{
					ContentID: 1,
					Rating:    5,
					Title:     "Test Title",
					Text:      "Test Text",
				},
				UserID: 1,
			},
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при добавлении отзыва"), errors.New("database error")),
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().AddReview(gomock.Any()).Return(nil, errors.New("database error"))
			},
		},
		{
			Name: "Review already exists",
			Review: dto.ReviewCreate{
				ReviewCreateRequest: dto.ReviewCreateRequest{
					ContentID: 1,
					Rating:    5,
					Title:     "Test Title",
					Text:      "Test Text",
				},
				UserID: 1,
			},
			ExpectedErr: usecase.ErrReviewAlreadyExists,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().AddReview(gomock.Any()).Return(nil, repository.ErrReviewAlreadyExists)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockReviewRepo := mockrepo.NewMockReview(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			service := NewReviewService(mockReviewRepo, mockUserRepo,
				mockContentRepo, mockStaticRepo)
			tc.SetupReviewRepoMock(mockReviewRepo)
			_, err := service.CreateReview(tc.Review)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestDeleteReview(t *testing.T) {
	// Setup
}

func TestVoteReview(t *testing.T) {
	// Setup
}

func TestIsVotedByUser(t *testing.T) {
	// Setup
}

func TestUnVoteReview(t *testing.T) {
	// Setup
}
