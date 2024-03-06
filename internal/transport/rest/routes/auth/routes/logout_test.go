package routes

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/session"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Logout(t *testing.T) {
	// Добавление существующей сессии
	existingSession := "existing_session"
	session.SessionsDB.Sessions[existingSession] = 1

	tests := []struct {
		name       string
		cookie     *http.Cookie
		wantStatus int
		setup      func()
	}{
		{
			name:       "Unauthenticated request",
			cookie:     nil,
			wantStatus: http.StatusForbidden,
			setup:      func() {},
		},
		{
			name: "Authenticated request",
			cookie: &http.Cookie{
				Name:  "session",
				Value: existingSession,
			},
			wantStatus: http.StatusOK,
			setup:      func() {},
		},
		{
			name: "Unauthenticated request after session deletion",
			cookie: &http.Cookie{
				Name:  "session",
				Value: existingSession,
			},
			wantStatus: http.StatusForbidden,
			setup: func() {
				session.SessionsDB.DeleteSession(existingSession)
			},
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
