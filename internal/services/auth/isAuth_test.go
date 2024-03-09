package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/session"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"testing"
	"time"
)

func Test_IsAuth(t *testing.T) {
	// добавление существующей сессии
	existingSession := "existing_session"
	session.SessionsDB.Sessions[existingSession] = 1

	tests := []struct {
		name    string
		cookie  string
		wantErr *exceptions.Exception
	}{
		{
			name:   "Unsuccessful auth - session does not exist",
			cookie: "nonexistent_session",
			wantErr: &exceptions.Exception{
				When:  time.Now(),
				What:  "Не авторизован",
				Layer: exceptions.Service,
				Type:  exceptions.Forbidden,
			},
		},
		{
			name:    "Successful auth",
			cookie:  existingSession,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := IsAuth(tt.cookie)
			if err != nil && tt.wantErr != nil {
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("IsAuth() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if err != tt.wantErr {
				t.Errorf("IsAuth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
