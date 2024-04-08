package postgres

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/lib/pq"
)

type CompilationDB struct {
	DB *sql.DB
}

// GetCompilationContentLength получает число контента в подборке
func (c CompilationDB) GetCompilationContentLength(id int) (int, error) {
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

// AddCompilation добавляет подборку в бд
// Должны быть заполнены поля CompilationTypeID, Title, PosterUploadID
// Возвращает указатель на подборку с заполненным ID
func (c CompilationDB) AddCompilation(compilation *entity.Compilation) (*entity.Compilation, error) {
	query, args, err := sq.Insert("compilation").
		Columns("compilation_type_id", "title", "poster").
		Values(compilation.CompilationTypeID, compilation.Title, compilation.PosterUploadID).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса AddCompilation"))
	}
	err = c.DB.QueryRow(query, args...).Scan(&compilation.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case entity.PSQLUniqueViolation:
				return nil, entity.NewClientError("подборка уже существует", entity.ErrAlreadyExists)
			case entity.PSQLForeignKeyViolation:
				return nil, entity.NewClientError("категория подборки или постера с таким id не существует", entity.ErrNotFound)
			case entity.PSQLCheckViolation:
				return nil, entity.NewClientError("неверный формат постера или названия подборки", entity.ErrBadRequest)

			}
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при добавлении подборки"))
	}
	return compilation, nil
}

// GetCompilation получает подборку из бд по ID
func (c CompilationDB) GetCompilation(id int) (*entity.Compilation, error) {
	query, args, err := sq.Select("id", "compilation_type_id", "title", "poster").
		From("compilation").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса GetCompilation"))
	}

	compilation := &entity.Compilation{}
	err = c.DB.QueryRow(query, args...).Scan(&compilation.ID, &compilation.CompilationTypeID, &compilation.Title, &compilation.PosterUploadID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewClientError("подборка не найдена", entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении подборки"))
	}
	return compilation, nil
}

// GetAllCompilationsByCompilationTypeID получает все подборки по ID категории подборок из бд отсортированные по алфавиту
func (c CompilationDB) GetAllCompilationsByCompilationTypeID(compilationTypeID, page, limit int) ([]*entity.Compilation, error) {
	query, args, err := sq.Select("id", "compilation_type_id", "title", "poster").
		From("compilation").
		Where(sq.Eq{"compilation_type_id": compilationTypeID}).
		OrderBy("title ASC").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса GetAllCompilationsByCompilationTypeID"))
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

// UpdateCompilation обновляет подборку в бд
// Должны быть заполнены поля ID, CompilationTypeID, Title, PosterUploadID
// Возвращает указатель на подборку с обновленными полями
func (c CompilationDB) UpdateCompilation(compilation *entity.Compilation) (*entity.Compilation, error) {
	query, args, err := sq.Update("compilation").
		Set("compilation_type_id", compilation.CompilationTypeID).
		Set("title", compilation.Title).
		Set("poster", compilation.PosterUploadID).
		Where(sq.Eq{"id": compilation.ID}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса UpdateCompilation"))
	}
	_, err = c.DB.Exec(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при обновлении подборки"))
	}
	return compilation, nil
}

// DeleteCompilation удаляет подборку из бд
func (c CompilationDB) DeleteCompilation(id int) error {
	query, args, err := sq.Delete("compilation").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса DeleteCompilation"))
	}
	_, err = c.DB.Exec(query, args...)
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при удалении подборки"))
	}
	return nil
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
