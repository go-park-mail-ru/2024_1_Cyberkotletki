package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
)

type OngoingContentDB struct {
	DB *sql.DB
}

type ScanOngoingContent struct {
	ID             int
	ContentType    string
	Title          string
	PosterStaticID sql.NullInt64
	ReleaseDate    time.Time
}

func scanAllFieldsOnoingContent(row sq.RowScanner) (*entity.OngoingContent, error) {
	ongoingContent := &entity.OngoingContent{}
	err := row.Scan(
		&ongoingContent.ID,
		&ongoingContent.Type,
		&ongoingContent.Title,
		&ongoingContent.ReleaseDate,
	)
	return ongoingContent, err
}

func NewOngoingContentRepository(db *sql.DB) repository.OngoingContent {
	return &OngoingContentDB{
		DB: db,
	}
}

func (oc *OngoingContentDB) getOngoingMovieData(id int) (*entity.OngoingMovie, error) {
	query, args, err := sq.Select("premiere, duration").
		From("ongoing_movie").
		Where(sq.Eq{"ongoing_content_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getOngoingMovieData"))
	}
	var movie entity.OngoingMovie
	var premiere sql.NullTime
	var duration sql.NullInt64
	err = oc.DB.QueryRow(query, args...).Scan(&premiere, &duration)
	if err != nil {
		return nil, entity.PSQLQueryErr("getOngoingMovieData", err)
	}
	movie.Premiere = premiere.Time
	movie.Duration = int(duration.Int64)
	return &movie, nil
}

func (oc *OngoingContentDB) getOngoingSeriesData(id int) (*entity.OngoingSeries, error) {
	query, args, err := sq.Select("year_start, year_end").
		From("ongoing_series").
		Where(sq.Eq{"ongoing_content_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getOngoingSeriesData"))
	}
	var series entity.OngoingSeries
	var yearStart sql.NullInt64
	var yearEnd sql.NullInt64
	err = oc.DB.QueryRow(query, args...).Scan(&yearStart, &yearEnd)
	if err != nil {
		return nil, entity.PSQLQueryErr("getOngoingSeriesData", err)
	}
	series.YearStart = int(yearStart.Int64)
	series.YearEnd = int(yearEnd.Int64)
	return &series, nil
}

func (oc *OngoingContentDB) getGenreByID(genreID int) (*entity.Genre, error) {
	query, args, err := sq.Select("name").
		From("genre").
		Where(sq.Eq{"id": genreID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var genre entity.Genre
	err = oc.DB.QueryRow(query, args...).Scan(&genre.Name)
	genre.ID = genreID
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при запросе getGenreByID"))
	}
	return &genre, nil
}

func (oc *OngoingContentDB) getOngoingContentGenres(contentID int) ([]entity.Genre, error) {
	query, args, err := sq.Select("genre_id").
		From("ongoing_content_genre").
		Where(sq.Eq{"ongoing_content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getOngoingContentGenres"))
	}
	rows, err := oc.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("getOngoingContentGenres", err)
	}
	defer rows.Close()
	genres := make([]entity.Genre, 0)
	for rows.Next() {
		var genreID int
		err := rows.Scan(&genreID)
		if err != nil {
			return nil, entity.PSQLQueryErr("getOngoingContentGenres при сканировании жанров", err)
		}
		genre, err := oc.getGenreByID(genreID)
		if err != nil {
			return nil, err
		}
		genres = append(genres, *genre)
	}
	return genres, nil
}

func (oc *OngoingContentDB) getOngoingContentInfo(contentID int) (*entity.OngoingContent, error) {
	query, args, err := sq.Select("id", "content_type, title, poster_static_id, release_date").
		From("ongoing_content").
		Where(sq.Eq{"id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getOngoingContentInfo"))
	}
	var scanOngoingContent ScanOngoingContent
	var ongoingContent entity.OngoingContent
	err = oc.DB.QueryRow(query, args...).Scan(&scanOngoingContent.ID,
		&scanOngoingContent.ContentType,
		&scanOngoingContent.Title,
		&scanOngoingContent.PosterStaticID,
		&scanOngoingContent.ReleaseDate,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrOngoingContentNotFound
		}
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении контента календаря релизов"))
	}
	ongoingContent.ID = scanOngoingContent.ID
	ongoingContent.Title = scanOngoingContent.Title
	ongoingContent.PosterStaticID = int(scanOngoingContent.PosterStaticID.Int64)
	ongoingContent.ReleaseDate = scanOngoingContent.ReleaseDate
	ongoingContent.Type = scanOngoingContent.ContentType

	return &ongoingContent, nil
}

func (oc *OngoingContentDB) GetOngoingContentByID(id int) (*entity.OngoingContent, error) {
	ongoingContent, err := oc.getOngoingContentInfo(id)
	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	var occurredErrorsChan = make(chan error, 2)
	wg.Add(1)
	go func() {
		defer wg.Done()
		genres, err := oc.getOngoingContentGenres(id)
		if err != nil {
			occurredErrorsChan <- err
			return
		}
		ongoingContent.Genres = genres
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		switch ongoingContent.Type {
		case entity.OngoingContentTypeMovie:
			movie, err := oc.getOngoingMovieData(id)
			if err != nil {
				occurredErrorsChan <- err
				return
			}
			ongoingContent.Movie = movie
		case entity.OngoingContentTypeSeries:
			series, err := oc.getOngoingSeriesData(id)
			if err != nil {
				occurredErrorsChan <- err
				return
			}
			ongoingContent.Series = series
		}
	}()
	wg.Wait()
	if len(occurredErrorsChan) > 0 {
		return nil, <-occurredErrorsChan
	}
	return ongoingContent, nil
}

// GetNearestOngoings возвращает ближайшие релизы
func (oc *OngoingContentDB) GetNearestOngoings(limit int) ([]*entity.OngoingContent, error) {
	query, args, err := selectAllFields().
		From("ongoing_content").
		Where(sq.Gt{"release_date": time.Now()}).
		OrderBy("release_date").
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса GetNearestOngoings"))
	}
	rows, err := oc.DB.Query(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrOngoingContentNotFound
		}
		return nil, entity.PSQLQueryErr("GetNearestOngoings", err)
	}
	defer rows.Close()
	ongoingContents := make([]*entity.OngoingContent, 0)
	for rows.Next() {
		ongoingContent, err := scanAllFieldsOnoingContent(rows)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, repository.ErrOngoingContentNotFound
			}
			return nil, entity.PSQLQueryErr("GetNearestOngoings при сканировании клендаря релизов", err)
		}
		ongoingContents = append(ongoingContents, ongoingContent)
	}
	if err = rows.Err(); err != nil {
		return nil, entity.PSQLQueryErr("GetNearestOngoings при получении результатов", err)
	}
	return ongoingContents, nil
}

// GetOngoingContentByMonthAndYear возвращает релизы по месяцу и году
func (oc *OngoingContentDB) GetOngoingContentByMonthAndYear(month, year int) ([]*entity.OngoingContent, error) {
	query, args, err := selectAllFields().
		From("ongoing_content").
		Where(sq.Eq{"extract(month from release_date)": month, "extract(year from release_date)": year}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса GetOngoingContentByMonthAndYear"))
	}
	rows, err := oc.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("GetOngoingContentByMonthAndYear", err)
	}
	defer rows.Close()
	ongoingContents := make([]*entity.OngoingContent, 0)
	for rows.Next() {
		ongoingContent, err := scanAllFieldsOnoingContent(rows)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, repository.ErrOngoingContentNotFound
			}
			return nil, entity.PSQLQueryErr("GetOngoingContentByMonthAndYear при сканировании клендаря релизов", err)
		}
		ongoingContents = append(ongoingContents, ongoingContent)
	}
	return ongoingContents, err
}

// GetAllReleaseYears возвращает все года релизов
func (oc *OngoingContentDB) GetAllReleaseYears() ([]int, error) {
	query, args, err := sq.Select("extract(year from release_date)").
		From("ongoing_content").
		GroupBy("extract(year from release_date)").
		OrderBy("extract(year from release_date)").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса GetAllReleaseYears"))
	}
	rows, err := oc.DB.Query(query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrOngoingContentYearsNotFound
		}
		return nil, entity.PSQLQueryErr("GetAllReleaseYears", err)
	}
	defer rows.Close()
	years := make([]int, 0)
	for rows.Next() {
		var year int
		err := rows.Scan(&year)
		if err != nil {
			return nil, entity.PSQLQueryErr("GetAllReleaseYears при сканировании года релизов", err)
		}
		years = append(years, year)
	}
	return years, nil
}

// IsOngoingConentFinished возвращает true, если контент завершен
func (oc *OngoingContentDB) IsOngoingConentFinished(contentID int) (bool, error) {
	query, args, err := sq.Select("release_date").
		From("ongoing_content").
		Where(sq.Eq{"id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return false, entity.PSQLWrap(err, errors.New("ошибка при составлении запроса IsOngoingConentFinished"))
	}
	var releaseDate time.Time
	err = oc.DB.QueryRow(query, args...).Scan(&releaseDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, repository.ErrOngoingContentNotFound
		}
		return false, entity.PSQLQueryErr("IsOngoingConentFinished", err)
	}
	if time.Now().After(releaseDate) {
		return true, nil
	}
	return false, nil
}
