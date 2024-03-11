package auth

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/user"
	u "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/user"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"testing"
)

func TestLoginData_Login(t *testing.T) {
	// Инициализация базы данных
	user.UsersDatabase.InitDB()

	salt, hash := u.HashPassword("correct_password")

	// Добавление пользователя с нужным паролем и хешем
	user.UsersDatabase.DB[4] = u.User{
		Id:           4,
		Name:         "Test",
		Email:        "test@example.com",
		PasswordHash: hash,
		PasswordSalt: salt,
	}

	tests := []struct {
		name    string
		input   LoginData
		wantErr error
	}{
		{
			name: "Unsuccessful login - wrong password",
			input: LoginData{
				Login:    "egor@example.com",
				Password: "wrong_password",
			},
			wantErr: exc.ForbiddenErr, // Ожидаем ошибку Forbidden, так как пароль неверный
		},
		{
			name: "Unsuccessful login - user does not exist",
			input: LoginData{
				Login:    "nonexistent@example.com",
				Password: "hashed_password1",
			},
			wantErr: exc.NotFoundErr, // userDB.UsersDatabase.GetUserByEmail(loginData.Login) вызывает exc.NotFoundErr
		},
		{
			name: "Successful login",
			input: LoginData{
				Login:    "test@example.com",
				Password: "correct_password",
			},
			wantErr: nil, // Ожидаем nil, так как логин и пароль верные
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Login(tt.input)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && !exc.Is(err, tt.wantErr) {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
