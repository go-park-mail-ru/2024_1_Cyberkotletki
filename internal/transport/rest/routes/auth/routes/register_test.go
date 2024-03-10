package routes

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/services/auth"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Register(t *testing.T) {
	tests := []struct {
		name       string
		body       *auth.RegisterData
		wantStatus int
		setup      func()
	}{
		{
			name:       "Registration with invalid data",
			body:       &auth.RegisterData{},
			wantStatus: http.StatusBadRequest,
			setup:      func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			body, _ := json.Marshal(tt.body)
			req, err := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Register)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}
		})
	}
}
