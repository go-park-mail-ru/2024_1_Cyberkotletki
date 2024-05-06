package postgres

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/stretchr/testify/require"
)

var releaseDateTmp = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

func getDriverValuesOngoing(args []any) []driver.Value {
	driverValues := make([]driver.Value, len(args))
	for i, v := range args {
		driverValues[i] = v
	}
	return driverValues
}

func setupGetOngoingContentSuccess(mock sqlmock.Sqlmock, contentID int, contentType string) {
	query, args, _ := sq.Select(
		"id",
		"content_type",
		"title",
		"poster_static_id",
		"release_date",
	).
		From("ongoing_content").
		Where(sq.Eq{"id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{
			"id",
			"content_type",
			"title",
			"poster_static_id",
			"release_date",
		}).AddRow(
			contentID,
			contentType,
			"title",
			1,
			releaseDateTmp,
		))
}
func setupGetOngoingContentSuccessTime(mock sqlmock.Sqlmock, contentID int, releaseDate time.Time) {
	query, args, _ := sq.Select(
		"id",
		"content_type",
		"title",
		"poster_static_id",
		"release_date",
	).
		From("ongoing_content").
		Where(sq.Eq{"id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{
			"id",
			"content_type",
			"title",
			"poster_static_id",
			"release_date",
		}).AddRow(
			contentID,
			"movie",
			"title",
			1,
			releaseDate,
		))
}

func setupGetOngoingContentGenresSuccess(mock sqlmock.Sqlmock, contentID int, genresID []int) {
	query, args, _ := sq.Select("genre_id").
		From("ongoing_content_genre").
		Where(sq.Eq{"ongoing_content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows := sqlmock.NewRows([]string{"genre_id"})
	for _, genreID := range genresID {
		rows.AddRow(genreID)
	}
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(rows)
}

func setupGetOngoingGenreByIDSuccess(mock sqlmock.Sqlmock, genreID int, genreName string) {
	query, args, _ := sq.Select("name").
		From("genre").
		Where(sq.Eq{"id": genreID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{"name"}).AddRow(genreName))
}

func TestOngoingContentRepository_GetOngoingContentByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Request     int
		ExpectedOut *entity.OngoingContent
		ExpectedErr error
		SetupMock   func(mock sqlmock.Sqlmock)
	}{
		{
			Name:    "Успешное получение фильма",
			Request: 1,
			ExpectedOut: &entity.OngoingContent{
				ID:             1,
				Title:          "title",
				PosterStaticID: 1,
				ReleaseDate:    releaseDateTmp,
				Type:           "movie",
				Genres:         []entity.Genre{{ID: 1, Name: "Action"}},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock) {
				// устанавливаем порядок вызовов и мокаем их
				setupGetOngoingContentSuccess(mock, 1, "movie")
				setupGetOngoingContentGenresSuccess(mock, 1, []int{1})
				setupGetOngoingGenreByIDSuccess(mock, 1, "Action")
			},
		},
		{
			Name:    "Успешное получение сериала",
			Request: 1,
			ExpectedOut: &entity.OngoingContent{
				ID:             1,
				Title:          "title",
				PosterStaticID: 1,
				ReleaseDate:    releaseDateTmp,
				Type:           "series",
				Genres:         []entity.Genre{{ID: 1, Name: "Action"}},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock) {
				// устанавливаем порядок вызовов и мокаем их
				setupGetOngoingContentSuccess(mock, 1, "series")
				setupGetOngoingContentGenresSuccess(mock, 1, []int{1})
				setupGetOngoingGenreByIDSuccess(mock, 1, "Action")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			mock.MatchExpectationsInOrder(false)
			require.NoError(t, err)
			repo := NewOngoingContentRepository(db)
			tc.SetupMock(mock)
			output, err := repo.GetOngoingContentByID(tc.Request)
			if !errors.Is(err, tc.ExpectedErr) {
				require.Fail(t, fmt.Errorf("unexpected error, expected: %v, got: %v", tc.ExpectedErr, err).Error())
			}
			require.Equal(t, tc.ExpectedOut, output)
		})
	}
}

/*
func TestOngoingContentRepository_GetNearestOngoings(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Request     int
		ExpectedOut []*entity.OngoingContent
		ExpectedErr error
		SetupMock   func(mock sqlmock.Sqlmock, query string, args []driver.Value)
	}{
		{
			Name:    "Успешное получение ближайших премьер",
			Request: 2,
			ExpectedOut: []*entity.OngoingContent{
				{
					ID:             1,
					Title:          "title",
					PosterStaticID: 1,
					ReleaseDate:    releaseDateTmp,
					Type:           "movie",
					Movie: &entity.OngoingMovie{
						Premiere: time.Time{},
						Duration: 100,
					},
					Series: nil,
					Genres: []entity.Genre{{ID: 1, Name: "Action"}},
				},
				{
					ID:             2,
					Title:          "title",
					PosterStaticID: 1,
					ReleaseDate:    releaseDateTmp,
					Type:           "series",
					Movie:          nil,
					Series: &entity.OngoingSeries{
						YearStart: 1980,
						YearEnd:   1981,
					},
					Genres: []entity.Genre{{ID: 1, Name: "Action"}},
				},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(
					sqlmock.NewRows([]string{
						"id",
						"content_type",
						"title",
						"poster_static_id",
						"release_date",
					}).AddRow(
						1,
						entity.OngoingContentTypeMovie,
						"title",
						1,
						releaseDateTmp,
					).AddRow(
						2,
						entity.OngoingContentTypeSeries,
						"title",
						1,
						releaseDateTmp,
					))
				setupGetOngoingMovieDataSuccess(mock, 1, time.Time{}, 100)
				setupGetOngoingSeriesDataSuccess(mock, 2, 1980, 1981)
				setupGetOngoingContentGenresSuccess(mock, 1, []int{1})
				setupGetOngoingContentGenresSuccess(mock, 2, []int{1})
				setupGetOngoingGenreByIDSuccess(mock, 1, "Action")
			},
		},
		{
			Name:        "Успешное получение ближайших премьер, но их нет",
			Request:     2,
			ExpectedOut: nil,
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(args...).WillReturnRows(
					sqlmock.NewRows([]string{}))
			},
		},
		{
			Name:        "Ошибка при получении ближайших премьер",
			Request:     2,
			ExpectedOut: nil,
			ExpectedErr: errors.New("error"),
			SetupMock: func(mock sqlmock.Sqlmock, query string, args []driver.Value) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(args...).
					WillReturnError(fmt.Errorf("ошибка"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := NewOngoingContentRepository(db)
			query, args, _ := selectAllFields().
				From("ongoing_content").
				Where(sq.Gt{"release_date": time.Time{}}).
				OrderBy("release_date").
				Where(sq.Gt{"release_date": time.Time{}}).
				OrderBy("release_date").
				Limit(uint64(tc.Request)).
				PlaceholderFormat(sq.Dollar).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			require.NoError(t, err)
			driverValues := make([]driver.Value, len(args))
			for i, v := range args {
				driverValues[i] = v
			}
			tc.SetupMock(mock, query, driverValues)
			output, err := repo.GetNearestOngoings(tc.Request)

			require.Equal(t, tc.ExpectedOut, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
*/

func TestOngoingContentDB_GetAllReleaseYears(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		ExpectedOut []int
		ExpectedErr error
		SetupMock   func(mock sqlmock.Sqlmock)
	}{
		{
			Name:        "Successful retrieval of release years",
			ExpectedOut: []int{1980, 1981, 1982},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"extract(year from release_date)"}).
					AddRow(1980).
					AddRow(1981).
					AddRow(1982)
				mock.ExpectQuery("^SELECT extract\\(year from release_date\\) FROM ongoing_content GROUP BY extract\\(year from release_date\\) ORDER BY extract\\(year from release_date\\)$").WillReturnRows(rows)
			},
		},

		{
			Name:        "No release years found",
			ExpectedOut: nil,
			ExpectedErr: repository.ErrOngoingContentYearsNotFound,
			SetupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("^SELECT extract\\(year from release_date\\) FROM ongoing_content GROUP BY extract\\(year from release_date\\) ORDER BY extract\\(year from release_date\\)$").WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			mock.MatchExpectationsInOrder(false)
			require.NoError(t, err)
			repo := NewOngoingContentRepository(db)
			tc.SetupMock(mock)
			output, err := repo.GetAllReleaseYears()
			if !errors.Is(err, tc.ExpectedErr) {
				require.Fail(t, fmt.Errorf("unexpected error, expected: %v, got: %v", tc.ExpectedErr, err).Error())
			}
			require.Equal(t, tc.ExpectedOut, output)
		})
	}
}
