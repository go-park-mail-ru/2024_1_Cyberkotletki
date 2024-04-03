package postgres

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestUsersDB_AddUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		user        *entity.User
		expectedErr error
		expectedOut *entity.User
		setupMock   func(mock sqlmock.Sqlmock, query string)
	}{
		{
			name: "Успешное добавление",
			user: &entity.User{
				Email:        "email@mail.ru",
				PasswordHash: []byte("hashed"),
				PasswordSalt: []byte("salt"),
			},
			expectedErr: nil,
			setupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("email@mail.ru", []byte("hashed"), []byte("salt")).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			expectedOut: &entity.User{
				ID:             1,
				Name:           "",
				Email:          "email@mail.ru",
				PasswordHash:   []byte("hashed"),
				PasswordSalt:   []byte("salt"),
				AvatarUploadID: 0,
			},
		},
		{
			name: "Email уже занят",
			user: &entity.User{
				Email:        "email@mail.ru",
				PasswordHash: []byte("hashed"),
				PasswordSalt: []byte("salt"),
			},
			expectedErr: entity.NewClientError("пользователь с таким email уже существует", entity.ErrAlreadyExists),
			setupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("email@mail.ru", []byte("hashed"), []byte("salt")).
					WillReturnError(&pq.Error{Code: entity.PSQLUniqueViolation})
			},
		},
		{
			name: "Нарушение ограничения CHECK",
			user: &entity.User{
				Email:        "тут мог бы быть длинный емейл длинной больше 255 символов, но его нет",
				PasswordHash: []byte("в этом хэше слишком много байт, из-за чего нарушается ограничение CHECK"),
				PasswordSalt: []byte("salt"),
			},
			expectedErr: entity.NewClientError(
				"одно или несколько полей заполнены некорректно",
				entity.ErrBadRequest,
				errors.New("ошибка при выполнении sql-запроса AddUser: нарушение целостности данных"),
			),
			setupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(
						"тут мог бы быть длинный емейл длинной больше 255 символов, но его нет",
						[]byte("в этом хэше слишком много байт, из-за чего нарушается ограничение CHECK"),
						[]byte("salt"),
					).
					WillReturnError(&pq.Error{Code: entity.PSQLCheckViolation})
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
			repo := &UsersDB{
				DB: db,
			}
			// Ожидаемый запрос
			query, _, err := sq.Insert("users").
				Columns("email", "password_hashed", "salt_password").
				Values(tc.user.Email, tc.user.PasswordHash, tc.user.PasswordSalt).
				Suffix("RETURNING id").
				PlaceholderFormat(sq.Dollar).
				ToSql()
			tc.setupMock(mock, query)
			user, err := repo.AddUser(tc.user)
			require.Equal(t, tc.expectedErr, err)
			require.EqualValues(t, tc.expectedOut, user)
		})
	}
}
