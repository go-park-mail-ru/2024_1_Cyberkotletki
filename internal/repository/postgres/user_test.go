package postgres

import (
	"database/sql"
	"database/sql/driver"
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

func TestUsersDB_GetUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Params      map[string]interface{}
		ExpectedErr error
		ExpectedOut *entity.User
		SetupMock   func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name: "Пользователь найден",
			Params: map[string]interface{}{
				"email": "email@email.com",
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "password_hashed", "salt_password", "avatar_upload_id"}).
						AddRow(1, "email@email.com", "", []byte("hashed"), []byte("salt"), 0))
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
			Name: "Пользователь не найден",
			Params: map[string]interface{}{
				"email": "email@email.com",
			},
			ExpectedErr: entity.NewClientError("пользователь не найден", entity.ErrNotFound),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(sql.ErrNoRows)
			},
			ExpectedOut: nil,
		},
		{
			Name: "Ошибка при выполнении запроса",
			Params: map[string]interface{}{
				"email": "email@email.com",
			},
			ExpectedErr: entity.PSQLWrap(errors.New("ошибка при выполнении sql-запроса GetUser"), errors.New("ошибка при выполнении sql-запроса GetUser")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(errors.New("ошибка при выполнении sql-запроса GetUser"))
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
				Select("id", "email", "name", "password_hashed", "salt_password", "avatar_upload_id").
				From("users").
				Where(tc.Params).
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
			user, err := repo.GetUser(tc.Params)
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
				AvatarUploadID: 0,
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
				AvatarUploadID: 0,
			},
			ExpectedErr: entity.PSQLWrap(errors.New("ошибка при выполнении sql-запроса UpdateUser"), errors.New("ошибка при выполнении sql-запроса UpdateUser")),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectExec(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(errors.New("ошибка при выполнении sql-запроса UpdateUser"))
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
			err = repo.UpdateUser(map[string]interface{}{"id": 1}, map[string]interface{}{
				"email":            tc.User.Email,
				"name":             tc.User.Name,
				"password_hashed":  tc.User.PasswordHash,
				"salt_password":    tc.User.PasswordSalt,
				"avatar_upload_id": tc.User.AvatarUploadID,
			})
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
