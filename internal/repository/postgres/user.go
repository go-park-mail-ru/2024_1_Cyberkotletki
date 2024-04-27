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

func NewUserRepository(db *sql.DB) repository.User {
	return &UsersDB{
		DB: db,
	}
}

// AddUser добавляет пользователя в базу данных.
// У переданного пользователя должен быть заполнен email, passwordHash, passwordSalt.
// Если операция происходит успешно, то в переданный по указателю user будет записан id и вернется указатель на
// этого же пользователя
func (u *UsersDB) AddUser(user *entity.User) (*entity.User, error) {
	query, args, err := sq.Insert("users").
		Columns("email", "password_hashed", "salt_password").
		Values(user.Email, user.PasswordHash, user.PasswordSalt).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Join(errors.New("ошибка при составлении запроса AddUser"), err)
	}
	result, err := u.DB.Exec(query, args...)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case entity.PSQLUniqueViolation:
			return nil, repository.ErrUserAlreadyExists
		case entity.PSQLCheckViolation:
			return nil, repository.ErrUserIncorrectData
		default:
			return nil, entity.PSQLWrap(pqErr, errors.New("ошибка при выполнении запроса AddUser"))
		}
	} else if err != nil {
		// неизвестная ошибка
		return nil, entity.PSQLWrap(err, errors.New("ошибка при выполнении запроса AddUser"))
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при выполнении запроса AddUser"))
	}
	user.ID = int(lastID)
	return user, nil
}

func (u *UsersDB) getUser(query string, args ...any) (*entity.User, error) {
	user := User{}
	err := u.DB.
		QueryRow(query, args...).
		Scan(&user.ID, &user.Email, &user.Name, &user.PasswordHash, &user.PasswordSalt, &user.AvatarUploadID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrUserNotFound
	} else if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при выполнении запроса getUser"))
	}
	return user.GetEntity(), nil
}

// GetUserByID возвращает пользователя из базы данных по айди
func (u *UsersDB) GetUserByID(userID int) (*entity.User, error) {
	query, args, err := sq.
		Select("id", "email", "name", "password_hashed", "salt_password", "avatar_upload_id").
		From("users").
		Where(map[string]any{"id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Join(errors.New("ошибка при составлении запроса GetUserByID"), err)
	}
	return u.getUser(query, args...)
}

// GetUserByEmail возвращает пользователя из базы данных по почте
func (u *UsersDB) GetUserByEmail(userEmail string) (*entity.User, error) {
	query, args, err := sq.
		Select("id", "email", "name", "password_hashed", "salt_password", "avatar_upload_id").
		From("users").
		Where(map[string]any{"email": userEmail}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Join(errors.New("ошибка при составлении запроса GetUserByEmail"), err)
	}
	return u.getUser(query, args...)
}

// UpdateUser обновляет пользователя в базе данных по переданным параметрам
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
	query, args, err := sq.Update("users").
		Where(map[string]any{"id": user.ID}).
		SetMap(setMap).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return errors.Join(errors.New("ошибка при составлении запроса UpdateUser"), err)
	}
	// ни одна колонка может быть не изменена, но будем считать, что это успешное обновление
	if _, err = u.DB.Exec(query, args...); err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при выполнении запроса UpdateUser"))
	}
	return nil
}
