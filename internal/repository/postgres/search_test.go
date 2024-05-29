package postgres

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestSearchDB_SearchContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Query       string
		ExpectedErr error
		ExpectedOut []int
		SetupMock   func(mock sqlmock.Sqlmock, query string)
	}{
		{
			Name:        "Успешный поиск",
			Query:       "Query",
			ExpectedOut: []int{1},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("Query", "Query").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			dbx := sqlx.NewDb(db, "sqlmock")
			if err != nil {
				t.Fatalf("не удалось создать мок: %s", err)
			}
			mock.MatchExpectationsInOrder(false)
			mockContent := NewContentRepository(dbx)
			repo := NewSearchRepository(dbx, mockContent)
			query, _, err := sq.Select("id").
				From("content").
				Where(sq.Or{
					sq.Expr("word_similarity(title, ?) > 0.3", tc.Query),
					sq.Expr("word_similarity(original_title, ?) > 0.3", tc.Query),
				}).
				OrderBy(fmt.Sprintf(
					`CASE WHEN word_similarity(title, '%s') > 0.3
THEN similarity(title, '%s') ELSE similarity(original_title, '%s') END DESC`, tc.Query, tc.Query, tc.Query)).
				Limit(5).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			if err != nil {
				t.Fatalf("ошибка при составлении запроса SearchContent: %s", err)
			}
			tc.SetupMock(mock, query)
			contents, err := repo.SearchContent(tc.Query)
			require.Equal(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOut, contents)
		})
	}
}

func TestSearchDB_SearchPerson(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Query       string
		ExpectedErr error
		ExpectedOut []int
		SetupMock   func(mock sqlmock.Sqlmock, query string)
	}{
		{
			Name:        "Успешный поиск",
			Query:       "Query",
			ExpectedOut: []int{1},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("Query", "Query").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			dbx := sqlx.NewDb(db, "sqlmock")
			require.NoError(t, err)
			mock.MatchExpectationsInOrder(false)
			mockContent := NewContentRepository(dbx)
			repo := NewSearchRepository(dbx, mockContent)
			query, _, err := sq.Select("id").
				From("person").
				Where(sq.Or{
					sq.Expr("word_similarity(Name, ?) > 0.3", tc.Query),
					sq.Expr("word_similarity(en_name, ?) > 0.3", tc.Query),
				}).
				OrderBy(fmt.Sprintf(
					`CASE WHEN word_similarity(Name, '%s') > 0.3
THEN similarity(Name, '%s') ELSE similarity(en_name, '%s') END DESC`, tc.Query, tc.Query, tc.Query)).
				Limit(5).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			if err != nil {
				t.Fatalf("ошибка при составлении запроса SearchPerson: %s", err)
			}
			tc.SetupMock(mock, query)
			persons, err := repo.SearchPerson(tc.Query)
			require.Equal(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOut, persons)
		})
	}
}
