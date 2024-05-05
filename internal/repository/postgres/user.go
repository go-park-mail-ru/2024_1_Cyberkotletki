package postgres

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/lib/pq"
)

type UsersDB struct {
	DB *sql.DB
}

type DBUser struct {
	ID             int
	Email          string
	Name           sql.NullString
	PasswordHash   []byte
	PasswordSalt   []byte
	AvatarUploadID sql.NullInt64
	Rating         int
}

func (u *DBUser) GetEntity() *entity.User {
	return &entity.User{
		ID:             u.ID,
		Email:          u.Email,
		Name:           u.Name.String,
		PasswordHash:   u.PasswordHash,
		PasswordSalt:   u.PasswordSalt,
		AvatarUploadID: int(u.AvatarUploadID.Int64),
	}
}

func NewUserRepository(db *sql.DB) repository.User {
	return &UsersDB{
		DB: db,
	}
}

// AddUser добавляет пользователя в базу данных.
// Если операция происходит успешно, то в переданный по указателю user будет записан id нового пользователя.
func (u *UsersDB) AddUser(email string, passwordHash, passwordSalt []byte) (*entity.User, error) {
	query, args, err := sq.Insert("\"user\"").
		Columns("email", "password_hashed", "salt_password").
		Values(email, passwordHash, passwordSalt).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(errors.New("ошибка при составлении запроса AddUser"), err)
	}
	var lastID int
	err = u.DB.QueryRow(query, args...).Scan(&lastID)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case entity.PSQLUniqueViolation:
			return nil, repository.ErrUserAlreadyExists
		case entity.PSQLCheckViolation:
			return nil, repository.ErrUserIncorrectData
		}
	}
	if err != nil {
		// неизвестная ошибка
		return nil, entity.PSQLQueryErr("AddUser", err)
	}
	user := new(entity.User)
	user.ID = lastID
	user.Email = email
	user.PasswordHash = passwordHash
	user.PasswordSalt = passwordSalt
	return user, nil
}

func (u *UsersDB) getUser(where map[string]any) (*entity.User, error) {
	query, args, err := sq.
		Select("\"id\"", "email", "name", "password_hashed", "salt_password", "avatar_upload_id", "rating").
		From("\"user\"").
		Where(where).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(errors.New("ошибка при составлении запроса getUser"), err)
	}
	user := DBUser{}
	err = u.DB.
		QueryRow(query, args...).
		Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.PasswordHash,
			&user.PasswordSalt,
			&user.AvatarUploadID,
			&user.Rating,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrUserNotFound
	}
	if err != nil {
		return nil, entity.PSQLQueryErr("getUser", err)
	}
	return user.GetEntity(), nil
}

// GetUserByID возвращает пользователя из базы данных по айди
func (u *UsersDB) GetUserByID(userID int) (*entity.User, error) {
	return u.getUser(map[string]any{"id": userID})
}

// GetUserByEmail возвращает пользователя из базы данных по почте
func (u *UsersDB) GetUserByEmail(userEmail string) (*entity.User, error) {
	return u.getUser(map[string]any{"email": userEmail})
}

// UpdateUser обновляет пользователя в базе данных по переданной структуре
func (u *UsersDB) UpdateUser(user *entity.User) error {
	setMap := make(map[string]any)
	if user.Email != "" {
		setMap["email"] = user.Email
	}
	if user.Name != "" {
		setMap["name"] = user.Name
	}
	if len(user.PasswordHash) > 0 {
		setMap["password_hashed"] = user.PasswordHash
	}
	if len(user.PasswordSalt) > 0 {
		setMap["salt_password"] = user.PasswordSalt
	}
	if user.AvatarUploadID != 0 {
		setMap["avatar_upload_id"] = user.AvatarUploadID
	}
	query, args, err := sq.Update("\"user\"").
		Where(map[string]any{"id": user.ID}).
		SetMap(setMap).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(errors.New("ошибка при составлении запроса UpdateUser"), err)
	}
	// ни одна колонка может быть не изменена, но будем считать, что это успешное обновление
	if _, err = u.DB.Exec(query, args...); err != nil {
		return entity.PSQLQueryErr("UpdateUser", err)
	}
	return nil
}
