package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	_ "github.com/lib/pq" // Драйвер для работы с PostgreSQL
	"os"
	"path/filepath"
)

type StaticDB struct {
	DB        *sql.DB
	basicPath string
	maxSize   int
}

func NewStaticRepository(database config.PostgresDatabase, basicPath string, maxSize int) (repository.Static, error) {
	db, err := sql.Open("postgres", database.ConnectURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &StaticDB{
		DB:        db,
		basicPath: basicPath,
		maxSize:   maxSize,
	}, nil
}

// GetStatic возвращает путь к статике по его ID
func (s StaticDB) GetStatic(staticID int) (string, error) {
	// no-lint
	query, args, _ := sq.
		Select("path", "name").
		From("static").
		Where(sq.Eq{"id": staticID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	var path, name string
	err := s.DB.QueryRow(query, args...).Scan(&path, &name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", entity.NewClientError("файл не найден", entity.ErrNotFound)
		}
		return "", entity.PSQLWrap(err, errors.New("ошибка при выполнении sql-запроса GetStatic"))
	}
	return fmt.Sprintf("%s/%s", path, name), nil
}

// UploadStatic загружает статику на сервер
func (s StaticDB) UploadStatic(path, filename string, data []byte) (int, error) {
	if len(data) > s.maxSize {
		return -1, entity.NewClientError("файл слишком большой", entity.ErrBadRequest)
	}
	// Создаем путь, если он еще не существует
	dir := filepath.Dir(fmt.Sprintf("%s/%s/", s.basicPath, path))
	err := os.MkdirAll(dir, 0755)
	fmt.Println(dir)
	if err != nil {
		return -1, entity.NewClientError(
			"произошла внутренняя ошибка",
			entity.ErrInternal,
			err,
			fmt.Errorf("не удалось создать папку для хранения статики"),
		)
	}
	// Создаем файл на диске
	dst, err := os.Create(fmt.Sprintf("%s/%s/%s", s.basicPath, path, filename))
	if err != nil {
		return -1, entity.NewClientError(
			"произошла внутренняя ошибка",
			entity.ErrInternal,
			err,
			fmt.Errorf("не удалось создать файл"),
		)
	}
	if _, err = dst.Write(data); err != nil {
		return -1, entity.NewClientError(
			"произошла внутренняя ошибка",
			entity.ErrInternal,
			err,
			fmt.Errorf("не удалось записать данные в файл"),
		)
	}
	if err = dst.Close(); err != nil {
		return -1, entity.NewClientError(
			"произошла внутренняя ошибка",
			entity.ErrInternal,
			err,
			fmt.Errorf("не удалось закрыть файл"),
		)
	}

	// Затем добавляем запись в базу данных
	// no-lint
	query, args, _ := sq.
		Insert("static").
		Columns("path", "name").
		Values(path, filename).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).ToSql()
	var id int
	if err = s.DB.QueryRow(query, args...).Scan(&id); err != nil {
		return -1, entity.PSQLWrap(err, errors.New("ошибка при выполнении sql-запроса UploadStatic"))
	}
	return id, nil
}
