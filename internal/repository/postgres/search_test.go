package postgres

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestSearchDB_SearchContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Query       string
		ExpectedErr error
		ExpectedOut []entity.Content
		SetupMock   func(mock sqlmock.Sqlmock, query string)
	}{
		{
			Name:  "Успешный поиск",
			Query: "Query",
			ExpectedOut: []entity.Content{
				{
					ID:             1,
					Title:          "title",
					OriginalTitle:  "original_title",
					Rating:         5.0,
					PosterStaticID: 1,
					Genres:         []entity.Genre{{ID: 1, Name: "genre"}},
					Country:        []entity.Country{{ID: 1, Name: "country"}},
					Actors:         []entity.Person{{ID: 1, Name: "name", EnName: "en_name", BirthDate: time.Time{}, DeathDate: time.Time{}, Sex: "M", Height: 100, PhotoStaticID: 1}},
					Directors:      []entity.Person{{ID: 1, Name: "name", EnName: "en_name", BirthDate: time.Time{}, DeathDate: time.Time{}, Sex: "M", Height: 100, PhotoStaticID: 1}},
					Type:           "movie",
					Movie:          &entity.Movie{Premiere: time.Time{}, Duration: 100},
				},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("Query", "Query").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, content_type, title, original_title, rating, poster_upload_id FROM content WHERE id = $1")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "content_type", "title", "original_title", "rating", "poster_upload_id"}).
						AddRow(1, "movie", "title", "original_title", 5.0, 1))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT genre_id FROM genre_content WHERE content_id = $1")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"genre_id"}).AddRow(1))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT country_id FROM country_content WHERE content_id = $1")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"country_id"}).AddRow(1))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT premiere, duration FROM movie WHERE content_id = $1")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"premiere", "duration"}).AddRow(time.Time{}, 100))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM role WHERE name_en = $1")).
					WithArgs("actor").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM genre WHERE id = $1")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("genre"))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id FROM role WHERE name_en = $1")).
					WithArgs("director").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT person_id FROM person_role WHERE content_id = $1 AND role_id = $2")).
					WithArgs(1, 1).
					WillReturnRows(sqlmock.NewRows([]string{"person_id"}).AddRow(1))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT person_id FROM person_role WHERE content_id = $1 AND role_id = $2")).
					WithArgs(1, 2).
					WillReturnRows(sqlmock.NewRows([]string{"person_id"}).AddRow(1))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT name FROM country WHERE id = $1")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("country"))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, en_name, birth_date, death_date, sex, height, photo_upload_id FROM person WHERE id = $1")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "en_name", "birth_date", "death_date", "sex", "height", "photo_upload_id"}).
						AddRow(1, "name", "en_name", time.Time{}, time.Time{}, "M", 100, 1))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, en_name, birth_date, death_date, sex, height, photo_upload_id FROM person WHERE id = $1")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "en_name", "birth_date", "death_date", "sex", "height", "photo_upload_id"}).
						AddRow(1, "name", "en_name", time.Time{}, time.Time{}, "M", 100, 1))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("не удалось создать мок: %s", err)
			}
			mock.MatchExpectationsInOrder(false)
			mockContent := NewContentRepository(db)
			repo := NewSearchRepository(db, mockContent)
			query, _, err := sq.Select("id").
				From("content").
				Where(sq.Or{
					sq.Expr("word_similarity(title, ?) > 0.3", tc.Query),
					sq.Expr("word_similarity(original_title, ?) > 0.3", tc.Query),
				}).
				OrderBy(fmt.Sprintf(
					`CASE WHEN word_similarity(title, '%s') > 0.3
THEN similarity(title, '%s') ELSE similarity(original_title, '%s') END DESC`, tc.Query, tc.Query, tc.Query)).
				Limit(5).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			if err != nil {
				t.Fatalf("ошибка при составлении запроса SearchContent: %s", err)
			}
			tc.SetupMock(mock, query)
			contents, err := repo.SearchContent(tc.Query)
			require.Equal(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOut, contents)
		})
	}
}

func TestSearchDB_SearchPerson(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		Query       string
		ExpectedErr error
		ExpectedOut []entity.Person
		SetupMock   func(mock sqlmock.Sqlmock, query string)
	}{
		{
			Name:  "Успешный поиск",
			Query: "Query",
			ExpectedOut: []entity.Person{
				{
					ID:            1,
					Name:          "name",
					EnName:        "en_name",
					Sex:           "M",
					Height:        100,
					PhotoStaticID: 1,
					BirthDate:     time.Time{},
					DeathDate:     time.Time{},
				},
			},
			ExpectedErr: nil,
			SetupMock: func(mock sqlmock.Sqlmock, query string) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs("Query", "Query").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, en_name, birth_date, death_date, sex, height, photo_upload_id FROM person WHERE id = $1")).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "en_name", "birth_date", "death_date", "sex", "height", "photo_upload_id"}).
						AddRow(1, "name", "en_name", time.Time{}, time.Time{}, "M", 100, 1))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("не удалось создать мок: %s", err)
			}
			mock.MatchExpectationsInOrder(false)
			mockContent := NewContentRepository(db)
			repo := NewSearchRepository(db, mockContent)
			query, _, err := sq.Select("id").
				From("person").
				Where(sq.Or{
					sq.Expr("word_similarity(Name, ?) > 0.3", tc.Query),
					sq.Expr("word_similarity(en_name, ?) > 0.3", tc.Query),
				}).
				OrderBy(fmt.Sprintf(
					`CASE WHEN word_similarity(Name, '%s') > 0.3
THEN similarity(Name, '%s') ELSE similarity(en_name, '%s') END DESC`, tc.Query, tc.Query, tc.Query)).
				Limit(5).
				PlaceholderFormat(sq.Dollar).
				ToSql()
			if err != nil {
				t.Fatalf("ошибка при составлении запроса SearchPerson: %s", err)
			}
			tc.SetupMock(mock, query)
			persons, err := repo.SearchPerson(tc.Query)
			require.Equal(t, tc.ExpectedErr, err)
			require.EqualValues(t, tc.ExpectedOut, persons)
		})
	}
}
