package postgres

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Драйвер для работы с PostgreSQL
	"io"
)

type StaticDB struct {
	DB         *sqlx.DB
	s3         *s3.S3
	bucketName string
	maxSize    int
}

func NewStaticRepository(db *sqlx.DB, s3 *s3.S3, bucketName string, maxSize int) repository.Static {
	return &StaticDB{
		DB:         db,
		s3:         s3,
		bucketName: bucketName,
		maxSize:    maxSize,
	}
}

func (s StaticDB) GetMaxSize() int {
	return s.maxSize
}

// GetStaticURL возвращает путь к статике по его ID
func (s StaticDB) GetStaticURL(staticID int) (string, error) {
	query, args, err := sq.
		Select("path", "Name").
		From("static").
		Where(sq.Eq{"id": staticID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return "", entity.PSQLWrap(errors.New("ошибка при составлении запроса GetStaticURL"), err)
	}
	var path, name string
	err = s.DB.QueryRow(query, args...).Scan(&path, &name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", repository.ErrStaticNotFound
		}
		return "", entity.PSQLQueryErr("GetStaticURL", err)
	}
	return fmt.Sprintf("%s/%s", path, name), nil
}

// UploadStatic загружает статику на сервер
func (s StaticDB) UploadStatic(path, filename string, reader io.ReadSeeker) (int, error) {
	// Проверка размера файла
	size, err := reader.Seek(0, io.SeekEnd)
	if err != nil {
		return -1, errors.Join(err, errors.New("ошибка при определении размера файла"))
	}
	if size > int64(s.maxSize) {
		return -1, repository.ErrStaticTooBigFile
	}
	_, err = reader.Seek(0, io.SeekStart) // Возвращаемся в начало файла
	if err != nil {
		return -1, errors.Join(err, errors.New("ошибка при возвращении io.ReadSeeker в начало файла"))
	}

	// Загрузка файла в S3
	_, err = s.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(fmt.Sprintf("%s/%s", path, filename)),
		Body:   reader,
	})
	if err != nil {
		return -1, err
	}

	// Добавляем запись в базу данных
	query, args, err := sq.
		Insert("static").
		Columns("path", "Name").
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

func (s StaticDB) GetStaticFile(staticURI string) (io.ReadSeeker, error) {
	// Получаем статику из S3
	output, err := s.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(staticURI),
	})
	if err != nil {
		return nil, repository.ErrStaticNotFound
	}
	buf, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, entity.S3Wrap(errors.New("ошибка при чтении файла из S3"), err)
	}
	return bytes.NewReader(buf), nil
}
