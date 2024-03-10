package user

import (
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	var userObj = NewUserEmpty()

	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{
			name:     "Password length less than 8",
			password: "1234567",
			wantErr:  exc.BadRequestErr, //  длина пароля < 8
		},
		{
			name:     "Password length more than 32",
			password: "1234567890123456789012345678901234",
			wantErr:  exc.BadRequestErr, // длина пароля > 32
		},
		{
			name:     "Valid password",
			password: "12345678",
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := userObj.ValidatePassword(tt.password)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !exc.Is(err, tt.wantErr) {
				t.Errorf("ValidatePassword() error = %v, wantErr %v", err, exc.BadRequestErr)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	var userObj = NewUserEmpty()

	tests := []struct {
		name    string
		email   string
		wantErr error
	}{
		{
			name:    "Email without @",
			email:   "testexample.com",
			wantErr: exc.BadRequestErr, //email не содержит @
		},
		{
			name:    "Email without . after @",
			email:   "test@examplecom",
			wantErr: exc.BadRequestErr, // после @ нет .
		},
		{
			name:    "Email without alphanumeric before @",
			email:   "@example.com",
			wantErr: exc.BadRequestErr, //перед @ нет буквенно-цифровых символов
		},
		{
			name:    "Valid email",
			email:   "test@example.com",
			wantErr: nil,
		},
		{
			name:    "Email with empty part after .",
			email:   "test@example.",
			wantErr: exc.BadRequestErr, // после . идет ""
		},
		{
			name:    "Email with non-alphanumeric character in host",
			email:   "test@e_ample.com",
			wantErr: exc.BadRequestErr, // в хосте есть не буквенно-цифровой символ
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := userObj.ValidateEmail(tt.email)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !exc.Is(err, tt.wantErr) {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, exc.BadRequestErr)
			}
		})
	}
}

func TestCheckPassword(t *testing.T) {
	var userObj = NewUserEmpty()
	userObj.PasswordHash = HashPassword("12345678")

	tests := []struct {
		name     string
		password string
		want     bool
	}{
		{
			name:     "Password matches hash",
			password: "12345678",
			want:     true,
		},
		{
			name:     "Password does not match hash",
			password: "87654321",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := userObj.CheckPassword(tt.password); got != tt.want {
				t.Errorf("CheckPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
