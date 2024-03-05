package db

import (
	filmDB "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	"testing"
)

func TestInitF(t *testing.T) {
	var film filmDB.FilmsDB
	film.InitDB()

	if len(film.DB) != 30 {
		t.Errorf("Expected length of FilmsDB to be 30, got %d", len(film.DB))
	}

	if film, ok := film.DB[1]; !ok || film.Content.Title != "1+1" {
		t.Errorf("Expected to find film with title 1+1, got %v", film)
	}
}

func TestGetFilm(t *testing.T) {
	var filmDB filmDB.FilmsDB
	filmDB.InitDB()

	// Добавьте фильм в базу данных перед тестированием функции GetFilm
	testFilm := content.Film{Content: content.Content{
		Id: 1,
	},
	// Заполните остальные поля, если это необходимо
	}
	filmDB.DB[testFilm.Id] = testFilm

	film, err := filmDB.GetFilm(1)

	if err != nil {
		t.Errorf("TestGetFilm failed, unexpected error: %v", err)
	}

	if film == nil {
		t.Errorf("TestGetFilm failed, film is nil")
	}

	if film.Id != testFilm.Id {
		t.Errorf("TestGetFilm failed, film does not match testFilm")
	}
}

func TestGetFilmsByGenre(t *testing.T) {
	var filmDB filmDB.FilmsDB
	filmDB.InitDB()

	// Добавьте фильмы в базу данных перед тестированием функции GetFilmsByGenre
	films, err := filmDB.GetFilmsByGenre(1)

	if err != nil {
		t.Errorf("TestGetFilmsByGenre failed, unexpected error: %v", err)
	}

	if films[0].Genres[0] != films[1].Genres[0] {
		t.Errorf("TestGetFilmsByGenre failed, films does not match expected films")
	}
}
