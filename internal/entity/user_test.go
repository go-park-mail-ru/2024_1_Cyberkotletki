package entity

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestValidatePassword(t *testing.T) {

	testCases := []struct {
		Name   string
		Input  string
		Output error
	}{
		{
			Name:   "Валидный пароль",
			Input:  "AmazingPassword1!",
			Output: nil,
		},
		{
			Name:   "Короткий пароль",
			Input:  "AmPa1!",
			Output: fmt.Errorf("пароль должен содержать не менее 8 символов"),
		},
		{
			Name:   "Длинный пароль",
			Input:  "AmazingPasswordWithALotOfWordsInThisSentence1!",
			Output: fmt.Errorf("пароль должен содержать не более 32 символов"),
		},
		{
			Name:   "Недопустимые символы",
			Input:  "КириллицаВПароле1!",
			Output: fmt.Errorf("пароль должен состоять из латинских букв, цифр и специальных символов !@#$%%^&*"),
		},
		{
			Name:   "Нет заглавных букв",
			Input:  "amazing*password1!",
			Output: fmt.Errorf("пароль должен содержать как минимум одну заглавную букву"),
		},
		{
			Name:   "Нет строчных букв",
			Input:  "AMAZING*PASSWORD1!",
			Output: fmt.Errorf("пароль должен содержать как минимум одну строчную букву"),
		},
		{
			Name:   "Нет цифры",
			Input:  "AmazingPassword!",
			Output: fmt.Errorf("пароль должен содержать как минимум одну цифру"),
		},
		{
			Name:   "Нет специального символа",
			Input:  "AmazingPassword1",
			Output: fmt.Errorf("пароль должен содержать как минимум один из специальных символов !@#$%%^&*"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			err := ValidatePassword(tc.Input)
			if err != nil {
				require.EqualError(t, err, err.Error())
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {

	testCases := []struct {
		Name   string
		Input  string
		Output error
	}{
		{
			Name:   "Валидная почта",
			Input:  "email@email.com",
			Output: nil,
		},
		{
			Name:   "Длинная почта",
			Input:  strings.Repeat("email@email.com", 100),
			Output: fmt.Errorf("почта не может быть длиннее 256 символов"),
		},
		{
			Name:   "Невалидная почта",
			Input:  "email@email.com!",
			Output: fmt.Errorf("невалидная почта"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			err := ValidatePassword(tc.Input)
			if err != nil {
				require.EqualError(t, err, err.Error())
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	salt, hash, err := HashPassword("AmazingPassword1!")
	require.NoError(t, err)
	require.Equal(t, 8, len(salt))
	require.Equal(t, 32, len(hash))
}

func TestCheckPassword(t *testing.T) {
	user := NewUserEmpty()
	salt, hash, err := HashPassword("AmazingPassword1!")
	user.PasswordSalt = salt
	user.PasswordHash = hash
	require.NoError(t, err)
	require.Equal(t, user.CheckPassword("AmazingPassword1!"), true)
	require.Equal(t, user.CheckPassword("AmazingPassword!"), false)
}
