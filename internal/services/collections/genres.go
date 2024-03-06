package collections

import exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"

type GenresData struct {
	Genres []string `json:"genres" example:"action,drama,comedian"`
}

func GetGenres() (GenresData, *exc.Exception) {

	return GenresData{
		Genres: []string{"action", "drama", "comedian"},
	}, nil
}
