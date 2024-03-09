package httputil

import (
	"encoding/json"
	"fmt"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, response any) {
	jsonData, err := json.Marshal(response)
	if err != nil {
		NewError(
			w,
			http.StatusInternalServerError,
			exc.New(
				exc.Transport,
				exc.Internal,
				"произошла непредвиденная ошибка",
				fmt.Sprintf("не удалось преобразовать в json: %v", response),
			),
		)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err = w.Write(jsonData)
	if err != nil {
		NewError(
			w,
			http.StatusInternalServerError,
			exc.New(
				exc.Transport,
				exc.Internal,
				"произошла непредвиденная ошибка",
				fmt.Sprintf("не удалось записать преобразованный json в ответ: %v", response),
			),
		)
	}
}
