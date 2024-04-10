package postgres

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestCompilationDB_GetCompilationsByTypeID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		RequestID      int
		ExpectedOutput []entity.Compilation
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:           "Успешное получение",
			RequestID:      1,
			ExpectedOutput: []entity.Compilation{{ID: 1, Title: "Test", CompilationTypeID: 1, PosterUploadID: 1}},
			ExpectedErr:    nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "compilation_type_id", "poster_upload_id"}).AddRow(1, "Test", 1, 1))
			},
		},
		{
			Name:           "Несуществующая подборка",
			RequestID:      2,
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLWrap(sql.ErrNoRows, errors.New("ошибка при получении подборок")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Name:           "Неизвестная ошибка",
			RequestID:      3,
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении подборок")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := &CompilationDB{
				DB: db,
			}
			query, args, _ := sq.Select("id", "title", "compilation_type_id", "poster_upload_id").
				From("compilation").
				Where(sq.Eq{"compilation_type_id": tc.RequestID}).
				OrderBy("title ASC").
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetCompilationsByTypeID(tc.RequestID)
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestCompilationDB_GetCompilationContentLength(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		RequestID      int
		ExpectedOutput int
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:           "Успешное получение",
			RequestID:      1,
			ExpectedOutput: 5,
			ExpectedErr:    nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
			},
		},
		{
			Name:           "Несуществующая подборка/ подборка без контента",
			RequestID:      1,
			ExpectedOutput: 0,
			ExpectedErr:    entity.PSQLWrap(sql.ErrNoRows, errors.New("ошибка при получении числа контента в подборке")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Name:           "Неизвестная ошибка",
			RequestID:      1,
			ExpectedOutput: 0,
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении числа контента в подборке")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := &CompilationDB{
				DB: db,
			}
			query, args, _ := sq.Select("count(*)").
				From("compilation_content").
				Where(sq.Eq{"compilation_id": tc.RequestID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetCompilationContentLength(tc.RequestID)
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestCompilationDB_GetCompilationContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		RequestID      int
		Page           int
		Limit          int
		ExpectedOutput []int
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:           "Успешное получение",
			RequestID:      1,
			Page:           1,
			Limit:          5,
			ExpectedOutput: []int{1, 2, 3, 4, 5},
			ExpectedErr:    nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"content_id"}).
						AddRow(1).
						AddRow(2).
						AddRow(3).
						AddRow(4).
						AddRow(5))
			},
		},
		{
			Name:           "Подборка без контента",
			RequestID:      1,
			Page:           1,
			Limit:          5,
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLWrap(sql.ErrNoRows, errors.New("контент подборки не найден"), entity.ErrNotFound),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			Name:           "Неизвестная ошибка",
			RequestID:      1,
			Page:           1,
			Limit:          5,
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении контента подборки")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка"))
			},
		},
		{
			Name:           "Страница больше общего числа страниц",
			RequestID:      1,
			Page:           3,
			Limit:          2,
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLWrap(sql.ErrNoRows, errors.New("контент подборки не найден"), entity.ErrNotFound),
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
			repo := &CompilationDB{
				DB: db,
			}
			query, args, _ := sq.Select("content_id").
				From("compilation_content").
				Where(sq.Eq{"compilation_id": tc.RequestID}).
				Limit(uint64(tc.Limit)).
				Offset(uint64((tc.Page - 1) * tc.Limit)).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetCompilationContent(tc.RequestID, tc.Page, tc.Limit)
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestCompilationDB_GetAllCompilationTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		ExpectedOutput []entity.CompilationType
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:           "Успешное получение",
			ExpectedOutput: []entity.CompilationType{{ID: 1, Type: "Test"}},
			ExpectedErr:    nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "type"}).AddRow(1, "Test"))
			},
		},
		{
			Name:           "Неизвестная ошибка",
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении категорий подборок")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(fmt.Errorf("ошибка"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := &CompilationDB{
				DB: db,
			}
			query, _, _ := sq.Select("id", "type").
				From("compilation_type").
				OrderBy("type ASC").
				PlaceholderFormat(sq.Dollar).
				ToSql()
			tc.SetupMock(mock, query, nil)
			output, err := repo.GetAllCompilationTypes()
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
