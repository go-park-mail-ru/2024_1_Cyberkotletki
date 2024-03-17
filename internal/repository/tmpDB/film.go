package tmpDB

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	contentrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/content"
	"sync"
	"sync/atomic"
	"time"
)

type ContentDB struct {
	sync.RWMutex
	DB          map[int]entity.Film
	filmsLastId atomic.Int64
}

func NewContentRepository() contentrepo.Film {
	f := &ContentDB{
		DB: make(map[int]entity.Film),
	}
	AddExamples(f)
	return f
}

func AddExamples(f *ContentDB) {
	f.Lock()
	defer f.Unlock()

	f.DB = make(map[int]entity.Film)

	// Заполнение базы данных DB
	f.filmsLastId.CompareAndSwap(f.filmsLastId.Load(), f.filmsLastId.Load()+1)
	f.DB[1] = entity.Film{
		Content: entity.Content{
			Id:            1,
			Title:         "1+1",
			OriginalTitle: "Intouchables",
			Country: []entity.Country{
				{
					Id:   2,
					Name: "Франция",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					Id:        1,
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
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/6712baa5-53be-44a0-a700-6f1042c7fc97.jpg",
			Playback:    "",
		},
		Duration: 112,
	}
	f.DB[2] = entity.Film{
		Content: entity.Content{
			Id:            2,
			Title:         "Волк с Уолл-стрит",
			OriginalTitle: "The Wolf from Wall Street",
			Country: []entity.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					Id:        4,
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
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/be085ecb-331c-444a-a693-a9a4f6aa3241.jpg",
			Playback:    "",
		},
		Duration: 180,
	}
	f.DB[3] = entity.Film{
		Content: entity.Content{
			Id:            3,
			Title:         "Брат",
			OriginalTitle: "",
			Country: []entity.Country{
				{
					Id:   3,
					Name: "Россия",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					Id:        7,
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
					Id:        8,
					FirstName: "Сергей",
					LastName:  "Бордов",
				},
				{
					Id:        9,
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
	f.DB[4] = entity.Film{
		Content: entity.Content{
			Id:            4,
			Title:         "Побег из Шоушенка",
			OriginalTitle: "The Shawshank Redemption",
			Country: []entity.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					Id:        10,
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
					Id:        11,
					FirstName: "Тим",
					LastName:  "Роббинс",
				},
				{
					Id:        12,
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
	f.DB[5] = entity.Film{
		Content: entity.Content{
			Id:            5,
			Title:         "Зеленая миля",
			OriginalTitle: "The Green Mile",
			Country: []entity.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					Id:        10,
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
					Id:        13,
					FirstName: "Том Хэнкс",
					LastName:  "Бодров",
				},
				{
					Id:        14,
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
	f.DB[6] = entity.Film{
		Content: entity.Content{
			Id:            6,
			Title:         "Форрест Гамп",
			OriginalTitle: "Forrest Gump",
			Country: []entity.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					Id:        15,
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
					Id:        13,
					FirstName: "Том",
					LastName:  "Хэнкс",
				},
				{
					Id:        16,
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
	f.DB[7] = entity.Film{
		Content: entity.Content{
			Id:            7,
			Title:         "Достучаться до небес",
			OriginalTitle: "Knockin' on Haven's Door",
			Country: []entity.Country{
				{
					Id:   4,
					Name: "Германия",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					Id:        17,
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
					Id:        18,
					FirstName: "Тиль",
					LastName:  "Швайгер",
				},
				{
					Id:        19,
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
	f.DB[8] = entity.Film{
		Content: entity.Content{
			Id:            8,
			Title:         "Дьявол носит Prada",
			OriginalTitle: "The Devil Wears Prada",
			Country: []entity.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					Id:        20,
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
					Id:        21,
					FirstName: "Мэрил",
					LastName:  "Стрип",
				},
				{
					Id:        22,
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
	f.DB[9] = entity.Film{
		Content: entity.Content{
			Id:            9,
			Title:         "Паразиты",
			OriginalTitle: "Gisaengchung",
			Country: []entity.Country{
				{
					Id:   5,
					Name: "Корея Южная",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					Id:        22,
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
					Id:        23,
					FirstName: "Сон",
					LastName:  "Кан-хо",
				},
				{
					Id:        24,
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
	f.DB[10] = entity.Film{
		Content: entity.Content{
			Id:            10,
			Title:         "Однажды в... Голливуде",
			OriginalTitle: "Once Upon a Time in... Hollywood",
			Country: []entity.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   1,
					Name: "драма",
				},
			},
			Directors: []entity.Person{
				{
					Id:        25,
					FirstName: "Квентин",
					LastName:  "Тарантино",
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
			Rating:           7.7,
			Actors: []entity.Person{
				{
					Id:        5,
					FirstName: "Леонардо",
					LastName:  "ДиКаприо",
				},
				{
					Id:        26,
					FirstName: "Брэд",
					LastName:  "Питт",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/22151415-7f1b-4d35-909c-454d0274c975.jpg",
			Playback:    "",
		},
		Duration: 161,
	}
	f.DB[11] = entity.Film{
		Content: entity.Content{
			Id:            11,
			Title:         "Леон",
			OriginalTitle: "Leon",
			Country: []entity.Country{
				{
					Id:   2,
					Name: "Франция",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   2,
					Name: "боевик",
				},
			},
			Directors: []entity.Person{
				{
					Id:        27,
					FirstName: "Люк",
					LastName:  "Бессон",
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
			Rating:           8.7,
			Actors: []entity.Person{
				{
					Id:        28,
					FirstName: "Жан",
					LastName:  "Рено",
				},
				{
					Id:        29,
					FirstName: "Натали",
					LastName:  "Портман",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/13fc2c3d-bbfd-4334-b69c-1fc2c5de9d4d.jpg",
			Playback:    "",
		},
		Duration: 133,
	}
	f.DB[12] = entity.Film{
		Content: entity.Content{
			Id:            12,
			Title:         "Карты, деньги, два ствола",
			OriginalTitle: "Lock, Stock and Two Smocking Bar",
			Country: []entity.Country{
				{
					Id:   6,
					Name: "Великобритания",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   2,
					Name: "боевик",
				},
			},
			Directors: []entity.Person{
				{
					Id:        30,
					FirstName: "Гай",
					LastName:  "Ричи",
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
			Premiere:         time.Date(1998, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1998, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.6,
			Actors: []entity.Person{
				{
					Id:        31,
					FirstName: "Джейсон",
					LastName:  "Флеминг",
				},
				{
					Id:        32,
					FirstName: "Декстер",
					LastName:  "Флетчер",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/3ad819ce-71d1-47bf-993f-505a78ecb5cf.jpg",
			Playback:    "",
		},
		Duration: 107,
	}
	f.DB[13] = entity.Film{
		Content: entity.Content{
			Id:            13,
			Title:         "Брат 2",
			OriginalTitle: "",
			Country: []entity.Country{
				{
					Id:   3,
					Name: "Россия",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   2,
					Name: "боевик",
				},
			},
			Directors: []entity.Person{
				{
					Id:        7,
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
			Premiere:         time.Date(2000, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(2000, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.2,
			Actors: []entity.Person{
				{
					Id:        8,
					FirstName: "Сергей",
					LastName:  "Бордов",
				},
				{
					Id:        9,
					FirstName: "Виктор",
					LastName:  "Сухоруков",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/fb91434e-7ff9-4930-87ed-01ecd502eeb7.jpg",
			Playback:    "",
		},
		Duration: 127,
	}
	f.DB[14] = entity.Film{
		Content: entity.Content{
			Id:            14,
			Title:         "Шерлок Холмс",
			OriginalTitle: "Sherlock Holmes",
			Country: []entity.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   2,
					Name: "боевик",
				},
			},
			Directors: []entity.Person{
				{
					Id:        30,
					FirstName: "Гай",
					LastName:  "Ричи",
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
			Premiere:         time.Date(2009, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(2009, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.1,
			Actors: []entity.Person{
				{
					Id:        32,
					FirstName: "Роберт",
					LastName:  "Дауни",
				},
				{
					Id:        33,
					FirstName: "Джуд",
					LastName:  "Лоу",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/d6182066-4c04-46b4-97b4-fbb43e8a3915.jpg",
			Playback:    "",
		},
		Duration: 128,
	}
	f.DB[15] = entity.Film{
		Content: entity.Content{
			Id:            15,
			Title:         "Законопослушный гражданин",
			OriginalTitle: "Law Abiding Citizen",
			Country: []entity.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   2,
					Name: "боевик",
				},
			},
			Directors: []entity.Person{
				{
					Id:        34,
					FirstName: "Ф. Гэри",
					LastName:  "Грей",
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
			Premiere:         time.Date(2009, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(2009, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.0,
			Actors: []entity.Person{
				{
					Id:        35,
					FirstName: "Джейми",
					LastName:  "Фокс",
				},
				{
					Id:        36,
					FirstName: "Джерард",
					LastName:  "Батлер",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/b969d3f8-bc34-4624-bd70-5a89014fe832.jpg",
			Playback:    "",
		},
		Duration: 108,
	}
	f.DB[16] = entity.Film{
		Content: entity.Content{
			Id:            16,
			Title:         "Бесславные Ублюдки",
			OriginalTitle: "Inglorious Bastards",
			Country: []entity.Country{
				{
					Id:   4,
					Name: "Германия",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   2,
					Name: "боевик",
				},
			},
			Directors: []entity.Person{
				{
					Id:        25,
					FirstName: "Квентин",
					LastName:  "Тарантино",
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
			Premiere:         time.Date(2009, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(2009, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.0,
			Actors: []entity.Person{
				{
					Id:        26,
					FirstName: "Брэд",
					LastName:  "Питт",
				},
				{
					Id:        37,
					FirstName: "Кристофер",
					LastName:  "Вальц",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/be649b3e-2610-40a6-b717-8591fe5e8911.jpg",
			Playback:    "",
		},
		Duration: 153,
	}
	f.DB[17] = entity.Film{
		Content: entity.Content{
			Id:            17,
			Title:         "Такси",
			OriginalTitle: "Taxi",
			Country: []entity.Country{
				{
					Id:   2,
					Name: "Франция",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   2,
					Name: "боевик",
				},
			},
			Directors: []entity.Person{
				{
					Id:        38,
					FirstName: "Жерар",
					LastName:  "Пирес",
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
			Premiere:         time.Date(1998, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1998, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.0,
			Actors: []entity.Person{
				{
					Id:        39,
					FirstName: "Сами",
					LastName:  "Насери",
				},
				{
					Id:        40,
					FirstName: "Фредерик",
					LastName:  "Дифенталь",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/97816692-ab94-4554-9b86-b928e56a4163.jpg",
			Playback:    "",
		},
		Duration: 189,
	}
	f.DB[18] = entity.Film{
		Content: entity.Content{
			Id:            18,
			Title:         "Бэтмен: Начало",
			OriginalTitle: "Batman Begins",
			Country: []entity.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   2,
					Name: "боевик",
				},
			},
			Directors: []entity.Person{
				{
					Id:        41,
					FirstName: "Кристофер",
					LastName:  "Нолан",
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
			Premiere:         time.Date(2005, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(2005, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           7.9,
			Actors: []entity.Person{
				{
					Id:        42,
					FirstName: "Кристиан",
					LastName:  "Бэйл",
				},
				{
					Id:        43,
					FirstName: "Кэти",
					LastName:  "Холмс",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/ec6d73d4-93b7-4ef4-8a5f-c0655e948f1d.jpg",
			Playback:    "",
		},
		Duration: 140,
	}
	f.DB[19] = entity.Film{
		Content: entity.Content{
			Id:            19,
			Title:         "Переводчик",
			OriginalTitle: "The Covenant",
			Country: []entity.Country{
				{
					Id:   6,
					Name: "Великобритания",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   2,
					Name: "боевик",
				},
			},
			Directors: []entity.Person{
				{
					Id:        30,
					FirstName: "Гай",
					LastName:  "Ричи",
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
			Premiere:         time.Date(2022, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(2022, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           7.9,
			Actors: []entity.Person{
				{
					Id:        44,
					FirstName: "Джейк",
					LastName:  "Джилленхол",
				},
				{
					Id:        45,
					FirstName: "Дар",
					LastName:  "Салим",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/738942dd-2e36-4c3a-970a-fc351368acd3.jpg",
			Playback:    "",
		},
		Duration: 123,
	}
	f.DB[20] = entity.Film{
		Content: entity.Content{
			Id:            20,
			Title:         "Безумный Макс: Дорога ярости",
			OriginalTitle: "Mad Max: Fury Road",
			Country: []entity.Country{
				{
					Id:   6,
					Name: "Австралия",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   2,
					Name: "боевик",
				},
			},
			Directors: []entity.Person{
				{
					Id:        46,
					FirstName: "Джордж",
					LastName:  "Миллер",
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
			Premiere:         time.Date(2015, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(2015, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           7.8,
			Actors: []entity.Person{
				{
					Id:        47,
					FirstName: "Том",
					LastName:  "Харди",
				},
				{
					Id:        48,
					FirstName: "Шарлиз",
					LastName:  "Терон",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/9bddc2a2-c1f8-4156-aa45-482ce4a2ca11.jpg",
			Playback:    "",
		},
		Duration: 120,
	}
	f.DB[21] = entity.Film{
		Content: entity.Content{
			Id:            21,
			Title:         "Один дома",
			OriginalTitle: "Home Alone",
			Country: []entity.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   3,
					Name: "комедия",
				},
			},
			Directors: []entity.Person{
				{
					Id:        49,
					FirstName: "Крис",
					LastName:  "Коламбус",
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
			Premiere:         time.Date(1990, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1990, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.3,
			Actors: []entity.Person{
				{
					Id:        50,
					FirstName: "Маколей",
					LastName:  "Калкин",
				},
				{
					Id:        51,
					FirstName: "Джо",
					LastName:  "Пеши",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/694860f4-e0fa-464b-928a-35857864ff7f.jpg",
			Playback:    "",
		},
		Duration: 103,
	}
	f.DB[22] = entity.Film{
		Content: entity.Content{
			Id:            22,
			Title:         "Иван Васильевич меняет профессию",
			OriginalTitle: "",
			Country: []entity.Country{
				{
					Id:   8,
					Name: "СССР",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   3,
					Name: "комедия",
				},
			},
			Directors: []entity.Person{
				{
					Id:        52,
					FirstName: "Леонид",
					LastName:  "Гайдай",
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
			Premiere:         time.Date(1973, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1973, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.8,
			Actors: []entity.Person{
				{
					Id:        53,
					FirstName: "Александр",
					LastName:  "Демьяненко",
				},
				{
					Id:        54,
					FirstName: "Юрий",
					LastName:  "Яковлев",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/f92f275b-b53f-4d97-9171-ce49cbb9e077.jpg",
			Playback:    "",
		},
		Duration: 88,
	}
	f.DB[23] = entity.Film{
		Content: entity.Content{
			Id:            23,
			Title:         "Операция  «Ы» и другие приклюяения Шурика",
			OriginalTitle: "",
			Country: []entity.Country{
				{
					Id:   7,
					Name: "СССР",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   3,
					Name: "комедия",
				},
			},
			Directors: []entity.Person{
				{
					Id:        52,
					FirstName: "Леонид",
					LastName:  "Гайдай",
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
			Premiere:         time.Date(1965, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1965, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.7,
			Actors: []entity.Person{
				{
					Id:        53,
					FirstName: "Александр",
					LastName:  "Демьяненко",
				},
				{
					Id:        56,
					FirstName: "Наталья",
					LastName:  "Селезнёва",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/4b1d2c78-5a85-4045-989b-e1c1c283c2e5.jpg",
			Playback:    "",
		},
		Duration: 95,
	}
	f.DB[24] = entity.Film{
		Content: entity.Content{
			Id:            24,
			Title:         "Холоп",
			OriginalTitle: "",
			Country: []entity.Country{
				{
					Id:   3,
					Name: "Россия",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   3,
					Name: "комедия",
				},
			},
			Directors: []entity.Person{
				{
					Id:        57,
					FirstName: "Клим",
					LastName:  "Шипенко",
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
			Rating:           7.1,
			Actors: []entity.Person{
				{
					Id:        58,
					FirstName: "Милош",
					LastName:  "Бикович",
				},
				{
					Id:        59,
					FirstName: "Александра",
					LastName:  "Бортич",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/24d80d45-f3db-46d2-abc3-da115d979a5e.jpg",
			Playback:    "",
		},
		Duration: 109,
	}
	f.DB[25] = entity.Film{
		Content: entity.Content{
			Id:            25,
			Title:         "Джентельмены удачи",
			OriginalTitle: "",
			Country: []entity.Country{
				{
					Id:   7,
					Name: "СССР",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   3,
					Name: "комедия",
				},
			},
			Directors: []entity.Person{
				{
					Id:        60,
					FirstName: "Александр",
					LastName:  "Серый",
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
			Premiere:         time.Date(1971, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1971, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.5,
			Actors: []entity.Person{
				{
					Id:        61,
					FirstName: "Евгений",
					LastName:  "Леонов",
				},
				{
					Id:        62,
					FirstName: "Георгий",
					LastName:  "Вицин",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/de2cd29e-a7c6-401d-aeb5-1bc51ebac02f.jpg",
			Playback:    "",
		},
		Duration: 84,
	}
	f.DB[26] = entity.Film{
		Content: entity.Content{
			Id:            26,
			Title:         "Бриллиантовая рука",
			OriginalTitle: "",
			Country: []entity.Country{
				{
					Id:   7,
					Name: "СССР",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   3,
					Name: "комедия",
				},
			},
			Directors: []entity.Person{
				{
					Id:        52,
					FirstName: "Леонид",
					LastName:  "Гайдай",
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
			Premiere:         time.Date(1968, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1968, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.5,
			Actors: []entity.Person{
				{
					Id:        63,
					FirstName: "Юрий",
					LastName:  "Никулин",
				},
				{
					Id:        64,
					FirstName: "Андрей",
					LastName:  "Миронов",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/e9137c38-42f7-4f10-a59a-40175b8ad909.jpg",
			Playback:    "",
		},
		Duration: 94,
	}
	f.DB[27] = entity.Film{
		Content: entity.Content{
			Id:            27,
			Title:         "Кавказская пленница, или Новые приключения Шурика",
			OriginalTitle: "",
			Country: []entity.Country{
				{
					Id:   7,
					Name: "СССР",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   3,
					Name: "комедия",
				},
			},
			Directors: []entity.Person{
				{
					Id:        52,
					FirstName: "Леонид",
					LastName:  "Гайдай",
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
			Premiere:         time.Date(1966, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1966, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.5,
			Actors: []entity.Person{
				{
					Id:        53,
					FirstName: "Александр",
					LastName:  "Демьяненко",
				},
				{
					Id:        65,
					FirstName: "Наталья",
					LastName:  "Варлей",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/d8f4a209-e92e-4536-9464-b8f33bcee617.jpg",
			Playback:    "",
		},
		Duration: 82,
	}
	f.DB[28] = entity.Film{
		Content: entity.Content{
			Id:            28,
			Title:         "Один дома 2: Затерянный в Нью-Йорке",
			OriginalTitle: "Home Alone 2: Lost in New York",
			Country: []entity.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   3,
					Name: "комедия",
				},
			},
			Directors: []entity.Person{
				{
					Id:        49,
					FirstName: "Крис",
					LastName:  "Коламбус",
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
			Premiere:         time.Date(1992, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(1992, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           8.0,
			Actors: []entity.Person{
				{
					Id:        50,
					FirstName: "Маколей",
					LastName:  "Калкин",
				},
				{
					Id:        51,
					FirstName: "Джо",
					LastName:  "Пеши",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/0db9ad6e-a213-4fa9-bf2e-a4ea1185963d.jpg",
			Playback:    "",
		},
		Duration: 119,
	}
	f.DB[29] = entity.Film{
		Content: entity.Content{
			Id:            29,
			Title:         "Круэлла",
			OriginalTitle: "Cruella",
			Country: []entity.Country{
				{
					Id:   1,
					Name: "США",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   3,
					Name: "комедия",
				},
			},
			Directors: []entity.Person{
				{
					Id:        66,
					FirstName: "Крэйг",
					LastName:  "Гиллеспи",
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
			Premiere:         time.Date(2021, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(2021, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           7.6,
			Actors: []entity.Person{
				{
					Id:        67,
					FirstName: "Эмма",
					LastName:  "Стоун",
				},
				{
					Id:        68,
					FirstName: "Эммма",
					LastName:  "Томпсон",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/ba9e666c-e9f0-45d7-90e5-4cc96e20626b.jpg",
			Playback:    "",
		},
		Duration: 134,
	}
	f.DB[30] = entity.Film{
		Content: entity.Content{
			Id:            30,
			Title:         "Батя",
			OriginalTitle: "",
			Country: []entity.Country{
				{
					Id:   3,
					Name: "Россия",
				},
			},
			Genres: []entity.Genre{
				{
					Id:   3,
					Name: "комедия",
				},
			},
			Directors: []entity.Person{
				{
					Id:        68,
					FirstName: "Дмитрий",
					LastName:  "Ефимович",
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
			Premiere:         time.Date(2020, time.January, 0, 0, 0, 0, 0, time.UTC),
			Release:          time.Date(2020, time.January, 0, 0, 0, 0, 0, time.UTC),
			AgeRestriction:   0,
			Rating:           7.7,
			Actors: []entity.Person{
				{
					Id:        69,
					FirstName: "Владимир",
					LastName:  "Вдовиченков",
				},
				{
					Id:        70,
					FirstName: "Андрей",
					LastName:  "Андреев",
				},
			},
			Dubbing:     []entity.Person{},
			Awards:      []entity.Award{},
			Description: "",
			Poster:      "/posters/cdcdaad3-ca72-4fee-971a-bf9d7866d0d1.jpg",
			Playback:    "",
		},
		Duration: 76,
	}
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
func (f *ContentDB) GetFilmsByGenre(genreId int) ([]entity.Film, error) {
	f.Lock()
	defer f.Unlock()

	var films []entity.Film
	for _, film := range f.DB {
		for _, genreObj := range film.Content.Genres {
			if genreObj.Id == genreId {
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
