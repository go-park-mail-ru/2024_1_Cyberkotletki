package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/lib/pq"
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

type ScanPerson struct {
	ID            int
	FirstName     string
	LastName      string
	BirthDate     sql.NullTime
	DeathDate     sql.NullTime
	StartCareer   sql.NullTime
	EndCareer     sql.NullTime
	Sex           string
	Height        sql.NullInt64
	Spouse        sql.NullString
	Children      sql.NullString
	PhotoStaticID sql.NullInt64
}

func NewContentRepository(database config.PostgresDatabase) (repository.Content, error) {
	db, err := sql.Open("postgres", database.ConnectURL)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// запросы занимают до 1 с, 1 соединение может обработать 1 req/sec или 60 req/min
	// оптимальным будет 25 соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute)

	return &ContentDB{
		DB: db,
	}, nil
}

// GetMovieDataByID возвращает информацию о фильме по его ID
func (c *ContentDB) GetMovieDataByID(id int) (*entity.Movie, error) {
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

// GetSeriesDataByID возвращает информацию о сериале по его ID
func (c *ContentDB) GetSeriesDataByID(id int) (*entity.Series, error) {
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

// GetRoleIDByName возвращает айди роли по ее названию
func (c *ContentDB) GetRoleIDByName(role string) (int, error) {
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
	roleID, err := c.GetRoleIDByName(role)
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

// GetGenreByID возвращает жанр по его ID
func (c *ContentDB) GetGenreByID(id int) (*entity.Genre, error) {
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

// GetGenreByName возвращает жанр по его названию
func (c *ContentDB) GetGenreByName(name string) (*entity.Genre, error) {
	query, args, _ := sq.Select("id").
		From("genre").
		Where(sq.Eq{"name": name}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var genre entity.Genre
	err := c.DB.QueryRow(query, args...).Scan(&genre.ID)
	genre.Name = name
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
		genre, err := c.GetGenreByID(genreID)
		if err != nil {
			return nil, err
		}
		genres = append(genres, *genre)
	}
	return genres, nil

}

// GetCountryByID возвращает страну по ее ID
func (c *ContentDB) GetCountryByID(id int) (*entity.Country, error) {
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

// GetCountryByName возвращает страну по ее названию
func (c *ContentDB) GetCountryByName(name string) (*entity.Country, error) {
	query, args, _ := sq.Select("id").
		From("country").
		Where(sq.Eq{"name": name}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	var country entity.Country
	err := c.DB.QueryRow(query, args...).Scan(&country.ID)
	country.Name = name
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
		country, err := c.GetCountryByID(countryID)
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

func (c *ContentDB) addPersonToContent(role string, contentID, personID int) error {
	role_ID, err := c.GetRoleIDByName(role)
	if err != nil {
		return err
	}
	query, args, _ := sq.Insert("person_role").
		Columns("role_id", "content_id", "person_id").
		Values(role_ID, contentID, personID).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err = c.DB.Exec(query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case entity.PSQLForeignKeyViolation:
				return entity.NewClientError("Контента или персоны с таким id нет Не найден", entity.ErrNotFound)
			case entity.PSQLUniqueViolation:
				return entity.NewClientError("Такая роль уже добавлена", entity.ErrAlreadyExists)
			default:
				return entity.PSQLWrap(err, fmt.Errorf("ошибка при добавлении персоны к контенту"))
			}
		}
		return entity.PSQLWrap(err, errors.New("ошибка при добавлении персоны к контенту"))
	}
	return err
}

func (c *ContentDB) addMovieData(id int, movie *entity.Movie) error {
	query, args, _ := sq.Insert("movie").
		Columns("content_id", "premiere", "release", "duration").
		Values(id, movie.Premiere, movie.Release, movie.Duration).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err := c.DB.Exec(query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case entity.PSQLCheckViolation:
				return entity.NewClientError("Некорректные данные фильма", entity.ErrBadRequest)
			default:
				return entity.PSQLWrap(err, fmt.Errorf("ошибка при добавлении информации о фильме"))
			}
		}
		return entity.PSQLWrap(err, errors.New("ошибка при добавлении информации о фильме"))
	}

	return nil
}

// addSeriesData добавляет информацию о сериале
func (c *ContentDB) addSeriesData(id int, series *entity.Series) error {
	query, args, _ := sq.Insert("series").
		Columns("content_id", "year_start", "year_end").
		Values(id, series.YearStart, series.YearEnd).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err := c.DB.Exec(query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case entity.PSQLCheckViolation:
				return entity.NewClientError("Некорректные данные сериала", entity.ErrBadRequest)
			default:
				return entity.PSQLWrap(err, fmt.Errorf("ошибка при добавлении информации о сериале"))
			}
		}
		return entity.PSQLWrap(err, errors.New("ошибка при добавлении информации о сериале"))
	}

	// Добавление сезонов сериала
	for _, season := range series.Seasons {
		err := c.addSeasonToSeries(id, &season)
		if err != nil {
			return err
		}
	}
	return nil
}

// addSeasonToSeries добавляет сезон к сериалу по его ID
func (c *ContentDB) addSeasonToSeries(seriesID int, season *entity.Season) error {
	// no-lint
	query, args, _ := sq.Insert("season").
		Columns("id", "year_start", "year_end", "content_id").
		Values(season.ID, season.YearStart, season.YearEnd, seriesID).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err := c.DB.Exec(query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case entity.PSQLCheckViolation:
				return entity.NewClientError("Некорректные данные сезона", entity.ErrBadRequest)
			default:
				return entity.PSQLWrap(err, fmt.Errorf("ошибка при добавлении сезона к сериалу"))
			}
		}
		return entity.PSQLWrap(err, errors.New("ошибка при добавлении сезона к сериалу"))
	}

	for _, episode := range season.Episodes {
		err = c.addEpisodeToSeason(season.ID, &episode)
		if err != nil {
			return err
		}
	}

	return nil
}

// addEpisodeToSeason добавляет эпизод к сезону по его ID
func (c *ContentDB) addEpisodeToSeason(seasonID int, episode *entity.Episode) error {
	// no-lint
	query, args, _ := sq.Insert("episode").
		Columns("id", "episode_number", "title", "season_id").
		Values(episode.ID, episode.EpisodeNumber, episode.Title, seasonID).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	_, err := c.DB.Exec(query, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case entity.PSQLForeignKeyViolation:
				return entity.NewClientError("Сезона с таким id нет", entity.ErrNotFound)
			case entity.PSQLCheckViolation:
				return entity.NewClientError("Некорректные данные эпизода", entity.ErrBadRequest)
			default:
				return entity.PSQLWrap(err, fmt.Errorf("ошибка при добавлении эпизода к сезону"))
			}
		}
		return entity.PSQLWrap(err, errors.New("ошибка при добавлении эпизода к сезону"))
	}

	return nil
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
		movie, err := c.GetMovieDataByID(id)
		if err != nil {
			return nil, err
		}
		content.Movie = movie
	case entity.ContentTypeSeries:
		series, err := c.GetSeriesDataByID(id)
		if err != nil {
			log.Printf("Rep GetCont 647 Ошибка при получении сериала %v", err)
			return nil, err
		}
		content.Series = series
	default:
		return nil, entity.NewClientError("Неизвестный тип контента, не", entity.ErrNotFound)
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
	var scanContent ScanContent
	var content entity.Content
	err := c.DB.QueryRow(query, args...).Scan(
		&scanContent.ID,
		&scanContent.Title,
		&scanContent.OriginalTitle,
		&scanContent.PosterStaticID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewClientError("Контент не найден", entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении краткой информации о контенте"))
	}
	content.ID = scanContent.ID
	content.Title = scanContent.Title
	content.OriginalTitle = scanContent.OriginalTitle.String
	content.PosterStaticID = scanContent.PosterStaticID
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
		movie, err := c.GetMovieDataByID(id)
		if err != nil {
			return nil, err
		}
		content.Movie = movie
	case entity.ContentTypeSeries:
		series, err := c.GetSeriesDataByID(id)
		if err != nil {
			return nil, err
		}
		content.Series = series
	default:
		return nil, entity.NewClientError("Неизвестный тип контента, не найден", entity.ErrNotFound)
	}
	return &content, nil
}

func (c *ContentDB) AddContent(content *entity.Content) (*entity.Content, error) {
	log.Printf("Добавление контента AddContent")
	// вставка контента
	query, args, _ := sq.Insert("content").
		Columns("title",
			"original_title",
			"slogan",
			"budget",
			"age_restriction",
			"audience",
			"imdb",
			"description",
			"poster_upload_id",
			"box_office",
			"marketing_budget").
		Values(
			content.Title,
			content.OriginalTitle,
			content.Slogan,
			content.Budget,
			content.AgeRestriction,
			content.Audience,
			content.IMDBRating,
			content.Description,
			content.PosterStaticID,
			content.BoxOffice,
			content.Marketing,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	err := c.DB.QueryRow(query, args...).Scan(&content.ID)
	if err != nil {
		log.Printf("Ошибка при добавлении контента %v", err)
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case entity.PSQLCheckViolation:
				return nil, entity.NewClientError("Некорректные данные контента", entity.ErrBadRequest)
			case entity.PSQLForeignKeyViolation:
				return nil, entity.NewClientError("Постера с таким id нет, не найден", entity.ErrNotFound)
			case entity.PSQLUniqueViolation:
				return nil, entity.NewClientError("Контент с таким id уже существует", entity.ErrAlreadyExists)
			default:
				return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при добавлении контента"))
			}
		}
		return nil, entity.PSQLWrap(err, errors.New("ошибка при добавлении контента"))
	}

	// вставка страны производства контента
	for _, country := range content.Country {
		query, args, _ = sq.Insert("country_content").
			Columns("content_id", "country_id").
			Values(content.ID, country.ID).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		_, err = c.DB.Exec(query, args...)
		if err != nil {
			log.Printf("Ошибка при добавлении страны производства контента %v", err)
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				switch pqErr.Code {
				case entity.PSQLForeignKeyViolation:
					return nil, entity.NewClientError("Контента или страны с таким id нет,не найден", entity.ErrNotFound)
				case entity.PSQLCheckViolation:
					return nil, entity.NewClientError("Некорректные данные страны производства контента", entity.ErrBadRequest)
				default:
					return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при добавлении страны производства контента"))
				}
			}
			return nil, entity.PSQLWrap(err, errors.New("ошибка при добавлении страны производства контента"))
		}
	}

	// вставка жанра контента
	for _, genre := range content.Genres {
		query, args, _ = sq.Insert("genre_content").
			Columns("content_id", "genre_id").
			Values(content.ID, genre.ID).
			PlaceholderFormat(sq.Dollar).
			ToSql()

		_, err = c.DB.Exec(query, args...)
		if err != nil {
			log.Printf("Ошибка при добавлении жанра контента %v", err)
			var pqErr *pq.Error
			if errors.As(err, &pqErr) {
				switch pqErr.Code {
				case entity.PSQLForeignKeyViolation:
					return nil, entity.NewClientError("Контента или жанра с таким id нет, не найден", entity.ErrNotFound)
				case entity.PSQLCheckViolation:
					return nil, entity.NewClientError("Некорректные данные жанра контента", entity.ErrBadRequest)
				default:
					return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при добавлении жанра контента"))
				}
			}
			return nil, entity.PSQLWrap(err, errors.New("ошибка при добавлении жанра контента"))
		}
	}

	// Вставка данных о людях, связанных с контентом
	roles := []struct {
		Role    string
		Persons []entity.Person
	}{
		{Role: entity.RoleActor, Persons: content.Actors},
		{Role: entity.RoleDirector, Persons: content.Directors},
		{Role: entity.RoleProducer, Persons: content.Producers},
		{Role: entity.RoleWriter, Persons: content.Writers},
		{Role: entity.RoleOperator, Persons: content.Operators},
		{Role: entity.RoleComposer, Persons: content.Composers},
		{Role: entity.RoleEditor, Persons: content.Editors},
	}
	for _, rolePersons := range roles {
		for _, person := range rolePersons.Persons {
			err := c.addPersonToContent(rolePersons.Role, content.ID, person.ID)
			if err != nil {
				log.Printf("Ошибка при добавлении персоны к контенту %v", err)
				return nil, err
			}
		}
	}

	// Вставка специфических данных для фильмов и сериалов
	switch content.Type {
	case entity.ContentTypeMovie:
		err := c.addMovieData(content.ID, content.Movie)
		if err != nil {
			log.Printf("Ошибка при добавлении информации о фильме %v", err)
			return nil, err
		}
	case entity.ContentTypeSeries:
		err := c.addSeriesData(content.ID, content.Series)
		if err != nil {
			log.Printf("Ошибка при добавлении информации о сериале %v", err)
			return nil, err
		}
	default:
		return nil, entity.NewClientError("Неизвестный тип контента, не найден", entity.ErrNotFound)
	}

	return nil, err

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
	var scanPerson ScanPerson
	var person entity.Person
	err := c.DB.QueryRow(query, args...).Scan(
		&scanPerson.ID,
		&scanPerson.FirstName,
		&scanPerson.LastName,
		&scanPerson.BirthDate,
		&scanPerson.DeathDate,
		&scanPerson.StartCareer,
		&scanPerson.EndCareer,
		&scanPerson.Sex,
		&scanPerson.Height,
		&scanPerson.Spouse,
		&scanPerson.Children,
		&scanPerson.PhotoStaticID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.NewClientError("Персона не найдена", entity.ErrNotFound)
		}
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении информации о персоне"))
	}
	person.ID = scanPerson.ID
	person.FirstName = scanPerson.FirstName
	person.LastName = scanPerson.LastName
	person.BirthDate = scanPerson.BirthDate.Time
	person.DeathDate = scanPerson.DeathDate.Time
	person.StartCareer = scanPerson.StartCareer.Time
	person.EndCareer = scanPerson.EndCareer.Time
	person.Sex = scanPerson.Sex
	if scanPerson.Height.Valid {
		person.Height = int(scanPerson.Height.Int64)
	}
	person.Spouse = scanPerson.Spouse.String
	person.Children = scanPerson.Children.String
	if scanPerson.PhotoStaticID.Valid {
		person.PhotoStaticID = int(scanPerson.PhotoStaticID.Int64)
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
