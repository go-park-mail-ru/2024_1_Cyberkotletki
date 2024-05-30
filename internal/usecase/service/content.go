package service

import (
	"log"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
)

type ContentService struct {
	contentRepo repository.Content
	reviewRepo  repository.Review
	staticRepo  repository.Static
}

func NewContentService(
	contentRepo repository.Content,
	reviewRepo repository.Review,
	staticRepo repository.Static,
) usecase.Content {
	return &ContentService{
		contentRepo: contentRepo,
		reviewRepo:  reviewRepo,
		staticRepo:  staticRepo,
	}
}

func CountryEntityToDTO(countryEntity []entity.Country) []string {
	countries := make([]string, len(countryEntity))
	for i, country := range countryEntity {
		countries[i] = country.Name
	}
	return countries
}

func GenreEntityToDTO(genreEntity []entity.Genre) []string {
	genres := make([]string, len(genreEntity))
	for i, genre := range genreEntity {
		genres[i] = genre.Name
	}
	return genres
}

func PersonEntityToPreviewDTO(personEntity []entity.Person) []dto.PersonPreview {
	persons := make([]dto.PersonPreview, len(personEntity))
	for i, person := range personEntity {
		persons[i] = dto.PersonPreview{
			ID:        person.ID,
			FirstName: person.FirstName,
			LastName:  person.LastName,
		}
	}
	return persons
}

func MovieEntityToDTO(movieEntity entity.Movie) dto.MovieContent {
	return dto.MovieContent{
		Premiere: movieEntity.Premiere,
		Release:  movieEntity.Release,
		Duration: movieEntity.Duration,
	}
}

func SeriesEntityToDTO(seriesEntity entity.Series) dto.SeriesContent {
	seasons := make([]dto.Season, len(seriesEntity.Seasons))
	for seasonIndex, season := range seriesEntity.Seasons {
		episodes := make([]dto.Episode, len(season.Episodes))
		for episodeIndex, episode := range season.Episodes {
			episodes[episodeIndex] = dto.Episode{
				ID:            episode.ID,
				EpisodeNumber: episode.EpisodeNumber,
				Title:         episode.Title,
			}
		}
		seasons[seasonIndex] = dto.Season{
			ID:        season.ID,
			YearStart: season.YearStart,
			YearEnd:   season.YearEnd,
			Episodes:  episodes,
		}
	}
	return dto.SeriesContent{
		YearStart: seriesEntity.YearStart,
		YearEnd:   seriesEntity.YearEnd,
		Seasons:   seasons,
	}

}

func (c *ContentService) PersonsDTOToentity(persons []dto.PersonPreview) []entity.Person {
	entityPersons := make([]entity.Person, len(persons))
	for i, person := range persons {
		personEntity, err := c.contentRepo.GetPerson(person.ID)
		if err != nil {
			log.Println("usecase PersonsDTOToentity", err)
			return []entity.Person{}
		}
		entityPersons[i] = *personEntity
	}
	return entityPersons
}

func (c *ContentService) CountriesDTOToEntity(countries []string) []entity.Country {
	entityCountries := make([]entity.Country, len(countries))
	for i, country := range countries {
		countryEntity, err := c.contentRepo.GetCountryByName(country)
		if err != nil {
			log.Println("usecase CountriesDTOToEntity", err)
			return []entity.Country{}
		}
		entityCountries[i] = *countryEntity
	}
	return entityCountries
}

func (c *ContentService) GenresDTOToEntity(genres []string) []entity.Genre {
	entityGenres := make([]entity.Genre, len(genres))
	for i, genre := range genres {
		genreEntity, err := c.contentRepo.GetGenreByName(genre)
		if err != nil {
			log.Println("usecase GenresDTOToEntity", err)
			return []entity.Genre{}
		}
		entityGenres[i] = *genreEntity
	}
	return entityGenres
}

// GetContentByID возвращает dto.Content по его ID
func (c *ContentService) GetContentByID(id int) (*dto.Content, error) {
	contentEntity, err := c.contentRepo.GetContent(id)
	if err != nil {
		return nil, err
	}
	posterURL, err := c.staticRepo.GetStatic(contentEntity.PosterStaticID)
	if err != nil {
		return nil, err
	}
	rating, err := c.reviewRepo.GetContentRating(id)
	if err != nil {
		return nil, err
	}
	contentDTO := dto.Content{
		ID:             contentEntity.ID,
		Title:          contentEntity.Title,
		OriginalTitle:  contentEntity.OriginalTitle,
		Slogan:         contentEntity.Slogan,
		Budget:         contentEntity.Budget,
		AgeRestriction: contentEntity.AgeRestriction,
		Audience:       contentEntity.Audience,
		Rating:         rating,
		IMDBRating:     contentEntity.IMDBRating,
		Description:    contentEntity.Description,
		PosterURL:      posterURL,
		BoxOffice:      contentEntity.BoxOffice,
		Marketing:      contentEntity.Marketing,
		Countries:      CountryEntityToDTO(contentEntity.Country),
		Genres:         GenreEntityToDTO(contentEntity.Genres),
		Actors:         PersonEntityToPreviewDTO(contentEntity.Actors),
		Directors:      PersonEntityToPreviewDTO(contentEntity.Directors),
		Producers:      PersonEntityToPreviewDTO(contentEntity.Producers),
		Writers:        PersonEntityToPreviewDTO(contentEntity.Writers),
		Operators:      PersonEntityToPreviewDTO(contentEntity.Operators),
		Composers:      PersonEntityToPreviewDTO(contentEntity.Composers),
		Editors:        PersonEntityToPreviewDTO(contentEntity.Editors),
		Type:           contentEntity.Type,
	}
	if contentEntity.Type == entity.ContentTypeMovie {
		contentDTO.Movie = MovieEntityToDTO(*contentEntity.Movie)
	} else if contentEntity.Type == entity.ContentTypeSeries {
		contentDTO.Series = SeriesEntityToDTO(*contentEntity.Series)
	}

	return &contentDTO, nil
}

func (c *ContentService) CreateContent(content dto.Content) (*dto.Content, error) {
	if err := entity.ValidateContent(
		content.Title,
		content.OriginalTitle,
		content.Slogan,
		content.Description,
		content.Budget,
		content.AgeRestriction,
		content.Audience,
		content.BoxOffice,
		content.Marketing,
		content.IMDBRating,
		content.Type,
	); err != nil {
		return nil, err

	}

	contentEntity := entity.Content{
		Title:          content.Title,
		OriginalTitle:  content.OriginalTitle,
		Slogan:         content.Slogan,
		Budget:         content.Budget,
		AgeRestriction: content.AgeRestriction,
		Audience:       content.Audience,
		IMDBRating:     content.IMDBRating,
		Description:    content.Description,
		BoxOffice:      content.BoxOffice,
		Marketing:      content.Marketing,
		Country:        c.CountriesDTOToEntity(content.Countries),
		Genres:         c.GenresDTOToEntity(content.Genres),
		Actors:         c.PersonsDTOToentity(content.Actors),
		Directors:      c.PersonsDTOToentity(content.Directors),
		Producers:      c.PersonsDTOToentity(content.Producers),
		Writers:        c.PersonsDTOToentity(content.Writers),
		Operators:      c.PersonsDTOToentity(content.Operators),
		Composers:      c.PersonsDTOToentity(content.Composers),
		Editors:        c.PersonsDTOToentity(content.Editors),
		Type:           content.Type,
	}
	if content.PosterURL != "" {
		staticID, err := c.staticRepo.AddStaticToDB("posters", content.PosterURL)
		if err != nil {
			log.Println("usecase AddStaticToDB", err)
			return nil, err
		}
		contentEntity.PosterStaticID = staticID
	}

	if content.Type == entity.ContentTypeMovie {
		contentEntity.Movie = &entity.Movie{
			Premiere: content.Movie.Premiere,
			Release:  content.Movie.Release,
			Duration: content.Movie.Duration,
		}
	} else if content.Type == entity.ContentTypeSeries {
		seasons := make([]entity.Season, len(content.Series.Seasons))
		for seasonIndex, season := range content.Series.Seasons {
			episodes := make([]entity.Episode, len(season.Episodes))
			for episodeIndex, episode := range season.Episodes {
				episodes[episodeIndex] = entity.Episode{
					EpisodeNumber: episode.EpisodeNumber,
					Title:         episode.Title,
				}
			}
			seasons[seasonIndex] = entity.Season{
				YearStart: season.YearStart,
				YearEnd:   season.YearEnd,
				Episodes:  episodes,
			}
		}
		contentEntity.Series = &entity.Series{
			YearStart: content.Series.YearStart,
			YearEnd:   content.Series.YearEnd,
			Seasons:   seasons,
		}
	}
	fullContent, err := c.contentRepo.AddContent(&contentEntity)
	if err != nil {
		log.Println("usecase AddContent 258", err)
		return nil, err
	}
	return c.GetContentByID(fullContent.ID)

}

// GetPersonByID возвращает dto.Person по его ID
func (c *ContentService) GetPersonByID(id int) (*dto.Person, error) {
	personEntity, err := c.contentRepo.GetPerson(id)
	if err != nil {
		return nil, err
	}
	personDTO := dto.Person{
		ID:          personEntity.ID,
		FirstName:   personEntity.FirstName,
		LastName:    personEntity.LastName,
		BirthDate:   personEntity.BirthDate,
		DeathDate:   personEntity.DeathDate,
		StartCareer: personEntity.StartCareer,
		EndCareer:   personEntity.EndCareer,
		Sex:         personEntity.Sex,
		BirthPlace:  personEntity.BirthPlace,
		Height:      personEntity.Height,
		Spouse:      personEntity.Spouse,
		Children:    personEntity.Children,
	}
	contentRoles, err := c.contentRepo.GetPersonRoles(personEntity.ID)
	if err != nil {
		return nil, err
	}
	personDTO.Roles = make([]dto.PreviewContentCard, len(contentRoles))
	for roleIndex, role := range contentRoles {
		personDTO.Roles[roleIndex] = dto.PreviewContentCard{
			ID:            role.ID,
			Title:         role.Title,
			OriginalTitle: role.OriginalTitle,
		}
		posterURL, err := c.staticRepo.GetStatic(role.PosterStaticID)
		if err != nil {
			return nil, err
		}
		personDTO.Roles[roleIndex].Poster = posterURL
	}
	if personEntity.PhotoStaticID != 0 {
		static, err := c.staticRepo.GetStatic(personEntity.PhotoStaticID)
		if err != nil {
			return nil, err
		}
		personDTO.PhotoURL = static
	}
	return &personDTO, nil
}
