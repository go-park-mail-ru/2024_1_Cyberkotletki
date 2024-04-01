package postgres

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/jackc/pgx/v5"
	"os"
)

type StaticDB struct {
	DB        *pgx.Conn
	ctx       context.Context
	basicPath string
	maxSize   int
}

func NewStaticRepository(database config.PostgresDatabase, basicPath string, maxSize int) (repository.Static, error) {
	ctx := context.Background()
	db, err := pgx.Connect(ctx, database.ConnectURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &StaticDB{
		DB:        db,
		ctx:       ctx,
		basicPath: basicPath,
		maxSize:   maxSize,
	}, nil
}

// GetStatic возвращает путь к статике по его ID
func (s StaticDB) GetStatic(staticID int) (string, error) {
	query, args, err := sq.
		Select("path", "name").
		From("static").
		Where(sq.Eq{"id": staticID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return "", entity.PSQLWrap(err)
	}
	var path, name string
	err = s.DB.QueryRow(s.ctx, query, args...).Scan(&path, &name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", entity.NewClientError("файл не найден", entity.ErrNotFound)
		}
		return "", entity.PSQLWrap(err)
	}
	return fmt.Sprintf("%s/%s", path, name), nil
}

// UploadStatic загружает статику на сервер
func (s StaticDB) UploadStatic(path, filename string, data []byte) (int, error) {
	// Сначала создаем файл на диске
	if len(data) > s.maxSize {
		return -1, entity.NewClientError("файл слишком большой", entity.ErrBadRequest)
	}
	dst, err := os.Create(fmt.Sprintf("%s/%s/%s", s.basicPath, path, filename))
	if err != nil {
		return -1, entity.NewClientError("произошла внутренняя ошибка", entity.ErrInternal, err)
	}
	if _, err = dst.Write(data); err != nil {
		return -1, entity.NewClientError("произошла внутренняя ошибка", entity.ErrInternal, err)
	}
	if err = dst.Close(); err != nil {
		return -1, entity.NewClientError("произошла внутренняя ошибка", entity.ErrInternal, err)
	}

	// Затем добавляем запись в базу данных
	query, args, err := sq.
		Insert("static").
		Columns("path", "name").
		Values(path, filename).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return -1, entity.PSQLWrap(err)
	}
	var id int
	if err = s.DB.QueryRow(s.ctx, query, args...).Scan(&id); err != nil {
		return -1, entity.PSQLWrap(err)
	}
	return id, nil
}
