package postgres

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
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

func setupGetContentSuccess(mock sqlmock.Sqlmock, contentID int) {
	query, args, _ := sq.Select(
		"id",
		"title",
		"original_title",
		"slogan",
		"budget",
		"age_restriction",
		"audience",
		"imdb",
		"description",
		"poster_upload_id",
		"box_office",
		"marketing_budget",
	).
		From("content").
		Where(sq.Eq{"id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{
			"id",
			"title",
			"original_title",
			"slogan",
			"budget",
			"age_restriction",
			"audience",
			"imdb",
			"description",
			"poster_upload_id",
			"box_office",
			"marketing_budget",
		}).AddRow(
			contentID,
			"title",
			"original title",
			"slogan",
			10,
			18,
			1000000,
			9,
			"description",
			1,
			10,
			10,
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
		Where(sq.Eq{"name": roleName}).
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

func setupGetPerson(mock sqlmock.Sqlmock, personID int, personFirstName, personLastName string) {
	query, args, _ := sq.Select(
		"id",
		"first_name",
		"last_name",
		"birth_date",
		"death_date",
		"start_career",
		"end_career",
		"sex",
		"height",
		"spouse",
		"children",
		"photo_upload_id",
	).
		From("person").
		Where(sq.Eq{"id": personID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{
			"id",
			"first_name",
			"last_name",
			"birth_date",
			"death_date",
			"start_career",
			"end_career",
			"sex",
			"height",
			"spouse",
			"children",
			"photo_upload_id",
		}).AddRow(
			personID,
			personFirstName,
			personLastName,
			time.Time{},
			time.Time{},
			time.Time{},
			time.Time{},
			"M",
			180,
			"spouse",
			"children",
			1,
		))
}

func setupGetMovieDataSuccess(mock sqlmock.Sqlmock, contentID int) {
	query, args, _ := sq.Select("premiere", "release", "duration").
		From("movie").
		Where(sq.Eq{"content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{
			"premiere",
			"release",
			"duration",
		}).AddRow(
			time.Time{},
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
		Where(sq.Eq{"content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows := sqlmock.NewRows([]string{"season_id"})
	for _, seasonID := range seasons {
		rows.AddRow(seasonID)
	}
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(rows)
}

func setupGetSeasonSuccess(mock sqlmock.Sqlmock, seasonID int, yearStart, yearEnd int) {
	query, args, _ := sq.Select("id", "year_start", "year_end").
		From("season").
		Where(sq.Eq{"id": seasonID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{"id", "year_start", "year_end"}).AddRow(seasonID, yearStart, yearEnd))
}

func setupGetEpisodesBySeasonIDSuccess(mock sqlmock.Sqlmock, seasonID int, episodes []int) {
	query, args, _ := sq.Select("id").
		From("episode").
		Where(sq.Eq{"season_id": seasonID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows := sqlmock.NewRows([]string{"episode_id"})
	for _, episodeID := range episodes {
		rows.AddRow(episodeID)
	}
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(rows)
}

func setupGetEpisodeSuccess(mock sqlmock.Sqlmock, episodeID int, episodeNumber int, title string) {
	query, args, _ := sq.Select("id", "episode_number", "title").
		From("episode").
		Where(sq.Eq{"id": episodeID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(getDriverValues(args)...).WillReturnRows(
		sqlmock.NewRows([]string{"id", "episode_number", "title"}).AddRow(episodeID, episodeNumber, title))
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
				ID:             1,
				Title:          "title",
				OriginalTitle:  "original title",
				Slogan:         "slogan",
				Budget:         10,
				AgeRestriction: 18,
				Audience:       1000000,
				IMDBRating:     9,
				Description:    "description",
				PosterStaticID: 1,
				BoxOffice:      10,
				Marketing:      10,
				Country:        []entity.Country{{ID: 1, Name: "Russia"}},
				Genres:         []entity.Genre{{ID: 1, Name: "Action"}},
				Actors:         []entity.Person{{ID: 10, FirstName: "Actor", LastName: "Actorov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Directors:      []entity.Person{{ID: 11, FirstName: "Director", LastName: "Directorov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Producers:      []entity.Person{{ID: 12, FirstName: "Producer", LastName: "Producerov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Writers:        []entity.Person{{ID: 13, FirstName: "Writer", LastName: "Writerov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Operators:      []entity.Person{{ID: 14, FirstName: "Operator", LastName: "Operatorov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Composers:      []entity.Person{{ID: 15, FirstName: "Composer", LastName: "Composerov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Editors:        []entity.Person{{ID: 16, FirstName: "Editor", LastName: "Editorov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Type:           "movie",
				Movie: &entity.Movie{
					Premiere: time.Time{},
					Release:  time.Time{},
					Duration: 100,
				},
				Series: nil,
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock) {
				// устанавливаем порядок вызовов и мокаем их
				setupGetContentSuccess(mock, 1)
				setupGetContentTypeSuccess(mock, 1, "movie")
				setupGetContentProductionsCountriesSuccess(mock, 1, []int{1})
				setupGetCountryByIDSuccess(mock, 1, "Russia")
				setupGetContentGenresSuccess(mock, 1, []int{1})
				setupGetGenreByIDSuccess(mock, 1, "Action")

				setupGetRoleIDByNameSuccess(mock, "actor", 1)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 1, 1, 10)
				setupGetPerson(mock, 10, "Actor", "Actorov")

				setupGetRoleIDByNameSuccess(mock, "director", 2)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 2, 1, 11)
				setupGetPerson(mock, 11, "Director", "Directorov")

				setupGetRoleIDByNameSuccess(mock, "producer", 3)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 3, 1, 12)
				setupGetPerson(mock, 12, "Producer", "Producerov")

				setupGetRoleIDByNameSuccess(mock, "writer", 4)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 4, 1, 13)
				setupGetPerson(mock, 13, "Writer", "Writerov")

				setupGetRoleIDByNameSuccess(mock, "operator", 5)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 5, 1, 14)
				setupGetPerson(mock, 14, "Operator", "Operatorov")

				setupGetRoleIDByNameSuccess(mock, "composer", 6)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 6, 1, 15)
				setupGetPerson(mock, 15, "Composer", "Composerov")

				setupGetRoleIDByNameSuccess(mock, "editor", 7)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 7, 1, 16)
				setupGetPerson(mock, 16, "Editor", "Editorov")

				setupGetMovieDataSuccess(mock, 1)
			},
		},
		{
			Name:    "Успешное получение сериала",
			Request: 1,
			ExpectedOut: &entity.Content{
				ID:             1,
				Title:          "title",
				OriginalTitle:  "original title",
				Slogan:         "slogan",
				Budget:         10,
				AgeRestriction: 18,
				Audience:       1000000,
				IMDBRating:     9,
				Description:    "description",
				PosterStaticID: 1,
				BoxOffice:      10,
				Marketing:      10,
				Country:        []entity.Country{{ID: 1, Name: "Russia"}},
				Genres:         []entity.Genre{{ID: 1, Name: "Action"}},
				Actors:         []entity.Person{{ID: 10, FirstName: "Actor", LastName: "Actorov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Directors:      []entity.Person{{ID: 11, FirstName: "Director", LastName: "Directorov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Producers:      []entity.Person{{ID: 12, FirstName: "Producer", LastName: "Producerov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Writers:        []entity.Person{{ID: 13, FirstName: "Writer", LastName: "Writerov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Operators:      []entity.Person{{ID: 14, FirstName: "Operator", LastName: "Operatorov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Composers:      []entity.Person{{ID: 15, FirstName: "Composer", LastName: "Composerov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Editors:        []entity.Person{{ID: 16, FirstName: "Editor", LastName: "Editorov", Sex: "M", Height: 180, Spouse: "spouse", Children: "children", PhotoStaticID: 1}},
				Type:           "series",
				Movie:          nil,
				Series: &entity.Series{
					Seasons: []entity.Season{
						{
							ID:        1,
							YearStart: 1980,
							YearEnd:   1981,
							Episodes: []entity.Episode{
								{
									ID:            1,
									EpisodeNumber: 10,
									Title:         "Episode",
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
				setupGetContentSuccess(mock, 1)
				setupGetContentTypeSuccess(mock, 1, "series")
				setupGetContentProductionsCountriesSuccess(mock, 1, []int{1})
				setupGetCountryByIDSuccess(mock, 1, "Russia")
				setupGetContentGenresSuccess(mock, 1, []int{1})
				setupGetGenreByIDSuccess(mock, 1, "Action")

				setupGetRoleIDByNameSuccess(mock, "actor", 1)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 1, 1, 10)
				setupGetPerson(mock, 10, "Actor", "Actorov")

				setupGetRoleIDByNameSuccess(mock, "director", 2)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 2, 1, 11)
				setupGetPerson(mock, 11, "Director", "Directorov")

				setupGetRoleIDByNameSuccess(mock, "producer", 3)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 3, 1, 12)
				setupGetPerson(mock, 12, "Producer", "Producerov")

				setupGetRoleIDByNameSuccess(mock, "writer", 4)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 4, 1, 13)
				setupGetPerson(mock, 13, "Writer", "Writerov")

				setupGetRoleIDByNameSuccess(mock, "operator", 5)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 5, 1, 14)
				setupGetPerson(mock, 14, "Operator", "Operatorov")

				setupGetRoleIDByNameSuccess(mock, "composer", 6)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 6, 1, 15)
				setupGetPerson(mock, 15, "Composer", "Composerov")

				setupGetRoleIDByNameSuccess(mock, "editor", 7)
				setupGetPersonsByRoleAndContentIDSuccess(mock, 7, 1, 16)
				setupGetPerson(mock, 16, "Editor", "Editorov")

				setupGetSeriesDataSuccess(mock, 1)
				setupGetSeasonsSuccess(mock, 1, []int{1})
				setupGetSeasonSuccess(mock, 1, 1980, 1981)
				setupGetEpisodesBySeasonIDSuccess(mock, 1, []int{1})
				setupGetEpisodeSuccess(mock, 1, 10, "Episode")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			repo := &ContentDB{
				DB: db,
			}
			tc.SetupMock(mock)
			output, err := repo.GetContent(tc.Request)
			require.Equal(t, tc.ExpectedOut, output)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
