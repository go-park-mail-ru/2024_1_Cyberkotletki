package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/user"
	"testing"
)

func TestRegisterData_Register(t *testing.T) {
	// Инициализация базы данных
	user.UsersDatabase.InitDB()

	tests := []struct {
		name    string
		input   RegisterData
		wantErr bool
	}{
		{
			name: "Successful registration",
			input: RegisterData{
				Email:    "newuser@example.com",
				Password: "NewUserPassword1!",
			},
			wantErr: false,
		},
		{
			name: "Unsuccessful registration - email already exists",
			input: RegisterData{
				Email:    "egor@example.com",
				Password: "EgorPassword1!",
			},
			wantErr: true,
		},
		{
			name: "Unsuccessful registration - invalid email",
			input: RegisterData{
				Email:    "invalid email",
				Password: "InvalidEmailPassword1!",
			},
			wantErr: true,
		},
		{
			name: "Unsuccessful registration - invalid password",
			input: RegisterData{
				Email:    "invalid email",
				Password: "+",
			},
			wantErr: true,
		},
		{
			name: "Unsuccessful registration - invalid password",
			input: RegisterData{
				Email:    "invalid email",
				Password: "thisisaverylongpasswordthisisaverylongpasswordthisisaverylongpasswordthisisaverylongpassword",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Register(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
