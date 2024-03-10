package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Logout(t *testing.T) {
	// Добавление существующей сессии
	// не удалось протестировать второй случай
	tests := []struct {
		name       string
		cookie     *http.Cookie
		wantStatus int
		setup      func()
	}{
		{
			name:       "Logout without session",
			cookie:     nil, // нет cookie "session", поэтому пользователь не вошел в систему
			wantStatus: http.StatusOK,
			setup:      func() {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			req, err := http.NewRequest("POST", "/auth/logout", nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.cookie != nil {
				req.AddCookie(tt.cookie)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(Logout)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}
		})
	}
}
