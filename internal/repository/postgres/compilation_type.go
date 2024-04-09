package postgres

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
)

type CompilationTypeDB struct {
	DB *sql.DB
}

func NewCompilationTypeRepository(database config.PostgresDatabase) (repository.CompilationType, error) {
	db, err := sql.Open("postgres", database.ConnectURL)
	if err != nil {
		return nil, err
	}
	return &CompilationTypeDB{
		DB: db,
	}, nil
}

// GetCompilationType получает категорию подборки по id
func (c CompilationTypeDB) GetCompilationType(id int) (*entity.CompilationType, error) {
	query, args, err := sq.Select("id", "type").
		From("compilation_type").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса GetCompilationType"))
	}
	compilationType := &entity.CompilationType{}
	err = c.DB.QueryRow(query, args...).Scan(&compilationType.ID, &compilationType.Type)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewClientError("категория подборки не найдена", entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении категории"))
	}
	return compilationType, nil
}

// GetAllCompilationTypes получает все категории подборок
// они отсортированы по алфавиту
func (c CompilationTypeDB) GetAllCompilationTypes() ([]*entity.CompilationType, error) {
	query, args, err := sq.Select("id", "type").
		From("compilation_type").
		OrderBy("type").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса GetAllCompilationTypes"))
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении подборок"))
	}
	// закрываем соединение
	defer rows.Close()
	// далее сканируем строки, создаем сущности и добавляем их в массив
	compilationTypes := make([]*entity.CompilationType, 0)
	for rows.Next() {
		compilationType := &entity.CompilationType{}
		err = rows.Scan(&compilationType.ID, &compilationType.Type)
		if err != nil {
			return nil, entity.PSQLWrap(err, errors.New("ошибка при сканировании категорий подборок"))
		}
		compilationTypes = append(compilationTypes, compilationType)
	}
	return compilationTypes, nil
}
