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

/*func TestCompilationDB_GetCompilationsByCompilationTypeID(t *testing.T) {
	t.Parallel()

	tastCases := []struct {
		Name                   string
		RequestCompilationType struct {
			CompilationType int
			Page            int
			Limit           int
		}
		ExpectedCompilations []*entity.Compilation
		ExpectedErr          error
		SetupMock            func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name: "Успешное получение подборок по ID категории подборок",
			RequestCompilationType: struct {
				CompilationType int
				Page            int
				Limit           int
			}{
				CompilationType: 1,
				Page:            1,
				Limit:           1,
			},
			ExpectedCompilations: []*entity.Compilation{
				{
					ID:                1,
					Title:             "title",
					CompilationTypeID: 1,
					PosterUploadID:    1,
				},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"с.id", "с.title", "с.compilation_type_id", "с.poster"}).
						AddRow(1, "title", 1, 1))
			},
		},
		{
			Name: "Подборок по ID категории подборок нет",
			RequestCompilationType: struct {
				CompilationType int
				Page            int
				Limit           int
			}{
				CompilationType: 1,
				Page:            1,
				Limit:           1,
			},
			ExpectedCompilations: []*entity.Compilation{},
			ExpectedErr:          nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(sqlmock.NewRows([]string{}))
			},
		},
		{
			Name: "Неизвестная ошибка",
			RequestCompilationType: struct {
				CompilationType int
				Page            int
				Limit           int
			}{
				CompilationType: 1,
				Page:            1,
				Limit:           1,
			},
			ExpectedCompilations: nil,
			ExpectedErr:          entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении подборок")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка"))
			},
		},
	}

	for _, to := range tastCases {
		to := to
		t.Run(to.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := &CompilationDB{
				DB: db,
			}
			query, args, err := sq.Select("с.id", "с.compilation_type_id", "с.title", "с.poster").
				From("compilation с").
				Where(sq.Eq{"с.compilation_type_id": to.RequestCompilationType.CompilationType}).
				OrderBy("с.title ASC").
				Limit(uint64(to.RequestCompilationType.Limit)).
				Offset(uint64((to.RequestCompilationType.Page - 1) * to.RequestCompilationType.Limit)).
				PlaceholderFormat(sq.Dollar).ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			to.SetupMock(mock, query, driverValues)
			output, err := repo.GetCompilationsByCompilationTypeID(
				to.RequestCompilationType.CompilationType, to.RequestCompilationType.Page, to.RequestCompilationType.Limit)
			require.Equal(t, to.ExpectedCompilations, output)
			require.Equal(t, to.ExpectedErr, err)
		})
	}
}

func TestCompilationDB_GetCompilation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		RequestID      int
		ExpectedOutput *entity.Compilation
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:      "Успех при получении подборки по ID",
			RequestID: 1,
			ExpectedOutput: &entity.Compilation{
				ID:                1,
				Title:             "title",
				CompilationTypeID: 1,
				PosterUploadID:    1,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "compilation_type_id", "poster_upload_id"}).
						AddRow(1, "title", 1, 1))
			},
		},
		{
			Name:           "Нет подборки по ID",
			RequestID:      1,
			ExpectedOutput: nil,
			ExpectedErr:    entity.NewClientError("подборка не найдена", entity.ErrNotFound),
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
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении подборки")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка")) // or any other unexpected error code
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
			query, args, err := sq.Select("id", "title", "compilation_type_id", "poster_upload_id").
				From("compilation").
				Where(sq.Eq{"id": tc.RequestID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetCompilation(tc.RequestID)
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}*/

func TestCompilationDB_GetCompilationContentLength(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		RequestID      int
		ExpectedLength int
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:           "Успешное получение количества контента в подборке по ID",
			RequestID:      1,
			ExpectedLength: 2,
			ExpectedErr:    nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
			},
		},
		{
			Name:           "Неизвестная ошибка",
			RequestID:      1,
			ExpectedLength: 0,
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
			query, args, err := sq.Select("count(*)").
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
			length, err := repo.GetCompilationContentLength(tc.RequestID)
			require.Equal(t, tc.ExpectedLength, length)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestCompilationDB_GetCompilationContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		RequestID      int
		ExpectedOutput []int
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:           "Успешное получение списка контента в подборке",
			RequestID:      1,
			ExpectedOutput: []int{1, 2},
			ExpectedErr:    nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"content_id"}).
						AddRow(1).
						AddRow(2))
			},
		},
		{
			Name:           "Подборка не найдена",
			RequestID:      1,
			ExpectedOutput: []int{},
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
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении контента подборки")),
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
			query, args, err := sq.Select("content_id").
				From("compilation_content").
				Where(sq.Eq{"compilation_id": tc.RequestID}).
				Limit(uint64(10)).
				Offset(uint64((1 - 1) * 10)).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetCompilationContent(tc.RequestID, 1, 10)
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
