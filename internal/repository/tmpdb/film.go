package tmpdb

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"sync"
	"time"
)

type ContentDB struct {
	sync.RWMutex
	DB map[int]entity.Film
	// filmsLastId atomic.Int64
}

func NewContentRepository() repository.Content {
	filmsDB := &ContentDB{
		DB: make(map[int]entity.Film),
	}
	// nolint: dupl
	filmsDB.DB[1] = entity.Film{
		Content: entity.Content{
			ID:            1,
			Title:         "1+1",
			OriginalTitle: "Intouchables",
			Country: []entity.Country{
				{
					ID:   2,
					Name: "Франция",
				},
			},
			Genres: []entity.Genre{
				{
					ID:   3,
					Name: "комедия",
				},
			},
			Directors: []entity.Person{
				{
					ID:        1,
					FirstName: "Оливье",
					LastName:  "Накаш",
				},
			},
			Writers:          []entity.Person{},
			Producers:        []entity.Person{},
			Cinematographers: []entity.Person{},
			Slogan:           "",
			Composers:        []entity.Person{},
			Artists:          []entity.Person{},
			Editors:          []entity.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []entity.BoxOffice{},
			Audiences:        []entity.Audience{},
			Premiere:         time.Date(2011, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(2011, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.8,
			Actors: []entity.Person{
				{
					ID:        2,
					FirstName: "Франсуа",
					LastName:  "Клюзе",
				},
				{
					ID:        3,
					FirstName: "Омар",
					LastName:  "Си",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/6712baa5-53be-44a0-a700-6f1042c7fc97.jpg",
			Playback:    "",
		},
		Duration: 112,
	}
	// nolint: dupl
	filmsDB.DB[2] = entity.Film{
		Content: entity.Content{
			ID:            2,
			Title:         "Волк с Уолл-стрит",
			OriginalTitle: "The Wolf from Wall Street",
			Country: []entity.Country{
				{
					ID:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					ID:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					ID:        4,
					FirstName: "Мартин",
					LastName:  "Скорсезе",
				},
			},
			Writers:          []entity.Person{},
			Producers:        []entity.Person{},
			Cinematographers: []entity.Person{},
			Slogan:           "",
			Composers:        []entity.Person{},
			Artists:          []entity.Person{},
			Editors:          []entity.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []entity.BoxOffice{},
			Audiences:        []entity.Audience{},
			Premiere:         time.Date(2013, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(2013, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.0,
			Actors: []entity.Person{
				{
					ID:        5,
					FirstName: "Леонардо",
					LastName:  "ДиКаприо",
				},
				{
					ID:        6,
					FirstName: "Джона",
					LastName:  "Хилл",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/be085ecb-331c-444a-a693-a9a4f6aa3241.jpg",
			Playback:    "",
		},
		Duration: 180,
	}
	// nolint: dupl
	filmsDB.DB[3] = entity.Film{
		Content: entity.Content{
			ID:            3,
			Title:         "Брат",
			OriginalTitle: "",
			Country: []entity.Country{
				{
					ID:   3,
					Name: "Россия",
				},
			},
			Genres: []entity.Genre{
				{
					ID:   2,
					Name: "боевик",
				},
			},
			Directors: []entity.Person{
				{
					ID:        7,
					FirstName: "Алексей",
					LastName:  "Балабанов",
				},
			},
			Writers:          []entity.Person{},
			Producers:        []entity.Person{},
			Cinematographers: []entity.Person{},
			Slogan:           "",
			Composers:        []entity.Person{},
			Artists:          []entity.Person{},
			Editors:          []entity.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []entity.BoxOffice{},
			Audiences:        []entity.Audience{},
			Premiere:         time.Date(1997, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1997, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.0,
			Actors: []entity.Person{
				{
					ID:        8,
					FirstName: "Сергей",
					LastName:  "Бордов",
				},
				{
					ID:        9,
					FirstName: "Виктор",
					LastName:  "Сухоруков",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/8d061467-a6f5-4198-80e9-02cd31157a00.jpg",
			Playback:    "",
		},
		Duration: 100,
	}
	// nolint: dupl
	filmsDB.DB[4] = entity.Film{
		Content: entity.Content{
			ID:            4,
			Title:         "Побег из Шоушенка",
			OriginalTitle: "The Shawshank Redemption",
			Country: []entity.Country{
				{
					ID:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					ID:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					ID:        10,
					FirstName: "Фрэнк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []entity.Person{},
			Producers:        []entity.Person{},
			Cinematographers: []entity.Person{},
			Slogan:           "",
			Composers:        []entity.Person{},
			Artists:          []entity.Person{},
			Editors:          []entity.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []entity.BoxOffice{},
			Audiences:        []entity.Audience{},
			Premiere:         time.Date(1994, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1994, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []entity.Person{
				{
					ID:        11,
					FirstName: "Тим",
					LastName:  "Роббинс",
				},
				{
					ID:        12,
					FirstName: "Морган",
					LastName:  "Фриман",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/4d685b8f-3762-41b1-af4a-3776f9a11057.jpg",
			Playback:    "",
		},
		Duration: 142,
	}
	// nolint: dupl
	filmsDB.DB[5] = entity.Film{
		Content: entity.Content{
			ID:            5,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []entity.Country{
				{
					ID:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					ID:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					ID:        10,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []entity.Person{},
			Producers:        []entity.Person{},
			Cinematographers: []entity.Person{},
			Slogan:           "",
			Composers:        []entity.Person{},
			Artists:          []entity.Person{},
			Editors:          []entity.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []entity.BoxOffice{},
			Audiences:        []entity.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []entity.Person{
				{
					ID:        13,
					FirstName: "Том Хэнкс",
					LastName:  "Бодров",
				},
				{
					ID:        14,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/69480722-3e05-4519-bb04-111ecfd5ef8c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	// nolint: dupl
	filmsDB.DB[6] = entity.Film{
		Content: entity.Content{
			ID:            6,
			Title:         "Форрест Гамп",
			OriginalTitle: "Forrest Gump",
			Country: []entity.Country{
				{
					ID:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					ID:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					ID:        15,
					FirstName: "Роберт",
					LastName:  "Земекис",
				},
			},
			Writers:          []entity.Person{},
			Producers:        []entity.Person{},
			Cinematographers: []entity.Person{},
			Slogan:           "",
			Composers:        []entity.Person{},
			Artists:          []entity.Person{},
			Editors:          []entity.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []entity.BoxOffice{},
			Audiences:        []entity.Audience{},
			Premiere:         time.Date(1994, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1994, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.9,
			Actors: []entity.Person{
				{
					ID:        13,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					ID:        16,
					FirstName: "Робин",
					LastName:  "Райт",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/c3bf8150-b240-4b2b-aa7d-e1963e90c558.jpg",
			Playback:    "",
		},
		Duration: 142,
	}
	// nolint: dupl
	filmsDB.DB[7] = entity.Film{
		Content: entity.Content{
			ID:            7,
			Title:         "Достучаться до небес",
			OriginalTitle: "Knockin' on Haven's Door",
			Country: []entity.Country{
				{
					ID:   4,
					Name: "Германия",
				},
			},
			Genres: []entity.Genre{
				{
					ID:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					ID:        17,
					FirstName: "Томас",
					LastName:  "Ян",
				},
			},
			Writers:          []entity.Person{},
			Producers:        []entity.Person{},
			Cinematographers: []entity.Person{},
			Slogan:           "",
			Composers:        []entity.Person{},
			Artists:          []entity.Person{},
			Editors:          []entity.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []entity.BoxOffice{},
			Audiences:        []entity.Audience{},
			Premiere:         time.Date(1997, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1997, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.6,
			Actors: []entity.Person{
				{
					ID:        18,
					FirstName: "Тиль",
					LastName:  "Швайгер",
				},
				{
					ID:        19,
					FirstName: "Ян",
					LastName:  "Йозеф Лиферс",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/423a8b50-0239-4e59-919f-9045046c0809.jpg",
			Playback:    "",
		},
		Duration: 87,
	}
	// nolint: dupl
	filmsDB.DB[8] = entity.Film{
		Content: entity.Content{
			ID:            8,
			Title:         "Дьявол носит Prada",
			OriginalTitle: "The Devil Wears Prada",
			Country: []entity.Country{
				{
					ID:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					ID:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					ID:        20,
					FirstName: "Дэвид",
					LastName:  "Фрэнкел",
				},
			},
			Writers:          []entity.Person{},
			Producers:        []entity.Person{},
			Cinematographers: []entity.Person{},
			Slogan:           "",
			Composers:        []entity.Person{},
			Artists:          []entity.Person{},
			Editors:          []entity.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []entity.BoxOffice{},
			Audiences:        []entity.Audience{},
			Premiere:         time.Date(2006, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(2006, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           7.7,
			Actors: []entity.Person{
				{
					ID:        21,
					FirstName: "Мэрил",
					LastName:  "Стрип",
				},
				{
					ID:        22,
					FirstName: "Энн",
					LastName:  "Хэтэуэй",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/d104ced8-fd3d-4a52-b7b9-11220e38ac3e.jpg",
			Playback:    "",
		},
		Duration: 109,
	}
	// nolint: dupl
	filmsDB.DB[9] = entity.Film{
		Content: entity.Content{
			ID:            9,
			Title:         "Паразиты",
			OriginalTitle: "Gisaengchung",
			Country: []entity.Country{
				{
					ID:   5,
					Name: "Корея Южная",
				},
			},
			Genres: []entity.Genre{
				{
					ID:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					ID:        22,
					FirstName: "Пон",
					LastName:  "Джун-хо",
				},
			},
			Writers:          []entity.Person{},
			Producers:        []entity.Person{},
			Cinematographers: []entity.Person{},
			Slogan:           "",
			Composers:        []entity.Person{},
			Artists:          []entity.Person{},
			Editors:          []entity.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []entity.BoxOffice{},
			Audiences:        []entity.Audience{},
			Premiere:         time.Date(2019, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(2019, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.0,
			Actors: []entity.Person{
				{
					ID:        23,
					FirstName: "Сон",
					LastName:  "Кан-хо",
				},
				{
					ID:        24,
					FirstName: "Ли",
					LastName:  "Сон-гюн",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/2c3258b0-7d6c-4256-8e76-5cf8a381dc1d.jpg",
			Playback:    "",
		},
		Duration: 131,
	}
	return filmsDB
}

func (f *ContentDB) GetFilm(id int) (*entity.Film, error) {
	f.Lock()
	defer f.Unlock()
	if filmObj, ok := f.DB[id]; ok {
		return &filmObj, nil
	}
	return nil, entity.NewClientError("фильм с таким id не найден", entity.ErrNotFound)
}

// GetFilmsByGenre возвращает фильмы определенного жанра
func (f *ContentDB) GetFilmsByGenre(genreID int) ([]entity.Film, error) {
	f.Lock()
	defer f.Unlock()

	var films []entity.Film
	for _, film := range f.DB {
		for _, genreObj := range film.Content.Genres {
			if genreObj.ID == genreID {
				films = append(films, film)
				break
			}
		}
	}
	if films == nil {
		return nil, entity.NewClientError("фильмы с таким жанром не найдены", entity.ErrNotFound)
	}
	return films, nil
}

func (f *ContentDB) GetContent(id int) (*entity.Content, error) {
	film, err := f.GetFilm(id)
	if err != nil {
		return nil, err
	}
	return &film.Content, nil
}
