package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
)

type ContentDB struct {
	DB *sql.DB
}

type ScanContent struct {
	ID             int
	Title          string
	OriginalTitle  sql.NullString
	Slogan         sql.NullString
	Budget         sql.NullInt64
	AgeRestriction int
	Audience       sql.NullInt64
	IMDBRating     float64
	Description    string
	PosterStaticID int
	BoxOffice      sql.NullInt64
	Marketing      sql.NullInt64
}

func NewContentRepository(database config.PostgresDatabase) (repository.Content, error) {
	db, err := sql.Open("postgres", database.ConnectURL)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &ContentDB{
		DB: db,
	}, nil
}

// getMovieData возвращает информацию о фильме по его ID
func (c *ContentDB) getMovieData(id int) (*entity.Movie, error) {
	// no-lint
	query, args, _ := sq.Select("premiere", "release", "duration").
		From("movie").
		Where(sq.Eq{"content_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var movie entity.Movie
	err := c.DB.QueryRow(query, args...).Scan(&movie.Premiere, &movie.Release, &movie.Duration)
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении информации о фильме"))
	}
	return &movie, nil
}

// getSeriesData возвращает информацию о сериале по его ID
func (c *ContentDB) getSeriesData(id int) (*entity.Series, error) {
	// no-lint
	query, args, _ := sq.Select("year_start", "year_end").
		From("series").
		Where(sq.Eq{"content_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var series entity.Series
	err := c.DB.QueryRow(query, args...).Scan(&series.YearStart, &series.YearEnd)
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении информации о сериале"))
	}
	seasons, err := c.getSeasonsByContentID(id)
	if err != nil {
		return nil, err
	}
	series.Seasons = seasons
	return &series, nil
}

// getEpisodeData возвращает информацию об эпизоде по его ID
func (c *ContentDB) getEpisodeData(id int) (*entity.Episode, error) {
	// no-lint
	query, args, _ := sq.Select("id", "episode_number", "title").
		From("episode").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var episode entity.Episode
	err := c.DB.QueryRow(query, args...).Scan(&episode.ID, &episode.EpisodeNumber, &episode.Title)
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении информации об эпизоде"))
	}
	return &episode, nil
}

// getEpisodesBySeasonID возвращает эпизоды сезона по его ID
// nolint: dupl
func (c *ContentDB) getEpisodesBySeasonID(seasonID int) ([]entity.Episode, error) {
	// no-lint
	query, args, _ := sq.Select("id").
		From("episode").
		Where(sq.Eq{"season_id": seasonID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении эпизодов сезона"))
	}
	defer rows.Close()
	var episodes []entity.Episode
	for rows.Next() {
		var episodeID int
		err := rows.Scan(&episodeID)
		if err != nil {
			return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении эпизодов сезона"))
		}
		episode, err := c.getEpisodeData(episodeID)
		if err != nil {
			return nil, err
		}
		episodes = append(episodes, *episode)
	}
	return episodes, nil
}

// getSeasonData возвращает информацию о сезоне по его ID
func (c *ContentDB) getSeasonData(id int) (*entity.Season, error) {
	// no-lint
	query, args, _ := sq.Select("id", "year_start", "year_end").
		From("season").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var season entity.Season
	err := c.DB.QueryRow(query, args...).Scan(&season.ID, &season.YearStart, &season.YearEnd)
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении информации о сезоне"))
	}
	return &season, nil
}

// getSeasonsByContentID возвращает сезоны контента по его ID
func (c *ContentDB) getSeasonsByContentID(contentID int) ([]entity.Season, error) {
	// no-lint
	query, args, _ := sq.Select("id").
		From("season").
		Where(sq.Eq{"content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении сезонов контента"))
	}
	defer rows.Close()
	var seasons []entity.Season
	for rows.Next() {
		var seasonID int
		err := rows.Scan(&seasonID)
		if err != nil {
			return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении сезонов контента"))
		}
		season, err := c.getSeasonData(seasonID)
		if err != nil {
			return nil, err
		}
		episodes, err := c.getEpisodesBySeasonID(seasonID)
		if err != nil {
			return nil, err
		}
		season.Episodes = episodes
		seasons = append(seasons, *season)
	}
	return seasons, nil
}

// getRoleIDByName возвращает айди роли по ее названию
func (c *ContentDB) getRoleIDByName(role string) (int, error) {
	// no-lint
	query, args, _ := sq.Select("id").
		From("role").
		Where(sq.Eq{"name": role}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var roleID int
	err := c.DB.QueryRow(query, args...).Scan(&roleID)
	if err != nil {
		return 0, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении айди роли"))
	}
	return roleID, nil
}

// getPersonsByRoleAndContentID возвращает персон контента по его ID и роли
func (c *ContentDB) getPersonsByRoleAndContentID(role string, contentID int) ([]entity.Person, error) {
	roleID, err := c.getRoleIDByName(role)
	if err != nil {
		return nil, err
	}
	// no-lint
	query, args, _ := sq.Select("person_id").
		From("person_role").
		Where(sq.Eq{"content_id": contentID, "role_id": roleID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении персон контента"))
	}
	defer rows.Close()
	var persons []entity.Person
	for rows.Next() {
		var personID int
		err := rows.Scan(&personID)
		if err != nil {
			return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении персон контента"))
		}
		person, err := c.GetPerson(personID)
		if err != nil {
			return nil, err
		}
		persons = append(persons, *person)
	}
	return persons, nil
}

// getGenreByID возвращает жанр по его ID
func (c *ContentDB) getGenreByID(id int) (*entity.Genre, error) {
	query, args, _ := sq.Select("name").
		From("genre").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var genre entity.Genre
	err := c.DB.QueryRow(query, args...).Scan(&genre.Name)
	genre.ID = id
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении жанра контента"))
	}
	return &genre, nil
}

// getContentGenres возвращает жанры контента по его ID
// nolint: dupl
func (c *ContentDB) getContentGenres(id int) ([]entity.Genre, error) {
	// no-lint
	query, args, _ := sq.Select("genre_id").
		From("genre_content").
		Where(sq.Eq{"content_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении жанров контента"))
	}
	defer rows.Close()
	var genres []entity.Genre
	for rows.Next() {
		var genreID int
		err := rows.Scan(&genreID)
		if err != nil {
			return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении жанров контента"))
		}
		genre, err := c.getGenreByID(genreID)
		if err != nil {
			return nil, err
		}
		genres = append(genres, *genre)
	}
	return genres, nil

}

// getCountryByID возвращает страну по ее ID
func (c *ContentDB) getCountryByID(id int) (*entity.Country, error) {
	// no-lint
	query, args, _ := sq.Select("name").
		From("country").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var country entity.Country
	err := c.DB.QueryRow(query, args...).Scan(&country.Name)
	country.ID = id
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении страны производства контента"))
	}
	return &country, nil
}

// getContentProductionCountries возвращает страны производства контента по его ID
// nolint: dupl
func (c *ContentDB) getContentProductionCountries(id int) ([]entity.Country, error) {
	// no-lint
	query, args, _ := sq.Select("country_id").
		From("country_content").
		Where(sq.Eq{"content_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении стран производства контента"))
	}
	defer rows.Close()
	var countries []entity.Country
	for rows.Next() {
		var countryID int
		err := rows.Scan(&countryID)
		if err != nil {
			return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении стран производства контента"))
		}
		country, err := c.getCountryByID(countryID)
		if err != nil {
			return nil, err
		}
		countries = append(countries, *country)
	}
	return countries, nil
}

// getContentType возвращает тип контента по его ID. Возможные значения: movie, series
func (c *ContentDB) getContentType(id int) (string, error) {
	// no-lint
	query, args, _ := sq.Select("type").
		From("content_type").
		Where(sq.Eq{"content_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var contentType string
	err := c.DB.QueryRow(query, args...).Scan(&contentType)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", entity.NewClientError("Контент не найден", entity.ErrNotFound)
		}
		return "", entity.PSQLWrap(err, fmt.Errorf("ошибка при получении типа контента"))
	}
	return contentType, nil
}

func (c *ContentDB) GetContent(id int) (*entity.Content, error) {
	// TODO: подзапросы можно (и нужно) распараллелить
	// no-lint
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
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var scanContent ScanContent
	var content entity.Content
	err := c.DB.QueryRow(query, args...).Scan(
		&scanContent.ID,
		&scanContent.Title,
		&scanContent.OriginalTitle,
		&scanContent.Slogan,
		&scanContent.Budget,
		&scanContent.AgeRestriction,
		&scanContent.Audience,
		&scanContent.IMDBRating,
		&scanContent.Description,
		&scanContent.PosterStaticID,
		&scanContent.BoxOffice,
		&scanContent.Marketing,
	)
	content.ID = scanContent.ID
	content.Title = scanContent.Title
	content.OriginalTitle = scanContent.OriginalTitle.String
	content.Slogan = scanContent.Slogan.String
	content.Budget = int(scanContent.Budget.Int64)
	content.AgeRestriction = scanContent.AgeRestriction
	content.Audience = int(scanContent.Audience.Int64)
	content.IMDBRating = scanContent.IMDBRating
	content.Description = scanContent.Description
	content.PosterStaticID = scanContent.PosterStaticID
	content.BoxOffice = int(scanContent.BoxOffice.Int64)
	content.Marketing = int(scanContent.Marketing.Int64)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewClientError("Контент не найден", entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении контента"))
	}
	contentType, err := c.getContentType(id)
	if err != nil {
		return nil, err
	}
	content.Type = contentType
	countries, err := c.getContentProductionCountries(id)
	if err != nil {
		return nil, err
	}
	content.Country = countries
	genres, err := c.getContentGenres(id)
	if err != nil {
		return nil, err
	}
	content.Genres = genres

	type RolePersons struct {
		Role    string
		Persons *[]entity.Person
	}
	roles := []RolePersons{
		{Role: entity.RoleActor, Persons: &content.Actors},
		{Role: entity.RoleDirector, Persons: &content.Directors},
		{Role: entity.RoleProducer, Persons: &content.Producers},
		{Role: entity.RoleWriter, Persons: &content.Writers},
		{Role: entity.RoleOperator, Persons: &content.Operators},
		{Role: entity.RoleComposer, Persons: &content.Composers},
		{Role: entity.RoleEditor, Persons: &content.Editors},
	}
	for _, rolePersons := range roles {
		personsResult, err := c.getPersonsByRoleAndContentID(rolePersons.Role, id)
		if err != nil {
			return nil, err
		}
		*rolePersons.Persons = personsResult
	}

	switch contentType {
	case entity.ContentTypeMovie:
		movie, err := c.getMovieData(id)
		if err != nil {
			return nil, err
		}
		content.Movie = movie
	case entity.ContentTypeSeries:
		series, err := c.getSeriesData(id)
		if err != nil {
			return nil, err
		}
		content.Series = series
	default:
		return nil, entity.NewClientError("Неизвестный тип контента", entity.ErrNotFound)
	}
	return &content, nil
}

// GetPreviewContent возвращает краткую информацию о контенте по его ID.
// Заполняет только поля id, title, original_title, poster_upload_id, countries, genres, actors, directors.
// Для фильма заполняет поля premiere, duration.
// Для сериала заполняет поля year_start, year_end, seasons.
func (c *ContentDB) GetPreviewContent(id int) (*entity.Content, error) {
	// TODO: подзапросы можно (и нужно) распараллелить
	// no-lint
	query, args, _ := sq.Select(
		"id",
		"title",
		"original_title",
		"poster_upload_id",
	).
		From("content").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var content entity.Content
	err := c.DB.QueryRow(query, args...).Scan(
		&content.ID,
		&content.Title,
		&content.OriginalTitle,
		&content.PosterStaticID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewClientError("Контент не найден", entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении краткой информации о контенте"))
	}
	countries, err := c.getContentProductionCountries(id)
	if err != nil {
		return nil, err
	}
	content.Country = countries
	genres, err := c.getContentGenres(id)
	if err != nil {
		return nil, err
	}
	content.Genres = genres
	actors, err := c.getPersonsByRoleAndContentID(entity.RoleActor, id)
	if err != nil {
		return nil, err
	}
	content.Actors = actors
	directors, err := c.getPersonsByRoleAndContentID(entity.RoleDirector, id)
	if err != nil {
		return nil, err
	}
	content.Directors = directors
	contentType, err := c.getContentType(id)
	if err != nil {
		return nil, err
	}
	switch contentType {
	case entity.ContentTypeMovie:
		movie, err := c.getMovieData(id)
		if err != nil {
			return nil, err
		}
		content.Movie = movie
	case entity.ContentTypeSeries:
		series, err := c.getSeriesData(id)
		if err != nil {
			return nil, err
		}
		content.Series = series
	default:
		return nil, entity.NewClientError("Неизвестный тип контента", entity.ErrNotFound)
	}
	return &content, nil
}

func (c *ContentDB) GetPerson(id int) (*entity.Person, error) {
	// no-lint
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
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var person entity.Person
	err := c.DB.QueryRow(query, args...).Scan(
		&person.ID,
		&person.FirstName,
		&person.LastName,
		&person.BirthDate,
		&person.DeathDate,
		&person.StartCareer,
		&person.EndCareer,
		&person.Sex,
		&person.Height,
		&person.Spouse,
		&person.Children,
		&person.PhotoStaticID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewClientError("Персона не найдена", entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении информации о персоне"))
	}
	return &person, nil
}

// GetPersonRoles возвращает список контента, в создании которого персона принимала участие по ID персоны.
// nolint: dupl
func (c *ContentDB) GetPersonRoles(id int) ([]entity.Content, error) {
	// no-lint
	query, args, _ := sq.Select("content_id").
		From("person_role").
		Where(sq.Eq{"person_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении ролей персоны"))
	}
	defer rows.Close()
	var contents []entity.Content
	for rows.Next() {
		var contentID int
		err := rows.Scan(&contentID)
		if err != nil {
			return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении ролей персоны"))
		}
		content, err := c.GetPreviewContent(contentID)
		if err != nil {
			return nil, err
		}
		contents = append(contents, *content)
	}
	return contents, nil
}
