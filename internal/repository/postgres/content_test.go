package postgres

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func getDriverValues(args []any) []driver.Value {
	driverValues := make([]driver.Value, len(args))
	for i, v := range args {
		driverValues[i] = v
	}
	return driverValues
}

func setupGetContentSuccess(mock sqlmock.Sqlmock, contentID int, contentType string) {
	query, args, _ := sq.Select(
		"id",
		"content_type",
		"title",
		"original_title",
		"slogan",
		"budget",
		"age_restriction",
		"imdb",
		"rating",
		"description",
		"poster_upload_id",
		"trailer_url",
		"backdrop_upload_id",
		"ongoing",
		"ongoing_date",
	).
		From("content").
		Where(sq.Eq{"id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{
			"id",
			"content_type",
			"title",
			"original_title",
			"slogan",
			"budget",
			"age_restriction",
			"imdb",
			"rating",
			"description",
			"poster_upload_id",
			"trailer_url",
			"backdrop_upload_id",
			"ongoing",
			"ongoing_date",
		}).AddRow(
			contentID,
			contentType,
			"title",
			"original title",
			"slogan",
			"10",
			18,
			9,
			8,
			"description",
			501,
			"trailer",
			500,
			false,
			nil,
		))
}

func setupGetContentTypeSuccess(mock sqlmock.Sqlmock, contentID int, contentType string) {
	query, args, _ := sq.Select("type").
		From("content_type").
		Where(sq.Eq{"content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{"type"}).AddRow(contentType))
}

func setupGetContentProductionsCountriesSuccess(mock sqlmock.Sqlmock, contentID int, countriesID []int) {
	query, args, _ := sq.Select("country_id").
		From("country_content").
		Where(sq.Eq{"content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows := sqlmock.NewRows([]string{"country_id"})
	for _, countryID := range countriesID {
		rows.AddRow(countryID)
	}
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(rows)
}

func setupGetCountryByIDSuccess(mock sqlmock.Sqlmock, countryID int, countryName string) {
	query, args, _ := sq.Select("name").
		From("country").
		Where(sq.Eq{"id": countryID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{"name"}).AddRow(countryName))
}

func setupGetContentGenresSuccess(mock sqlmock.Sqlmock, contentID int, genresID []int) {
	query, args, _ := sq.Select("genre_id").
		From("genre_content").
		Where(sq.Eq{"content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows := sqlmock.NewRows([]string{"genre_id"})
	for _, genreID := range genresID {
		rows.AddRow(genreID)
	}
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(rows)
}

func setupGetGenreByIDSuccess(mock sqlmock.Sqlmock, genreID int, genreName string) {
	query, args, _ := sq.Select("name").
		From("genre").
		Where(sq.Eq{"id": genreID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{"name"}).AddRow(genreName))
}

func setupGetRoleIDByNameSuccess(mock sqlmock.Sqlmock, roleName string, roleID int) {
	query, args, _ := sq.Select("id").
		From("role").
		Where(sq.Eq{"name_en": roleName}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(roleID))
}

func setupGetPersonsByRoleAndContentIDSuccess(mock sqlmock.Sqlmock, roleID, contentID int, personID int) {
	query, args, _ := sq.Select("person_id").
		From("person_role").
		Where(sq.Eq{"content_id": contentID, "role_id": roleID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{"person_id"}).AddRow(personID))
}

func setupGetPerson(mock sqlmock.Sqlmock) {
	query, args, _ := sq.Select(
		"id",
		"name",
		"en_name",
		"birth_date",
		"death_date",
		"sex",
		"height",
		"photo_upload_id",
	).
		From("person").
		Where(sq.Eq{"id": 1}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{
			"id",
			"name",
			"en_name",
			"birth_date",
			"death_date",
			"sex",
			"height",
			"photo_upload_id",
		}).AddRow(
			1,
			"Имя",
			"Name",
			time.Time{},
			time.Time{},
			"M",
			175,
			1,
		))
}

func setupGetMovieDataSuccess(mock sqlmock.Sqlmock, contentID int) {
	query, args, _ := sq.Select("premiere", "duration").
		From("movie").
		Where(sq.Eq{"content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{
			"premiere",
			"duration",
		}).AddRow(
			time.Time{},
			100,
		))
}

func setupGetSeriesDataSuccess(mock sqlmock.Sqlmock, contentID int) {
	query, args, _ := sq.Select("year_start", "year_end").
		From("series").
		Where(sq.Eq{"content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{
			"year_start",
			"year_end",
		}).AddRow(
			1980,
			1981,
		))
}

func setupGetSeasonsSuccess(mock sqlmock.Sqlmock, contentID int, seasons []int) {
	query, args, _ := sq.Select("id").
		From("season").
		Where(sq.Eq{"series_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows := sqlmock.NewRows([]string{"id"})
	for _, seasonID := range seasons {
		rows.AddRow(seasonID)
	}
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(rows)
}

func setupGetSeasonSuccess(mock sqlmock.Sqlmock, seasonID int, title string) {
	query, args, _ := sq.Select("id", "title").
		From("season").
		Where(sq.Eq{"id": seasonID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{"id", "title"}).AddRow(seasonID, title))
}

func setupGetEpisodesBySeasonIDSuccess(mock sqlmock.Sqlmock, seasonID int, episodes []int) {
	query, args, _ := sq.Select("id").
		From("episode").
		Where(sq.Eq{"season_id": seasonID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows := sqlmock.NewRows([]string{"id"})
	for _, episodeID := range episodes {
		rows.AddRow(episodeID)
	}
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(rows)
}

func setupGetEpisodeSuccess(mock sqlmock.Sqlmock, episodeID int, episodeNumber int, title string) {
	query, args, _ := sq.Select("id", "episode_number", "title", "duration").
		From("episode").
		Where(sq.Eq{"id": episodeID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{"id", "episode_number", "title", "duration"}).AddRow(episodeID, episodeNumber, title, 100))
}

func setupGetFactsSuccess(mock sqlmock.Sqlmock, contentID int, facts []string) {
	query, args, _ := sq.Select("fact").
		From("content_fact").
		Where(sq.Eq{"content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows := sqlmock.NewRows([]string{"fact"})
	for _, fact := range facts {
		rows.AddRow(fact)
	}
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(rows)
}

func setupGetPicturesSuccess(mock sqlmock.Sqlmock, contentID int, pictures []int) {
	query, args, _ := sq.Select("static_id").
		From("content_image").
		Where(sq.Eq{"content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows := sqlmock.NewRows([]string{"static_id"})
	for _, picture := range pictures {
		rows.AddRow(picture)
	}
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(rows)
}

func TestContentDB_GetContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Request     int
		ExpectedOut *entity.Content
		ExpectedErr error
		SetupMock   func(mock sqlmock.Sqlmock)
	}{
		{
			Name:    "Успешное получение фильма",
			Request: 1,
			ExpectedOut: &entity.Content{
				ID:               1,
				Title:            "title",
				OriginalTitle:    "original title",
				Slogan:           "slogan",
				Budget:           "10",
				AgeRestriction:   18,
				Rating:           8,
				IMDBRating:       9,
				Description:      "description",
				PosterStaticID:   501,
				TrailerLink:      "trailer",
				BackdropStaticID: 500,
				PicturesStaticID: []int{900},
				Facts:            []string{"fact"},
				Country:          []entity.Country{{ID: 1, Name: "Russia"}},
				Genres:           []entity.Genre{{ID: 1, Name: "Action"}},
				Actors:           []entity.Person{entity.GetExamplePerson()},
				Directors:        []entity.Person{entity.GetExamplePerson()},
				Producers:        []entity.Person{entity.GetExamplePerson()},
				Writers:          []entity.Person{entity.GetExamplePerson()},
				Operators:        []entity.Person{entity.GetExamplePerson()},
				Composers:        []entity.Person{entity.GetExamplePerson()},
				Editors:          []entity.Person{entity.GetExamplePerson()},
				Type:             "movie",
				Movie: &entity.Movie{
					Premiere: time.Time{},
					Duration: 100,
				},
				Series: nil,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock) {
				// устанавливаем порядок вызовов и мокаем их
				setupGetContentSuccess(mock, 1, entity.ContentTypeMovie)
				setupGetContentProductionsCountriesSuccess(mock, 1, []int{1})
				setupGetCountryByIDSuccess(mock, 1, "Russia")
				setupGetContentGenresSuccess(mock, 1, []int{1})
				setupGetGenreByIDSuccess(mock, 1, "Action")

				setupGetRoleIDByNameSuccess(mock, "actor", 1)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 1, 1, 1)
				setupGetPerson(mock)

				setupGetRoleIDByNameSuccess(mock, "director", 2)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 2, 1, 1)
				setupGetPerson(mock)

				setupGetRoleIDByNameSuccess(mock, "producer", 3)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 3, 1, 1)
				setupGetPerson(mock)

				setupGetRoleIDByNameSuccess(mock, "writer", 4)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 4, 1, 1)
				setupGetPerson(mock)

				setupGetRoleIDByNameSuccess(mock, "operator", 5)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 5, 1, 1)
				setupGetPerson(mock)

				setupGetRoleIDByNameSuccess(mock, "composer", 6)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 6, 1, 1)
				setupGetPerson(mock)

				setupGetRoleIDByNameSuccess(mock, "editor", 7)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 7, 1, 1)
				setupGetPerson(mock)

				setupGetMovieDataSuccess(mock, 1)
				setupGetFactsSuccess(mock, 1, []string{"fact"})
				setupGetPicturesSuccess(mock, 1, []int{900})
			},
		},
		{
			Name:    "Успешное получение сериала",
			Request: 1,
			ExpectedOut: &entity.Content{
				ID:               1,
				Title:            "title",
				OriginalTitle:    "original title",
				Slogan:           "slogan",
				Budget:           "10",
				AgeRestriction:   18,
				Rating:           8,
				IMDBRating:       9,
				Description:      "description",
				PosterStaticID:   501,
				TrailerLink:      "trailer",
				BackdropStaticID: 500,
				PicturesStaticID: []int{900},
				Facts:            []string{"fact"},
				Country:          []entity.Country{{ID: 1, Name: "Russia"}},
				Genres:           []entity.Genre{{ID: 1, Name: "Action"}},
				Actors:           []entity.Person{entity.GetExamplePerson()},
				Directors:        []entity.Person{entity.GetExamplePerson()},
				Producers:        []entity.Person{entity.GetExamplePerson()},
				Writers:          []entity.Person{entity.GetExamplePerson()},
				Operators:        []entity.Person{entity.GetExamplePerson()},
				Composers:        []entity.Person{entity.GetExamplePerson()},
				Editors:          []entity.Person{entity.GetExamplePerson()},
				Type:             "series",
				Movie:            nil,
				Series: &entity.Series{
					Seasons: []entity.Season{
						{
							ID:    1,
							Title: "Season",
							Episodes: []entity.Episode{
								{
									ID:            1,
									EpisodeNumber: 10,
									Title:         "Episode",
									Duration:      100,
								},
							},
						},
					},
					YearStart: 1980,
					YearEnd:   1981,
				},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock) {
				// устанавливаем порядок вызовов и мокаем их
				setupGetContentSuccess(mock, 1, entity.ContentTypeSeries)
				setupGetContentTypeSuccess(mock, 1, "series")
				setupGetContentProductionsCountriesSuccess(mock, 1, []int{1})
				setupGetCountryByIDSuccess(mock, 1, "Russia")
				setupGetContentGenresSuccess(mock, 1, []int{1})
				setupGetGenreByIDSuccess(mock, 1, "Action")

				setupGetRoleIDByNameSuccess(mock, "actor", 1)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 1, 1, 1)
				setupGetPerson(mock)

				setupGetRoleIDByNameSuccess(mock, "director", 2)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 2, 1, 1)
				setupGetPerson(mock)

				setupGetRoleIDByNameSuccess(mock, "producer", 3)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 3, 1, 1)
				setupGetPerson(mock)

				setupGetRoleIDByNameSuccess(mock, "writer", 4)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 4, 1, 1)
				setupGetPerson(mock)

				setupGetRoleIDByNameSuccess(mock, "operator", 5)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 5, 1, 1)
				setupGetPerson(mock)

				setupGetRoleIDByNameSuccess(mock, "composer", 6)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 6, 1, 1)
				setupGetPerson(mock)

				setupGetRoleIDByNameSuccess(mock, "editor", 7)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 7, 1, 1)
				setupGetPerson(mock)

				setupGetSeriesDataSuccess(mock, 1)
				setupGetSeasonsSuccess(mock, 1, []int{1})
				setupGetSeasonSuccess(mock, 1, "Season")
				setupGetEpisodesBySeasonIDSuccess(mock, 1, []int{1})
				setupGetEpisodeSuccess(mock, 1, 10, "Episode")

				setupGetFactsSuccess(mock, 1, []string{"fact"})
				setupGetPicturesSuccess(mock, 1, []int{900})
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			dbx := sqlx.NewDb(db, "sqlmock")
			mock.MatchExpectationsInOrder(false)
			require.NoError(t, err)
			repo := NewContentRepository(dbx)
			tc.SetupMock(mock)
			output, err := repo.GetContent(tc.Request)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedOut, output)
		})
	}
}
