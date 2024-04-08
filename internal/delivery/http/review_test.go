package http

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockusecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReviewHandler_GetReview(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		ReviewID               string
		ExpectedErr            error
		Cookies                *http.Cookie
		SetupReviewUsecaseMock func(usecase *mockusecase.MockReview)
	}{
		{
			Name:        "Успешное получение отзыва",
			ReviewID:    "1",
			ExpectedErr: nil,
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetReview(1).Return(nil, nil)
			},
		},
		{
			Name:     "Отзыв не найден",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "отзыв не найден",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetReview(1).Return(nil, entity.NewClientError("отзыв не найден", entity.ErrNotFound))
			},
		},
		{
			Name:        "Неожиданная ошибка",
			ReviewID:    "1",
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "internal"},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetReview(1).Return(nil, entity.NewClientError("internal", entity.ErrInternal))
			},
		},
		{
			Name:                   "Невалидный айди",
			ReviewID:               "ogo!",
			ExpectedErr:            &echo.HTTPError{Code: 400, Message: "невалидный id рецензии"},
			Cookies:                nil,
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			e := echo.New()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockReviewUsecase := mockusecase.NewMockReview(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			tc.SetupReviewUsecaseMock(mockReviewUsecase)
			reviewHandler := NewReviewEndpoints(mockReviewUsecase, mockAuthUsecase)
			req := httptest.NewRequest(http.MethodGet, "/review/", nil)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/review/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.ReviewID)
			err := reviewHandler.GetReview(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewEndpoints_GetMyContentReview(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		ContentID              string
		ExpectedErr            error
		ExpectedOutput         *dto.ReviewResponse
		Cookies                *http.Cookie
		SetupReviewUsecaseMock func(usecase *mockusecase.MockReview)
		SetupAuthUsecaseMock   func(usecase *mockusecase.MockAuth)
	}{
		{
			Name:        "Успешное получение отзыва",
			ContentID:   "1",
			ExpectedErr: nil,
			ExpectedOutput: &dto.ReviewResponse{
				Review: dto.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    5,
					Title:     "Title",
					Text:      "i like it",
					CreatedAt: "2022-01-02T15:04:05Z",
					Likes:     5,
					Dislikes:  5,
				},
				AuthorName:   "Author",
				AuthorAvatar: "avatars/avatar.jpg",
				ContentName:  "Content",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetContentReviewByAuthor(1, 1).Return(&dto.ReviewResponse{
					Review: dto.Review{
						ID:        1,
						AuthorID:  1,
						ContentID: 1,
						Rating:    5,
						Title:     "Title",
						Text:      "i like it",
						CreatedAt: "2022-01-02T15:04:05Z",
						Likes:     5,
						Dislikes:  5,
					},
					AuthorName:   "Author",
					AuthorAvatar: "avatars/avatar.jpg",
					ContentName:  "Content",
				}, nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:      "Невалидный айди",
			ContentID: "ogo!",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный id контента",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock:   func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name:      "Отзыв не найден",
			ContentID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "отзыв не найден",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetContentReviewByAuthor(1, 1).Return(nil, entity.NewClientError("отзыв не найден", entity.ErrNotFound))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:                   "Пользователь не авторизован",
			ContentID:              "1",
			ExpectedErr:            &echo.HTTPError{Code: 401, Message: "необходима авторизация"},
			Cookies:                &http.Cookie{Name: "session", Value: "xxx"},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession(gomock.Any()).Return(-1, entity.NewClientError("необходима авторизация", entity.ErrUnauthorized))
			},
		},
		{
			Name:        "Неожиданная ошибка",
			ContentID:   "1",
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "internal"},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetContentReviewByAuthor(1, 1).Return(nil, entity.NewClientError("internal", entity.ErrInternal))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
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
			mockReviewUsecase := mockusecase.NewMockReview(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			tc.SetupReviewUsecaseMock(mockReviewUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			reviewHandler := NewReviewEndpoints(mockReviewUsecase, mockAuthUsecase)
			req := httptest.NewRequest(http.MethodGet, "/review/myReview", nil)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/review/myReview")
			c.QueryParams().Add("id", tc.ContentID)
			err := reviewHandler.GetMyContentReview(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewEndpoints_CreateReview(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		Body                   func() io.Reader
		ExpectedErr            error
		ExpectedOutput         *dto.ReviewResponse
		Cookies                *http.Cookie
		SetupReviewUsecaseMock func(usecase *mockusecase.MockReview)
		SetupAuthUsecaseMock   func(usecase *mockusecase.MockAuth)
	}{
		{
			Name: "Успешное создание отзыва",
			Body: func() io.Reader {
				return strings.NewReader(`{"contentID":1,"rating":5,"title":"Title","text":"i like it"}`)
			},
			ExpectedOutput: &dto.ReviewResponse{
				Review: dto.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    5,
					Title:     "Title",
					Text:      "i like it",
					CreatedAt: "2022-01-02T15:04:05Z",
					Likes:     5,
					Dislikes:  5,
				},
				AuthorName:   "Author",
				AuthorAvatar: "avatars/avatar.jpg",
				ContentName:  "Content",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().CreateReview(dto.ReviewCreate{
					ReviewCreateRequest: dto.ReviewCreateRequest{
						ContentID: 1,
						Rating:    5,
						Title:     "Title",
						Text:      "i like it",
					},
					UserID: 1,
				}).Return(&dto.ReviewResponse{
					Review: dto.Review{
						ID:        1,
						AuthorID:  1,
						ContentID: 1,
						Rating:    5,
						Title:     "Title",
						Text:      "i like it",
						CreatedAt: "2022-01-02T15:04:05Z",
						Likes:     5,
						Dislikes:  5,
					},
					AuthorName:   "Author",
					AuthorAvatar: "avatars/avatar.jpg",
					ContentName:  "Content",
				}, nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name: "Невалидный запрос",
			Body: func() io.Reader {
				return strings.NewReader("kek")
			},
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный запрос",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock:   func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name: "Невалидный айди",
			Body: func() io.Reader {
				return strings.NewReader(`{"contentID":-100,"rating":5,"title":"Title","text":"i like it"}`)
			},
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный id контента",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().CreateReview(dto.ReviewCreate{
					ReviewCreateRequest: dto.ReviewCreateRequest{
						ContentID: -100,
						Rating:    5,
						Title:     "Title",
						Text:      "i like it",
					},
					UserID: 1,
				}).Return(nil, entity.NewClientError("невалидный id контента", entity.ErrBadRequest))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name: "Пользователь не авторизован",
			Body: func() io.Reader {
				return strings.NewReader(`{"contentID":1,"rating":5,"title":"Title","text":"i like it"}`)
			},
			ExpectedErr:            &echo.HTTPError{Code: 401, Message: "необходима авторизация"},
			Cookies:                &http.Cookie{Name: "session", Value: "xxx"},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession(gomock.Any()).Return(-1, entity.NewClientError("необходима авторизация", entity.ErrUnauthorized))
			},
		},
		{
			Name: "Неожиданная ошибка",
			Body: func() io.Reader {
				return strings.NewReader(`{"contentID":1,"rating":5,"title":"Title","text":"i like it"}`)
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "internal"},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().CreateReview(dto.ReviewCreate{
					ReviewCreateRequest: dto.ReviewCreateRequest{
						ContentID: 1,
						Rating:    5,
						Title:     "Title",
						Text:      "i like it",
					},
					UserID: 1,
				}).Return(nil, entity.NewClientError("internal", entity.ErrInternal))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
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
			mockReviewUsecase := mockusecase.NewMockReview(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			tc.SetupReviewUsecaseMock(mockReviewUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			reviewHandler := NewReviewEndpoints(mockReviewUsecase, mockAuthUsecase)
			req := httptest.NewRequest(http.MethodPost, "/review/", tc.Body())
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/review/")
			c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			err := reviewHandler.CreateReview(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewEndpoints_UpdateReview(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		Body                   func() io.Reader
		ExpectedErr            error
		ExpectedOutput         *dto.ReviewResponse
		Cookies                *http.Cookie
		SetupReviewUsecaseMock func(usecase *mockusecase.MockReview)
		SetupAuthUsecaseMock   func(usecase *mockusecase.MockAuth)
	}{
		{
			Name: "Успешное обновление отзыва",
			Body: func() io.Reader {
				return strings.NewReader(`{"reviewID":1,"rating":5,"title":"Title","text":"i like it"}`)
			},
			ExpectedOutput: &dto.ReviewResponse{
				Review: dto.Review{
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Rating:    5,
					Title:     "Title",
					Text:      "i like it",
					CreatedAt: "2022-01-02T15:04:05Z",
					Likes:     5,
					Dislikes:  5,
				},
				AuthorName:   "Author",
				AuthorAvatar: "avatars/avatar.jpg",
				ContentName:  "Content",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().EditReview(dto.ReviewUpdate{
					ReviewUpdateRequest: dto.ReviewUpdateRequest{
						ReviewID: 1,
						Rating:   5,
						Title:    "Title",
						Text:     "i like it",
					},
					UserID: 1,
				}).Return(&dto.ReviewResponse{
					Review: dto.Review{
						ID:        1,
						AuthorID:  1,
						ContentID: 1,
						Rating:    5,
						Title:     "Title",
						Text:      "i like it",
						CreatedAt: "2022-01-02T15:04:05Z",
						Likes:     5,
						Dislikes:  5,
					},
					AuthorName:   "Author",
					AuthorAvatar: "avatars/avatar.jpg",
					ContentName:  "Content",
				}, nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name: "Невалидный запрос",
			Body: func() io.Reader {
				return strings.NewReader("kek")
			},
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный запрос",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock:   func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name: "Отзыв не найден",
			Body: func() io.Reader {
				return strings.NewReader(`{"reviewID":1,"rating":5,"title":"Title","text":"i like it"}`)
			},
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "отзыв не найден",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().EditReview(dto.ReviewUpdate{
					ReviewUpdateRequest: dto.ReviewUpdateRequest{
						ReviewID: 1,
						Rating:   5,
						Title:    "Title",
						Text:     "i like it",
					},
					UserID: 1,
				}).Return(nil, entity.NewClientError("отзыв не найден", entity.ErrNotFound))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name: "Нет доступа к чужой рецензии",
			Body: func() io.Reader {
				return strings.NewReader(`{"reviewID":1,"rating":5,"title":"Title","text":"i like it"}`)
			},
			ExpectedErr: &echo.HTTPError{
				Code:    403,
				Message: "нет доступа к чужой рецензии",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().EditReview(dto.ReviewUpdate{
					ReviewUpdateRequest: dto.ReviewUpdateRequest{
						ReviewID: 1,
						Rating:   5,
						Title:    "Title",
						Text:     "i like it",
					},
					UserID: 1,
				}).Return(nil, entity.NewClientError("нет доступа к чужой рецензии", entity.ErrForbidden))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name: "Невалидный айди",
			Body: func() io.Reader {
				return strings.NewReader(`{"reviewID":-100,"rating":5,"title":"Title","text":"i like it"}`)
			},
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный id рецензии",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().EditReview(dto.ReviewUpdate{
					ReviewUpdateRequest: dto.ReviewUpdateRequest{
						ReviewID: -100,
						Rating:   5,
						Title:    "Title",
						Text:     "i like it",
					},
					UserID: 1,
				}).Return(nil, entity.NewClientError("невалидный id рецензии", entity.ErrBadRequest))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name: "Пользователь не авторизован",
			Body: func() io.Reader {
				return strings.NewReader(`{"reviewID":1,"rating":5,"title":"Title","text":"i like it"}`)
			},
			ExpectedErr:            &echo.HTTPError{Code: 401, Message: "необходима авторизация"},
			Cookies:                &http.Cookie{Name: "session", Value: "xxx"},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession(gomock.Any()).Return(-1, entity.NewClientError("необходима авторизация", entity.ErrUnauthorized))
			},
		},
		{
			Name: "Неожиданная ошибка",
			Body: func() io.Reader {
				return strings.NewReader(`{"reviewID":1,"rating":5,"title":"Title","text":"i like it"}`)
			},
			ExpectedErr: &echo.HTTPError{Code: 500, Message: "internal"},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().EditReview(dto.ReviewUpdate{
					ReviewUpdateRequest: dto.ReviewUpdateRequest{
						ReviewID: 1,
						Rating:   5,
						Title:    "Title",
						Text:     "i like it",
					},
					UserID: 1,
				}).Return(nil, entity.NewClientError("internal", entity.ErrInternal))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
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
			mockReviewUsecase := mockusecase.NewMockReview(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			tc.SetupReviewUsecaseMock(mockReviewUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			reviewHandler := NewReviewEndpoints(mockReviewUsecase, mockAuthUsecase)
			req := httptest.NewRequest(http.MethodPut, "/review/", tc.Body())
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/review/")
			c.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			err := reviewHandler.UpdateReview(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewEndpoints_DeleteReview(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		ReviewID               string
		ExpectedErr            error
		Cookies                *http.Cookie
		SetupReviewUsecaseMock func(usecase *mockusecase.MockReview)
		SetupAuthUsecaseMock   func(usecase *mockusecase.MockAuth)
	}{
		{
			Name:        "Успешное удаление отзыва",
			ReviewID:    "1",
			ExpectedErr: nil,
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().DeleteReview(1, 1).Return(nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:     "Отзыв не найден",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "отзыв не найден",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().DeleteReview(1, 1).Return(entity.NewClientError("отзыв не найден", entity.ErrNotFound))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:     "Нет доступа к чужому отзыву",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    403,
				Message: "нет доступа к чужому отзыву",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().DeleteReview(1, 1).Return(entity.NewClientError("нет доступа к чужому отзыву", entity.ErrForbidden))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:     "Невалидный айди",
			ReviewID: "ogo!",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный id рецензии",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock:   func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name:     "Пользователь не авторизован",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    401,
				Message: "необходима авторизация",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession(gomock.Any()).Return(-1, entity.NewClientError("необходима авторизация", entity.ErrUnauthorized))
			},
		},
		{
			Name:     "Неожиданная ошибка",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    500,
				Message: "internal",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().DeleteReview(1, 1).Return(entity.NewClientError("internal", entity.ErrInternal))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
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
			mockReviewUsecase := mockusecase.NewMockReview(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			tc.SetupReviewUsecaseMock(mockReviewUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			reviewHandler := NewReviewEndpoints(mockReviewUsecase, mockAuthUsecase)
			req := httptest.NewRequest(http.MethodDelete, "/review/", nil)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/review/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.ReviewID)
			err := reviewHandler.DeleteReview(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewEndpoints_GetRecentReviews(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		ExpectedErr            error
		ExpectedOutput         *dto.ReviewResponseList
		SetupReviewUsecaseMock func(usecase *mockusecase.MockReview)
	}{
		{
			Name:        "Успешное получение отзывов",
			ExpectedErr: nil,
			ExpectedOutput: &dto.ReviewResponseList{
				Reviews: []dto.ReviewResponse{
					{
						Review: dto.Review{
							ID:        1,
							AuthorID:  1,
							ContentID: 1,
							Rating:    5,
							Title:     "Title",
							Text:      "i like it",
							CreatedAt: "2022-01-02T15:04:05Z",
							Likes:     5,
							Dislikes:  5,
						},
						AuthorName:   "Author",
						AuthorAvatar: "avatars/avatar.jpg",
						ContentName:  "Content",
					},
				},
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetLatestReviews(3).Return(&dto.ReviewResponseList{
					Reviews: []dto.ReviewResponse{
						{
							Review: dto.Review{
								ID:        1,
								AuthorID:  1,
								ContentID: 1,
								Rating:    5,
								Title:     "Title",
								Text:      "i like it",
								CreatedAt: "2022-01-02T15:04:05Z",
								Likes:     5,
								Dislikes:  5,
							},
							AuthorName:   "Author",
							AuthorAvatar: "avatars/avatar.jpg",
							ContentName:  "Content",
						},
					},
				}, nil)
			},
		},
		{
			Name:           "Неожиданная ошибка",
			ExpectedErr:    &echo.HTTPError{Code: 500, Message: "internal"},
			ExpectedOutput: nil,
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetLatestReviews(3).Return(nil, entity.NewClientError("internal", entity.ErrInternal))
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
			mockReviewUsecase := mockusecase.NewMockReview(ctrl)
			tc.SetupReviewUsecaseMock(mockReviewUsecase)
			reviewHandler := NewReviewEndpoints(mockReviewUsecase, nil)
			req := httptest.NewRequest(http.MethodGet, "/review/recent", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			err := reviewHandler.GetRecentReviews(c)
			require.Equal(t, tc.ExpectedErr, err)
			if tc.ExpectedErr == nil {
				var reviews dto.ReviewResponseList
				err = json.NewDecoder(rec.Body).Decode(&reviews)
				require.NoError(t, err)
				require.Equal(t, *tc.ExpectedOutput, reviews)
			}
		})
	}
}

func TestReviewEndpoints_GetUserLatestReviews(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		UserID                 string
		ExpectedErr            error
		ExpectedOutput         *dto.ReviewResponseList
		SetupReviewUsecaseMock func(usecase *mockusecase.MockReview)
	}{
		{
			Name:        "Успешное получение отзывов",
			UserID:      "1",
			ExpectedErr: nil,
			ExpectedOutput: &dto.ReviewResponseList{
				Reviews: []dto.ReviewResponse{
					{
						Review: dto.Review{
							ID:        1,
							AuthorID:  1,
							ContentID: 1,
							Rating:    5,
							Title:     "Title",
							Text:      "i like it",
							CreatedAt: "2022-01-02T15:04:05Z",
							Likes:     5,
							Dislikes:  5,
						},
						AuthorName:   "Author",
						AuthorAvatar: "avatars/avatar.jpg",
						ContentName:  "Content",
					},
				},
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetUserReviews(1, 3, 1).Return(&dto.ReviewResponseList{
					Reviews: []dto.ReviewResponse{
						{
							Review: dto.Review{
								ID:        1,
								AuthorID:  1,
								ContentID: 1,
								Rating:    5,
								Title:     "Title",
								Text:      "i like it",
								CreatedAt: "2022-01-02T15:04:05Z",
								Likes:     5,
								Dislikes:  5,
							},
							AuthorName:   "Author",
							AuthorAvatar: "avatars/avatar.jpg",
							ContentName:  "Content",
						},
					},
				}, nil)
			},
		},
		{
			Name:   "Неожиданная ошибка",
			UserID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    500,
				Message: "internal",
			},
			ExpectedOutput: nil,
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetUserReviews(1, 3, 1).Return(nil, entity.NewClientError("internal", entity.ErrInternal))
			},
		},
		{
			Name:   "Невалидный айди",
			UserID: "ogo!",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный id пользователя",
			},
			ExpectedOutput:         nil,
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
		},
		{
			Name:   "Пользователь не найден",
			UserID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "пользователь не найден",
			},
			ExpectedOutput: nil,
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetUserReviews(1, 3, 1).Return(nil, entity.NewClientError("пользователь не найден", entity.ErrNotFound))
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
			mockReviewUsecase := mockusecase.NewMockReview(ctrl)
			tc.SetupReviewUsecaseMock(mockReviewUsecase)
			reviewHandler := NewReviewEndpoints(mockReviewUsecase, nil)
			req := httptest.NewRequest(http.MethodGet, "/review/user", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/review/user/:id/recent")
			c.SetParamNames("id")
			c.SetParamValues(tc.UserID)
			err := reviewHandler.GetUserLatestReviews(c)
			require.Equal(t, tc.ExpectedErr, err)
			if tc.ExpectedErr == nil {
				var reviews dto.ReviewResponseList
				err = json.NewDecoder(rec.Body).Decode(&reviews)
				require.NoError(t, err)
				require.Equal(t, *tc.ExpectedOutput, reviews)
			}
		})
	}
}

func TestReviewEndpoints_GetUserReviews(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		UserID                 string
		Page                   string
		ExpectedErr            error
		ExpectedOutput         *dto.ReviewResponseList
		SetupReviewUsecaseMock func(usecase *mockusecase.MockReview)
	}{
		{
			Name:        "Успешное получение отзывов",
			UserID:      "1",
			Page:        "1",
			ExpectedErr: nil,
			ExpectedOutput: &dto.ReviewResponseList{
				Reviews: []dto.ReviewResponse{
					{
						Review: dto.Review{
							ID:        1,
							AuthorID:  1,
							ContentID: 1,
							Rating:    5,
							Title:     "Title",
							Text:      "i like it",
							CreatedAt: "2022-01-02T15:04:05Z",
							Likes:     5,
							Dislikes:  5,
						},
						AuthorName:   "Author",
						AuthorAvatar: "avatars/avatar.jpg",
						ContentName:  "Content",
					},
				},
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetUserReviews(1, 10, 1).Return(&dto.ReviewResponseList{
					Reviews: []dto.ReviewResponse{
						{
							Review: dto.Review{
								ID:        1,
								AuthorID:  1,
								ContentID: 1,
								Rating:    5,
								Title:     "Title",
								Text:      "i like it",
								CreatedAt: "2022-01-02T15:04:05Z",
								Likes:     5,
								Dislikes:  5,
							},
							AuthorName:   "Author",
							AuthorAvatar: "avatars/avatar.jpg",
							ContentName:  "Content",
						},
					},
				}, nil)
			},
		},
		{
			Name:   "Неожиданная ошибка",
			UserID: "1",
			Page:   "1",
			ExpectedErr: &echo.HTTPError{
				Code:    500,
				Message: "internal",
			},
			ExpectedOutput: nil,
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetUserReviews(1, 10, 1).Return(nil, entity.NewClientError("internal", entity.ErrInternal))
			},
		},
		{
			Name:   "Невалидный айди",
			UserID: "ogo!",
			Page:   "1",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный id пользователя",
			},
			ExpectedOutput:         nil,
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
		},
		{
			Name:   "Невалидный номер страницы",
			UserID: "1",
			Page:   "ogo!",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный номер страницы",
			},
			ExpectedOutput:         nil,
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
		},
		{
			Name:   "Пользователь не найден",
			UserID: "1",
			Page:   "1",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "пользователь не найден",
			},
			ExpectedOutput: nil,
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetUserReviews(1, 10, 1).Return(nil, entity.NewClientError("пользователь не найден", entity.ErrNotFound))
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
			mockReviewUsecase := mockusecase.NewMockReview(ctrl)
			tc.SetupReviewUsecaseMock(mockReviewUsecase)
			reviewHandler := NewReviewEndpoints(mockReviewUsecase, nil)
			req := httptest.NewRequest(http.MethodGet, "/review/user", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/review/user/:id/:page")
			c.SetParamNames("id", "page")
			c.SetParamValues(tc.UserID, tc.Page)
			err := reviewHandler.GetUserReviews(c)
			require.Equal(t, tc.ExpectedErr, err)
			if tc.ExpectedErr == nil {
				var reviews dto.ReviewResponseList
				err = json.NewDecoder(rec.Body).Decode(&reviews)
				require.NoError(t, err)
				require.Equal(t, *tc.ExpectedOutput, reviews)
			}
		})
	}
}

func TestReviewEndpoints_GetContentReviews(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		ContentID              string
		Page                   string
		ExpectedErr            error
		ExpectedOutput         *dto.ReviewResponseList
		SetupReviewUsecaseMock func(usecase *mockusecase.MockReview)
	}{
		{
			Name:        "Успешное получение отзывов",
			ContentID:   "1",
			Page:        "1",
			ExpectedErr: nil,
			ExpectedOutput: &dto.ReviewResponseList{
				Reviews: []dto.ReviewResponse{
					{
						Review: dto.Review{
							ID:        1,
							AuthorID:  1,
							ContentID: 1,
							Rating:    5,
							Title:     "Title",
							Text:      "i like it",
							CreatedAt: "2022-01-02T15:04:05Z",
							Likes:     5,
							Dislikes:  5,
						},
						AuthorName:   "Author",
						AuthorAvatar: "avatars/avatar.jpg",
						ContentName:  "Content",
					},
				},
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetContentReviews(1, 10, 1).Return(&dto.ReviewResponseList{
					Reviews: []dto.ReviewResponse{
						{
							Review: dto.Review{
								ID:        1,
								AuthorID:  1,
								ContentID: 1,
								Rating:    5,
								Title:     "Title",
								Text:      "i like it",
								CreatedAt: "2022-01-02T15:04:05Z",
								Likes:     5,
								Dislikes:  5,
							},
							AuthorName:   "Author",
							AuthorAvatar: "avatars/avatar.jpg",
							ContentName:  "Content",
						},
					},
				}, nil)
			},
		},
		{
			Name:      "Неожиданная ошибка",
			ContentID: "1",
			Page:      "1",
			ExpectedErr: &echo.HTTPError{
				Code:    500,
				Message: "internal",
			},
			ExpectedOutput: nil,
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetContentReviews(1, 10, 1).Return(nil, entity.NewClientError("internal", entity.ErrInternal))
			},
		},
		{
			Name:      "Невалидный айди",
			ContentID: "ogo!",
			Page:      "1",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный id контента",
			},
			ExpectedOutput:         nil,
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
		},
		{
			Name:      "Невалидный номер страницы",
			ContentID: "1",
			Page:      "ogo!",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный номер страницы",
			},
			ExpectedOutput:         nil,
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
		},
		{
			Name:      "Контент не найден",
			ContentID: "1",
			Page:      "1",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "контент не найден",
			},
			ExpectedOutput: nil,
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().GetContentReviews(1, 10, 1).Return(nil, entity.NewClientError("контент не найден", entity.ErrNotFound))
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
			mockReviewUsecase := mockusecase.NewMockReview(ctrl)
			tc.SetupReviewUsecaseMock(mockReviewUsecase)
			reviewHandler := NewReviewEndpoints(mockReviewUsecase, nil)
			req := httptest.NewRequest(http.MethodGet, "/review/content", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/review/content/:id/:page")
			c.SetParamNames("id", "page")
			c.SetParamValues(tc.ContentID, tc.Page)
			err := reviewHandler.GetContentReviews(c)
			require.Equal(t, tc.ExpectedErr, err)
			if tc.ExpectedErr == nil {
				var reviews dto.ReviewResponseList
				err = json.NewDecoder(rec.Body).Decode(&reviews)
				require.NoError(t, err)
				require.Equal(t, *tc.ExpectedOutput, reviews)
			}
		})
	}
}

func TestReviewEndpoints_LikeReview(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		ReviewID               string
		ExpectedErr            error
		Cookies                *http.Cookie
		SetupReviewUsecaseMock func(usecase *mockusecase.MockReview)
		SetupAuthUsecaseMock   func(usecase *mockusecase.MockAuth)
	}{
		{
			Name:        "Успешный лайк на отзыв",
			ReviewID:    "1",
			ExpectedErr: nil,
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().LikeReview(1, 1).Return(nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:     "Неожиданная ошибка",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    500,
				Message: "internal",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().LikeReview(1, 1).Return(entity.NewClientError("internal", entity.ErrInternal))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:     "Невалидный айди",
			ReviewID: "ogo!",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный id рецензии",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock:   func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name:     "Отзыв не найден",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "отзыв не найден",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().LikeReview(1, 1).Return(entity.NewClientError("отзыв не найден", entity.ErrNotFound))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:     "Пользователь не авторизован",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    401,
				Message: "необходима авторизация",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession(gomock.Any()).Return(-1, entity.NewClientError("необходима авторизация", entity.ErrUnauthorized))
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
			mockReviewUsecase := mockusecase.NewMockReview(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			tc.SetupReviewUsecaseMock(mockReviewUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			reviewHandler := NewReviewEndpoints(mockReviewUsecase, mockAuthUsecase)
			req := httptest.NewRequest(http.MethodPost, "/review/like", nil)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/review/like/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.ReviewID)
			err := reviewHandler.LikeReview(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewEndpoints_DislikeReview(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		ReviewID               string
		ExpectedErr            error
		Cookies                *http.Cookie
		SetupReviewUsecaseMock func(usecase *mockusecase.MockReview)
		SetupAuthUsecaseMock   func(usecase *mockusecase.MockAuth)
	}{
		{
			Name:        "Успешный дизлайк на отзыв",
			ReviewID:    "1",
			ExpectedErr: nil,
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().DislikeReview(1, 1).Return(nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:     "Неожиданная ошибка",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    500,
				Message: "internal",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().DislikeReview(1, 1).Return(entity.NewClientError("internal", entity.ErrInternal))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:     "Невалидный айди",
			ReviewID: "ogo!",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный id рецензии",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock:   func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name:     "Отзыв не найден",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "отзыв не найден",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().DislikeReview(1, 1).Return(entity.NewClientError("отзыв не найден", entity.ErrNotFound))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:     "Пользователь не авторизован",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    401,
				Message: "необходима авторизация",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession(gomock.Any()).Return(-1, entity.NewClientError("необходима авторизация", entity.ErrUnauthorized))
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
			mockReviewUsecase := mockusecase.NewMockReview(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			tc.SetupReviewUsecaseMock(mockReviewUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			reviewHandler := NewReviewEndpoints(mockReviewUsecase, mockAuthUsecase)
			req := httptest.NewRequest(http.MethodPost, "/review/dislike", nil)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/review/dislike/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.ReviewID)
			err := reviewHandler.DislikeReview(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewEndpoints_UnlikeReview(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                   string
		ReviewID               string
		ExpectedErr            error
		Cookies                *http.Cookie
		SetupReviewUsecaseMock func(usecase *mockusecase.MockReview)
		SetupAuthUsecaseMock   func(usecase *mockusecase.MockAuth)
	}{
		{
			Name:        "Успешное удаление лайка с отзыва",
			ReviewID:    "1",
			ExpectedErr: nil,
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().UnlikeReview(1, 1).Return(nil)
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:     "Неожиданная ошибка",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    500,
				Message: "internal",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().UnlikeReview(1, 1).Return(entity.NewClientError("internal", entity.ErrInternal))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:     "Невалидный айди",
			ReviewID: "ogo!",
			ExpectedErr: &echo.HTTPError{
				Code:    400,
				Message: "невалидный id рецензии",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock:   func(usecase *mockusecase.MockAuth) {},
		},
		{
			Name:     "Отзыв не найден",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    404,
				Message: "отзыв не найден",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {
				usecase.EXPECT().UnlikeReview(1, 1).Return(entity.NewClientError("отзыв не найден", entity.ErrNotFound))
			},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession("xxx").Return(1, nil)
			},
		},
		{
			Name:     "Пользователь не авторизован",
			ReviewID: "1",
			ExpectedErr: &echo.HTTPError{
				Code:    401,
				Message: "необходима авторизация",
			},
			Cookies: &http.Cookie{
				Name:  "session",
				Value: "xxx",
			},
			SetupReviewUsecaseMock: func(usecase *mockusecase.MockReview) {},
			SetupAuthUsecaseMock: func(usecase *mockusecase.MockAuth) {
				usecase.EXPECT().GetUserIDBySession(gomock.Any()).Return(-1, entity.NewClientError("необходима авторизация", entity.ErrUnauthorized))
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
			mockReviewUsecase := mockusecase.NewMockReview(ctrl)
			mockAuthUsecase := mockusecase.NewMockAuth(ctrl)
			tc.SetupReviewUsecaseMock(mockReviewUsecase)
			tc.SetupAuthUsecaseMock(mockAuthUsecase)
			reviewHandler := NewReviewEndpoints(mockReviewUsecase, mockAuthUsecase)
			req := httptest.NewRequest(http.MethodPost, "/review/unlike", nil)
			if tc.Cookies != nil {
				req.AddCookie(tc.Cookies)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/review/unlike/:id")
			c.SetParamNames("id")
			c.SetParamValues(tc.ReviewID)
			err := reviewHandler.UnlikeReview(c)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
