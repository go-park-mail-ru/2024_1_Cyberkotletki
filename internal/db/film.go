package db

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/small_models"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type FilmsDB struct {
	DB          map[int]content.Film
	dbMutex     sync.RWMutex
	filmsLastId atomic.Int64
}

// InitFilmsDB Инициализирует небольшую таблицу фильмов
func (f *FilmsDB) InitFilmsDB() {
	f.dbMutex.Lock()
	defer f.dbMutex.Unlock()

	f.DB = make(map[int]content.Film)

	// Заполнение базы данных DB
	f.DB[1] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "1+1",
			OriginalTitle: "Intouchables",
			Country: []small_models.Country{
				{
					Id:   2,
					Name: "Франция",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Оливье",
					LastName:  "Накаш",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(2011, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           8.8,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Франсуа",
					LastName:  "Клюзе",
				},
				{
					Id:        3,
					FirstName: "Омар",
					LastName:  "Си",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/50d7618a-0366-43dc-97f8-627ad0335409.jpg",
			Playback:    "",
		},
		Duration: 112,
	}
	f.DB[2] = content.Film{
		Content: content.Content{
			Id:            2,
			Title:         "Волк с Уолл-стрит",
			OriginalTitle: "The Wolf jf Wall Stret",
			Country: []small_models.Country{
				{
					Id:   2,
					Name: "Франция",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        4,
					FirstName: "Мартин",
					LastName:  "Скорсезе",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(2013, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           8.0,
			Actors: []person.Person{
				{
					Id:        5,
					FirstName: "Леонардо",
					LastName:  "ДиКаприо",
				},
				{
					Id:        6,
					FirstName: "Джона",
					LastName:  "Хилл",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/50d7618a-0366-43dc-97f8-627ad0335409.jpg",
			Playback:    "",
		},
		Duration: 180,
	}
	f.DB[3] = content.Film{
		Content: content.Content{
			Id:            3,
			Title:         "Побег из Шоушенка",
			OriginalTitle: "The Shaushank Redemimtion",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        7,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1994, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        8,
					FirstName: "Тим",
					LastName:  "Роббинс",
				},
				{
					Id:        9,
					FirstName: "Морган",
					LastName:  "Фриман",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/79107706-c75a-4048-83ac-8f86f5d1050a.jpg",
			Playback:    "",
		},
		Duration: 142,
	}
	f.DB[4] = content.Film{
		Content: content.Content{
			Id:            4,
			Title:         "Форест Гамп",
			OriginalTitle: "Forrest Gump",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        10,
					FirstName: "Роберт",
					LastName:  "Земекис",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1994, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        11,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        12,
					FirstName: "Робин",
					LastName:  "Райт",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 142,
	}
	f.DB[5] = content.Film{
		Content: content.Content{
			Id:            5,
			Title:         "Брат",
			OriginalTitle: "",
			Country: []small_models.Country{
				{
					Id:   3,
					Name: "Россия",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        13,
					FirstName: "Алексей",
					LastName:  "Балабанов",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1997, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           8.3,
			Actors: []person.Person{
				{
					Id:        14,
					FirstName: "Сергей",
					LastName:  "Бодров",
				},
				{
					Id:        15,
					FirstName: "Виктор",
					LastName:  "Сухоруков",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 100,
	}
	f.DB[6] = content.Film{
		Content: content.Content{
			Id:            6,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        7,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        11,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        12,
					FirstName: "Девид",
					LastName:  " Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[7] = content.Film{
		Content: content.Content{
			Id:            7,
			Title:         "Король и Шут",
			OriginalTitle: "",
			Country: []small_models.Country{
				{
					Id:   3,
					Name: "Россия",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        15,
					FirstName: "Дэвид",
					LastName:  "Финчер",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(2023, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           8.2,
			Actors: []person.Person{
				{
					Id:        16,
					FirstName: "Эдвард",
					LastName:  "Нортон",
				},
				{
					Id:        17,
					FirstName: "Бредд",
					LastName:  "Питт",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 50,
	}
	f.DB[8] = content.Film{
		Content: content.Content{
			Id:            8,
			Title:         "Властелин Колец",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[9] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[10] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[11] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[12] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[13] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[14] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[15] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[16] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[17] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[18] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[19] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[20] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[21] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[22] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[23] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[24] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[25] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[26] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[27] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[28] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[29] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[30] = content.Film{
		Content: content.Content{
			Id:            1,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []person.Person{
				{
					Id:        1,
					FirstName: "Френк",
					LastName:  "Дарабонт",
				},
			},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           0,
			Marketing:        0,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Date(1999, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           9.1,
			Actors: []person.Person{
				{
					Id:        2,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        3,
					FirstName: "Дэвид",
					LastName:  "Морс",
				},
			},
			Dubbing:     []person.Person{},
			Awards:      []small_models.Award{},
			Description: "",
			Poster:      "/assets/examples/static/posters/1cb669a7-ad7d-4757-b75e-068d4d54240c.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
}

func (f *FilmsDB) GetFilm(id int) (*content.Film, error) {
	f.dbMutex.Lock()
	defer f.dbMutex.Unlock()

	film_obj, ok := f.DB[id]
	if ok {
		return &film_obj, nil
	}
	err := errors.New("film with this id not found")
	return nil, err
}

// возвращает фильмы определенного жанра
func (f *FilmsDB) GetFilmsByGenre(genreId int) ([]content.Film, error) {
	f.dbMutex.Lock()
	defer f.dbMutex.Unlock()

	var films []content.Film
	for _, film := range f.DB {
		for _, genre := range film.Content.Genres {
			if genre.Id == genreId {
				films = append(films, film)
				break
			}
		}
	}
	if films == nil {
		err := errors.New("films not found")
		return nil, err
	}
	return films, nil
}

// возвращает фильмы, отсортированные по дате релиза
func (f *FilmsDB) GetFilmsByReleaseDate() ([]content.Film, error) {
	f.dbMutex.Lock()
	defer f.dbMutex.Unlock()

	films := make([]content.Film, 0, len(f.DB))
	for _, film := range f.DB {
		films = append(films, film)
	}
	sort.Slice(films, func(i, j int) bool {
		return films[i].Content.Release.Before(films[j].Content.Release)
	})

	if films == nil {
		err := errors.New("films not found")
		return nil, err
	}
	return films, nil
}
