package postgres

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
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
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs("email@mail.ru", []byte("hashed"), []byte("salt")).
					WillReturnResult(sqlmock.NewResult(1, 1))
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
			expectedErr: repository.ErrUserAlreadyExists,
			setupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
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
			expectedErr: repository.ErrUserIncorrectData,
			setupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(
						"тут мог бы быть длинный емейл длинной больше 255 символов, но его нет",
						[]byte("в этом хэше слишком много байт, из-за чего нарушается ограничение CHECK"),
						[]byte("salt"),
					).
					WillReturnError(&pq.Error{Code: entity.PSQLCheckViolation})
			},
		},
		{
			name: "Неизвестная ошибка psql",
			user: &entity.User{
				Email:        "тут мог бы быть длинный емейл длинной больше 255 символов, но его нет",
				PasswordHash: []byte("в этом хэше слишком много байт, из-за чего нарушается ограничение CHECK"),
				PasswordSalt: []byte("salt"),
			},
			expectedErr: entity.PSQLQueryErr("AddUser", &pq.Error{Code: "123"}),
			setupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(
						"тут мог бы быть длинный емейл длинной больше 255 символов, но его нет",
						[]byte("в этом хэше слишком много байт, из-за чего нарушается ограничение CHECK"),
						[]byte("salt"),
					).
					WillReturnError(&pq.Error{Code: "123"})
			},
		},
		{
			name: "Неизвестная ошибка sql",
			user: &entity.User{
				Email:        "тут мог бы быть длинный емейл длинной больше 255 символов, но его нет",
				PasswordHash: []byte("в этом хэше слишком много байт, из-за чего нарушается ограничение CHECK"),
				PasswordSalt: []byte("salt"),
			},
			expectedErr: entity.PSQLQueryErr("AddUser", sql.ErrConnDone),
			setupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(
						"тут мог бы быть длинный емейл длинной больше 255 символов, но его нет",
						[]byte("в этом хэше слишком много байт, из-за чего нарушается ограничение CHECK"),
						[]byte("salt"),
					).
					WillReturnError(sql.ErrConnDone)
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
			repo := NewUserRepository(db)
			// Ожидаемый запрос
			query, _, err := sq.Insert("users").
				Columns("email", "password_hashed", "salt_password").
				Values(tc.user.Email, tc.user.PasswordHash, tc.user.PasswordSalt).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			tc.setupMock(mock, query)
			user, err := repo.AddUser(tc.user.Email, tc.user.PasswordHash, tc.user.PasswordSalt)
			require.Equal(t, tc.expectedErr, err)
			require.EqualValues(t, tc.expectedOut, user)
		})
	}
}

func TestUsersDB_GetUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Request     int
		ExpectedErr error
		ExpectedOut *entity.User
		SetupMock   func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:        "Пользователь найден",
			Request:     1,
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "password_hashed", "salt_password", "avatar_upload_id", "rating"}).
						AddRow(1, "email@email.com", "", []byte("hashed"), []byte("salt"), 0, 0))
			},
			ExpectedOut: &entity.User{
				ID:             1,
				Email:          "email@email.com",
				Name:           "",
				PasswordHash:   []byte("hashed"),
				PasswordSalt:   []byte("salt"),
				AvatarUploadID: 0,
			},
		},
		{
			Name:        "Пользователь не найден",
			Request:     1,
			ExpectedErr: repository.ErrUserNotFound,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(sql.ErrNoRows)
			},
			ExpectedOut: nil,
		},
		{
			Name:        "Ошибка при выполнении запроса",
			Request:     1,
			ExpectedErr: entity.PSQLWrap(errors.New("ошибка при выполнении запроса getUser"), errors.New("ошибка при выполнении запроса GetUser")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(errors.New("ошибка при выполнении запроса GetUser"))
			},
			ExpectedOut: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
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
			query, args, err := sq.
				Select("id", "email", "name", "password_hashed", "salt_password", "avatar_upload_id", "rating").
				From("users").
				Where(map[string]any{"id": tc.Request}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			if err != nil {
				t.Fatalf("ошибка при формировании sql-запроса GetUser: %s", err)
			}
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			user, err := repo.GetUserByID(tc.Request)
			require.Equal(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOut, user)
		})
	}
}

func TestUsersDB_GetUserByEmail(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Request     string
		ExpectedErr error
		ExpectedOut *entity.User
		SetupMock   func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:        "Пользователь найден",
			Request:     "email@email.com",
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "password_hashed", "salt_password", "avatar_upload_id", "rating"}).
						AddRow(1, "email@email.com", "", []byte("hashed"), []byte("salt"), 0, 0))
			},
			ExpectedOut: &entity.User{
				ID:             1,
				Email:          "email@email.com",
				Name:           "",
				PasswordHash:   []byte("hashed"),
				PasswordSalt:   []byte("salt"),
				AvatarUploadID: 0,
			},
		},
		{
			Name:        "Пользователь не найден",
			Request:     "email@email.com",
			ExpectedErr: repository.ErrUserNotFound,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(sql.ErrNoRows)
			},
			ExpectedOut: nil,
		},
		{
			Name:        "Ошибка при выполнении запроса",
			Request:     "email@email.com",
			ExpectedErr: entity.PSQLQueryErr("getUser", sql.ErrConnDone),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(sql.ErrConnDone)
			},
			ExpectedOut: nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
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
			query, args, err := sq.
				Select("id", "email", "name", "password_hashed", "salt_password", "avatar_upload_id", "rating").
				From("users").
				Where(map[string]any{"email": tc.Request}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			if err != nil {
				t.Fatalf("ошибка при формировании sql-запроса GetUser: %s", err)
			}
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			user, err := repo.GetUserByEmail(tc.Request)
			require.Equal(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOut, user)
		})
	}
}

func TestUsersDB_UpdateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		User        *entity.User
		ExpectedErr error
		SetupMock   func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name: "Успешное обновление",
			User: &entity.User{
				ID:             1,
				Email:          "email@email.com",
				Name:           "name",
				PasswordHash:   []byte("hashed"),
				PasswordSalt:   []byte("salt"),
				AvatarUploadID: 1,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
		},
		{
			Name: "Ошибка при выполнении запроса",
			User: &entity.User{
				ID:             1,
				Email:          "email@email.com",
				Name:           "name",
				PasswordHash:   []byte("hashed"),
				PasswordSalt:   []byte("salt"),
				AvatarUploadID: 1,
			},
			ExpectedErr: entity.PSQLQueryErr("UpdateUser", sql.ErrConnDone),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(sql.ErrConnDone)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
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
			query, args, err := sq.Update("users").
				Where(map[string]interface{}{"id": 1}).
				SetMap(map[string]interface{}{
					"email":            tc.User.Email,
					"name":             tc.User.Name,
					"password_hashed":  tc.User.PasswordHash,
					"salt_password":    tc.User.PasswordSalt,
					"avatar_upload_id": tc.User.AvatarUploadID,
				}).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			err = repo.UpdateUser(tc.User)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
