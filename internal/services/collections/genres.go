package collections

type GenresData struct {
	Genres []string `json:"genres" example:"action,drama,comedian"`
}

func GetGenres() (GenresData, error) {

	return GenresData{
		Genres: []string{"action", "drama", "comedian"},
	}, nil
}
