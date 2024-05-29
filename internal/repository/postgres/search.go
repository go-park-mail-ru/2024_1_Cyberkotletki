package postgres

import (
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/jmoiron/sqlx"
)

type SearchDB struct {
	ContentDB repository.Content
	DB        *sqlx.DB
}

func NewSearchRepository(db *sqlx.DB, contentDB repository.Content) repository.Search {
	return &SearchDB{
		DB:        db,
		ContentDB: contentDB,
	}
}

// SearchContent ищет контент по запросу
// nolint: dupl
func (s SearchDB) SearchContent(query string) ([]int, error) {
	sqlQuery, args, err := sq.
		Select("id").
		From("content").
		Where(sq.Or{
			sq.Expr("word_similarity(title, ?) > 0.3", query),
			sq.Expr("word_similarity(original_title, ?) > 0.3", query),
		}).
		OrderBy(fmt.Sprintf(
			`CASE WHEN word_similarity(title, '%s') > 0.3 
THEN similarity(title, '%s') ELSE similarity(original_title, '%s') END DESC`, query, query, query)).
		Limit(5).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса SearchContent"))
	}

	rows, err := s.DB.Query(sqlQuery, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("SearchContent", err)
	}
	defer rows.Close()

	var contentIDs []int
	for rows.Next() {
		var contentID int
		err = rows.Scan(&contentID)
		if err != nil {
			return nil, entity.PSQLWrap(err, errors.New("ошибка при сканировании строк в SearchContent"))
		}
		contentIDs = append(contentIDs, contentID)
	}

	return contentIDs, nil
}

// SearchPerson ищет персону по запросу
// nolint: dupl
func (s SearchDB) SearchPerson(query string) ([]int, error) {
	sqlQuery, args, err := sq.
		Select("id").
		From("person").
		Where(sq.Or{
			sq.Expr("word_similarity(Name, ?) > 0.3", query),
			sq.Expr("word_similarity(en_name, ?) > 0.3", query),
		}).
		OrderBy(fmt.Sprintf(
			`CASE WHEN word_similarity(Name, '%s') > 0.3 
THEN similarity(Name, '%s') ELSE similarity(en_name, '%s') END DESC`, query, query, query)).
		Limit(5).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса SearchPerson"))
	}

	rows, err := s.DB.Query(sqlQuery, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("SearchPerson", err)
	}
	defer rows.Close()

	var personIDs []int
	for rows.Next() {
		var personID int
		err = rows.Scan(&personID)
		if err != nil {
			return nil, entity.PSQLWrap(err, errors.New("ошибка при сканировании строк в SearchPerson"))
		}
		personIDs = append(personIDs, personID)
	}

	return personIDs, nil
}
