package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/sync/errgroup"
)

type ContentDB struct {
	DB *sqlx.DB
}

type ScanContent struct {
	ID               int
	ContentType      string
	Title            string
	OriginalTitle    sql.NullString
	Slogan           sql.NullString
	Budget           sql.NullString
	AgeRestriction   int
	IMDBRating       float64
	Rating           float64
	Description      string
	PosterStaticID   sql.NullInt64
	TrailerURL       sql.NullString
	BackdropStaticID sql.NullInt64
	Ongoing          bool
	OngoingDate      sql.NullTime
	StreamURL        sql.NullString
}

func NewContentRepository(db *sqlx.DB) repository.Content {
	return &ContentDB{
		DB: db,
	}
}

// getMovieData возвращает информацию о фильме по его ID
func (c *ContentDB) getMovieData(id int) (*entity.Movie, error) {
	query, args, err := sq.Select("premiere", "duration").
		From("movie").
		Where(sq.Eq{"content_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getMovieData"))
	}
	var movie entity.Movie
	var premiere sql.NullTime
	var duration sql.NullInt64
	err = c.DB.QueryRow(query, args...).Scan(&premiere, &duration)
	if err != nil {
		return nil, entity.PSQLQueryErr("getMovieData", err)
	}
	movie.Premiere = premiere.Time
	movie.Duration = int(duration.Int64)
	return &movie, nil
}

// getSeriesData возвращает информацию о сериале по его ID
func (c *ContentDB) getSeriesData(id int) (*entity.Series, error) {
	query, args, err := sq.Select("year_start", "year_end").
		From("series").
		Where(sq.Eq{"content_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getSeriesData"))
	}
	var series entity.Series
	var yearStart, yearEnd sql.NullInt64
	err = c.DB.QueryRow(query, args...).Scan(&yearStart, &yearEnd)
	if err != nil {
		return nil, entity.PSQLQueryErr("getSeriesData", err)
	}
	series.YearStart = int(yearStart.Int64)
	series.YearEnd = int(yearEnd.Int64)
	seasons, err := c.getSeasonsByContentID(id)
	if err != nil {
		return nil, err
	}
	series.Seasons = seasons
	return &series, nil
}

// getEpisodeData возвращает информацию об эпизоде по его ID
func (c *ContentDB) getEpisodeData(id int) (*entity.Episode, error) {
	query, args, err := sq.Select("id", "episode_number", "title", "duration").
		From("episode").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getEpisodeData"))
	}
	var duration sql.NullInt64
	var episode entity.Episode
	err = c.DB.QueryRow(query, args...).Scan(&episode.ID, &episode.EpisodeNumber, &episode.Title, &duration)
	episode.Duration = int(duration.Int64)
	if err != nil {
		return nil, entity.PSQLQueryErr("getEpisodeData", err)
	}
	return &episode, nil
}

// getEpisodesBySeasonID возвращает эпизоды сезона по его ID
// nolint: dupl
func (c *ContentDB) getEpisodesBySeasonID(seasonID int) ([]entity.Episode, error) {
	query, args, err := sq.Select("id").
		From("episode").
		Where(sq.Eq{"season_id": seasonID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getEpisodesBySeasonID"))
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("getEpisodesBySeasonID", err)
	}
	defer rows.Close()
	var episodes []entity.Episode
	for rows.Next() {
		var episodeID int
		err := rows.Scan(&episodeID)
		if err != nil {
			return nil, entity.PSQLQueryErr("getEpisodesBySeasonID при сканировании", err)
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
	query, args, err := sq.Select("id", "title").
		From("season").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getSeasonData"))
	}
	var season entity.Season
	err = c.DB.QueryRow(query, args...).Scan(&season.ID, &season.Title)
	if err != nil {
		return nil, entity.PSQLQueryErr("getSeasonData при сканировании", err)
	}
	return &season, nil
}

// getSeasonsByContentID возвращает сезоны контента по его ID
func (c *ContentDB) getSeasonsByContentID(contentID int) ([]entity.Season, error) {
	query, args, err := sq.Select("id").
		From("season").
		Where(sq.Eq{"series_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getSeasonsByContentID"))
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("getSeasonsByContentID", err)
	}
	defer rows.Close()
	var seasons []entity.Season
	for rows.Next() {
		var seasonID int
		err := rows.Scan(&seasonID)
		if err != nil {
			return nil, entity.PSQLQueryErr("getSeasonsByContentID при сканировании", err)
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
	query, args, err := sq.Select("id").
		From("role").
		Where(sq.Eq{"name_en": role}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return 0, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getRoleIDByName"))
	}
	var roleID int
	err = c.DB.QueryRow(query, args...).Scan(&roleID)
	if err != nil {
		return 0, entity.PSQLQueryErr("getRoleIDByName", err)
	}
	return roleID, nil
}

// getPersonsByRoleAndContentID возвращает персон контента по его ID и роли
func (c *ContentDB) getPersonsByRoleAndContentID(role string, contentID int) ([]entity.Person, error) {
	roleID, err := c.getRoleIDByName(role)
	if err != nil {
		return nil, err
	}
	query, args, err := sq.Select("person_id").
		From("person_role").
		Where(sq.Eq{"content_id": contentID, "role_id": roleID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil,
			entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getPersonsByRoleAndContentID"))
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("getPersonsByRoleAndContentID", err)
	}
	defer rows.Close()
	var persons []entity.Person
	for rows.Next() {
		var personID int
		err := rows.Scan(&personID)
		if err != nil {
			return nil, entity.PSQLQueryErr("getPersonsByRoleAndContentID при сканировании", err)
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
		return nil, entity.PSQLQueryErr("getGenreByID", err)
	}
	return &genre, nil
}

// getContentGenres возвращает жанры контента по его ID
// nolint: dupl
func (c *ContentDB) getContentGenres(id int) ([]entity.Genre, error) {
	query, args, err := sq.Select("genre_id").
		From("genre_content").
		Where(sq.Eq{"content_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getContentGenres"))
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("getContentGenres", err)
	}
	defer rows.Close()
	var genres []entity.Genre
	for rows.Next() {
		var genreID int
		err := rows.Scan(&genreID)
		if err != nil {
			return nil, entity.PSQLQueryErr("getContentGenres при сканировании", err)
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
	query, args, err := sq.Select("name").
		From("country").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getCountryByID"))
	}
	var country entity.Country
	err = c.DB.QueryRow(query, args...).Scan(&country.Name)
	country.ID = id
	if err != nil {
		return nil, entity.PSQLQueryErr("getCountryByID", err)
	}
	return &country, nil
}

// getContentProductionCountries возвращает страны производства контента по его ID
// nolint: dupl
func (c *ContentDB) getContentProductionCountries(id int) ([]entity.Country, error) {
	query, args, err := sq.Select("country_id").
		From("country_content").
		Where(sq.Eq{"content_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil,
			entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getContentProductionCountries"))
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("getContentProductionCountries", err)
	}
	defer rows.Close()
	var countries []entity.Country
	for rows.Next() {
		var countryID int
		err := rows.Scan(&countryID)
		if err != nil {
			return nil, entity.PSQLQueryErr("getContentProductionCountries при сканировании", err)
		}
		country, err := c.getCountryByID(countryID)
		if err != nil {
			return nil, err
		}
		countries = append(countries, *country)
	}
	return countries, nil
}

// getContentFacts возвращает факты о контенте по его ID
func (c *ContentDB) getContentFacts(id int) ([]string, error) {
	query, args, err := sq.Select("fact").
		From("content_fact").
		Where(sq.Eq{"content_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getContentFacts"))
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("getContentFacts", err)
	}
	defer rows.Close()
	var facts []string
	for rows.Next() {
		var fact string
		err := rows.Scan(&fact)
		if err != nil {
			return nil, entity.PSQLQueryErr("getContentFacts", err)
		}
		facts = append(facts, fact)
	}
	return facts, nil
}

// getContentPictures возвращает изображения контента по его ID
func (c *ContentDB) getContentPictures(id int) ([]int, error) {
	query, args, err := sq.Select("static_id").
		From("content_image").
		Where(sq.Eq{"content_id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getContentPictures"))
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("getContentPictures", err)
	}
	defer rows.Close()
	var pictures []int
	for rows.Next() {
		var picture int
		err := rows.Scan(&picture)
		if err != nil {
			return nil, entity.PSQLQueryErr("getContentPictures при сканировании", err)
		}
		pictures = append(pictures, picture)
	}
	return pictures, nil
}

func (c *ContentDB) getContentInfo(id int) (*entity.Content, error) {
	query, args, err := sq.Select(
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
		"streaming_url",
	).
		From("content").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getContentInfo"))
	}
	var scanContent ScanContent
	var content entity.Content
	err = c.DB.QueryRow(query, args...).Scan(
		&scanContent.ID,
		&scanContent.ContentType,
		&scanContent.Title,
		&scanContent.OriginalTitle,
		&scanContent.Slogan,
		&scanContent.Budget,
		&scanContent.AgeRestriction,
		&scanContent.IMDBRating,
		&scanContent.Rating,
		&scanContent.Description,
		&scanContent.PosterStaticID,
		&scanContent.TrailerURL,
		&scanContent.BackdropStaticID,
		&content.Ongoing,
		&scanContent.OngoingDate,
		&scanContent.StreamURL,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrContentNotFound
		}
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении контента"))
	}
	content.ID = scanContent.ID
	content.Title = scanContent.Title
	content.OriginalTitle = scanContent.OriginalTitle.String
	content.Slogan = scanContent.Slogan.String
	content.Budget = scanContent.Budget.String
	content.AgeRestriction = scanContent.AgeRestriction
	content.IMDBRating = scanContent.IMDBRating
	content.Rating = scanContent.Rating
	content.Description = scanContent.Description
	content.PosterStaticID = int(scanContent.PosterStaticID.Int64)
	content.TrailerLink = scanContent.TrailerURL.String
	content.BackdropStaticID = int(scanContent.BackdropStaticID.Int64)
	content.Type = scanContent.ContentType
	if scanContent.OngoingDate.Valid {
		content.OngoingDate = &scanContent.OngoingDate.Time
	}
	if scanContent.StreamURL.Valid {
		content.StreamURL = scanContent.StreamURL.String
	}
	return &content, nil
}

func (c *ContentDB) GetContent(id int) (*entity.Content, error) {
	group, _ := errgroup.WithContext(context.Background())

	// получаем основную информацию о контенте
	content, err := c.getContentInfo(id)
	if err != nil {
		return nil, err
	}

	// запрашиваем дополнительные данные
	group.Go(func() error {
		pictures, err := c.getContentPictures(id)
		if err != nil {
			return err
		}
		content.PicturesStaticID = pictures
		return nil
	})

	group.Go(func() error {
		facts, err := c.getContentFacts(id)
		if err != nil {
			return err
		}
		content.Facts = facts
		return nil
	})

	group.Go(func() error {
		countries, err := c.getContentProductionCountries(id)
		if err != nil {
			return err
		}
		content.Country = countries
		return nil
	})

	group.Go(func() error {
		genres, err := c.getContentGenres(id)
		if err != nil {
			return err
		}
		content.Genres = genres
		return nil
	})

	group.Go(func() error {
		actors, err := c.getPersonsByRoleAndContentID(entity.RoleActor, id)
		if err != nil {
			return err
		}
		content.Actors = actors
		return nil
	})

	group.Go(func() error {
		directors, err := c.getPersonsByRoleAndContentID(entity.RoleDirector, id)
		if err != nil {
			return err
		}
		content.Directors = directors
		return nil
	})

	group.Go(func() error {
		producers, err := c.getPersonsByRoleAndContentID(entity.RoleProducer, id)
		if err != nil {
			return err
		}
		content.Producers = producers
		return nil
	})

	group.Go(func() error {
		writers, err := c.getPersonsByRoleAndContentID(entity.RoleWriter, id)
		if err != nil {
			return err
		}
		content.Writers = writers
		return nil
	})

	group.Go(func() error {
		operators, err := c.getPersonsByRoleAndContentID(entity.RoleOperator, id)
		if err != nil {
			return err
		}
		content.Operators = operators
		return nil
	})

	group.Go(func() error {
		composers, err := c.getPersonsByRoleAndContentID(entity.RoleComposer, id)
		if err != nil {
			return err
		}
		content.Composers = composers
		return nil
	})

	group.Go(func() error {
		editors, err := c.getPersonsByRoleAndContentID(entity.RoleEditor, id)
		if err != nil {
			return err
		}
		content.Editors = editors
		return nil
	})

	group.Go(func() error {
		switch content.Type {
		case entity.ContentTypeMovie:
			movie, err := c.getMovieData(id)
			if err != nil {
				return err
			}
			content.Movie = movie
		case entity.ContentTypeSeries:
			series, err := c.getSeriesData(id)
			if err != nil {
				return err
			}
			content.Series = series
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return content, nil
}

func (c *ContentDB) getPreviewContentInfo(id int) (*entity.Content, error) {
	query, args, err := sq.Select(
		"id",
		"content_type",
		"title",
		"original_title",
		"rating",
		"poster_upload_id",
		"ongoing",
		"ongoing_date",
	).
		From("content").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса getPreviewContentInfo"))
	}
	var scanContent ScanContent
	var content entity.Content
	err = c.DB.QueryRow(query, args...).Scan(
		&scanContent.ID,
		&scanContent.ContentType,
		&scanContent.Title,
		&scanContent.OriginalTitle,
		&scanContent.Rating,
		&scanContent.PosterStaticID,
		&scanContent.Ongoing,
		&scanContent.OngoingDate,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrContentNotFound
		}
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при получении контента"))
	}
	content.ID = scanContent.ID
	content.Type = scanContent.ContentType
	content.Title = scanContent.Title
	content.OriginalTitle = scanContent.OriginalTitle.String
	content.Rating = scanContent.Rating
	content.PosterStaticID = int(scanContent.PosterStaticID.Int64)
	content.Ongoing = scanContent.Ongoing
	if scanContent.OngoingDate.Valid {
		content.OngoingDate = &scanContent.OngoingDate.Time
	}
	return &content, nil
}

// GetPreviewContent возвращает краткую информацию о контенте по его ID.
// Заполняет только поля id, title, original_title, poster_upload_id, countries, genres, actors, directors.
// Для фильма заполняет поля premiere, duration.
// Для сериала заполняет поля year_start, year_end, seasons.
func (c *ContentDB) GetPreviewContent(id int) (*entity.Content, error) {
	group, _ := errgroup.WithContext(context.Background())

	// запрашиваем краткую информацию о контенте
	content, err := c.getPreviewContentInfo(id)
	if err != nil {
		return nil, err
	}

	// запрашиваем дополнительные данные
	group.Go(func() error {
		countries, err := c.getContentProductionCountries(id)
		if err != nil {
			return err
		}
		content.Country = countries
		return nil
	})

	group.Go(func() error {
		genres, err := c.getContentGenres(id)
		if err != nil {
			return err
		}
		content.Genres = genres
		return nil
	})

	group.Go(func() error {
		actors, err := c.getPersonsByRoleAndContentID(entity.RoleActor, id)
		if err != nil {
			return err
		}
		content.Actors = actors
		return nil
	})

	group.Go(func() error {
		directors, err := c.getPersonsByRoleAndContentID(entity.RoleDirector, id)
		if err != nil {
			return err
		}
		content.Directors = directors
		return nil
	})

	group.Go(func() error {
		switch content.Type {
		case entity.ContentTypeMovie:
			movie, err := c.getMovieData(id)
			if err != nil {
				return err
			}
			content.Movie = movie
		case entity.ContentTypeSeries:
			series, err := c.getSeriesData(id)
			if err != nil {
				return err
			}
			content.Series = series
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		return nil, err
	}

	return content, nil
}

func (c *ContentDB) GetPerson(id int) (*entity.Person, error) {
	query, args, err := sq.Select(
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
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса GetPerson"))
	}
	person := new(entity.Person)
	err = c.DB.QueryRowx(query, args...).StructScan(person)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrPersonNotFound
		}
		return nil, entity.PSQLQueryErr("GetPerson", err)
	}

	return person, nil
}

// GetPersonRoles возвращает список контента, в создании которого персона принимала участие по ID персоны.
func (c *ContentDB) GetPersonRoles(personID int) ([]entity.PersonRole, error) {
	// получаем всевозможные роли
	var roles []entity.Role
	query, args, err := sq.Select("id", "name", "name_en").From("role").ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса GetPersonRoles"))
	}
	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("GetPersonRoles при получении списка ролей", err)
	}
	for rows.Next() {
		var role entity.Role
		err = rows.Scan(&role.ID, &role.Name, &role.EnName)
		if err != nil {
			return nil, entity.PSQLQueryErr("GetPersonRoles при сканировании списка ролей", err)
		}
		roles = append(roles, role)
	}
	rows.Close()

	// получаем роли персоны
	personRoles := make([]entity.PersonRole, 0)
	for _, role := range roles {
		query, args, err := sq.Select("content_id").
			From("person_role").
			Where(sq.Eq{"person_id": personID, "role_id": role.ID}).
			PlaceholderFormat(sq.Dollar).
			ToSql()
		if err != nil {
			return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса GetPersonRoles"))
		}
		rows, err := c.DB.Query(query, args...)
		if err != nil {
			return nil, entity.PSQLQueryErr("GetPersonRoles при получении ролей персоны", err)
		}
		for rows.Next() {
			var contentID int
			err := rows.Scan(&contentID)
			if err != nil {
				return nil, entity.PSQLQueryErr("GetPersonRoles при сканировании ролей персоны", err)
			}
			personRoles = append(personRoles, entity.PersonRole{
				PersonID:  personID,
				Role:      role,
				ContentID: contentID,
			})
		}
		rows.Close()
	}
	return personRoles, nil
}

func (c *ContentDB) GetSimilarContent(id int) ([]entity.Content, error) {
	// ну тут уж придётся обойтись без squirrel :)
	query := `
    WITH target_persons AS (
        SELECT pr.person_id
        FROM person_role pr
        WHERE pr.content_id = $1
    ),
    target_genres AS (
        SELECT gc.genre_id
        FROM genre_content gc
        WHERE gc.content_id = $1
    ),
    matched_contents AS (
        SELECT pr.content_id, COUNT(*) AS person_match_count
        FROM person_role pr
        JOIN target_persons tp ON pr.person_id = tp.person_id
        WHERE pr.content_id != $1
        GROUP BY pr.content_id
    ),
    genre_matched_contents AS (
        SELECT gc.content_id, COUNT(*) AS genre_match_count
        FROM genre_content gc
        JOIN target_genres tg ON gc.genre_id = tg.genre_id
        WHERE gc.content_id != $1
        GROUP BY gc.content_id
    )
    SELECT mc.content_id
    FROM matched_contents mc
    JOIN genre_matched_contents gmc ON mc.content_id = gmc.content_id
    ORDER BY (mc.person_match_count + gmc.genre_match_count) DESC
    LIMIT 10;
    `

	rows, err := c.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Собираем результаты
	var contents []entity.Content
	for rows.Next() {
		var contentID int
		err := rows.Scan(&contentID)
		if err != nil {
			return nil, err
		}
		content, err := c.GetPreviewContent(contentID)
		if err != nil {
			return nil, err
		}
		contents = append(contents, *content)
	}

	return contents, nil
}

func (c *ContentDB) GetNearestOngoings(limit int) ([]int, error) {
	query, args, err := sq.Select("id").
		From("content").
		Where(sq.Eq{"ongoing": true}).
		OrderBy("ongoing_date ASC").
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса GetNearestOngoings"))
	}

	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("GetNearestOngoings", err)
	}
	defer rows.Close()

	var contentIDs []int
	for rows.Next() {
		var contentID int
		err := rows.Scan(&contentID)
		if err != nil {
			return nil, entity.PSQLQueryErr("GetNearestOngoings при сканировании", err)
		}
		contentIDs = append(contentIDs, contentID)
	}

	return contentIDs, nil
}

func (c *ContentDB) GetOngoingContentByMonthAndYear(month, year int) ([]int, error) {
	query, args, err := sq.Select("id").
		From("content").
		Where(sq.Eq{"ongoing": true}).
		Where(sq.Expr("EXTRACT(MONTH FROM ongoing_date) = ?", month)).
		Where(sq.Expr("EXTRACT(YEAR FROM ongoing_date) = ?", year)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса GetOngoingContentByMonthAndYear"))
	}

	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("GetOngoingContentByMonthAndYear", err)
	}
	defer rows.Close()

	var contentIDs []int
	for rows.Next() {
		var contentID int
		err := rows.Scan(&contentID)
		if err != nil {
			return nil, entity.PSQLQueryErr("GetOngoingContentByMonthAndYear при сканировании", err)
		}
		contentIDs = append(contentIDs, contentID)
	}

	return contentIDs, nil
}

func (c *ContentDB) GetAllOngoingsYears() ([]int, error) {
	query, args, err := sq.Select("EXTRACT(YEAR FROM ongoing_date)").
		From("content").
		Where(sq.Eq{"ongoing": true}).
		GroupBy("EXTRACT(YEAR FROM ongoing_date)").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса GetAllOngoingsYears"))
	}

	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("GetAllOngoingsYears", err)
	}
	defer rows.Close()

	var years []int
	for rows.Next() {
		var year int
		err := rows.Scan(&year)
		if err != nil {
			return nil, entity.PSQLQueryErr("GetAllOngoingsYears при сканировании", err)
		}
		years = append(years, year)
	}

	return years, nil
}

func (c *ContentDB) IsOngoingContentReleased(contentID int) (bool, error) {
	query, args, err := sq.Select("ongoing").
		From("content").
		Where(sq.Eq{"id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return false, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса IsOngoingContentReleased"))
	}

	var finished bool
	err = c.DB.QueryRow(query, args...).Scan(&finished)
	if err != nil {
		return false, entity.PSQLQueryErr("IsOngoingContentReleased", err)
	}

	return !finished, nil
}

func (c *ContentDB) SetReleasedState(contentID int, released bool) error {
	query, args, err := sq.Update("content").
		Set("ongoing", !released).
		Where(sq.Eq{"id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса SetReleasedState"))
	}

	rowsAffected, err := c.DB.Exec(query, args...)
	if err != nil {
		return entity.PSQLQueryErr("SetReleasedState", err)
	}
	if totalAffected, err := rowsAffected.RowsAffected(); err != nil || totalAffected == 0 {
		return repository.ErrContentNotFound
	}

	return nil
}

func (c *ContentDB) SubscribeOnContent(userID, contentID int) error {
	query, args, err := sq.Insert("ongoing_subscribe").
		Columns("user_id", "content_id").
		Values(userID, contentID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса SubscribeOnContent"))
	}

	_, err = c.DB.Exec(query, args...)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == entity.PSQLUniqueViolation {
		return nil
	}
	if err != nil {
		return entity.PSQLQueryErr("SubscribeOnContent", err)
	}

	return nil
}

func (c *ContentDB) UnsubscribeFromContent(userID, contentID int) error {
	query, args, err := sq.Delete("ongoing_subscribe").
		Where(sq.Eq{"user_id": userID, "content_id": contentID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса UnsubscribeFromContent"))
	}

	// Неважно сколько rows affected, т.к. пользователь может не быть подписан на контент
	_, err = c.DB.Exec(query, args...)
	if err != nil {
		return entity.PSQLQueryErr("UnsubscribeFromContent", err)
	}

	return nil
}

func (c *ContentDB) GetSubscribedContentIDs(userID int) ([]int, error) {
	query, args, err := sq.Select("content_id").
		From("ongoing_subscribe").
		Where(sq.Eq{"user_id": userID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса GetSubscribedContentIDs"))
	}

	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("GetSubscribedContentIDs", err)
	}
	defer rows.Close()

	contentIDs := make([]int, 0)
	for rows.Next() {
		var contentID int
		err := rows.Scan(&contentID)
		if err != nil {
			return nil, entity.PSQLQueryErr("GetSubscribedContentIDs при сканировании", err)
		}
		contentIDs = append(contentIDs, contentID)
	}

	return contentIDs, nil
}

func (c *ContentDB) GetAvailableToWatch(page, limit int) ([]int, error) {
	query, args, err := sq.Select("id").
		From("content").
		Where(sq.NotEq{"streaming_url": ""}).
		Offset(uint64((page - 1) * limit)).
		Limit(uint64(limit)).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, entity.PSQLWrap(err, fmt.Errorf("ошибка при формировании запроса GetAvailableToWatch"))
	}

	rows, err := c.DB.Query(query, args...)
	if err != nil {
		return nil, entity.PSQLQueryErr("GetAvailableToWatch", err)
	}
	defer rows.Close()

	contentIDs := make([]int, 0)
	for rows.Next() {
		var contentID int
		err := rows.Scan(&contentID)
		if err != nil {
			return nil, entity.PSQLQueryErr("GetAvailableToWatch при сканировании", err)
		}
		contentIDs = append(contentIDs, contentID)
	}

	return contentIDs, nil
}
