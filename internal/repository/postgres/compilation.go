package postgres

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/jmoiron/sqlx"
)

type CompilationDB struct {
	DB *sqlx.DB
}

// NewCompilationRepository создает новый репозиторий подборок
func NewCompilationRepository(db *sqlx.DB) repository.Compilation {
	return &CompilationDB{
		DB: db,
	}
}

// GetCompilation получает подборку по ID
func (c *CompilationDB) GetCompilation(id int) (*entity.Compilation, error) {
	query, args, err := sq.Select("id", "title", "compilation_type_id", "poster_upload_id").
		From("compilation").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Join(errors.New("ошибка при составлении запроса GetCompilation"), err)
	}
	row := c.DB.QueryRow(query, args...)
	compilation := entity.Compilation{}
	var posterUploadID sql.NullInt64
	err = row.Scan(
		&compilation.ID,
		&compilation.Title,
		&compilation.CompilationTypeID,
		&posterUploadID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrCompilationNotFound
		}
		return nil, entity.PSQLQueryErr("GetCompilation при сканировании", err)
	}
	if posterUploadID.Valid {
		compilation.PosterUploadID = int(posterUploadID.Int64)
	}
	return &compilation, nil
}

// GetCompilationsByTypeID получает все подборки по ID категории
func (c *CompilationDB) GetCompilationsByTypeID(compilationTypeID int) ([]*entity.Compilation, error) {
	query, args, err := sq.Select("id", "title", "compilation_type_id", "poster_upload_id").
		From("compilation").
		Where(sq.Eq{"compilation_type_id": compilationTypeID}).
		OrderBy("id ASC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Join(errors.New("ошибка при составлении запроса GetCompilationsByTypeID"))
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("GetCompilationsByTypeID", err)
	}
	defer rows.Close()
	compilations := make([]*entity.Compilation, 0)
	for rows.Next() {
		var posterUploadID sql.NullInt64
		compilation := entity.Compilation{}
		err = rows.Scan(
			&compilation.ID,
			&compilation.Title,
			&compilation.CompilationTypeID,
			&posterUploadID,
		)
		if err != nil {
			return nil, entity.PSQLQueryErr("GetCompilationsByTypeID при сканировании", err)
		}
		if posterUploadID.Valid {
			compilation.PosterUploadID = int(posterUploadID.Int64)
		}
		compilations = append(compilations, &compilation)
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
		return 0, entity.PSQLQueryErr("GetCompilationContentLength", err)
	}
	return length, nil
}

// GetCompilationContent получает список id контента из бд у конкретной подборки по ID
func (c *CompilationDB) GetCompilationContent(id, page, limit int) ([]int, error) {
	query, args, err := sq.Select("content_id").
		From("compilation_content").
		Join("content ON compilation_content.content_id = content.id").
		Where(sq.Eq{"compilation_id": id}).
		OrderBy("content.rating DESC", "id ASC").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Join(errors.New("ошибка при составлении запроса GetCompilationContent"), err)
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrCompilationNotFound
		}
		return nil, entity.PSQLQueryErr("GetCompilationContent", err)
	}
	defer rows.Close()
	contentIDs := make([]int, 0, limit)
	for rows.Next() {
		var contentID int
		err = rows.Scan(&contentID)
		if err != nil {
			return nil, entity.PSQLQueryErr("GetCompilationContent при сканировании", err)
		}
		contentIDs = append(contentIDs, contentID)
	}
	return contentIDs, nil
}

// GetAllCompilationTypes получает все категории подборок
func (c *CompilationDB) GetAllCompilationTypes() ([]entity.CompilationType, error) {
	query, args, err := sq.Select("id", "type").
		From("compilation_type").
		OrderBy("id ASC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, errors.Join(errors.New("ошибка при составлении запроса GetAllCompilationTypes"), err)
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("GetAllCompilationTypes", err)
	}
	defer rows.Close()
	compilationTypes := make([]entity.CompilationType, 0)
	for rows.Next() {
		compilationType := entity.CompilationType{}
		err = rows.Scan(&compilationType.ID, &compilationType.Name)
		if err != nil {
			return nil, entity.PSQLQueryErr("GetAllCompilationTypes при сканировании", err)
		}
		compilationTypes = append(compilationTypes, compilationType)
	}
	return compilationTypes, nil
}
