package postgres

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestReviewDB_AddReview(t *testing.T) {
	t.Parallel()

	fixedTime := time.Now()

	testCases := []struct {
		Name           string
		RequestReview  *entity.Review
		ExpectedOutput *entity.Review
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name: "Успешное добавление",
			RequestReview: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 5,
			},
			ExpectedOutput: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 5,
				CreatedAt:     fixedTime,
				UpdatedAt:     fixedTime,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(1, fixedTime, fixedTime))
			},
		},
		{
			Name: "Невалидные данные",
			RequestReview: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 11,
			},
			ExpectedOutput: nil,
			ExpectedErr:    repository.ErrReviewBadRequest,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: entity.PSQLCheckViolation})
			},
		},
		{
			Name: "Рецензия уже существует",
			RequestReview: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 5,
				UpdatedAt:     fixedTime,
			},
			ExpectedOutput: nil,
			ExpectedErr:    repository.ErrReviewAlreadyExists,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: entity.PSQLUniqueViolation})
			},
		},
		{
			Name: "Несуществующий content_id",
			RequestReview: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     -1,
				Title:         "title",
				Text:          "text",
				ContentRating: 5,
			},
			ExpectedOutput: nil,
			ExpectedErr:    repository.ErrReviewViolation,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: entity.PSQLForeignKeyViolation})
			},
		},
		{
			Name: "Несуществующий user_id",
			RequestReview: &entity.Review{
				ID:            1,
				AuthorID:      -1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 5,
			},
			ExpectedOutput: nil,
			ExpectedErr:    repository.ErrReviewViolation,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: entity.PSQLForeignKeyViolation})
			},
		},
		{
			Name: "Неизвестная ошибка postgres",
			RequestReview: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 5,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLQueryErr("AddReview", &pq.Error{Code: "42P01"}),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: "42P01"}) // или любой другой непредусмотренный код ошибки
			},
		},
		{
			Name: "Неизвестная не-postgres ошибка",
			RequestReview: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 5,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLQueryErr("AddReview", fmt.Errorf("ошибка")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка")) // или любой другой непредусмотренный код ошибки
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := &ReviewDB{
				DB: db,
			}
			query, args, err := sq.Insert("review").
				Columns("user_id", "content_id", "title", "text", "content_rating").
				Values(tc.RequestReview.AuthorID, tc.RequestReview.ContentID, tc.RequestReview.Title, tc.RequestReview.Text, tc.RequestReview.ContentRating).
				Suffix("RETURNING id, created_at, updated_at").
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)

			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.AddReview(tc.RequestReview)
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_GetReviewByID(t *testing.T) {
	t.Parallel()

	fixedTime := time.Now()

	testCases := []struct {
		Name           string
		RequestID      int
		ExpectedOutput *entity.Review
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:      "Успешное получение",
			RequestID: 1,
			ExpectedOutput: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 5,
				CreatedAt:     fixedTime,
				UpdatedAt:     fixedTime,
				Rating:        -95,
				Likes:         5,
				Dislikes:      100,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"user_id",
						"content_id",
						"title",
						"text",
						"content_rating",
						"created_at",
						"updated_at",
						"likes",
						"dislikes",
						"rating",
					}).
						AddRow(1, 1, 1, "title", "text", 5, fixedTime, fixedTime, 5, 100, -95))
			},
		},
		{
			Name:           "Несуществующая рецензия",
			RequestID:      1,
			ExpectedOutput: nil,
			ExpectedErr:    repository.ErrReviewNotFound,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Name:           "Неизвестная ошибка",
			RequestID:      1,
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLQueryErr("GetReviewByID", fmt.Errorf("ошибка")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка")) // или любой другой непредусмотренный код ошибки
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := NewReviewRepository(db)
			query, args, err := selectAllFields().
				From("review").
				Where(sq.Eq{"id": tc.RequestID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetReviewByID(tc.RequestID)
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_GetReviewsCountByContentID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name             string
		RequestContentID int
		ExpectedOutput   int
		ExpectedErr      error
		SetupMock        func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:             "Успешное получение",
			RequestContentID: 1,
			ExpectedOutput:   1,
			ExpectedErr:      nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow(1))
			},
		},
		{
			Name:             "Неизвестная ошибка",
			RequestContentID: 1,
			ExpectedOutput:   0,
			ExpectedErr:      entity.PSQLQueryErr("GetReviewsCountByContentID", fmt.Errorf("ошибка")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка")) // или любой другой непредусмотренный код ошибки
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := NewReviewRepository(db)
			query, args, err := sq.Select("COUNT(*)").
				From("review").
				Where(sq.Eq{"content_id": tc.RequestContentID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetReviewsCountByContentID(tc.RequestContentID)
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_GetReviewsByContentID(t *testing.T) {
	t.Parallel()

	fixedTime := time.Now()

	testCases := []struct {
		Name    string
		Request struct {
			ContentID int
			Page      int
			Limit     int
		}
		ExpectedOutput []*entity.Review
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name: "Успешное получение",
			Request: struct {
				ContentID int
				Page      int
				Limit     int
			}{
				ContentID: 1,
				Page:      1,
				Limit:     1,
			},
			ExpectedOutput: []*entity.Review{
				{
					ID:            1,
					AuthorID:      1,
					ContentID:     1,
					Title:         "title",
					Text:          "text",
					ContentRating: 5,
					CreatedAt:     fixedTime,
					UpdatedAt:     fixedTime,
					Likes:         10,
					Dislikes:      100,
					Rating:        -90,
				},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"user_id",
						"content_id",
						"title",
						"text",
						"content_rating",
						"created_at",
						"updated_at",
						"likes",
						"dislikes",
						"rating",
					}).
						AddRows([]driver.Value{1, 1, 1, "title", "text", 5, fixedTime, fixedTime, 10, 100, -90}))
			},
		},
		{
			Name: "Рецензий нет",
			Request: struct {
				ContentID int
				Page      int
				Limit     int
			}{
				ContentID: 1,
				Page:      1,
				Limit:     1,
			},
			ExpectedOutput: []*entity.Review{},
			ExpectedErr:    nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(sqlmock.NewRows([]string{}))
			},
		},
		{
			Name: "Неизвестная ошибка",
			Request: struct {
				ContentID int
				Page      int
				Limit     int
			}{
				ContentID: 1,
				Page:      1,
				Limit:     1,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLQueryErr("GetReviewsByContentID", fmt.Errorf("ошибка")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка")) // или любой другой непредусмотренный код ошибки
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := &ReviewDB{
				DB: db,
			}
			query, args, err := selectAllFields().
				From("review").
				Where(sq.Eq{"content_id": tc.Request.ContentID}).
				OrderBy("rating DESC").
				Limit(uint64(tc.Request.Limit)).
				Offset(uint64((tc.Request.Page - 1) * tc.Request.Limit)).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetReviewsByContentID(tc.Request.ContentID, tc.Request.Page, tc.Request.Limit)
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_UpdateReview(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		RequestReview *entity.Review
		ExpectedErr   error
		SetupMock     func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name: "Успешное обновление",
			RequestReview: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 5,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			Name: "Невалидные данные",
			RequestReview: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 11,
			},
			ExpectedErr: repository.ErrReviewBadRequest,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: entity.PSQLCheckViolation})
			},
		},
		{
			Name: "Рецензия или автор рецензии не найдены",
			RequestReview: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 5,
			},
			ExpectedErr: repository.ErrReviewNotFound,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: entity.PSQLForeignKeyViolation})
			},
		},
		{
			Name: "Неизвестная ошибка postgres",
			RequestReview: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 5,
			},
			ExpectedErr: entity.PSQLQueryErr("UpdateReview", &pq.Error{Code: "42P01"}),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: "42P01"}) // или любой другой непредусмотренный код ошибки
			},
		},
		{
			Name: "Неизвестная не-postgres ошибка",
			RequestReview: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 5,
			},
			ExpectedErr: entity.PSQLQueryErr("UpdateReview", fmt.Errorf("ошибка")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка")) // или любой другой непредусмотренный код ошибки
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := &ReviewDB{
				DB: db,
			}
			query, args, err := sq.Update("review").
				Set("title", tc.RequestReview.Title).
				Set("text", tc.RequestReview.Text).
				Set("content_rating", tc.RequestReview.ContentRating).
				Where(sq.Eq{"id": tc.RequestReview.ContentID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			err = repo.UpdateReview(tc.RequestReview)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_DeleteReviewByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		RequestID   int
		ExpectedErr error
		SetupMock   func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:        "Успешное удаление",
			RequestID:   1,
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnResult(driver.ResultNoRows)
			},
		},
		{
			Name:        "Рецензия не найдена",
			RequestID:   1,
			ExpectedErr: repository.ErrReviewNotFound,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Name:        "Неизвестная ошибка",
			RequestID:   1,
			ExpectedErr: entity.PSQLQueryErr("DeleteReviewByID", fmt.Errorf("ошибка")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnError(fmt.Errorf("ошибка"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := &ReviewDB{
				DB: db,
			}
			query, args, err := sq.Delete("review").
				Where(sq.Eq{"id": tc.RequestID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			err = repo.DeleteReviewByID(tc.RequestID)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_GetReviewsByAuthorID(t *testing.T) {
	t.Parallel()

	fixedTime := time.Now()

	testCases := []struct {
		Name    string
		Request struct {
			AuthorID int
			Page     int
			Limit    int
		}
		ExpectedOutput []*entity.Review
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name: "Успешное получение",
			Request: struct {
				AuthorID int
				Page     int
				Limit    int
			}{
				AuthorID: 1,
				Page:     1,
				Limit:    1,
			},
			ExpectedOutput: []*entity.Review{
				{
					ID:            1,
					AuthorID:      1,
					ContentID:     1,
					Title:         "title",
					Text:          "text",
					ContentRating: 5,
					CreatedAt:     fixedTime,
					UpdatedAt:     fixedTime,
					Likes:         10,
					Dislikes:      100,
					Rating:        -90,
				},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"user_id",
						"content_id",
						"title",
						"text",
						"content_rating",
						"created_at",
						"updated_at",
						"likes",
						"dislikes",
						"rating",
					}).
						AddRows([]driver.Value{1, 1, 1, "title", "text", 5, fixedTime, fixedTime, 10, 100, -90}))
			},
		},
		{
			Name: "Рецензий нет",
			Request: struct {
				AuthorID int
				Page     int
				Limit    int
			}{
				AuthorID: 1,
				Page:     1,
				Limit:    1,
			},
			ExpectedOutput: []*entity.Review{},
			ExpectedErr:    nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(sqlmock.NewRows([]string{}))
			},
		},
		{
			Name: "Неизвестная ошибка",
			Request: struct {
				AuthorID int
				Page     int
				Limit    int
			}{
				AuthorID: 1,
				Page:     1,
				Limit:    1,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLQueryErr("GetReviewsByAuthorID", fmt.Errorf("ошибка")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка")) // или любой другой непредусмотренный код ошибки
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := &ReviewDB{
				DB: db,
			}
			query, args, err := selectAllFields().
				From("review").
				Where(sq.Eq{"user_id": tc.Request.AuthorID}).
				OrderBy("created_at DESC").
				Limit(uint64(tc.Request.Limit)).
				Offset(uint64((tc.Request.Page - 1) * tc.Request.Limit)).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetReviewsByAuthorID(tc.Request.AuthorID, tc.Request.Page, tc.Request.Limit)
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_GetContentReviewByAuthor(t *testing.T) {
	t.Parallel()

	fixedTime := time.Now()

	testCases := []struct {
		Name    string
		Request struct {
			AuthorID  int
			ContentID int
		}
		ExpectedOutput *entity.Review
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name: "Успешное получение",
			Request: struct {
				AuthorID  int
				ContentID int
			}{
				AuthorID:  1,
				ContentID: 1,
			},
			ExpectedOutput: &entity.Review{
				ID:            1,
				AuthorID:      1,
				ContentID:     1,
				Title:         "title",
				Text:          "text",
				ContentRating: 5,
				CreatedAt:     fixedTime,
				UpdatedAt:     fixedTime,
				Likes:         10,
				Dislikes:      100,
				Rating:        -90,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"user_id",
						"content_id",
						"title",
						"text",
						"content_rating",
						"created_at",
						"updated_at",
						"likes",
						"dislikes",
						"rating",
					}).
						AddRows([]driver.Value{1, 1, 1, "title", "text", 5, fixedTime, fixedTime, 10, 100, -90}))
			},
		},
		{
			Name: "Рецензия не найдена",
			Request: struct {
				AuthorID  int
				ContentID int
			}{
				AuthorID:  1,
				ContentID: 1,
			},
			ExpectedOutput: nil,
			ExpectedErr:    repository.ErrReviewNotFound,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Name: "Неизвестная ошибка",
			Request: struct {
				AuthorID  int
				ContentID int
			}{
				AuthorID:  1,
				ContentID: 1,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLQueryErr("GetContentReviewByAuthor", fmt.Errorf("ошибка")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка")) // или любой другой непредусмотренный код ошибки
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := &ReviewDB{
				DB: db,
			}
			query, args, err := selectAllFields().
				From("review").
				Where(sq.Eq{"user_id": tc.Request.AuthorID, "content_id": tc.Request.ContentID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetContentReviewByAuthor(tc.Request.AuthorID, tc.Request.ContentID)
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_GetLatestReviews(t *testing.T) {
	t.Parallel()

	fixedTime := time.Now()

	testCases := []struct {
		Name           string
		RequestLimit   int
		ExpectedOutput []*entity.Review
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:         "Успешное получение",
			RequestLimit: 1,
			ExpectedOutput: []*entity.Review{
				{
					ID:            1,
					AuthorID:      1,
					ContentID:     1,
					Title:         "title",
					Text:          "text",
					ContentRating: 5,
					CreatedAt:     fixedTime,
					UpdatedAt:     fixedTime,
					Likes:         10,
					Dislikes:      100,
					Rating:        -90,
				},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{
						"id",
						"user_id",
						"content_id",
						"title",
						"text",
						"content_rating",
						"created_at",
						"updated_at",
						"likes",
						"dislikes",
						"rating",
					}).
						AddRows([]driver.Value{1, 1, 1, "title", "text", 5, fixedTime, fixedTime, 10, 100, -90}))
			},
		},
		{
			Name:           "Рецензий нет",
			RequestLimit:   1,
			ExpectedOutput: []*entity.Review{},
			ExpectedErr:    nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(sqlmock.NewRows([]string{}))
			},
		},
		{
			Name:           "Неизвестная ошибка",
			RequestLimit:   1,
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLQueryErr("GetLatestReviews", fmt.Errorf("ошибка")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка")) // или любой другой непредусмотренный код ошибки
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := &ReviewDB{
				DB: db,
			}
			query, args, err := selectAllFields().
				From("review").
				OrderBy("created_at DESC").
				Limit(uint64(tc.RequestLimit)).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetLatestReviews(tc.RequestLimit)
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_VoteReview(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name    string
		Request struct {
			ReviewID int
			UserID   int
			Value    bool
		}
		ExpectedErr error
		SetupMock   func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name: "Успешное добавление лайка",
			Request: struct {
				ReviewID int
				UserID   int
				Value    bool
			}{
				ReviewID: 1,
				UserID:   1,
				Value:    true,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnResult(driver.ResultNoRows)
			},
		},
		{
			Name: "Успешное добавление дизлайка",
			Request: struct {
				ReviewID int
				UserID   int
				Value    bool
			}{
				ReviewID: 1,
				UserID:   1,
				Value:    false,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnResult(driver.ResultNoRows)
			},
		},
		{
			Name: "Несуществующая рецензия",
			Request: struct {
				ReviewID int
				UserID   int
				Value    bool
			}{
				ReviewID: 1,
				UserID:   1,
				Value:    false,
			},
			ExpectedErr: repository.ErrReviewNotFound,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnError(&pq.Error{Code: entity.PSQLForeignKeyViolation})
			},
		},
		{
			Name: "Оценка уже стоит",
			Request: struct {
				ReviewID int
				UserID   int
				Value    bool
			}{
				ReviewID: 1,
				UserID:   1,
				Value:    false,
			},
			ExpectedErr: repository.ErrReviewVoteAlreadyExists,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnError(&pq.Error{Code: entity.PSQLUniqueViolation})
			},
		},
		{
			Name: "Неизвестная ошибка",
			Request: struct {
				ReviewID int
				UserID   int
				Value    bool
			}{
				ReviewID: 1,
				UserID:   1,
				Value:    true,
			},
			ExpectedErr: entity.PSQLQueryErr("VoteReview", fmt.Errorf("ошибка")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка")) // или любой другой непредусмотренный код ошибки
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := &ReviewDB{
				DB: db,
			}
			query, args, err := sq.Insert("review_vote").
				Columns("review_id", "user_id", "value").
				Values(tc.Request.ReviewID, tc.Request.UserID, tc.Request.Value).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			err = repo.VoteReview(tc.Request.ReviewID, tc.Request.UserID, tc.Request.Value)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_UnVoteReview(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name    string
		Request struct {
			ReviewID int
			UserID   int
		}
		ExpectedErr error
		SetupMock   func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name: "Успешное удаление оценки",
			Request: struct {
				ReviewID int
				UserID   int
			}{
				ReviewID: 1,
				UserID:   1,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			Name: "Оценка на рецензию не найдена",
			Request: struct {
				ReviewID int
				UserID   int
			}{
				ReviewID: 1,
				UserID:   1,
			},
			ExpectedErr: repository.ErrReviewVoteNotFound,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: entity.PSQLForeignKeyViolation})
			},
		},
		{
			Name: "Неизвестная ошибка",
			Request: struct {
				ReviewID int
				UserID   int
			}{
				ReviewID: 1,
				UserID:   1,
			},
			ExpectedErr: entity.PSQLQueryErr("UnVoteReview", fmt.Errorf("ошибка")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка")) // или любой другой непредусмотренный код ошибки
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := NewReviewRepository(db)
			query, args, err := sq.Delete("review_vote").
				Where(sq.Eq{"review_id": tc.Request.ReviewID, "user_id": tc.Request.UserID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			err = repo.UnVoteReview(tc.Request.ReviewID, tc.Request.UserID)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_IsVotedByUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name    string
		Request struct {
			ReviewID int
			UserID   int
		}
		Expected    int
		ExpectedErr error
		SetupMock   func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name: "Стоит лайк",
			Request: struct {
				ReviewID int
				UserID   int
			}{
				ReviewID: 1,
				UserID:   1,
			},
			Expected:    1,
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
		},
		{
			Name: "Стоит дизлайк",
			Request: struct {
				ReviewID int
				UserID   int
			}{
				ReviewID: 1,
				UserID:   1,
			},
			Expected:    -1,
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
		},
		{
			Name: "Лайка нет",
			Request: struct {
				ReviewID int
				UserID   int
			}{
				ReviewID: 1,
				UserID:   1,
			},
			Expected:    0,
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Name: "Неизвестная ошибка",
			Request: struct {
				ReviewID int
				UserID   int
			}{
				ReviewID: 1,
				UserID:   1,
			},
			Expected:    0,
			ExpectedErr: entity.PSQLQueryErr("IsVotedByUser", fmt.Errorf("ошибка")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnError(fmt.Errorf("ошибка"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := NewReviewRepository(db)
			query, args, err := sq.Select("value").
				From("review_vote").
				Where(sq.Eq{"review_id": tc.Request.ReviewID, "user_id": tc.Request.UserID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.IsVotedByUser(tc.Request.ReviewID, tc.Request.UserID)
			require.Equal(t, tc.Expected, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
