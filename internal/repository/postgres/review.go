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

type ReviewDB struct {
	DB *sql.DB
}

func NewReviewRepository(database config.PostgresDatabase) (repository.Review, error) {
	db, err := sql.Open("postgres", database.ConnectURL)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &ReviewDB{
		DB: db,
	}, nil
}

// AddReview добавляет отзыв в базу данных.
// У переданного entity.Review должны быть заполнены поля: AuthorID, ContentID, Title, Text, Rating.
// Если операция происходит успешно, то в переданный по указателю review будут записаны ID и CreatedAt, затем
// вернется указатель на эту же рецензию
func (r *ReviewDB) AddReview(review *entity.Review) (*entity.Review, error) {
	query, args, err := sq.Insert("review").
		Columns("user_id", "content_id", "title", "text", "content_rating").
		Values(review.AuthorID, review.ContentID, review.Title, review.Text, review.Rating).
		Suffix("RETURNING id, updated_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса AddReview"))
	}

	err = r.DB.QueryRow(query, args...).Scan(&review.ID, &review.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case entity.PSQLCheckViolation:
				return nil, entity.NewClientError("указаны некорректные данные", entity.ErrBadRequest)
			case entity.PSQLUniqueViolation:
				return nil, entity.NewClientError("рецензия уже существует", entity.ErrAlreadyExists)
			case entity.PSQLForeignKeyViolation:
				return nil, entity.NewClientError("контент с таким id не существует", entity.ErrNotFound)
			default:
				return nil, entity.PSQLWrap(err, errors.New("ошибка при добавлении рецензии"), pqErr)
			}
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при добавлении рецензии"))
	}
	return review, nil
}

// GetReviewByID возвращает рецензию по ее ID
func (r *ReviewDB) GetReviewByID(id int) (*entity.Review, error) {
	query, args, err := sq.Select(
		"id",
		"user_id",
		"content_id",
		"title",
		"text",
		"content_rating",
		"created_at",
		"updated_at",
	).
		From("review").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса GetReviewByID"))
	}

	review := &entity.Review{}
	err = r.DB.QueryRow(query, args...).
		Scan(
			&review.ID,
			&review.AuthorID,
			&review.ContentID,
			&review.Title,
			&review.Text,
			&review.Rating,
			&review.CreatedAt,
			&review.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewClientError("рецензия не найдена", entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении рецензии"))
	}
	return review, nil
}

// GetReviewsByContentID возвращает все рецензии по ID контента, сортируя их по рейтингу
func (r *ReviewDB) GetReviewsByContentID(contentID, page, limit int) ([]*entity.Review, error) {
	// 3НДФ никого не щадит...
	// можно было бы денормализировать таблицы, то по тз курса СУБД нельзя
	// для оптимизации используется индекс на review_id в таблице review_like
	query, args, err := sq.Select("r.id", "r.user_id", "r.content_id", "r.title", "r.text", "r.content_rating").
		Column(sq.Expr("SUM(CASE WHEN rl.value THEN 1 ELSE 0 END) as likes")).
		Column(sq.Expr("SUM(CASE WHEN rl.value THEN 0 ELSE 1 END) as dislikes")).
		Column(sq.Expr("SUM(CASE WHEN rl.value THEN 1 ELSE -1 END) as rating")).
		From("review r").
		LeftJoin("review_like rl ON r.id = rl.review_id").
		Where(sq.Eq{"r.content_id": contentID}).
		GroupBy("r.id").
		OrderBy("rating DESC").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса GetReviewsByContentID"))
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении рецензий"))
	}
	defer rows.Close()

	reviews := make([]*entity.Review, 0, limit)
	for rows.Next() {
		review := &entity.Review{}
		err = rows.Scan(
			&review.ID,
			&review.AuthorID,
			&review.ContentID,
			&review.Title,
			&review.Text,
			&review.Rating,
			&review.Likes,
			&review.Dislikes,
		)
		if err != nil {
			return nil, entity.PSQLWrap(err, errors.New("ошибка при сканировании рецензий"))
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

// UpdateReview обновляет отзыв в базе данных.
// У переданного entity.Review должны быть заполнены поля: ID, ContentID, Title, Text, Rating.
// Обратно возвращается указатель на эту же рецензию без изменений
func (r *ReviewDB) UpdateReview(review *entity.Review) (*entity.Review, error) {
	query, args, err := sq.Update("review").
		Set("content_id", review.ContentID).
		Set("title", review.Title).
		Set("text", review.Text).
		Set("content_rating", review.Rating).
		Where(sq.Eq{"id": review.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса UpdateReview"))
	}

	_, err = r.DB.Exec(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при обновлении рецензии"))
	}
	return review, nil
}

// DeleteReviewByID удаляет отзыв по его ID
func (r *ReviewDB) DeleteReviewByID(id int) error {
	query, args, err := sq.Delete("review").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса DeleteReviewByID"))
	}

	_, err = r.DB.Exec(query, args...)
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при удалении рецензии"))
	}
	return nil
}

// GetReviewsByAuthorID возвращает отзывы по ID автора
func (r *ReviewDB) GetReviewsByAuthorID(authorID, page, limit int) ([]*entity.Review, error) {
	query, args, err := sq.Select(
		"id", "user_id", "content_id", "title", "text", "content_rating", "created_at",
	).
		From("review").
		OrderBy("id DESC").
		Where(sq.Eq{"user_id": authorID}).
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса GetReviewsByAuthorID"))
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении рецензий"))
	}
	defer rows.Close()

	reviews := make([]*entity.Review, 0)
	for rows.Next() {
		review := &entity.Review{}
		err = rows.Scan(
			&review.ID,
			&review.AuthorID,
			&review.ContentID,
			&review.Title,
			&review.Text,
			&review.Rating,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, entity.PSQLWrap(err, errors.New("ошибка при сканировании рецензий"))
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

// GetAuthorRating возвращает общий рейтинг автора на основе всех его отзывов
func (r *ReviewDB) GetAuthorRating(authorID int) (int, error) {
	query, args, err := sq.Select("r.user_id").
		Column(sq.Expr("SUM(CASE WHEN rl.value THEN 1 ELSE -1 END) as user_rating")).
		From("review r").
		LeftJoin("review_like rl ON r.id = rl.review_id").
		Where(sq.Eq{"r.user_id": authorID}).
		GroupBy("r.user_id").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса GetAuthorRating"))
	}

	var rating int
	err = r.DB.QueryRow(query, args...).Scan(&rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, entity.PSQLWrap(err, errors.New("ошибка при получении рейтинга автора"))
	}
	return rating, nil
}

// GetLatestReviews возвращает последние n отзывов
func (r *ReviewDB) GetLatestReviews(limit int) ([]*entity.Review, error) {
	query, args, err := sq.Select("id", "user_id", "content_id", "title", "text", "content_rating", "created_at").
		From("review").
		OrderBy("id DESC").
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса GetLatestReviews"))
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении рецензий"))
	}
	defer rows.Close()

	reviews := make([]*entity.Review, 0, limit)
	for rows.Next() {
		review := &entity.Review{}
		err = rows.Scan(
			&review.ID,
			&review.AuthorID,
			&review.ContentID,
			&review.Title,
			&review.Text,
			&review.Rating,
			&review.CreatedAt,
		)
		if err != nil {
			return nil, entity.PSQLWrap(err, errors.New("ошибка при сканировании рецензий"))
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

// LikeReview добавляет лайк к отзыву
func (r *ReviewDB) LikeReview(reviewID, userID int, like bool) error {
	query, args, err := sq.Insert("review_like").
		Columns("review_id", "user_id", "value").
		Values(reviewID, userID, like).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса LikeReview"))
	}

	_, err = r.DB.Exec(query, args...)
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при добавлении лайка"))
	}
	return nil
}

// UnlikeReview удаляет лайк отзыва
func (r *ReviewDB) UnlikeReview(reviewID, userID int) error {
	query, args, err := sq.Delete("review_like").
		Where(sq.Eq{"review_id": reviewID, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса UnlikeReview"))
	}

	_, err = r.DB.Exec(query, args...)
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при удалении лайка"))
	}
	return nil
}

// IsLikedByUser возвращает 1, если отзыв лайкнут, -1, если дизлайкнут, 0, если не оценен
func (r *ReviewDB) IsLikedByUser(reviewID, userID int) (int, error) {
	query, args, err := sq.Select("value").
		From("review_like").
		Where(sq.Eq{"review_id": reviewID, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса IsLikedOrDisliked"))
	}

	var value bool
	err = r.DB.QueryRow(query, args...).Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, entity.PSQLWrap(err, errors.New("ошибка при получении лайка"))
	}
	if value {
		return 1, nil
	}
	return -1, nil
}

// DeleteReview удаляет отзыв по его ID
func (r *ReviewDB) DeleteReview(reviewID int) error {
	query, args, err := sq.Delete("review").
		Where(sq.Eq{"id": reviewID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса DeleteReview"))
	}

	_, err = r.DB.Exec(query, args...)
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при удалении рецензии"))
	}
	return nil
}

// GetContentRating возвращает рейтинг контента на основе всех рецензий
func (r *ReviewDB) GetContentRating(contentID int) (int, error) {
	query, args, err := sq.Select("AVG(content_rating)").
		From("review").
		Where(sq.Eq{"content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, entity.PSQLWrap(err, errors.New("ошибка при формировании sql-запроса GetContentRating"))
	}

	var rating int
	err = r.DB.QueryRow(query, args...).Scan(&rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, entity.PSQLWrap(err, errors.New("ошибка при получении рейтинга контента"))
	}
	return rating, nil
}
