package httputil

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_NewError(t *testing.T) {
	tests := []struct {
		name   string
		status int
		err    exceptions.Exception
	}{
		{
			name:   "Test with 404 status and NotFound exceptions",
			status: http.StatusNotFound,
			err: exceptions.Exception{
				When:  time.Now(),
				What:  "Resource not found",
				Layer: exceptions.Service,
				Type:  exceptions.NotFound,
			},
		},
		{
			name:   "Test with 500 status and Server exceptions",
			status: http.StatusInternalServerError,
			err: exceptions.Exception{
				When:  time.Now(),
				What:  "Internal server error",
				Layer: exceptions.Server,
				Type:  exceptions.Untyped,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// соединение
			w := httptest.NewRecorder()
			NewError(w, tt.status, tt.err)
			if w.Code != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, w.Code)
			}
			expectedBody := fmt.Sprintf(`{"code": %d, "message": "%s"}`, tt.status, tt.err.What)
			if strings.TrimSpace(w.Body.String()) != expectedBody {
				t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
			}

		})
	}
}
