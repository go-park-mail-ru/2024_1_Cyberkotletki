package collections

import exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"

type CompilationData struct {
	Genre              string  `json:"genre" example:"action"`
	ContentIdentifiers []int64 `json:"ids" example:"1,2,3"`
}

func GetCompilation(genre string) (CompilationData, *exc.Exception) {

	return CompilationData{
		Genre:              genre,
		ContentIdentifiers: []int64{1, 2, 3},
	}, nil
}
