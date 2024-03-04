package collections

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/content"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"time"
)

type CompilationData struct {
	Genre              string `json:"genre" example:"action"`
	ContentIdentifiers []int  `json:"ids" example:"1,2,3"`
}

func GetCompilation(genre string) (*CompilationData, *exc.Exception) {
	var genreId int
	switch genre {
	case "drama":
		genreId = 1
	case "action":
		genreId = 2
	case "comedian":
		genreId = 3
	default:
		return nil, &exc.Exception{
			When:  time.Now(),
			What:  "Такого жанра не существует",
			Layer: exc.Service,
			Type:  exc.NotFound,
		}
	}
	if films, err := content.FilmsDatabase.GetFilmsByGenre(genreId); err != nil {
		return nil, err
	} else {
		var filmsIds []int
		for _, film := range films {
			filmsIds = append(filmsIds, film.Id)
		}
		return &CompilationData{
			Genre:              genre,
			ContentIdentifiers: filmsIds,
		}, nil
	}
}
