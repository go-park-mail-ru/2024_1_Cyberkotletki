package postgres

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	_ "github.com/lib/pq" // Драйвер для работы с PostgreSQL
	"io"
	"os"
	"path/filepath"
)

type StaticDB struct {
	DB        *sql.DB
	basicPath string
	maxSize   int
}

func NewStaticRepository(db *sql.DB, basicPath string, maxSize int) repository.Static {
	return &StaticDB{
		DB:        db,
		basicPath: basicPath,
		maxSize:   maxSize,
	}
}

func (s StaticDB) GetBasicPath() string {
	return s.basicPath
}

func (s StaticDB) GetMaxSize() int {
	return s.maxSize
}

// GetStatic возвращает путь к статике по его ID
func (s StaticDB) GetStatic(staticID int) (string, error) {
	query, args, err := sq.
		Select("path", "name").
		From("static").
		Where(sq.Eq{"id": staticID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return "", entity.PSQLWrap(errors.New("ошибка при составлении запроса GetStatic"), err)
	}
	var path, name string
	err = s.DB.QueryRow(query, args...).Scan(&path, &name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", repository.ErrStaticNotFound
		}
		return "", entity.PSQLQueryErr("GetStatic", err)
	}
	return fmt.Sprintf("%s/%s", path, name), nil
}

// UploadStatic загружает статику на сервер
func (s StaticDB) UploadStatic(path, filename string, buf bytes.Buffer) (int, error) {
	data := make([]byte, s.maxSize)
	bytesCount, err := buf.Read(data)
	if err != nil && err != io.EOF {
		return -1, errors.Join(err, errors.New("ошибка при чтении данных"))
	}
	if bytesCount > s.maxSize {
		return -1, repository.ErrStaticTooBigFile
	}
	data = data[:bytesCount]

	// Создаем путь, если он еще не существует
	dir := filepath.Dir(fmt.Sprintf("%s/%s/", s.basicPath, path))
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return -1, errors.Join(err, errors.New("не удалось создать папку для хранения статики"))
	}

	// Создаем файл на диске
	dst, err := os.Create(fmt.Sprintf("%s/%s/%s", s.basicPath, path, filename))
	if err != nil {
		return -1, errors.Join(err, errors.New("не удалось создать файл"))
	}
	if _, err = dst.Write(data); err != nil {
		return -1, errors.Join(err, errors.New("не удалось записать данные в созданный файл"))
	}
	if err = dst.Close(); err != nil {
		return -1, errors.Join(err, errors.New("не удалось закрыть файл после записи информации"))
	}

	// Добавляем запись в базу данных
	query, args, err := sq.
		Insert("static").
		Columns("path", "name").
		Values(path, filename).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return -1, entity.PSQLWrap(errors.New("ошибка при составлении запроса UploadStatic"), err)
	}
	var id int
	if err = s.DB.QueryRow(query, args...).Scan(&id); err != nil {
		return -1, entity.PSQLQueryErr("UploadStatic", err)
	}
	return id, nil
}
