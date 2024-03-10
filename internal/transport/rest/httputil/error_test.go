package httputil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_NewError(t *testing.T) {
	tests := []struct {
		name   string
		status int
		err    error
	}{
		{
			name:   "Test with 404 status and NotFound exceptions",
			status: http.StatusNotFound,
			err:    exc.NotFoundErr,
		},
		{
			name:   "Test with 500 status and Server exceptions",
			status: http.StatusInternalServerError,
			err:    exc.UntypedErr,
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
			var e exc.Error
			if errors.As(tt.err, &e) {
				expectedBody := fmt.Sprintf(`{"code": %d, "message": "%s"}`, tt.status, e.ClientMsg.Error())
				if strings.TrimSpace(w.Body.String()) != expectedBody {
					t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
				}
			}
		})
	}
}

func Test_NewError_WithType(t *testing.T) {
	w := httptest.NewRecorder()
	err := exc.New(exc.Service, exc.Forbidden, "тестовое сообщение")
	NewError(w, http.StatusForbidden, err)

	resp := w.Result()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected status %d, got %d", http.StatusForbidden, resp.StatusCode)
	}
	expectedBody := fmt.Sprintf(`{"code": 403, "message": "тестовое сообщение"}`)

	// буфер для хранения сжатого JSON
	var buf bytes.Buffer

	// сжатие expectedBody, чтоб убрать ненужные пробелы (для корректного сравнения)
	if err := json.Compact(&buf, []byte(expectedBody)); err != nil {
		t.Fatalf("Failed to compact JSON: %v", err)
	}
	expectedBody = buf.String()
	// чтение записанного HTTP-ответа
	body, _ := io.ReadAll(resp.Body)
	if strings.TrimSpace(string(body)) != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, string(body))
	}
}
