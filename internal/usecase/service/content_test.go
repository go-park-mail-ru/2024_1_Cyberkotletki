package service

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestContentService_GetContentByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		RequestID            int
		ExpectedContent      *dto.Content
		ExpectedError        error
		SetupContentRepoMock func(contentRepo *mockrepo.MockContent)
		SetupReviewRepoMock  func(reviewRepo *mockrepo.MockReview)
		SetupStaticRepoMock  func(staticRepo *mockrepo.MockStatic)
	}{
		{
			Name:      "Успешное получение фильма",
			RequestID: 1,
			ExpectedContent: &dto.Content{
				ID:             1,
				Title:          "Бэтмен",
				OriginalTitle:  "Batman",
				Slogan:         "I'm Batman",
				Budget:         1000000,
				AgeRestriction: 18,
				Audience:       1000000,
				Rating:         9.1,
				IMDBRating:     8.2,
				Description:    "Описание фильма или сериала",
				PosterURL:      "/static/poster.jpg",
				BoxOffice:      10000,
				Marketing:      10000,
				Countries:      []string{"Россия", "США"},
				Genres:         []string{"Боевик", "Драма"},
				Actors: []dto.PersonPreview{
					{
						ID:        1,
						FirstName: "Киану",
						LastName:  "Ривз",
					},
					{
						ID:        2,
						FirstName: "Кристиан",
						LastName:  "Бэйл",
					},
				},
				Directors: []dto.PersonPreview{
					{
						ID:        3,
						FirstName: "Квентин",
						LastName:  "Тарантино",
					},
				},
				Producers: []dto.PersonPreview{
					{
						ID:        4,
						FirstName: "Кристофер",
						LastName:  "Нолан",
					},
				},
				Writers: []dto.PersonPreview{
					{
						ID:        5,
						FirstName: "Джонатан",
						LastName:  "Нолан",
					},
				},
				Operators: []dto.PersonPreview{
					{
						ID:        6,
						FirstName: "Родриго",
						LastName:  "Прито",
					},
				},
				Composers: []dto.PersonPreview{
					{
						ID:        7,
						FirstName: "Ханс",
						LastName:  "Циммер",
					},
				},
				Editors: []dto.PersonPreview{
					{
						ID:        8,
						FirstName: "Ли",
						LastName:  "Смит",
					},
				},
				Type: "movie",
				Movie: dto.MovieContent{
					Premiere: time.Time{},
					Release:  time.Time{},
					Duration: 100,
				},
			},
			ExpectedError: nil,
			SetupContentRepoMock: func(contentRepo *mockrepo.MockContent) {
				contentRepo.EXPECT().GetContent(1).Return(&entity.Content{
					ID:             1,
					Title:          "Бэтмен",
					OriginalTitle:  "Batman",
					Slogan:         "I'm Batman",
					Budget:         1000000,
					AgeRestriction: 18,
					Audience:       1000000,
					IMDBRating:     8.2,
					Description:    "Описание фильма или сериала",
					PosterStaticID: 1,
					BoxOffice:      10000,
					Marketing:      10000,
					Country:        []entity.Country{{Name: "Россия"}, {Name: "США"}},
					Genres:         []entity.Genre{{Name: "Боевик"}, {Name: "Драма"}},
					Actors:         []entity.Person{{ID: 1, FirstName: "Киану", LastName: "Ривз"}, {ID: 2, FirstName: "Кристиан", LastName: "Бэйл"}},
					Directors:      []entity.Person{{ID: 3, FirstName: "Квентин", LastName: "Тарантино"}},
					Producers:      []entity.Person{{ID: 4, FirstName: "Кристофер", LastName: "Нолан"}},
					Writers:        []entity.Person{{ID: 5, FirstName: "Джонатан", LastName: "Нолан"}},
					Operators:      []entity.Person{{ID: 6, FirstName: "Родриго", LastName: "Прито"}},
					Composers:      []entity.Person{{ID: 7, FirstName: "Ханс", LastName: "Циммер"}},
					Editors:        []entity.Person{{ID: 8, FirstName: "Ли", LastName: "Смит"}},
					Type:           entity.ContentTypeMovie,
					Movie: &entity.Movie{
						Premiere: time.Time{},
						Release:  time.Time{},
						Duration: 100,
					},
				}, nil)
			},
			SetupReviewRepoMock: func(reviewRepo *mockrepo.MockReview) {
				reviewRepo.EXPECT().GetContentRating(1).Return(9.1, nil)
			},
			SetupStaticRepoMock: func(staticRepo *mockrepo.MockStatic) {
				staticRepo.EXPECT().GetStatic(1).Return("/static/poster.jpg", nil)
			},
		},
		{
			Name:      "Успешное получение сериала",
			RequestID: 2,
			ExpectedContent: &dto.Content{
				ID:             2,
				Title:          "Бэтмен",
				OriginalTitle:  "Batman",
				Slogan:         "I'm Batman",
				Budget:         1000000,
				AgeRestriction: 18,
				Audience:       1000000,
				Rating:         9.1,
				IMDBRating:     8.2,
				Description:    "Описание фильма или сериала",
				PosterURL:      "/static/poster.jpg",
				BoxOffice:      10000,
				Marketing:      10000,
				Countries:      []string{"Россия", "США"},
				Genres:         []string{"Боевик", "Драма"},
				Actors: []dto.PersonPreview{
					{
						ID:        1,
						FirstName: "Киану",
						LastName:  "Ривз",
					},
					{
						ID:        2,
						FirstName: "Кристиан",
						LastName:  "Бэйл",
					},
				},
				Directors: []dto.PersonPreview{
					{
						ID:        3,
						FirstName: "Квентин",
						LastName:  "Тарантино",
					},
				},
				Producers: []dto.PersonPreview{
					{
						ID:        4,
						FirstName: "Кристофер",
						LastName:  "Нолан",
					},
				},
				Writers: []dto.PersonPreview{
					{
						ID:        5,
						FirstName: "Джонатан",
						LastName:  "Нолан",
					},
				},
				Operators: []dto.PersonPreview{
					{
						ID:        6,
						FirstName: "Родриго",
						LastName:  "Прито",
					},
				},
				Composers: []dto.PersonPreview{
					{
						ID:        7,
						FirstName: "Ханс",
						LastName:  "Циммер",
					},
				},
				Editors: []dto.PersonPreview{
					{
						ID:        8,
						FirstName: "Ли",
						LastName:  "Смит",
					},
				},
				Type: "series",
				Series: dto.SeriesContent{
					YearStart: 2000,
					YearEnd:   2001,
					Seasons: []dto.Season{
						{
							ID:        1,
							YearStart: 2000,
							YearEnd:   2000,
							Episodes: []dto.Episode{
								{
									ID:            1,
									EpisodeNumber: 1,
									Title:         "Название серии",
								},
							},
						},
						{
							ID:        2,
							YearStart: 2001,
							YearEnd:   2001,
							Episodes: []dto.Episode{
								{
									ID:            2,
									EpisodeNumber: 1,
									Title:         "Название серии",
								},
								{
									ID:            3,
									EpisodeNumber: 2,
									Title:         "Название серии",
								},
							},
						},
					},
				},
			},
			ExpectedError: nil,
			SetupContentRepoMock: func(contentRepo *mockrepo.MockContent) {
				contentRepo.EXPECT().GetContent(2).Return(&entity.Content{
					ID:             2,
					Title:          "Бэтмен",
					OriginalTitle:  "Batman",
					Slogan:         "I'm Batman",
					Budget:         1000000,
					AgeRestriction: 18,
					Audience:       1000000,
					IMDBRating:     8.2,
					Description:    "Описание фильма или сериала",
					PosterStaticID: 1,
					BoxOffice:      10000,
					Marketing:      10000,
					Country:        []entity.Country{{Name: "Россия"}, {Name: "США"}},
					Genres:         []entity.Genre{{Name: "Боевик"}, {Name: "Драма"}},
					Actors:         []entity.Person{{ID: 1, FirstName: "Киану", LastName: "Ривз"}, {ID: 2, FirstName: "Кристиан", LastName: "Бэйл"}},
					Directors:      []entity.Person{{ID: 3, FirstName: "Квентин", LastName: "Тарантино"}},
					Producers:      []entity.Person{{ID: 4, FirstName: "Кристофер", LastName: "Нолан"}},
					Writers:        []entity.Person{{ID: 5, FirstName: "Джонатан", LastName: "Нолан"}},
					Operators:      []entity.Person{{ID: 6, FirstName: "Родриго", LastName: "Прито"}},
					Composers:      []entity.Person{{ID: 7, FirstName: "Ханс", LastName: "Циммер"}},
					Editors:        []entity.Person{{ID: 8, FirstName: "Ли", LastName: "Смит"}},
					Type:           entity.ContentTypeSeries,
					Series: &entity.Series{
						YearStart: 2000,
						YearEnd:   2001,
						Seasons: []entity.Season{
							{
								ID:        1,
								YearStart: 2000,
								YearEnd:   2000,
								Episodes: []entity.Episode{
									{
										ID:            1,
										EpisodeNumber: 1,
										Title:         "Название серии",
									},
								},
							},
							{
								ID:        2,
								YearStart: 2001,
								YearEnd:   2001,
								Episodes: []entity.Episode{
									{
										ID:            2,
										EpisodeNumber: 1,
										Title:         "Название серии",
									},
									{
										ID:            3,
										EpisodeNumber: 2,
										Title:         "Название серии",
									},
								},
							},
						},
					},
				}, nil)
			},
			SetupReviewRepoMock: func(reviewRepo *mockrepo.MockReview) {
				reviewRepo.EXPECT().GetContentRating(2).Return(9.1, nil)
			},
			SetupStaticRepoMock: func(staticRepo *mockrepo.MockStatic) {
				staticRepo.EXPECT().GetStatic(1).Return("/static/poster.jpg", nil)
			},
		},
		{
			Name:            "Ошибка при получении контента",
			RequestID:       3,
			ExpectedContent: nil,
			ExpectedError:   entity.ErrNotFound,
			SetupContentRepoMock: func(contentRepo *mockrepo.MockContent) {
				contentRepo.EXPECT().GetContent(3).Return(nil, entity.ErrNotFound)
			},
			SetupReviewRepoMock: func(reviewRepo *mockrepo.MockReview) {},
			SetupStaticRepoMock: func(staticRepo *mockrepo.MockStatic) {},
		},
		{
			Name:            "Ошибка при получении статики",
			RequestID:       4,
			ExpectedContent: nil,
			ExpectedError:   entity.ErrNotFound,
			SetupContentRepoMock: func(contentRepo *mockrepo.MockContent) {
				contentRepo.EXPECT().GetContent(4).Return(&entity.Content{
					ID:             4,
					Title:          "Бэтмен",
					OriginalTitle:  "Batman",
					Slogan:         "I'm Batman",
					Budget:         1000000,
					AgeRestriction: 18,
					Audience:       1000000,
					IMDBRating:     8.2,
					Description:    "Описание фильма или сериала",
					PosterStaticID: 1,
					BoxOffice:      10000,
					Marketing:      10000,
					Country:        []entity.Country{{Name: "Россия"}, {Name: "США"}},
					Genres:         []entity.Genre{{Name: "Боевик"}, {Name: "Драма"}},
					Actors:         []entity.Person{{ID: 1, FirstName: "Киану", LastName: "Ривз"}, {ID: 2, FirstName: "Кристиан", LastName: "Бэйл"}},
					Directors:      []entity.Person{{ID: 3, FirstName: "Квентин", LastName: "Тарантино"}},
					Producers:      []entity.Person{{ID: 4, FirstName: "Кристофер", LastName: "Нолан"}},
					Writers:        []entity.Person{{ID: 5, FirstName: "Джонатан", LastName: "Нолан"}},
					Operators:      []entity.Person{{ID: 6, FirstName: "Родриго", LastName: "Прито"}},
					Composers:      []entity.Person{{ID: 7, FirstName: "Ханс", LastName: "Циммер"}},
					Editors:        []entity.Person{{ID: 8, FirstName: "Ли", LastName: "Смит"}},
					Type:           entity.ContentTypeMovie,
					Movie: &entity.Movie{
						Premiere: time.Time{},
						Release:  time.Time{},
						Duration: 100,
					},
				}, nil)
			},
			SetupReviewRepoMock: func(reviewRepo *mockrepo.MockReview) {},
			SetupStaticRepoMock: func(staticRepo *mockrepo.MockStatic) {
				staticRepo.EXPECT().GetStatic(1).Return("", entity.ErrNotFound)
			},
		},
		{
			Name:            "Ошибка при получении рейтинга",
			RequestID:       5,
			ExpectedContent: nil,
			ExpectedError:   entity.ErrNotFound,
			SetupContentRepoMock: func(contentRepo *mockrepo.MockContent) {
				contentRepo.EXPECT().GetContent(5).Return(&entity.Content{
					ID:             5,
					Title:          "Бэтмен",
					OriginalTitle:  "Batman",
					Slogan:         "I'm Batman",
					Budget:         1000000,
					AgeRestriction: 18,
					Audience:       1000000,
					IMDBRating:     8.2,
					Description:    "Описание фильма или сериала",
					PosterStaticID: 1,
					BoxOffice:      10000,
					Marketing:      10000,
					Country:        []entity.Country{{Name: "Россия"}, {Name: "США"}},
					Genres:         []entity.Genre{{Name: "Боевик"}, {Name: "Драма"}},
					Actors:         []entity.Person{{ID: 1, FirstName: "Киану", LastName: "Ривз"}, {ID: 2, FirstName: "Кристиан", LastName: "Бэйл"}},
					Directors:      []entity.Person{{ID: 3, FirstName: "Квентин", LastName: "Тарантино"}},
					Producers:      []entity.Person{{ID: 4, FirstName: "Кристофер", LastName: "Нолан"}},
					Writers:        []entity.Person{{ID: 5, FirstName: "Джонатан", LastName: "Нолан"}},
					Operators:      []entity.Person{{ID: 6, FirstName: "Родриго", LastName: "Прито"}},
					Composers:      []entity.Person{{ID: 7, FirstName: "Ханс", LastName: "Циммер"}},
					Editors:        []entity.Person{{ID: 8, FirstName: "Ли", LastName: "Смит"}},
					Type:           entity.ContentTypeMovie,
					Movie: &entity.Movie{
						Premiere: time.Time{},
						Release:  time.Time{},
						Duration: 100,
					},
				}, nil)
			},
			SetupReviewRepoMock: func(reviewRepo *mockrepo.MockReview) {
				reviewRepo.EXPECT().GetContentRating(5).Return(0.0, entity.ErrNotFound)
			},
			SetupStaticRepoMock: func(staticRepo *mockrepo.MockStatic) {
				staticRepo.EXPECT().GetStatic(1).Return("/static/poster.jpg", nil)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockReviewRepo := mockrepo.NewMockReview(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			contentService := NewContentService(mockContentRepo, mockReviewRepo, mockStaticRepo)
			tc.SetupContentRepoMock(mockContentRepo)
			tc.SetupReviewRepoMock(mockReviewRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			output, err := contentService.GetContentByID(tc.RequestID)
			require.Equal(t, tc.ExpectedError, err)
			require.Equal(t, tc.ExpectedContent, output)
		})
	}
}

func TestContentService_GetPersonByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		RequestID            int
		ExpectedPerson       *dto.Person
		ExpectedError        error
		SetupContentRepoMock func(contentRepo *mockrepo.MockContent)
		SetupStaticRepoMock  func(staticRepo *mockrepo.MockStatic)
	}{
		{
			Name:      "Успешное получение персоны",
			RequestID: 1,
			ExpectedPerson: &dto.Person{
				ID:          1,
				FirstName:   "Киану",
				LastName:    "Ривз",
				BirthDate:   time.Time{},
				DeathDate:   time.Time{},
				StartCareer: time.Time{},
				EndCareer:   time.Time{},
				Sex:         "M",
				PhotoURL:    "/static/photo.jpg",
				BirthPlace:  "Бейрут",
				Height:      185,
				Spouse:      "Алисия Викандер",
				Children:    "Homer, Bart, Lisa, Maggie",
				Roles: []dto.PreviewContentCard{
					{
						ID:            1,
						Title:         "Бэтмен",
						OriginalTitle: "Batman",
						Poster:        "/static/poster1.jpg",
					},
					{
						ID:            2,
						Title:         "Джон Уик",
						OriginalTitle: "John Wick",
						Poster:        "/static/poster2.jpg",
					},
				},
			},
			ExpectedError: nil,
			SetupContentRepoMock: func(contentRepo *mockrepo.MockContent) {
				contentRepo.EXPECT().GetPerson(1).Return(&entity.Person{
					ID:            1,
					FirstName:     "Киану",
					LastName:      "Ривз",
					BirthDate:     time.Time{},
					DeathDate:     time.Time{},
					StartCareer:   time.Time{},
					EndCareer:     time.Time{},
					Sex:           "M",
					PhotoStaticID: 1,
					BirthPlace:    "Бейрут",
					Height:        185,
					Spouse:        "Алисия Викандер",
					Children:      "Homer, Bart, Lisa, Maggie",
				}, nil)
				contentRepo.EXPECT().GetPersonRoles(1).Return([]entity.Content{
					{
						ID:             1,
						Title:          "Бэтмен",
						OriginalTitle:  "Batman",
						PosterStaticID: 2,
					},
					{
						ID:             2,
						Title:          "Джон Уик",
						OriginalTitle:  "John Wick",
						PosterStaticID: 3,
					},
				}, nil)
			},
			SetupStaticRepoMock: func(staticRepo *mockrepo.MockStatic) {
				staticRepo.EXPECT().GetStatic(1).Return("/static/photo.jpg", nil)
				staticRepo.EXPECT().GetStatic(2).Return("/static/poster1.jpg", nil)
				staticRepo.EXPECT().GetStatic(3).Return("/static/poster2.jpg", nil)
			},
		},
		{
			Name:           "Ошибка при получении персоны",
			RequestID:      2,
			ExpectedPerson: nil,
			ExpectedError:  entity.ErrNotFound,
			SetupContentRepoMock: func(contentRepo *mockrepo.MockContent) {
				contentRepo.EXPECT().GetPerson(2).Return(nil, entity.ErrNotFound)
			},
			SetupStaticRepoMock: func(staticRepo *mockrepo.MockStatic) {},
		},
		{
			Name:           "Ошибка при получении ролей",
			RequestID:      3,
			ExpectedPerson: nil,
			ExpectedError:  entity.ErrNotFound,
			SetupContentRepoMock: func(contentRepo *mockrepo.MockContent) {
				contentRepo.EXPECT().GetPerson(3).Return(&entity.Person{
					ID:        3,
					FirstName: "Киану",
					LastName:  "Ривз",
				}, nil)
				contentRepo.EXPECT().GetPersonRoles(3).Return(nil, entity.ErrNotFound)
			},
			SetupStaticRepoMock: func(staticRepo *mockrepo.MockStatic) {},
		},
		{
			Name:           "Ошибка при получении статики постера",
			RequestID:      4,
			ExpectedPerson: nil,
			ExpectedError:  entity.ErrNotFound,
			SetupContentRepoMock: func(contentRepo *mockrepo.MockContent) {
				contentRepo.EXPECT().GetPerson(4).Return(&entity.Person{
					ID:            4,
					FirstName:     "Киану",
					LastName:      "Ривз",
					PhotoStaticID: 2,
				}, nil)
				contentRepo.EXPECT().GetPersonRoles(4).Return([]entity.Content{
					{
						ID:             1,
						Title:          "Бэтмен",
						OriginalTitle:  "Batman",
						PosterStaticID: 1,
					},
				}, nil)
			},
			SetupStaticRepoMock: func(staticRepo *mockrepo.MockStatic) {
				staticRepo.EXPECT().GetStatic(1).Return("", entity.ErrNotFound)
			},
		},
		{
			Name:           "Ошибка при получении статики фото",
			RequestID:      5,
			ExpectedPerson: nil,
			ExpectedError:  entity.ErrNotFound,
			SetupContentRepoMock: func(contentRepo *mockrepo.MockContent) {
				contentRepo.EXPECT().GetPerson(5).Return(&entity.Person{
					ID:            5,
					FirstName:     "Киану",
					LastName:      "Ривз",
					PhotoStaticID: 1,
				}, nil)
				contentRepo.EXPECT().GetPersonRoles(5).Return([]entity.Content{
					{
						ID:             1,
						Title:          "Бэтмен",
						OriginalTitle:  "Batman",
						PosterStaticID: 2,
					},
				}, nil)
			},
			SetupStaticRepoMock: func(staticRepo *mockrepo.MockStatic) {
				staticRepo.EXPECT().GetStatic(1).Return("", entity.ErrNotFound)
				staticRepo.EXPECT().GetStatic(2).Return("/static/poster.jpg", nil)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			contentService := NewContentService(mockContentRepo, nil, mockStaticRepo)
			tc.SetupContentRepoMock(mockContentRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			output, err := contentService.GetPersonByID(tc.RequestID)
			require.Equal(t, tc.ExpectedError, err)
			require.Equal(t, tc.ExpectedPerson, output)
		})
	}
}
