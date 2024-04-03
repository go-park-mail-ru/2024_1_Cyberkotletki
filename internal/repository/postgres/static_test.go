package postgres

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestUsersDB_GetStatic(t *testing.T) {
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
			expectedErr:    entity.NewClientError("файл не найден", entity.ErrNotFound),
			setupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Создаем мок подключения к базе данных
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("не удалось создать мок: %s", err)
			}
			// Создаем репозиторий с моком вместо реального подключения
			repo := &StaticDB{
				DB: db,
			}
			// Ожидаемый запрос
			query, _, err := sq.
				Select("path", "name").
				From("static").
				Where(sq.Eq{"id": tc.requestId}).
				PlaceholderFormat(sq.Dollar).ToSql()
			if err != nil {
				t.Fatalf("ошибка при формировании sql-запроса GetStatic: %s", err)
			}
			tc.setupMock(mock, query)
			// Вызываем тестируемый метод
			result, err := repo.GetStatic(tc.requestId)
			// Проверяем, что запросы выполнен без ошибок, если они не ожидались, и наоборот
			require.Equal(t, tc.expectedErr, err)
			// Проверяем, что результат соответствует ожидаемому
			require.Equal(t, tc.expectedOutput, result)
		})
	}
}
