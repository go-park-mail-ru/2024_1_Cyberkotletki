package postgres

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
)

type CompilationDB struct {
	DB *sql.DB
}

// NewCompilationRepository создает новый репозиторий подборок
func NewCompilationRepository(database config.PostgresDatabase) (repository.Compilation, error) {
	db, err := sql.Open("postgres", database.ConnectURL)
	if err != nil {
		return nil, err
	}
	return &CompilationDB{
		DB: db,
	}, nil
}

// GetAllCompilationsByCompilationTypeID получает все подборки по ID категории
// подборок из бд отсортированные по алфавиту
func (c *CompilationDB) GetCompilationsByCompilationTypeID(compilationTypeID,
	page, limit int) ([]*entity.Compilation, error) {
	query, args, err := sq.Select("с.id", "с.compilation_type_id", "с.title", "с.poster").
		From("compilation с").
		Where(sq.Eq{"с.compilation_type_id": compilationTypeID}).
		OrderBy("с.title ASC").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err,
			errors.New("ошибка при формировании sql-запроса GetAllCompilationsByCompilationTypeID"))
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении подборок"))
	}
	// закрываем соединение
	defer rows.Close()
	// далее сканируем строки, создаем сущности и добавляем их в массив
	compilations := make([]*entity.Compilation, 0, limit)
	for rows.Next() {
		compilation := &entity.Compilation{}
		err = rows.Scan(
			&compilation.ID,
			&compilation.CompilationTypeID,
			&compilation.Title,
			&compilation.PosterUploadID,
		)
		if err != nil {
			return nil, entity.PSQLWrap(err, errors.New("ошибка при сканировании подборок"))
		}
		compilations = append(compilations, compilation)
	}
	return compilations, nil
}

// GetCompilation получает подборку из бд по ID
func (c *CompilationDB) GetCompilation(id int) (*entity.Compilation, error) {
	query, args, err := sq.Select("id", "title", "compilation_type_id", " poster_upload_id").
		From("compilation").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err,
			errors.New("ошибка при формировании sql-запроса GetCompilation"))
	}

	compilation := &entity.Compilation{}
	err = c.DB.QueryRow(query, args...).Scan(&compilation.ID,
		&compilation.CompilationTypeID, &compilation.Title, &compilation.PosterUploadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewClientError("подборка не найдена", entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении подборки"))
	}
	return compilation, nil
}

// GetCompilationContentLength получает число контента в подборке
func (c *CompilationDB) GetCompilationContentLength(id int) (int, error) {
	query, args, err := sq.Select("count(*)").
		From("compilation_content").
		Where(sq.Eq{"compilation_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, entity.PSQLWrap(err,
			errors.New("ошибка при формировании sql-запроса GetCompilationContentLength"))
	}
	var length int
	err = c.DB.QueryRow(query, args...).Scan(&length)
	if err != nil {
		return 0, entity.PSQLWrap(err,
			errors.New("ошибка при получении числа контента в подборке"))
	}
	return length, nil

}

// GetCompilationContent получает список id контента из бд у конкретной подборки по ID
func (c *CompilationDB) GetCompilationContent(id, page, limit int) ([]int, error) {
	query, args, err := sq.Select("content_id").
		From("compilation_content").
		Where(sq.Eq{"compilation_id": id}).
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []int{}, entity.PSQLWrap(err, errors.New("контент подборки не найден"), entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err,
			errors.New("ошибка при получении контента подборки"))
	}
	defer rows.Close()

	contentIDs := make([]int, 0, limit)
	for rows.Next() {
		var contentID int
		err = rows.Scan(&contentID)
		if err != nil {
			return nil, err
		}
		contentIDs = append(contentIDs, contentID)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return contentIDs, nil
}
