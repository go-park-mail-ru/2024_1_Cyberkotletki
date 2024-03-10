package content

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"testing"
)

func TestFilmsDB_Init(t *testing.T) {
	var film FilmsDB
	film.InitDB()

	if len(film.DB) != 30 {
		t.Errorf("Expected length of FilmsDB to be 30, got %d", len(film.DB))
	}

	if film, ok := film.DB[1]; !ok || film.Content.Title != "1+1" {
		t.Errorf("Expected to find film with title 1+1, got %v", film)
	}
}

func TestFilmsDB_GetFilm(t *testing.T) {
	var film FilmsDB
	film.InitDB()

	testFilm := content.Film{Content: content.Content{
		Id: 1,
	}}
	film.DB[testFilm.Id] = testFilm

	tests := []struct {
		name    string
		filmId  int
		wantErr error
	}{
		{
			name:    "Successful get film",
			filmId:  1,
			wantErr: nil,
		},
		{
			name:    "Unsuccessful get film - film does not exist",
			filmId:  1000,
			wantErr: exc.NotFoundErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := film.GetFilm(tt.filmId)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("GetFilm() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && !exc.Is(err, tt.wantErr) {
				t.Errorf("GetFilm() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFilmsDB_GetFilmsByGenre(t *testing.T) {
	var films FilmsDB
	films.InitDB()

	tests := []struct {
		name    string
		genreId int
		wantErr error
	}{
		{
			name:    "Successful get films by genre",
			genreId: 1,
			wantErr: nil,
		},
		{
			name:    "Unsuccessful get films by genre - genre does not exist",
			genreId: 1000,
			wantErr: exc.NotFoundErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := films.GetFilmsByGenre(tt.genreId)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("GetFilmsByGenre() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && !exc.Is(err, tt.wantErr) {
				t.Errorf("GetFilmsByGenre() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
