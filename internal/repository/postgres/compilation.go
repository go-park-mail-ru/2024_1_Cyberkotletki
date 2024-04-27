package postgres

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
)

type CompilationDB struct {
	DB *sql.DB
}

// NewCompilationRepository создает новый репозиторий подборок
func NewCompilationRepository(db *sql.DB) repository.Compilation {
	return &CompilationDB{
		DB: db,
	}
}

// GetCompilationsByTypeID получает все подборки по ID категории
// подборок из бд отсортированные по алфавиту
func (c *CompilationDB) GetCompilationsByTypeID(compilationTypeID int) ([]entity.Compilation, error) {
	query, args, err := sq.Select("id", "title", "compilation_type_id", "poster_upload_id").
		From("compilation").
		Where(sq.Eq{"compilation_type_id": compilationTypeID}).
		OrderBy("title ASC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Join(errors.New("ошибка при составлении запроса GetCompilationsByTypeID"))
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении подборок"))
	}
	defer rows.Close()
	compilations := make([]entity.Compilation, 0)
	for rows.Next() {
		compilation := entity.Compilation{}
		err = rows.Scan(
			&compilation.ID,
			&compilation.Title,
			&compilation.CompilationTypeID,
			&compilation.PosterUploadID,
		)
		if err != nil {
			return nil, entity.PSQLWrap(err, errors.New("ошибка при сканировании подборок"))
		}
		compilations = append(compilations, compilation)
	}
	return compilations, nil
}

// GetCompilationContentLength получает число контента в подборке
func (c *CompilationDB) GetCompilationContentLength(id int) (int, error) {
	query, args, err := sq.Select("count(*)").
		From("compilation_content").
		Where(sq.Eq{"compilation_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса GetCompilationContentLength"))
	}
	var length int
	err = c.DB.QueryRow(query, args...).Scan(&length)
	if err != nil {
		return 0, entity.PSQLWrap(err, errors.New("ошибка при получении числа контента в подборке"))
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
			return nil, entity.PSQLWrap(err, errors.New("контент подборки не найден"), entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении контента подборки"))
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
		return nil, entity.PSQLWrap(err, errors.New("ошибка при сканировании контента"))
	}
	return contentIDs, nil
}

// GetAllCompilationTypes получает все категории подборок в алфавитном порядке
func (c *CompilationDB) GetAllCompilationTypes() ([]entity.CompilationType, error) {
	query, args, err := sq.Select("id", "type").
		From("compilation_type").
		OrderBy("type ASC").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, errors.Join(errors.New("ошибка при составлении запроса GetAllCompilationTypes"), err)
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении категорий подборок"))
	}
	defer rows.Close()
	compilationTypes := make([]entity.CompilationType, 0)
	for rows.Next() {
		compilationType := entity.CompilationType{}
		err = rows.Scan(&compilationType.ID, &compilationType.Type)
		if err != nil {
			return nil, entity.PSQLWrap(err, errors.New("ошибка при сканировании категорий подборок"))
		}
		compilationTypes = append(compilationTypes, compilationType)
	}
	return compilationTypes, nil
}
