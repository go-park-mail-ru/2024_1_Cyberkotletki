package postgres

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestStaticDB_GetStatic(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		requestId      int
		expectedOutput string
		expectedErr    error
		setupMock      func(mock sqlmock.Sqlmock, query string)
	}{
		{
			name:           "Существующий файл",
			requestId:      1,
			expectedOutput: "path/name",
			expectedErr:    nil,
			setupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"path", "name"}).AddRow("path", "name"))
			},
		},
		{
			name:           "Файл не найден",
			requestId:      2,
			expectedOutput: "",
			expectedErr:    repository.ErrStaticNotFound,
			setupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name:           "Неизвестная ошибка",
			requestId:      1,
			expectedOutput: "",
			expectedErr:    entity.PSQLQueryErr("GetStatic", sql.ErrConnDone),
			setupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(1).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("не удалось создать мок: %s", err)
			}
			repo := NewStaticRepository(db, "", 1)
			query, _, err := sq.
				Select("path", "name").
				From("static").
				Where(sq.Eq{"id": tc.requestId}).
				PlaceholderFormat(sq.Dollar).ToSql()
			if err != nil {
				t.Fatalf("ошибка при формировании sql-запроса GetStatic: %s", err)
			}
			tc.setupMock(mock, query)
			result, err := repo.GetStatic(tc.requestId)
			require.Equal(t, tc.expectedErr, err)
			require.Equal(t, tc.expectedOutput, result)
		})
	}
}
