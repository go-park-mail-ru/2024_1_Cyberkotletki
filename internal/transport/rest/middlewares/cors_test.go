package middlewares

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_SetCORS(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		allowedOrigin  string
		expectedStatus int
	}{
		{
			name:           "OPTIONS request",
			method:         http.MethodOptions,
			allowedOrigin:  "http://example.com",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET request",
			method:         http.MethodGet,
			allowedOrigin:  "http://example.com",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//новый HTTP-запрос с заданным методом
			req, err := http.NewRequest(tt.method, "/", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rec := httptest.NewRecorder()

			// новый объект CORSConfig с заданным разрешенным источником
			corsConfig := &CORSConfig{AllowedOrigin: tt.allowedOrigin}
			// обработчик
			handler := corsConfig.SetCORS(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
			handler.ServeHTTP(rec, req)
			if status := rec.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			// соответствие заголовка
			if origin := rec.Header().Get("Access-Control-Allow-Origin"); origin != tt.allowedOrigin {
				t.Errorf("handler returned wrong Access-Control-Allow-Origin header: got %v want %v", origin, tt.allowedOrigin)
			}
		})
	}
}

func Test_SetAllowedOriginsFromConfig(t *testing.T) {
	tests := []struct {
		name          string
		cors          string
		expectedValue string
	}{
		{
			name:          "Test with example.com",
			cors:          "http://example.com",
			expectedValue: "http://example.com",
		},
		{
			name:          "Test with localhost",
			cors:          "http://localhost:8080",
			expectedValue: "http://localhost:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corsConfig := &CORSConfig{}
			// InitParams с заданным CORS
			params := config.InitParams{CORS: tt.cors}
			corsConfig.SetAllowedOriginsFromConfig(params)
			//  соответствие AllowedOrigin в CORSConfig
			if corsConfig.AllowedOrigin != tt.expectedValue {
				t.Errorf("Expected AllowedOrigin %s, got %s", tt.expectedValue, corsConfig.AllowedOrigin)
			}
		})
	}
}
