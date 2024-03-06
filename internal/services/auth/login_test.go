package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/user"
	"testing"
)

func TestLoginData_Login(t *testing.T) {
	// Инициализация базы данных
	user.UsersDatabase.InitDB()
	tests := []struct {
		name    string
		input   LoginData
		wantErr bool
	}{
		{
			name: "Unsuccessful login - wrong password",
			input: LoginData{
				Login:    "egor@example.com",
				Password: "wrong_password",
			},
			wantErr: true,
		},
		{
			name: "Unsuccessful login - user does not exist",
			input: LoginData{
				Login:    "nonexistent@example.com",
				Password: "hashed_password1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Login(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
