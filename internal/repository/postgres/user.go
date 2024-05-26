package postgres

import (
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/lib/pq"
)

type UsersDB struct {
	DB *sql.DB
}

type User struct {
	ID             int
	Email          string
	Name           sql.NullString
	PasswordHash   []byte
	PasswordSalt   []byte
	AvatarUploadID sql.NullInt64
}

func (u *User) GetEntity() *entity.User {
	return &entity.User{
		ID:             u.ID,
		Email:          u.Email,
		Name:           u.Name.String,
		PasswordHash:   u.PasswordHash,
		PasswordSalt:   u.PasswordSalt,
		AvatarUploadID: int(u.AvatarUploadID.Int64),
	}
}

func NewUserRepository(database config.PostgresDatabase) (repository.User, error) {
	db, err := sql.Open("postgres", database.ConnectURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// почти при каждом запросе проверяется авторизация
	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Second * 5)

	return &UsersDB{
		DB: db,
	}, nil
}

// AddUser добавляет пользователя в базу данных.
// У переданного пользователя должен быть заполнен email, passwordHash, passwordSalt.
// Если операция происходит успешно, то в переданный по указателю user будет записан id и вернется указатель на
// этого же пользователя
func (u *UsersDB) AddUser(user *entity.User) (*entity.User, error) {
	// no-lint
	query, args, _ := sq.Insert("users").
		Columns("email", "password_hashed", "salt_password").
		Values(user.Email, user.PasswordHash, user.PasswordSalt).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	err := u.DB.QueryRow(query, args...).Scan(&user.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case entity.PSQLUniqueViolation:
				return nil, entity.NewClientError("пользователь с таким email уже существует", entity.ErrAlreadyExists)
			case entity.PSQLCheckViolation:
				return nil, entity.NewClientError(
					"одно или несколько полей заполнены некорректно",
					entity.ErrBadRequest,
					errors.New("ошибка при выполнении sql-запроса AddUser: нарушение целостности данных"),
				)
			default:
				return nil, entity.PSQLWrap(err, errors.New("ошибка при выполнении sql-запроса AddUser"), pqErr)
			}
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при выполнении sql-запроса AddUser"))
	}

	return user, nil
}

// GetUser возвращает пользователя из базы данных по переданным параметрам.
// Если пользователь не найден, то возвращается ошибка ErrNotFound
func (u *UsersDB) GetUser(params map[string]interface{}) (*entity.User, error) {
	// no-lint
	query, args, _ := sq.
		Select("id", "email", "name", "password_hashed", "salt_password", "avatar_upload_id").
		From("users").
		Where(params).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	user := User{}
	err := u.DB.QueryRow(query, args...).
		Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.PasswordSalt, &user.AvatarUploadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewClientError("пользователь не найден", entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при выполнении sql-запроса GetUser"))
	}
	return user.GetEntity(), nil
}

// UpdateUser обновляет пользователя в базе данных по переданным параметрам
func (u *UsersDB) UpdateUser(params map[string]interface{}, values map[string]interface{}) error {
	// no-lint
	query, args, _ := sq.Update("users").
		Where(params).
		SetMap(values).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if _, err := u.DB.Exec(query, args...); err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при выполнении sql-запроса UpdateUser"))
	}
	return nil
}
