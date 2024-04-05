package service

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestReviewService_CreateReview(t *testing.T) {
	t.Parallel()

	fixedTime := time.Now()

	testCases := []struct {
		Name                 string
		Input                *dto.ReviewCreate
		ExpectedErr          error
		ExpectedOutput       *dto.ReviewResponse
		SetupReviewRepoMock  func(repo *mockrepo.MockReview)
		SetupUserRepoMock    func(repo *mockrepo.MockUser)
		SetupContentRepoMock func(repo *mockrepo.MockContent)
		SetupStaticRepoMock  func(repo *mockrepo.MockStatic)
	}{
		{
			Name: "Успешное создание",
			Input: &dto.ReviewCreate{
				UserID:    1,
				ContentID: 1,
				Rating:    10,
				Title:     "title",
				Text:      "text",
			},
			ExpectedErr: nil,
			ExpectedOutput: &dto.ReviewResponse{
				Review: dto.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    10,
					Title:     "title",
					Text:      "text",
					CreatedAt: fixedTime.String(),
					Likes:     0,
					Dislikes:  0,
				},
				AuthorName:   "email",
				AuthorAvatar: "path",
				ContentName:  "movie",
			},
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().AddReview(gomock.Any()).Return(&entity.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    10,
					Title:     "title",
					Text:      "text",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Likes:     0,
					Dislikes:  0,
				}, nil)
			},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUser(gomock.Any()).Return(&entity.User{
					ID:             1,
					Email:          "email",
					AvatarUploadID: 1,
				}, nil)
			},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetContent(1).Return(&entity.Content{
					ID:    1,
					Title: "movie",
				}, nil)
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("path", nil)
			},
		},
		{
			Name: "Ошибка создания",
			Input: &dto.ReviewCreate{
				UserID:    1,
				ContentID: 1,
				Rating:    10,
				Title:     "title",
				Text:      "text",
			},
			ExpectedErr:    fmt.Errorf("ошибка создания отзыва"),
			ExpectedOutput: nil,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().AddReview(gomock.Any()).Return(nil, fmt.Errorf("ошибка создания отзыва"))
			},
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {},
			SetupStaticRepoMock:  func(repo *mockrepo.MockStatic) {},
		},
		{
			Name: "Ошибка получения автора",
			Input: &dto.ReviewCreate{
				UserID:    1,
				ContentID: 1,
				Rating:    10,
				Title:     "title",
				Text:      "text",
			},
			ExpectedErr:    fmt.Errorf("ошибка получения автора"),
			ExpectedOutput: nil,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().AddReview(gomock.Any()).Return(&entity.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    10,
					Title:     "title",
					Text:      "text",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Likes:     0,
					Dislikes:  0,
				}, nil)
			},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUser(gomock.Any()).Return(nil, fmt.Errorf("ошибка получения автора"))
			},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {},
			SetupStaticRepoMock:  func(repo *mockrepo.MockStatic) {},
		},
		{
			Name: "Ошибка получения контента",
			Input: &dto.ReviewCreate{
				UserID:    1,
				ContentID: 1,
				Rating:    10,
				Title:     "title",
				Text:      "text",
			},
			ExpectedErr:    fmt.Errorf("ошибка получения контента"),
			ExpectedOutput: nil,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().AddReview(gomock.Any()).Return(&entity.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    10,
					Title:     "title",
					Text:      "text",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Likes:     0,
					Dislikes:  0,
				}, nil)
			},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUser(gomock.Any()).Return(&entity.User{
					ID:             1,
					Email:          "email",
					AvatarUploadID: 1,
				}, nil)
			},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetContent(1).Return(nil, fmt.Errorf("ошибка получения контента"))
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("path", nil)
			},
		},
		{
			Name: "Ошибка получения аватара",
			Input: &dto.ReviewCreate{
				UserID:    1,
				ContentID: 1,
				Rating:    10,
				Title:     "title",
				Text:      "text",
			},
			ExpectedErr:    fmt.Errorf("не удалось получить аватар"),
			ExpectedOutput: nil,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().AddReview(gomock.Any()).Return(&entity.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    10,
					Title:     "title",
					Text:      "text",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Likes:     0,
					Dislikes:  0,
				}, nil)
			},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUser(gomock.Any()).Return(&entity.User{
					ID:             1,
					Email:          "email",
					AvatarUploadID: 1,
				}, nil)
			},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("", fmt.Errorf("не удалось получить аватар"))
			},
		},
		{
			Name: "Ошибка валидации c пустым текстом",
			Input: &dto.ReviewCreate{
				UserID:    1,
				ContentID: 1,
				Rating:    10,
				Title:     "title",
				Text:      "", // пустой текст
			},
			ExpectedErr:          entity.NewClientError("Количество символов в тексте рецензии должно быть от 1 до 10000", entity.ErrBadRequest),
			ExpectedOutput:       nil,
			SetupReviewRepoMock:  func(repo *mockrepo.MockReview) {},
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {},
			SetupStaticRepoMock:  func(repo *mockrepo.MockStatic) {},
		},
		{
			Name: "Ошибка валидации c пустым заголовком",
			Input: &dto.ReviewCreate{
				UserID:    1,
				ContentID: 1,
				Rating:    10,
				Title:     "", // пустой заголовок
			},
			ExpectedErr:          entity.NewClientError("Количество символов в заголовке рецензии должно быть от 1 до 50", entity.ErrBadRequest),
			ExpectedOutput:       nil,
			SetupReviewRepoMock:  func(repo *mockrepo.MockReview) {},
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {},
			SetupStaticRepoMock:  func(repo *mockrepo.MockStatic) {},
		},
		{
			Name: "Ошибка валидации с рейтингом больше 10 или меньше 1",
			Input: &dto.ReviewCreate{
				UserID:    1,
				ContentID: 1,
				Rating:    11, // рейтинг больше 10
				Title:     "title",
			},
			ExpectedErr:          entity.NewClientError("Рейтинг должен быть в диапазоне от 1 до 10", entity.ErrBadRequest),
			ExpectedOutput:       nil,
			SetupReviewRepoMock:  func(repo *mockrepo.MockReview) {},
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {},
			SetupStaticRepoMock:  func(repo *mockrepo.MockStatic) {},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockReviewRepo := mockrepo.NewMockReview(ctrl)
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			reviewService := NewReviewService(mockReviewRepo, mockUserRepo, mockContentRepo, mockStaticRepo)
			tc.SetupReviewRepoMock(mockReviewRepo)
			tc.SetupContentRepoMock(mockContentRepo)
			tc.SetupUserRepoMock(mockUserRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			output, err := reviewService.CreateReview(*tc.Input)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedOutput, output)
		})
	}
}

func TestReviewService_GetReview(t *testing.T) {
	t.Parallel()

	fixedTime := time.Now()

	testCases := []struct {
		Name                 string
		Input                int
		ExpectedErr          error
		ExpectedOutput       *dto.ReviewResponse
		SetupReviewRepoMock  func(repo *mockrepo.MockReview)
		SetupUserRepoMock    func(repo *mockrepo.MockUser)
		SetupContentRepoMock func(repo *mockrepo.MockContent)
		SetupStaticRepoMock  func(repo *mockrepo.MockStatic)
	}{
		{
			Name:        "Успешное получение",
			Input:       1,
			ExpectedErr: nil,
			ExpectedOutput: &dto.ReviewResponse{
				Review: dto.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    10,
					Title:     "title",
					Text:      "text",
					CreatedAt: fixedTime.String(),
					Likes:     0,
					Dislikes:  0,
				},
				AuthorName:   "email",
				AuthorAvatar: "path",
				ContentName:  "movie",
			},
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetReviewByID(1).Return(&entity.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    10,
					Title:     "title",
					Text:      "text",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Likes:     0,
					Dislikes:  0,
				}, nil)
			},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUser(gomock.Any()).Return(&entity.User{
					ID:             1,
					Email:          "email",
					AvatarUploadID: 1,
				}, nil)
			},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetContent(1).Return(&entity.Content{
					ID:    1,
					Title: "movie",
				}, nil)
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {
				repo.EXPECT().GetStatic(1).Return("path", nil)
			},
		},
		{
			Name:           "Ошибка получения",
			Input:          1,
			ExpectedErr:    fmt.Errorf("ошибка получения отзыва"),
			ExpectedOutput: nil,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetReviewByID(1).Return(nil, fmt.Errorf("ошибка получения отзыва"))
			},
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {},
			SetupStaticRepoMock:  func(repo *mockrepo.MockStatic) {},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockReviewRepo := mockrepo.NewMockReview(ctrl)
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			reviewService := NewReviewService(mockReviewRepo, mockUserRepo, mockContentRepo, mockStaticRepo)
			tc.SetupReviewRepoMock(mockReviewRepo)
			tc.SetupContentRepoMock(mockContentRepo)
			tc.SetupUserRepoMock(mockUserRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			output, err := reviewService.GetReview(tc.Input)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedOutput, output)
		})
	}
}

func TestReviewService_EditReview(t *testing.T) {
	t.Parallel()

	fixedTime := time.Now()

	testCases := []struct {
		Name                 string
		Input                *dto.ReviewUpdate
		ExpectedErr          error
		ExpectedOutput       *dto.ReviewResponse
		SetupReviewRepoMock  func(repo *mockrepo.MockReview)
		SetupUserRepoMock    func(repo *mockrepo.MockUser)
		SetupContentRepoMock func(repo *mockrepo.MockContent)
		SetupStaticRepoMock  func(repo *mockrepo.MockStatic)
	}{
		{
			Name: "Успешное редактирование",
			Input: &dto.ReviewUpdate{
				UserID:   1,
				ReviewID: 1,
				Rating:   5,
				Title:    "title",
				Text:     "text",
			},
			ExpectedErr: nil,
			ExpectedOutput: &dto.ReviewResponse{
				Review: dto.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    5,
					Title:     "title",
					Text:      "text",
					CreatedAt: fixedTime.String(),
					Likes:     0,
					Dislikes:  0,
				},
				AuthorName:   "name",
				AuthorAvatar: "",
				ContentName:  "movie",
			},
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetReviewByID(1).Return(&entity.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    10,
					Title:     "title",
					Text:      "text",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Likes:     0,
					Dislikes:  0,
				}, nil)
				repo.EXPECT().UpdateReview(gomock.Any()).Return(&entity.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    5,
					Title:     "title",
					Text:      "text",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Likes:     0,
					Dislikes:  0,
				}, nil)
			},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUser(gomock.Any()).Return(&entity.User{
					ID:             1,
					Name:           "name",
					Email:          "email",
					AvatarUploadID: 0,
				}, nil)
			},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetContent(1).Return(&entity.Content{
					ID:    1,
					Title: "movie",
				}, nil)
			},
			SetupStaticRepoMock: func(repo *mockrepo.MockStatic) {},
		},
		{
			Name: "Ошибка редактирования на уровне БД",
			Input: &dto.ReviewUpdate{
				UserID:   1,
				ReviewID: 1,
				Rating:   5,
				Title:    "title",
				Text:     "text",
			},
			ExpectedErr:    fmt.Errorf("ошибка редактирования отзыва"),
			ExpectedOutput: nil,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetReviewByID(1).Return(&entity.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    10,
					Title:     "title",
					Text:      "text",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Likes:     0,
					Dislikes:  0,
				}, nil)
				repo.EXPECT().UpdateReview(gomock.Any()).Return(nil, fmt.Errorf("ошибка редактирования отзыва"))
			},
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {},
			SetupStaticRepoMock:  func(repo *mockrepo.MockStatic) {},
		},
		{
			Name: "Ошибка редактирования на уровне валидации",
			Input: &dto.ReviewUpdate{
				UserID:   1,
				ReviewID: 1,
				Rating:   11, // рейтинг больше 10
				Title:    "title",
				Text:     "text",
			},
			ExpectedErr:    entity.NewClientError("Рейтинг должен быть в диапазоне от 1 до 10", entity.ErrBadRequest),
			ExpectedOutput: nil,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetReviewByID(1).Return(&entity.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    10,
					Title:     "title",
					Text:      "text",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Likes:     0,
					Dislikes:  0,
				}, nil)
			},
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {},
			SetupStaticRepoMock:  func(repo *mockrepo.MockStatic) {},
		},
		{
			Name: "Ошибка доступа",
			Input: &dto.ReviewUpdate{
				UserID:   2, // другой пользователь
				ReviewID: 1,
				Rating:   5,
				Title:    "title",
				Text:     "text",
			},
			ExpectedErr:    entity.NewClientError("Вы не можете редактировать чужой отзыв", entity.ErrForbidden),
			ExpectedOutput: nil,
			SetupReviewRepoMock: func(repo *mockrepo.MockReview) {
				repo.EXPECT().GetReviewByID(1).Return(&entity.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    10,
					Title:     "title",
					Text:      "text",
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Likes:     0,
					Dislikes:  0,
				}, nil)
			},
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {},
			SetupStaticRepoMock:  func(repo *mockrepo.MockStatic) {},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockReviewRepo := mockrepo.NewMockReview(ctrl)
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			reviewService := NewReviewService(mockReviewRepo, mockUserRepo, mockContentRepo, mockStaticRepo)
			tc.SetupReviewRepoMock(mockReviewRepo)
			tc.SetupContentRepoMock(mockContentRepo)
			tc.SetupUserRepoMock(mockUserRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			output, err := reviewService.EditReview(*tc.Input)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedOutput, output)
		})
	}
}
