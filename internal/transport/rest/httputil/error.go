package httputil

import (
	"encoding/json"
	"errors"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"net/http"
)

func NewError(w http.ResponseWriter, status int, err error) {
	httpError := HTTPError{
		Code: status,
	}
	var e exc.Error
	if errors.As(err, &e) {
		httpError.Message = e.ClientMsg.Error()
	} else {
		httpError.Message = err.Error()
	}

	if result, e := json.Marshal(httpError); e != nil {
		http.Error(w, "Произошла непредвиденная ошибка", http.StatusInternalServerError)
	} else {
		http.Error(w, string(result), status)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
}

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

var BadJSON = exc.New(exc.Transport, exc.Unprocessable, "", "невалидный JSON")
