package httputil

import (
	"encoding/json"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"net/http"
)

func NewError(w http.ResponseWriter, status int, err exceptions.Exception) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	w.WriteHeader(status)
	if result, e := json.Marshal(er); e != nil {
		http.Error(w, string(result), status)
	} else {
		http.Error(w, fmt.Sprintf(`{"code": %d, "message": "%s"}`, status, err.What), status)
	}
}

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
