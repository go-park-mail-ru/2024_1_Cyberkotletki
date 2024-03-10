package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/user"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"testing"
)

func TestRegisterData_Register(t *testing.T) {
	// Инициализация бд юзеров
	user.UsersDatabase.InitDB()

	tests := []struct {
		name    string
		input   RegisterData
		wantErr error
	}{
		{
			name: "Successful registration",
			input: RegisterData{
				Email:    "newuser@example.com",
				Password: "NewUserPassword1!",
			},
			wantErr: nil,
		},
		{
			name: "Unsuccessful registration - email already exists",
			input: RegisterData{
				Email:    "egor@example.com",
				Password: "EgorPassword1!",
			},
			wantErr: exc.AlreadyExistsErr,
		},
		{
			name: "Unsuccessful registration - invalid email",
			input: RegisterData{
				Email:    "invalid email",
				Password: "InvalidEmailPassword1!",
			},
			wantErr: exc.BadRequestErr,
		},
		{
			name: "Unsuccessful registration - invalid password",
			input: RegisterData{
				Email:    "egor@example.com",
				Password: "+",
			},
			wantErr: exc.BadRequestErr,
		},
		{
			name: "Unsuccessful registration - invalid password",
			input: RegisterData{
				Email:    "egor@example.com",
				Password: "thisisaverylongpasswordthisisaverylongpasswordthisisaverylongpasswordthisisaverylongpassword",
			},
			wantErr: exc.BadRequestErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Register(tt.input)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && !exc.Is(err, tt.wantErr) {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
