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

func TestCompilationTypeDB_GetCompilationType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		RequestID      int
		ExpectedOutput *entity.CompilationType
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:           "Успех",
			RequestID:      1,
			ExpectedOutput: &entity.CompilationType{ID: 1, Type: "TestType"},
			ExpectedErr:    nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"id", "type"}).
						AddRow(1, "TestType"))
			},
		},
		{
			Name:           "Не найден",
			RequestID:      1,
			ExpectedOutput: nil,
			ExpectedErr:    entity.NewClientError("категория подборки не найдена", entity.ErrNotFound),
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
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении категории")),
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
			repo := &CompilationTypeDB{
				DB: db,
			}
			query, args, err := sq.Select("id", "type").
				From("compilation_type").
				Where(sq.Eq{"id": tc.RequestID}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetCompilationType(tc.RequestID)
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestCompilationTypeDB_GetAllCompilationTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name           string
		ExpectedOutput []*entity.CompilationType
		ExpectedErr    error
		SetupMock      func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name: "Успешное получение",
			ExpectedOutput: []*entity.CompilationType{
				{ID: 1, Type: "TestType1"},
				{ID: 2, Type: "TestType2"},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "type"}).
						AddRow(1, "TestType1").
						AddRow(2, "TestType2"))
			},
		},
		{
			Name:           "Не найдено",
			ExpectedOutput: []*entity.CompilationType{},
			ExpectedErr:    nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(sqlmock.NewRows([]string{}))
			},
		},
		{
			Name:           "Неизвестная ошибка",
			ExpectedOutput: nil,
			ExpectedErr:    entity.PSQLWrap(fmt.Errorf("ошибка"), errors.New("ошибка при получении категорий")),
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
			repo := &CompilationTypeDB{
				DB: db,
			}
			query, args, err := sq.Select("id", "type").
				From("compilation_type").
				OrderBy("type").
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetAllCompilationTypes()
			require.Equal(t, tc.ExpectedOutput, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
