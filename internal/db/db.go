package db

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/small_models"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/user"
	"sort"
	"sync"
	"time"
)

var (
	UsersDB map[int]user.User
	FilmsDB map[int]content.Film
	dbMutex sync.Mutex
)

// Инициализирует небольшую таблицу пользователей
func InitUsersDB() map[int]user.User {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	UsersDB = make(map[int]user.User)

	// Заполнение базы данных UsersDB
	UsersDB[1] = user.User{
		Id:               1,
		Name:             "Egor",
		Email:            "egor@example.com",
		PasswordHash:     "hashed_password1",
		BirthDate:        time.Now(),
		SavedFilms:       []content.Film{},
		SavedSeries:      []content.Series{},
		SavedPersons:     []person.Person{},
		Friends:          []user.User{},
		ExpectedFilms:    []content.Film{},
		RegistrationDate: time.Now(),
	}
	UsersDB[2] = user.User{
		Id:               2,
		Name:             "Sasha",
		Email:            "sasha@example.com",
		PasswordHash:     "hashed_password2",
		BirthDate:        time.Now(),
		SavedFilms:       []content.Film{},
		SavedSeries:      []content.Series{},
		SavedPersons:     []person.Person{},
		Friends:          []user.User{},
		ExpectedFilms:    []content.Film{},
		RegistrationDate: time.Now(),
	}

	UsersDB[3] = user.User{
		Id:               3,
		Name:             "Kristina",
		Email:            "kristina@example.com",
		PasswordHash:     "hashed_password3",
		BirthDate:        time.Now(),
		SavedFilms:       []content.Film{},
		SavedSeries:      []content.Series{},
		SavedPersons:     []person.Person{},
		Friends:          []user.User{},
		ExpectedFilms:    []content.Film{},
		RegistrationDate: time.Now(),
	}
	return UsersDB
}

// Инициализирует небольшую таблицу фильмов
func InitFilmsDB() map[int]content.Film {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	FilmsDB = make(map[int]content.Film)

	// Заполнение базы данных FilmsDB
	FilmsDB[1] = content.Film{
		Content: content.Content{
			Id:    1,
			Title: "Игра",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "Россия",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "Комедия",
				},
			},
			Directors:        []person.Person{},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           123,
			Marketing:        321,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Now(),
			Release:          time.Now(),
			AgeRestriction:   21,
			Rating:           10,
			Actors:           []person.Person{},
			Dubbing:          []person.Person{},
			Awards:           []small_models.Award{},
			Description:      "blabla",
			Poster:           "",
			Playback:         "",
		},
		Duration: 120,
	}

	FilmsDB[2] = content.Film{
		Content: content.Content{
			Id:    2,
			Title: "Победа",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "Россия",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   2,
					Name: "Детектив",
				},
				{
					Id:   3,
					Name: "Боевик",
				},
			},
			Directors:        []person.Person{},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           123,
			Marketing:        321,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Now(),
			Release:          time.Now(),
			AgeRestriction:   6,
			Rating:           10,
			Actors:           []person.Person{},
			Dubbing:          []person.Person{},
			Awards:           []small_models.Award{},
			Description:      "blabla!",
			Poster:           "",
			Playback:         "",
		},
		Duration: 120,
	}
	FilmsDB[3] = content.Film{
		Content: content.Content{
			Id:    3,
			Title: "Вызов",
			Country: []small_models.Country{
				{
					Id:   1,
					Name: "Россия",
				},
			},
			Genres: []small_models.Genre{
				{
					Id:   1,
					Name: "Комедия",
				},
			},
			Directors:        []person.Person{},
			Writers:          []person.Person{},
			Producers:        []person.Person{},
			Cinematographers: []person.Person{},
			Slogan:           "",
			Composers:        []person.Person{},
			Artists:          []person.Person{},
			Editors:          []person.Person{},
			Budget:           444,
			Marketing:        555,
			BoxOffices:       []small_models.BoxOffice{},
			Audiences:        []small_models.Audience{},
			Premiere:         time.Now(),
			Release:          time.Now(),
			AgeRestriction:   0,
			Rating:           4,
			Actors:           []person.Person{},
			Dubbing:          []person.Person{},
			Awards:           []small_models.Award{},
			Description:      "",
			Poster:           "",
			Playback:         "",
		},
		Duration: 150,
	}
	return FilmsDB
}

func AddUser(user user.User) {
	dbMutex.Lock()
	UsersDB[user.Id] = user
	dbMutex.Unlock()
}

func AddFilm(film content.Film) {
	dbMutex.Lock()
	FilmsDB[film.Id] = film
	dbMutex.Unlock()
}

func GetUser(id int) (user.User, bool) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	user_obj, ok := UsersDB[id]
	return user_obj, ok
}

func GetFilm(id int) (content.Film, bool) {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	film, ok := FilmsDB[id]
	return film, ok
}

// возвращает фильмы определенного жанра
func GetFilmsByGenre(genreId int) []content.Film {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	var films []content.Film
	for _, film := range FilmsDB {
		for _, genre := range film.Content.Genres {
			if genre.Id == genreId {
				films = append(films, film)
				break
			}
		}
	}
	return films
}

// возвращает фильмы, отсортированные по дате релиза
func GetFilmsByReleaseDate() []content.Film {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	films := make([]content.Film, 0, len(FilmsDB))
	for _, film := range FilmsDB {
		films = append(films, film)
	}
	sort.Slice(films, func(i, j int) bool {
		return films[i].Content.Release.Before(films[j].Content.Release)
	})

	return films
}
