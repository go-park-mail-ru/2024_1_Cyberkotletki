package postgres

import (
	"context"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UsersDB struct {
	DB  *pgx.Conn
	ctx context.Context
}

func NewUserRepository(database config.PostgresDatabase) (repository.User, error) {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, database.ConnectURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &UsersDB{
		DB:  db,
		ctx: ctx,
	}, nil
}

// HasUser проверяет наличие пользователя в базе данных по email
func (u *UsersDB) HasUser(user *entity.User) (bool, error) {
	query, args, err := sq.
		Select("*").
		From("users").
		Where(sq.Eq{"email": user.Email}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return false, entity.PSQLWrap(err)
	}
	var result entity.User
	err = u.DB.QueryRow(u.ctx, query, args...).Scan(&result.ID, &result.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		// запись не найдена == пользователя нет
		return false, nil
	} else if err != nil {
		return false, entity.PSQLWrap(err)
	}

	return true, nil
}

// AddUser добавляет пользователя в базу данных.
// У переданного пользователя должен быть заполнен email, passwordHash, passwordSalt.
// Если операция происходит успешно, то в переданный по указателю user будет записан id и вернется указатель на
// этого же пользователя
func (u *UsersDB) AddUser(user *entity.User) (*entity.User, error) {
	query, args, err := sq.Insert("users").
		Columns("email", "password_hashed", "salt_password").
		Values(user.Email, user.PasswordHash, user.PasswordSalt).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err)
	}
	var pgErr *pgconn.PgError
	err = u.DB.QueryRow(u.ctx, query, args...).Scan(&user.ID)
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, entity.NewClientError("пользователь с таким email уже существует", entity.ErrAlreadyExists)
		}
		return nil, entity.PSQLWrap(err)
	}

	return user, nil
}

// GetUser возвращает пользователя из базы данных по переданным параметрам.
// Если пользователь не найден, то возвращается ошибка ErrNotFound
func (u *UsersDB) GetUser(params map[string]interface{}) (*entity.User, error) {
	query, args, err := sq.
		Select("id", "email", "name", "password_hashed", "salt_password").
		From("users").
		Where(params).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err)
	}
	var result entity.User
	err = u.DB.QueryRow(u.ctx, query, args...).
		Scan(&result.ID, &result.Email, &result.Name, &result.PasswordHash, &result.PasswordSalt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.NewClientError("пользователь не найден", entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err)
	}
	return &result, nil
}

// UpdateUser обновляет пользователя в базе данных по переданным параметрам
func (u *UsersDB) UpdateUser(params map[string]interface{}, values map[string]interface{}) error {
	query, args, err := sq.Update("users").
		SetMap(values).
		Where(params).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err)
	}
	err = u.DB.QueryRow(u.ctx, query, args...).Scan()
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.NewClientError("пользователь не найден", entity.ErrNotFound)
		}
		return entity.PSQLWrap(err)
	}
	return nil
}
