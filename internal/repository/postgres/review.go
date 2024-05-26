package postgres

import (
	"database/sql"
	"errors"
	"time"

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

	// запрос на изменение рецензии заннимает примерно 200 мс, 1 соединение может
	// обработать 5 req/sec или 300 req/min
	// запрос на вывод рецензий выполняется одноврменно с фильмом, занимает 40 мс.
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Second * 10)

	return &ReviewDB{
		DB: db,
	}, nil
}

// AddReview добавляет отзыв в базу данных.
// У переданного entity.Review должны быть заполнены поля: AuthorID, ContentID, Title, Text, Rating.
// Если операция происходит успешно, то в переданный по указателю review будут записаны ID, CreatedAt, UpdatedAt, затем
// вернется указатель на эту же рецензию
func (r *ReviewDB) AddReview(review *entity.Review) (*entity.Review, error) {
	// no-lint
	query, args, _ := sq.Insert("review").
		Columns("user_id", "content_id", "title", "text", "content_rating").
		Values(review.AuthorID, review.ContentID, review.Title, review.Text, review.Rating).
		Suffix("RETURNING id, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	err := r.DB.QueryRow(query, args...).Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case entity.PSQLCheckViolation:
				return nil, entity.NewClientError("указаны некорректные данные", entity.ErrBadRequest)
			case entity.PSQLUniqueViolation:
				return nil, entity.NewClientError("рецензия уже существует", entity.ErrAlreadyExists)
			case entity.PSQLForeignKeyViolation:
				return nil, entity.NewClientError(
					"контент с таким id не существует, либо такого пользователя не существует",
					entity.ErrNotFound,
				)
			default:
				return nil, entity.PSQLWrap(pqErr, errors.New("ошибка при добавлении рецензии"))
			}
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при добавлении рецензии"))
	}
	return review, nil
}

// GetReviewByID возвращает рецензию по ее ID
func (r *ReviewDB) GetReviewByID(id int) (*entity.Review, error) {
	// no-lint
	query, args, _ := sq.Select(
		"r.id",
		"r.user_id",
		"r.content_id",
		"r.title",
		"r.text",
		"r.content_rating",
		"r.created_at",
		"r.updated_at",
	).
		Column(sq.Expr("SUM(CASE WHEN rl.value IS TRUE THEN 1 ELSE 0 END) as likes")).
		Column(sq.Expr("SUM(CASE WHEN rl.value IS FALSE THEN 1 ELSE 0 END) as dislikes")).
		From("review r").
		LeftJoin("review_like rl ON r.id = rl.review_id").
		Where(sq.Eq{"r.id": id}).
		GroupBy("r.id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	review := &entity.Review{}
	err := r.DB.QueryRow(query, args...).
		Scan(
			&review.ID,
			&review.AuthorID,
			&review.ContentID,
			&review.Title,
			&review.Text,
			&review.Rating,
			&review.CreatedAt,
			&review.UpdatedAt,
			&review.Likes,
			&review.Dislikes,
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
	// Можно было бы денормализировать таблицы, но по тз курса СУБД нельзя
	// для оптимизации используется индекс на review_id в таблице review_like
	// no-lint
	query, args, _ := sq.Select(
		"r.id",
		"r.user_id",
		"r.content_id",
		"r.title",
		"r.text",
		"r.content_rating",
		"r.created_at",
		"r.updated_at",
	).
		Column(sq.Expr("SUM(CASE WHEN rl.value IS TRUE THEN 1 ELSE 0 END) as likes")).
		Column(sq.Expr("SUM(CASE WHEN rl.value IS FALSE THEN 1 ELSE 0 END) as dislikes")).
		Column(sq.Expr("SUM(CASE WHEN rl.value IS TRUE THEN 1 WHEN rl.value IS FALSE THEN -1 ELSE 0 END) as rating")).
		From("review r").
		LeftJoin("review_like rl ON r.id = rl.review_id").
		Where(sq.Eq{"r.content_id": contentID}).
		GroupBy("r.id").
		OrderBy("rating DESC").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении рецензий"))
	}
	defer rows.Close()

	reviews := make([]*entity.Review, 0, limit)
	for rows.Next() {
		var rating int
		review := &entity.Review{}
		err = rows.Scan(
			&review.ID,
			&review.AuthorID,
			&review.ContentID,
			&review.Title,
			&review.Text,
			&review.Rating,
			&review.CreatedAt,
			&review.UpdatedAt,
			&review.Likes,
			&review.Dislikes,
			&rating,
		)
		if err != nil {
			return nil, entity.PSQLWrap(err, errors.New("ошибка при сканировании рецензий"))
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

// UpdateReview обновляет отзыв в базе данных.
// У переданного entity.Review должны быть заполнены поля: ID, Title, Text, Rating.
func (r *ReviewDB) UpdateReview(review *entity.Review) (*entity.Review, error) {
	// no-lint
	query, args, _ := sq.Update("review").
		Set("title", review.Title).
		Set("text", review.Text).
		Set("content_rating", review.Rating).
		Where(sq.Eq{"id": review.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err := r.DB.Exec(query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case entity.PSQLCheckViolation:
				return nil, entity.NewClientError("указаны некорректные данные", entity.ErrBadRequest)
			case entity.PSQLForeignKeyViolation:
				return nil, entity.NewClientError(
					"контент с таким id не существует, либо такого пользователя не существует",
					entity.ErrNotFound,
				)
			default:
				return nil, entity.PSQLWrap(pqErr, errors.New("ошибка при обновлении рецензии"))
			}
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при обновлении рецензии"))
	}
	review, err = r.GetReviewByID(review.ID)
	if err != nil {
		return nil, err
	}
	return review, nil
}

// DeleteReviewByID удаляет отзыв по его ID
func (r *ReviewDB) DeleteReviewByID(id int) error {
	// no-lint
	query, args, _ := sq.Delete("review").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err := r.DB.Exec(query, args...)
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при удалении рецензии"))
	}
	return nil
}

// GetReviewsByAuthorID возвращает отзывы по ID автора, сортируя их по количеству лайки-дизлайки
func (r *ReviewDB) GetReviewsByAuthorID(authorID, page, limit int) ([]*entity.Review, error) {
	// no-lint
	query, args, _ := sq.Select(
		"r.id",
		"r.user_id",
		"r.content_id",
		"r.title",
		"r.text",
		"r.content_rating",
		"r.created_at",
		"r.updated_at",
	).
		Column(sq.Expr("SUM(CASE WHEN rl.value IS TRUE THEN 1 ELSE 0 END) as likes")).
		Column(sq.Expr("SUM(CASE WHEN rl.value IS FALSE THEN 1 ELSE 0 END) as dislikes")).
		Column(sq.Expr("SUM(CASE WHEN rl.value IS TRUE THEN 1 WHEN rl.value IS FALSE THEN -1 ELSE 0 END) as rating")).
		From("review r").
		LeftJoin("review_like rl ON r.id = rl.review_id").
		Where(sq.Eq{"r.user_id": authorID}).
		GroupBy("r.id").
		OrderBy("rating DESC").
		Limit(uint64(limit)).
		Offset(uint64((page - 1) * limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewClientError("рецензии не найдены", entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении рецензий"))
	}
	defer rows.Close()

	reviews := make([]*entity.Review, 0)
	for rows.Next() {
		var rating int
		review := &entity.Review{}
		err = rows.Scan(
			&review.ID,
			&review.AuthorID,
			&review.ContentID,
			&review.Title,
			&review.Text,
			&review.Rating,
			&review.CreatedAt,
			&review.UpdatedAt,
			&review.Likes,
			&review.Dislikes,
			&rating,
		)
		if err != nil {
			return nil, entity.PSQLWrap(err, errors.New("ошибка при сканировании рецензий"))
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}

func (r *ReviewDB) GetContentReviewByAuthor(authorID, contentID int) (*entity.Review, error) {
	// no-lint
	query, args, _ := sq.Select(
		"r.id",
		"r.user_id",
		"r.content_id",
		"r.title",
		"r.text",
		"r.content_rating",
		"r.created_at",
		"r.updated_at",
	).
		Column(sq.Expr("SUM(CASE WHEN rl.value IS TRUE THEN 1 ELSE 0 END) as likes")).
		Column(sq.Expr("SUM(CASE WHEN rl.value IS FALSE THEN 1 ELSE 0 END) as dislikes")).
		From("review r").
		LeftJoin("review_like rl ON r.id = rl.review_id").
		Where(sq.Eq{"r.user_id": authorID, "r.content_id": contentID}).
		GroupBy("r.id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	review := &entity.Review{}
	err := r.DB.QueryRow(query, args...).
		Scan(
			&review.ID,
			&review.AuthorID,
			&review.ContentID,
			&review.Title,
			&review.Text,
			&review.Rating,
			&review.CreatedAt,
			&review.UpdatedAt,
			&review.Likes,
			&review.Dislikes,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewClientError("рецензия не найдена", entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при получении рецензии"))
	}
	return review, nil
}

// GetAuthorRating возвращает общий рейтинг автора на основе всех его отзывов
func (r *ReviewDB) GetAuthorRating(authorID int) (int, error) {
	// no-lint
	query, args, _ := sq.Select().
		Column(sq.Expr("SUM(CASE WHEN rl.value IS TRUE THEN 1 WHEN rl.value IS FALSE THEN -1 ELSE 0 END) as user_rating")).
		From("review r").
		LeftJoin("review_like rl ON r.id = rl.review_id").
		Where(sq.Eq{"r.user_id": authorID}).
		GroupBy("r.user_id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var rating int
	err := r.DB.QueryRow(query, args...).Scan(&rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, entity.PSQLWrap(err, errors.New("ошибка при получении рейтинга автора"))
	}
	return rating, nil
}

// GetLatestReviews возвращает последние n отзывов. Не сортирует их и не даёт информацию о лайках и дизлайках
func (r *ReviewDB) GetLatestReviews(limit int) ([]*entity.Review, error) {
	// no-lint
	query, args, _ := sq.Select(
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
		OrderBy("id DESC").
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()

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
			&review.UpdatedAt,
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
	// no-lint
	query, args, _ := sq.Insert("review_like").
		Columns("review_id", "user_id", "value").
		Values(reviewID, userID, like).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err := r.DB.Exec(query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == entity.PSQLForeignKeyViolation {
				return entity.NewClientError("рецензия не найдена", entity.ErrNotFound)
			}
		}
		return entity.PSQLWrap(err, errors.New("ошибка при добавлении лайка"))
	}
	return nil
}

// UnlikeReview удаляет лайк отзыва
func (r *ReviewDB) UnlikeReview(reviewID, userID int) error {
	// no-lint
	query, args, _ := sq.Delete("review_like").
		Where(sq.Eq{"review_id": reviewID, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	rows, err := r.DB.Exec(query, args...)
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при удалении лайка"))
	}
	affected, err := rows.RowsAffected()
	if err != nil {
		return entity.PSQLWrap(err, errors.New("ошибка при получении количества затронутых строк"))
	}
	if affected == 0 {
		return entity.NewClientError("лайк не найден", entity.ErrNotFound)
	}
	return nil
}

// IsLikedByUser возвращает 1, если отзыв лайкнут, -1, если дизлайкнут, 0, если не оценен
func (r *ReviewDB) IsLikedByUser(reviewID, userID int) (int, error) {
	// no-lint
	query, args, _ := sq.Select("value").
		From("review_like").
		Where(sq.Eq{"review_id": reviewID, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var value bool
	err := r.DB.QueryRow(query, args...).Scan(&value)
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

// GetContentRating возвращает рейтинг контента на основе всех рецензий. Если нет ни одной рецензии, использует IMDB
func (r *ReviewDB) GetContentRating(contentID int) (float64, error) {
	query, args, _ := sq.Select("AVG(content_rating)").
		From("review").
		Where(sq.Eq{"content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var rating sql.NullFloat64
	err := r.DB.QueryRow(query, args...).Scan(&rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, entity.PSQLWrap(err, errors.New("ошибка при получении рейтинга контента"))
	}
	// Если рецензий нет, то используем IMDB
	if rating.Float64 == .0 {
		return r.getIMDBRating(contentID)
	}
	return rating.Float64, nil
}

func (r *ReviewDB) getIMDBRating(contentID int) (float64, error) {
	query, args, _ := sq.Select("imdb").
		From("content").
		Where(sq.Eq{"id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	var rating float64
	err := r.DB.QueryRow(query, args...).Scan(&rating)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, entity.NewClientError("контент не найден", entity.ErrNotFound)
		}
		return 0, entity.PSQLWrap(err, errors.New("ошибка при получении рейтинга контента"))
	}
	return rating, nil
}
