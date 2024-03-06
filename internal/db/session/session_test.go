package session

import (
	"testing"
)

func Test_SessionsDB(t *testing.T) {
	tests := []struct {
		name      string
		userId    int
		wantCheck bool
		wantDel   bool
	}{
		{
			name:      "Test with new session",
			userId:    1,
			wantCheck: true,
			wantDel:   true,
		},
		{
			name:      "Test with nonexistent session",
			userId:    2,
			wantCheck: false,
			wantDel:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			sessionId := SessionsDB.NewSession(tt.userId)
			// 2 сессия не должна существовать
			if tt.userId == 2 {
				SessionsDB.DeleteSession(sessionId)
			}

			userId, ok := SessionsDB.CheckSession(sessionId)
			if ok != tt.wantCheck {
				t.Errorf("CheckSession() = %v, want %v", ok, tt.wantCheck)
			}
			if ok && userId != tt.userId {
				t.Errorf("CheckSession() userId = %v, want %v", userId, tt.userId)
			}

			ok = SessionsDB.DeleteSession(sessionId)
			if ok != tt.wantDel {
				t.Errorf("DeleteSession() = %v, want %v", ok, tt.wantDel)
			}
		})
	}
}
