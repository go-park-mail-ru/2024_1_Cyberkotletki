package db

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/user"
	"testing"
	"time"
)

func TestInitUsersDB(t *testing.T) {
	usersDB := InitUsersDB()

	if len(usersDB) != 3 {
		t.Errorf("Expected length of UsersDB to be 3, got %d", len(usersDB))
	}

	if user, ok := usersDB[1]; !ok || user.Name != "Egor" {
		t.Errorf("Expected to find user with name Egor, got %v", user)
	}
}

func TestInitFilmsDB(t *testing.T) {
	filmsDB := InitFilmsDB()

	if len(filmsDB) != 3 {
		t.Errorf("Expected length of FilmsDB to be 2, got %d", len(filmsDB))
	}

	if film, ok := filmsDB[1]; !ok || film.Content.Title != "Игра" {
		t.Errorf("Expected to find film with title Игра, got %v", film)
	}
}

func TestAddUser(t *testing.T) {
	InitUsersDB()
	newUser := user.User{
		Id:               2,
		Name:             "Новое имя",
		Email:            "new_email@example.com",
		PasswordHash:     "new_hashed_password",
		BirthDate:        time.Now(),
		SavedFilms:       []content.Film{},
		SavedSeries:      []content.Series{},
		SavedPersons:     []person.Person{},
		Friends:          []user.User{},
		ExpectedFilms:    []content.Film{},
		RegistrationDate: time.Now(),
	}

	AddUser(newUser)

	if _, ok := GetUser(2); !ok {
		t.Errorf("TestAddUser failed, user not added")
	}
}

func TestAddFilm(t *testing.T) {
	InitFilmsDB()
	newFilm := content.Film{
		Content: content.Content{
			Id:    4,
			Title: "Новый фильм",
			// остальные поля...
		},
		Duration: 120,
	}

	AddFilm(newFilm)

	if _, ok := GetFilm(4); !ok {
		t.Errorf("TestAddFilm failed, film not added")
	}
}

func TestGetUser(t *testing.T) {
	InitUsersDB()
	user, ok := GetUser(1)

	if !ok || user.Id != 1 {
		t.Errorf("TestGetUser failed, expected user with id 1, got %v", user)
	}
}

func TestGetFilm(t *testing.T) {
	InitFilmsDB()
	film, ok := GetFilm(1)

	if !ok || film.Content.Id != 1 {
		t.Errorf("TestGetFilm failed, expected film with id 1, got %v", film)
	}
}

func TestGetFilmsByGenre(t *testing.T) {
	InitFilmsDB()

	films := GetFilmsByGenre(1)
	if len(films) == 0 {
		t.Errorf("Expected to find films with genre id 1, but found none")
	}
	for _, film := range films {
		found := false
		for _, genre := range film.Content.Genres {
			if genre.Id == 1 {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected all films to have genre id 1, but found a film without it: %v", film)
		}
	}
}

// Проверяем, что фильмы возвращаются в отсортированном порядке
func TestGetFilmsByReleaseDate(t *testing.T) {
	InitFilmsDB()

	films := GetFilmsByReleaseDate()
	if len(films) < 2 {
		t.Errorf("Expected to find at least two films for testing, but found less")
	}
	for i := 1; i < len(films); i++ {
		if films[i-1].Content.Release.After(films[i].Content.Release) {
			t.Errorf("Expected films to be sorted by release date, but found out of order films: %v, %v", films[i-1], films[i])
		}
	}
}
