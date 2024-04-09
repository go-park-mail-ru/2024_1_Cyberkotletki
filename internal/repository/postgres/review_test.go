package postgres

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
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
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
			},
			ExpectedOutput: &entity.Review{
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
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
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    11,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.NewClientError("указаны некорректные данные", entity.ErrBadRequest),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: entity.PSQLCheckViolation})
			},
		},
		{
			Name: "Отзыв уже добавлен",
			RequestReview: &entity.Review{
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
				UpdatedAt: fixedTime,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.NewClientError("рецензия уже существует", entity.ErrAlreadyExists),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: entity.PSQLUniqueViolation})
			},
		},
		{
			Name: "Несуществующий content_id",
			RequestReview: &entity.Review{
				ID:        1,
				AuthorID:  1,
				ContentID: -1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.NewClientError("контент с таким id не существует, либо такого пользователя не существует", entity.ErrNotFound),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: entity.PSQLForeignKeyViolation})
			},
		},
		{
			Name: "Несуществующий user_id",
			RequestReview: &entity.Review{
				ID:        1,
				AuthorID:  -1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.NewClientError("контент с таким id не существует, либо такого пользователя не существует", entity.ErrNotFound),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: entity.PSQLForeignKeyViolation})
			},
		},
		{
			Name: "Неизвестная ошибка postgres",
			RequestReview: &entity.Review{
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLWrap(&pq.Error{Code: "42P01"}, errors.New("ошибка при добавлении рецензии")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: "42P01"}) // или любой другой непредусмотренный код ошибки
			},
		},
		{
			Name: "Неизвестная не-postgres ошибка",
			RequestReview: &entity.Review{
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при добавлении рецензии")),
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
				Values(tc.RequestReview.AuthorID, tc.RequestReview.ContentID, tc.RequestReview.Title, tc.RequestReview.Text, tc.RequestReview.Rating).
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
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Likes:     5,
				Dislikes:  100,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "content_id", "title", "text", "content_rating", "created_at", "updated_at", "likes", "dislikes"}).
						AddRow(1, 1, 1, "title", "text", 5, fixedTime, fixedTime, 5, 100))
			},
		},
		{
			Name:           "Несуществующая рецензия",
			RequestID:      1,
			ExpectedOutput: nil,
			ExpectedErr:    entity.NewClientError("рецензия не найдена", entity.ErrNotFound),
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
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении рецензии")),
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
			query, args, err := sq.Select(
				"r.id",
				"r.user_id",
				"r.content_id",
				"r.title",
				"r.text",
				"r.content_rating",
				"r.created_at",
				"r.updated_at",
				"likes",
				"dislikes",
			).
				Column(sq.Expr("SUM(CASE WHEN rl.value THEN 1 ELSE 0 END) as likes")).
				Column(sq.Expr("SUM(CASE WHEN rl.value THEN 0 ELSE 1 END) as dislikes")).
				From("review r").
				LeftJoin("review_like rl ON r.id = rl.review_id").
				Where(sq.Eq{"r.id": tc.RequestID}).
				GroupBy("r.id").
				PlaceholderFormat(sq.Dollar).
				ToSql()
			fmt.Println(query)
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
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Title:     "title",
					Text:      "text",
					Rating:    5,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Likes:     10,
					Dislikes:  100,
				},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{
						"r.id",
						"r.user_id",
						"r.content_id",
						"r.title",
						"r.text",
						"r.content_rating",
						"r.created_at",
						"r.updated_at",
						"likes",
						"dislikes",
					}).
						AddRows([]driver.Value{1, 1, 1, "title", "text", 5, fixedTime, fixedTime, 10, 100}))
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
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении рецензий")),
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
			query, args, err := sq.Select(
				"r.id",
				"r.user_id",
				"r.content_id",
				"r.title",
				"r.text",
				"r.content_rating",
				"r.created_at",
				"r.updated_at",
				"likes",
				"dislikes",
			).
				Column(sq.Expr("SUM(CASE WHEN rl.value THEN 1 ELSE 0 END) as likes")).
				Column(sq.Expr("SUM(CASE WHEN rl.value THEN 0 ELSE 1 END) as dislikes")).
				Column(sq.Expr("SUM(CASE WHEN rl.value THEN 1 ELSE -1 END) as rating")).
				From("review r").
				LeftJoin("review_like rl ON r.id = rl.review_id").
				Where(sq.Eq{"r.content_id": tc.Request.ContentID}).
				GroupBy("r.id").
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
		Name           string
		RequestReview  *entity.Review
		ExpectedOutput *entity.Review
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name: "Успешное обновление",
			RequestReview: &entity.Review{
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
			},
			ExpectedOutput: &entity.Review{
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnResult(driver.ResultNoRows)
			},
		},
		{
			Name: "Невалидные данные",
			RequestReview: &entity.Review{
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    11,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.NewClientError("указаны некорректные данные", entity.ErrBadRequest),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: entity.PSQLCheckViolation})
			},
		},
		{
			Name: "Рецензия или автор рецензии не найдены",
			RequestReview: &entity.Review{
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.NewClientError("контент с таким id не существует, либо такого пользователя не существует", entity.ErrNotFound),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: entity.PSQLForeignKeyViolation})
			},
		},
		{
			Name: "Неизвестная ошибка postgres",
			RequestReview: &entity.Review{
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLWrap(&pq.Error{Code: "42P01"}, errors.New("ошибка при обновлении рецензии")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(&pq.Error{Code: "42P01"}) // или любой другой непредусмотренный код ошибки
			},
		},
		{
			Name: "Неизвестная не-postgres ошибка",
			RequestReview: &entity.Review{
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
			},
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при обновлении рецензии")),
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
				Set("content_id", tc.RequestReview.ContentID).
				Set("title", tc.RequestReview.Title).
				Set("text", tc.RequestReview.Text).
				Set("content_rating", tc.RequestReview.Rating).
				Where(sq.Eq{"id": tc.RequestReview.ID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.UpdateReview(tc.RequestReview)
			require.Equal(t, tc.ExpectedOutput, output)
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
			Name:        "Неизвестная ошибка",
			RequestID:   1,
			ExpectedErr: entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при удалении рецензии")),
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
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Title:     "title",
					Text:      "text",
					Rating:    5,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					Likes:     10,
					Dislikes:  100,
				},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{
						"r.id",
						"r.user_id",
						"r.content_id",
						"r.title",
						"r.text",
						"r.content_rating",
						"r.created_at",
						"r.updated_at",
						"likes",
						"dislikes",
					}).
						AddRows([]driver.Value{1, 1, 1, "title", "text", 5, fixedTime, fixedTime, 10, 100}))
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
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении рецензий")),
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
			query, args, err := sq.Select(
				"r.id",
				"r.user_id",
				"r.content_id",
				"r.title",
				"r.text",
				"r.content_rating",
				"r.created_at",
				"r.updated_at",
				"likes",
				"dislikes",
			).
				Column(sq.Expr("SUM(CASE WHEN rl.value THEN 1 ELSE 0 END) as likes")).
				Column(sq.Expr("SUM(CASE WHEN rl.value THEN 0 ELSE 1 END) as dislikes")).
				Column(sq.Expr("SUM(CASE WHEN rl.value THEN 1 ELSE -1 END) as rating")).
				From("review r").
				LeftJoin("review_like rl ON r.id = rl.review_id").
				Where(sq.Eq{"r.user_id": tc.Request.AuthorID}).
				GroupBy("r.id").
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
				ID:        1,
				AuthorID:  1,
				ContentID: 1,
				Title:     "title",
				Text:      "text",
				Rating:    5,
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
				Likes:     10,
				Dislikes:  100,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{
						"r.id",
						"r.user_id",
						"r.content_id",
						"r.title",
						"r.text",
						"r.content_rating",
						"r.created_at",
						"r.updated_at",
						"likes",
						"dislikes",
					}).
						AddRows([]driver.Value{1, 1, 1, "title", "text", 5, fixedTime, fixedTime, 10, 100}))
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
			ExpectedErr:    entity.NewClientError("рецензия не найдена", entity.ErrNotFound),
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
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении рецензии")),
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
			query, args, err := sq.Select(
				"r.id",
				"r.user_id",
				"r.content_id",
				"r.title",
				"r.text",
				"r.content_rating",
				"r.created_at",
				"r.updated_at",
				"likes",
				"dislikes",
			).
				Column(sq.Expr("SUM(CASE WHEN rl.value THEN 1 ELSE 0 END) as likes")).
				Column(sq.Expr("SUM(CASE WHEN rl.value THEN 0 ELSE 1 END) as dislikes")).
				From("review r").
				LeftJoin("review_like rl ON r.id = rl.review_id").
				Where(sq.Eq{"r.user_id": tc.Request.AuthorID, "r.content_id": tc.Request.ContentID}).
				GroupBy("r.id").
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

func TestReviewDB_GetAuthorRating(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		RequestID   int
		Expected    int
		ExpectedErr error
		SetupMock   func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:        "Успешное получение",
			RequestID:   1,
			Expected:    5,
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"rating"}).AddRow(5))
			},
		},
		{
			Name:        "Неизвестная ошибка",
			RequestID:   1,
			Expected:    0,
			ExpectedErr: entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении рейтинга автора")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnError(fmt.Errorf("ошибка"))
			},
		},
		{
			Name:        "Автора никто не оценивал",
			RequestID:   1,
			Expected:    0,
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(sql.ErrNoRows)
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
			query, args, err := sq.Select("r.user_id").
				Column(sq.Expr("SUM(CASE WHEN rl.value THEN 1 ELSE -1 END) as user_rating")).
				From("review r").
				LeftJoin("review_like rl ON r.id = rl.review_id").
				Where(sq.Eq{"r.user_id": tc.RequestID}).
				GroupBy("r.user_id").
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetAuthorRating(tc.RequestID)
			require.Equal(t, tc.Expected, output)
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
					ID:        1,
					AuthorID:  1,
					ContentID: 1,
					Title:     "title",
					Text:      "text",
					Rating:    5,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
				},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{
						"r.id",
						"r.user_id",
						"r.content_id",
						"r.title",
						"r.text",
						"r.content_rating",
						"r.created_at",
						"r.updated_at",
					}).
						AddRows([]driver.Value{1, 1, 1, "title", "text", 5, fixedTime, fixedTime}))
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
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении рецензий")),
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
			query, args, err := sq.Select("id", "user_id", "content_id", "title", "text", "content_rating", "created_at", "updated_at").
				From("review").
				OrderBy("id DESC").
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

func TestReviewDB_LikeReview(t *testing.T) {
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
			ExpectedErr: entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при добавлении лайка")),
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
			query, args, err := sq.Insert("review_like").
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
			err = repo.LikeReview(tc.Request.ReviewID, tc.Request.UserID, tc.Request.Value)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_UnlikeReview(t *testing.T) {
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
			Name: "Успешное удаление лайка",
			Request: struct {
				ReviewID int
				UserID   int
			}{
				ReviewID: 1,
				UserID:   1,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnResult(driver.ResultNoRows)
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
			ExpectedErr: entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при удалении лайка")),
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
			query, args, err := sq.Delete("review_like").
				Where(sq.Eq{"review_id": tc.Request.ReviewID, "user_id": tc.Request.UserID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			err = repo.UnlikeReview(tc.Request.ReviewID, tc.Request.UserID)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_IsLikedByUser(t *testing.T) {
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
			ExpectedErr: entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении лайка")),
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
			repo := &ReviewDB{
				DB: db,
			}
			query, args, err := sq.Select("value").
				From("review_like").
				Where(sq.Eq{"review_id": tc.Request.ReviewID, "user_id": tc.Request.UserID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.IsLikedByUser(tc.Request.ReviewID, tc.Request.UserID)
			require.Equal(t, tc.Expected, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestReviewDB_GetContentRating(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		RequestID   int
		Expected    int
		ExpectedErr error
		SetupMock   func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:        "Успешное получение",
			RequestID:   1,
			Expected:    5,
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"rating"}).AddRow(5))
			},
		},
		{
			Name:        "Неизвестная ошибка",
			RequestID:   1,
			Expected:    0,
			ExpectedErr: entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении рейтинга контента")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnError(fmt.Errorf("ошибка"))
			},
		},
		{
			Name:        "Контента никто не оценивал",
			RequestID:   1,
			Expected:    0,
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(sql.ErrNoRows)
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
			query, args, err := sq.Select("AVG(content_rating)").
				From("review").
				Where(sq.Eq{"content_id": tc.RequestID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetContentRating(tc.RequestID)
			require.Equal(t, tc.Expected, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
