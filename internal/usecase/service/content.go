package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
)

type ContentService struct {
	contentRepo repository.Content
	staticUC    usecase.Static
}

func NewContentService(contentRepo repository.Content, staticUC usecase.Static) usecase.Content {
	return &ContentService{
		contentRepo: contentRepo,
		staticUC:    staticUC,
	}
}

func countryEntityToDTO(countryEntity []entity.Country) []string {
	countries := make([]string, len(countryEntity))
	for i, country := range countryEntity {
		countries[i] = country.Name
	}
	return countries
}

func genreEntityToDTO(genreEntity []entity.Genre) []string {
	genres := make([]string, len(genreEntity))
	for i, genre := range genreEntity {
		genres[i] = genre.Name
	}
	return genres
}

func personEntityToPreviewDTO(personEntity []entity.Person) []dto.PersonPreview {
	persons := make([]dto.PersonPreview, len(personEntity))
	for i, person := range personEntity {
		persons[i] = dto.PersonPreview{
			ID:     person.ID,
			Name:   person.Name,
			EnName: person.EnName,
		}
	}
	return persons
}

func movieEntityToDTO(movieEntity entity.Movie) dto.MovieContent {
	premiere := movieEntity.Premiere
	if premiere.IsZero() {
		return dto.MovieContent{
			Duration: movieEntity.Duration,
		}
	}
	return dto.MovieContent{
		Premiere: &premiere,
		Duration: movieEntity.Duration,
	}
}

func seriesEntityToDTO(seriesEntity entity.Series) dto.SeriesContent {
	seasons := make([]dto.Season, len(seriesEntity.Seasons))
	for seasonIndex, season := range seriesEntity.Seasons {
		episodes := make([]dto.Episode, len(season.Episodes))
		for episodeIndex, episode := range season.Episodes {
			episodes[episodeIndex] = dto.Episode{
				ID:            episode.ID,
				EpisodeNumber: episode.EpisodeNumber,
				Title:         episode.Title,
				Duration:      episode.Duration,
			}
		}
		seasons[seasonIndex] = dto.Season{
			ID:       season.ID,
			Episodes: episodes,
		}
	}
	return dto.SeriesContent{
		YearStart: seriesEntity.YearStart,
		YearEnd:   seriesEntity.YearEnd,
		Seasons:   seasons,
	}

}

// GetContentByID возвращает dto.Content по его ID
func (c *ContentService) GetContentByID(id int) (*dto.Content, error) {
	contentEntity, err := c.contentRepo.GetContent(id)
	switch {
	case errors.Is(err, repository.ErrContentNotFound):
		return nil, usecase.ErrContentNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении контента"), err)
	}
	posterURL, err := c.staticUC.GetStatic(contentEntity.PosterStaticID)
	switch {
	case errors.Is(err, usecase.ErrStaticNotFound):
		// Если постер не найден, возвращаем пустую строку
		posterURL = ""
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении постера"), err)
	}
	backdropURL, err := c.staticUC.GetStatic(contentEntity.BackdropStaticID)
	switch {
	case errors.Is(err, usecase.ErrStaticNotFound):
		// Если фоновое изображение не найдено, возвращаем пустую строку
		backdropURL = ""
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении фонового изображения"), err)
	}
	pictures := make([]string, len(contentEntity.PicturesStaticID))
	for index, pictureID := range contentEntity.PicturesStaticID {
		pictureURL, err := c.staticUC.GetStatic(pictureID)
		switch {
		case errors.Is(err, usecase.ErrStaticNotFound):
			// Если изображение не найдено, возвращаем пустую строку
			pictures[index] = ""
		case err != nil:
			return nil, entity.UsecaseWrap(errors.New("ошибка при получении изображения"), err)
		}
		pictures[index] = pictureURL
	}
	contentDTO := dto.Content{
		ID:             contentEntity.ID,
		Title:          contentEntity.Title,
		OriginalTitle:  contentEntity.OriginalTitle,
		Slogan:         contentEntity.Slogan,
		Budget:         contentEntity.Budget,
		AgeRestriction: contentEntity.AgeRestriction,
		Rating:         contentEntity.Rating,
		IMDBRating:     contentEntity.IMDBRating,
		Description:    contentEntity.Description,
		PosterURL:      posterURL,
		TrailerLink:    contentEntity.TrailerLink,
		BackdropURL:    backdropURL,
		PicturesURL:    pictures,
		Facts:          contentEntity.Facts,
		Countries:      countryEntityToDTO(contentEntity.Country),
		Genres:         genreEntityToDTO(contentEntity.Genres),
		Actors:         personEntityToPreviewDTO(contentEntity.Actors),
		Directors:      personEntityToPreviewDTO(contentEntity.Directors),
		Producers:      personEntityToPreviewDTO(contentEntity.Producers),
		Writers:        personEntityToPreviewDTO(contentEntity.Writers),
		Operators:      personEntityToPreviewDTO(contentEntity.Operators),
		Composers:      personEntityToPreviewDTO(contentEntity.Composers),
		Editors:        personEntityToPreviewDTO(contentEntity.Editors),
		Type:           contentEntity.Type,
	}
	switch contentEntity.Type {
	case entity.ContentTypeMovie:
		contentDTO.Movie = movieEntityToDTO(*contentEntity.Movie)
	case entity.ContentTypeSeries:
		contentDTO.Series = seriesEntityToDTO(*contentEntity.Series)
	}

	return &contentDTO, nil
}

// GetPersonByID возвращает dto.Person по его ID
func (c *ContentService) GetPersonByID(id int) (*dto.Person, error) {
	personEntity, err := c.contentRepo.GetPerson(id)
	switch {
	case errors.Is(err, repository.ErrPersonNotFound):
		return nil, usecase.ErrPersonNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении персоны"), err)
	}
	photo, err := c.staticUC.GetStatic(personEntity.GetPhotoStaticID())
	switch {
	case errors.Is(err, usecase.ErrStaticNotFound):
		// Если фото не найдено, возвращаем пустую строку
		photo = ""
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении фото"), err)
	}
	personDTO := dto.Person{
		ID:       personEntity.ID,
		Name:     personEntity.Name,
		EnName:   personEntity.EnName,
		Sex:      personEntity.Sex,
		PhotoURL: photo,
	}
	if personEntity.BirthDate.Valid {
		personDTO.BirthDate = &personEntity.BirthDate.Time
	}
	if personEntity.DeathDate.Valid {
		personDTO.DeathDate = &personEntity.DeathDate.Time
	}
	if personEntity.Height.Valid {
		personDTO.Height = int(personEntity.Height.Int64)
	}
	contentRoles, err := c.contentRepo.GetPersonRoles(personEntity.ID)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении ролей персоны"), err)
	}
	personDTO.Roles = make(map[string][]dto.PreviewContentCardVertical, len(contentRoles))
	for _, role := range contentRoles {
		content, err := c.contentRepo.GetPreviewContent(role.ContentID)
		switch {
		case errors.Is(err, repository.ErrContentNotFound):
			continue
		case err != nil:
			return nil, entity.UsecaseWrap(errors.New("ошибка при получении контента"), err)
		}
		posterURL, err := c.staticUC.GetStatic(content.PosterStaticID)
		switch {
		case errors.Is(err, usecase.ErrStaticNotFound):
			// Если постер не найден, возвращаем пустую строку
			posterURL = ""
		case err != nil:
			return nil, entity.UsecaseWrap(errors.New("ошибка при получении постера"), err)
		}
		roleContent := dto.PreviewContentCardVertical{
			ID:     content.ID,
			Title:  content.Title,
			Genres: genreEntityToDTO(content.Genres),
			Poster: posterURL,
			Rating: content.Rating,
			Type:   content.Type,
		}
		switch content.Type {
		case entity.ContentTypeMovie:
			roleContent.ReleaseYear = content.Movie.Premiere.Year()
		case entity.ContentTypeSeries:
			roleContent.YearStart = content.Series.YearStart
			roleContent.YearEnd = content.Series.YearEnd
		}
		if _, ok := personDTO.Roles[role.Role.Name]; !ok {
			personDTO.Roles[role.Role.Name] = make([]dto.PreviewContentCardVertical, 0)
		}
		personDTO.Roles[role.Role.Name] = append(personDTO.Roles[role.Role.Name], roleContent)
	}
	return &personDTO, nil
}

// GetPreviewContentByID возвращает dto.PreviewContent по его ID
func (c *ContentService) GetPreviewContentByID(id int) (*dto.PreviewContent, error) {
	contentEntity, err := c.contentRepo.GetPreviewContent(id)
	switch {
	case errors.Is(err, repository.ErrContentNotFound):
		return nil, usecase.ErrContentNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении контента"), err)
	}
	posterURL, err := c.staticUC.GetStatic(contentEntity.PosterStaticID)
	switch {
	case errors.Is(err, usecase.ErrStaticNotFound):
		// Если постер не найден, возвращаем пустую строку
		posterURL = ""
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении постера"), err)
	}
	genres := genreEntityToDTO(contentEntity.Genres)
	genre := ""
	if len(genres) > 0 {
		genre = genres[0]
	}
	countries := countryEntityToDTO(contentEntity.Country)
	country := ""
	if len(countries) > 0 {
		country = countries[0]
	}
	directors := personEntityToPreviewDTO(contentEntity.Directors)
	director := ""
	if len(directors) > 0 {
		director = directors[0].Name
	}
	actorEntities := personEntityToPreviewDTO(contentEntity.Actors)
	actors := make([]string, len(actorEntities))
	for index, actor := range actorEntities {
		actors[index] = actor.Name
	}
	duration := 0
	releaseYear := 0
	yearStart := 0
	yearEnd := 0
	if contentEntity.Type == entity.ContentTypeMovie && contentEntity.Movie != nil {
		duration = contentEntity.Movie.Duration
		releaseYear = contentEntity.Movie.Premiere.Year()
	} else if contentEntity.Type == entity.ContentTypeSeries && contentEntity.Series != nil {
		yearStart = contentEntity.Series.YearStart
		yearEnd = contentEntity.Series.YearEnd
	}
	previewContentDTO := dto.PreviewContent{
		ID:            contentEntity.ID,
		Title:         contentEntity.Title,
		OriginalTitle: contentEntity.OriginalTitle,
		Country:       country,
		Genre:         genre,
		Director:      director,
		Actors:        actors,
		Poster:        posterURL,
		Rating:        contentEntity.Rating,
		Type:          contentEntity.Type,
		Duration:      duration,
		ReleaseYear:   releaseYear,
		YearEnd:       yearEnd,
		YearStart:     yearStart,
	}
	return &previewContentDTO, nil
}
