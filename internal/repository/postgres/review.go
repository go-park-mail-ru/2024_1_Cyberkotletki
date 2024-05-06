package postgres

import (
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ReviewDB struct {
	DB *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) repository.Review {
	return &ReviewDB{
		DB: db,
	}
}

func selectAllFields() sq.SelectBuilder {
	return sq.Select(
		"id",
		"user_id",
		"content_id",
		"title",
		"text",
		"content_rating",
		"created_at",
		"updated_at",
		"likes",
		"dislikes",
		"rating",
	)
}

func scanAllFields(row sq.RowScanner) (*entity.Review, error) {
	review := &entity.Review{}
	err := row.Scan(
		&review.ID,
		&review.AuthorID,
		&review.ContentID,
		&review.Title,
		&review.Text,
		&review.ContentRating,
		&review.CreatedAt,
		&review.UpdatedAt,
		&review.Likes,
		&review.Dislikes,
		&review.Rating,
	)
	return review, err
}

func scanRows(rows *sqlx.Rows) ([]*entity.Review, error) {
	reviews := make([]*entity.Review, 0)
	for rows.Next() {
		reviewEntity := new(entity.Review)
		err := rows.StructScan(reviewEntity)
		if err != nil {
			return nil, entity.PSQLQueryErr("scanRows при сканировании рецензий", err)
		}
		reviews = append(reviews, reviewEntity)
	}
	return reviews, nil
}

// GetLatestReviews возвращает последние добавленные рецензии
func (r *ReviewDB) GetLatestReviews(limit int) ([]*entity.Review, error) {
	query, args, err := selectAllFields().
		From("review").
		OrderBy("created_at DESC", "id ASC").
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса GetLatestReviews"))
	}
	rows, err := r.DB.Queryx(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("GetLatestReviews", err)
	}
	defer rows.Close()
	return scanRows(rows)
}

// AddReview добавляет отзыв в базу данных.
// У переданного entity.Review должны быть заполнены поля: AuthorID, ContentID, Title, Text, ContentRating.
// Если операция происходит успешно, то в переданный по указателю review будут записаны ID, CreatedAt, UpdatedAt, затем
// вернется указатель на эту же рецензию
func (r *ReviewDB) AddReview(review *entity.Review) (*entity.Review, error) {
	query, args, err := sq.Insert("review").
		Columns("user_id", "content_id", "title", "text", "content_rating").
		Values(review.AuthorID, review.ContentID, review.Title, review.Text, review.ContentRating).
		Suffix("RETURNING id, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса AddReview"))
	}
	err = r.DB.QueryRow(query, args...).Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case entity.PSQLCheckViolation:
			return nil, repository.ErrReviewBadRequest
		case entity.PSQLUniqueViolation:
			return nil, repository.ErrReviewAlreadyExists
		case entity.PSQLForeignKeyViolation:
			return nil, repository.ErrReviewViolation
		}
	}
	if err != nil {
		return nil, entity.PSQLQueryErr("AddReview", err)
	}
	return review, nil
}

// GetReviewByID возвращает рецензию по ее ID
func (r *ReviewDB) GetReviewByID(id int) (*entity.Review, error) {
	query, args, err := selectAllFields().
		From("review").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса GetReviewByID"))
	}
	review, err := scanAllFields(r.DB.QueryRow(query, args...))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrReviewNotFound
		}
		return nil, entity.PSQLQueryErr("GetReviewByID", err)
	}
	return review, nil
}

// GetReviewsCountByContentID возвращает количество рецензий по ID контента
func (r *ReviewDB) GetReviewsCountByContentID(contentID int) (int, error) {
	query, args, err := sq.Select("COUNT(*)").
		From("review").
		Where(sq.Eq{"content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса GetReviewsCountByContentID"))
	}
	var count int
	err = r.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, entity.PSQLQueryErr("GetReviewsCountByContentID", err)
	}
	return count, nil
}

// GetReviewsByContentID возвращает все рецензии по ID контента, сортируя их по рейтингу
func (r *ReviewDB) GetReviewsByContentID(contentID, page, limit int) ([]*entity.Review, error) {
	query, args, err := selectAllFields().
		From("review").
		Where(sq.Eq{"content_id": contentID}).
		OrderBy("rating DESC", "id ASC").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса GetReviewsByContentID"))
	}
	rows, err := r.DB.Queryx(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("GetReviewsByContentID", err)
	}
	defer rows.Close()
	return scanRows(rows)
}

// UpdateReview обновляет отзыв в базе данных.
// У переданного entity.Review должны быть заполнены поля: ID, Title, Text, ContentRating.
// Возвращает repository.ErrReviewBadRequest, если обновлённых строк нет
func (r *ReviewDB) UpdateReview(review *entity.Review) error {
	query, args, err := sq.Update("review").
		Set("title", review.Title).
		Set("text", review.Text).
		Set("content_rating", review.ContentRating).
		Where(sq.Eq{"id": review.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при составлении запроса UpdateReview"))
	}
	_, err = r.DB.Exec(query, args...)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case entity.PSQLCheckViolation:
			return repository.ErrReviewBadRequest
		case entity.PSQLForeignKeyViolation:
			return repository.ErrReviewNotFound
		}
	}
	if err != nil {
		return entity.PSQLQueryErr("UpdateReview", err)
	}
	return nil
}

// DeleteReviewByID удаляет отзыв по его ID
func (r *ReviewDB) DeleteReviewByID(id int) error {
	query, args, err := sq.Delete("review").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при составлении запроса DeleteReviewByID"))
	}
	_, err = r.DB.Exec(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return repository.ErrReviewNotFound
		}
		return entity.PSQLQueryErr("DeleteReviewByID", err)
	}
	return nil
}

// GetReviewsCountByAuthorID возвращает количество отзывов по ID автора
func (r *ReviewDB) GetReviewsCountByAuthorID(authorID int) (int, error) {
	query, args, err := sq.Select("COUNT(*)").
		From("review").
		Where(sq.Eq{"user_id": authorID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса GetReviewsCountByAuthorID"))
	}
	var count int
	err = r.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, entity.PSQLQueryErr("GetReviewsCountByAuthorID", err)
	}
	return count, nil
}

// GetReviewsByAuthorID возвращает отзывы по ID автора, сортируя их по дате добавления
func (r *ReviewDB) GetReviewsByAuthorID(authorID, page, limit int) ([]*entity.Review, error) {
	query, args, err := selectAllFields().
		From("review").
		Where(sq.Eq{"user_id": authorID}).
		OrderBy("created_at DESC", "id ASC").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса GetReviewsByAuthorID"))
	}
	rows, err := r.DB.Queryx(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("GetReviewsByAuthorID", err)
	}
	defer rows.Close()
	return scanRows(rows)
}

func (r *ReviewDB) GetContentReviewByAuthor(authorID, contentID int) (*entity.Review, error) {
	query, args, err := selectAllFields().
		From("review").
		Where(sq.Eq{"user_id": authorID, "content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса GetContentReviewByAuthor"))
	}
	review, err := scanAllFields(r.DB.QueryRow(query, args...))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrReviewNotFound
		}
		return nil, entity.PSQLQueryErr("GetContentReviewByAuthor", err)
	}
	return review, nil
}

// VoteReview добавляет оценку к отзыву
func (r *ReviewDB) VoteReview(reviewID, userID int, vote bool) error {
	query, args, err := sq.Insert("review_vote").
		Columns("review_id", "user_id", "value").
		Values(reviewID, userID, vote).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при составлении запроса VoteReview"))
	}
	_, err = r.DB.Exec(query, args...)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case entity.PSQLForeignKeyViolation:
			return repository.ErrReviewNotFound
		case entity.PSQLUniqueViolation:
			return repository.ErrReviewVoteAlreadyExists
		}
	}
	if err != nil {
		return entity.PSQLQueryErr("VoteReview", err)
	}
	return nil
}

// UnVoteReview удаляет оценку с отзыва
func (r *ReviewDB) UnVoteReview(reviewID, userID int) error {
	query, args, err := sq.Delete("review_vote").
		Where(sq.Eq{"review_id": reviewID, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при составлении запроса UnVoteReview"))
	}
	_, err = r.DB.Exec(query, args...)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case entity.PSQLForeignKeyViolation:
			return repository.ErrReviewVoteNotFound
		}
	}
	if err != nil {
		return entity.PSQLQueryErr("UnVoteReview", err)
	}
	return nil
}

// IsVotedByUser возвращает 1, если на отзыв поставлен лайк, и возвращает -1, если на отзыв поставлен дизлайк.
// Если пользователь не оценивал отзыв, возвращает 0
func (r *ReviewDB) IsVotedByUser(reviewID, userID int) (int, error) {
	query, args, err := sq.Select("value").
		From("review_vote").
		Where(sq.Eq{"review_id": reviewID, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса IsVotedByUser"))
	}
	var value bool
	err = r.DB.QueryRow(query, args...).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, entity.PSQLQueryErr("IsVotedByUser", err)
	}
	if value {
		return 1, nil
	}
	return -1, nil
}
