package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type FavouriteDB struct {
	DB *sqlx.DB
}

func NewFavouriteRepository(db *sqlx.DB) repository.Favourite {
	return &FavouriteDB{
		DB: db,
	}
}

func (f FavouriteDB) CreateFavourite(userID, contentID int, category string) error {
	query, args, err := sq.
		Insert("favourite").
		Columns("user_id", "content_id", "category").
		Values(userID, contentID, category).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при формировании запроса CreateFavourite"))
	}
	_, err = f.DB.Exec(query, args...)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case entity.PSQLUniqueViolation:
			// такой кейс игнорируем
			return nil
		case entity.PSQLForeignKeyViolation:
			fmt.Println(userID)
			fmt.Println(pqErr)
			return repository.ErrFavouriteContentNotFound
		}
	}
	if err != nil {
		return entity.PSQLQueryErr("CreateFavourite", err)
	}
	return nil
}

func (f FavouriteDB) DeleteFavourite(userID, contentID int) error {
	query, args, err := sq.
		Delete("favourite").
		Where(sq.Eq{"user_id": userID, "content_id": contentID}).
		Suffix("RETURNING content_id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при формировании запроса DeleteFavourite"))
	}
	row := f.DB.QueryRow(query, args...)
	var deletedContentID int
	err = row.Scan(&deletedContentID)
	if errors.Is(err, sql.ErrNoRows) {
		return repository.ErrFavouriteNotFound
	}
	if err != nil {
		return entity.PSQLQueryErr("DeleteFavourite", err)
	}
	return nil
}

func (f FavouriteDB) GetFavourites(userID int) ([]*entity.Favourite, error) {
	query, args, err := sq.
		Select("content_id", "category").
		From("favourite").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании запроса GetFavourites"))
	}
	rows, err := f.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("GetFavourites", err)
	}
	defer rows.Close()
	var favourites []*entity.Favourite
	for rows.Next() {
		var favourite entity.Favourite
		err = rows.Scan(&favourite.ContentID, &favourite.Category)
		if err != nil {
			return nil, entity.PSQLQueryErr("GetFavourites при сканировании", err)
		}
		favourite.UserID = userID
		favourites = append(favourites, &favourite)
	}
	return favourites, nil
}

func (f FavouriteDB) GetFavourite(userID, contentID int) (*entity.Favourite, error) {
	query, args, err := sq.
		Select("category").
		From("favourite").
		Where(sq.Eq{"user_id": userID, "content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании запроса GetStatus"))
	}
	row := f.DB.QueryRow(query, args...)
	var favourite entity.Favourite
	err = row.Scan(&favourite.Category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrFavouriteNotFound
		}
		return nil, entity.PSQLQueryErr("GetStatus", err)
	}
	favourite.UserID = userID
	favourite.ContentID = contentID
	return &favourite, nil
}
