package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/session"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"testing"
)

func Test_Logout(t *testing.T) {
	existingSession := "existing_session"
	session.SessionsDB.Sessions[existingSession] = 1

	tests := []struct {
		name    string
		cookie  string
		wantErr error
	}{
		{
			name:    "Session does not exist",
			cookie:  "nonexistent_session",
			wantErr: exc.UntypedErr, // сессия не существует
		},
		{
			name:    "Successful auth",
			cookie:  existingSession,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Logout(tt.cookie)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && !exc.Is(err, tt.wantErr) {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
