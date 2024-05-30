package entity

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name   string
		Input  string
		Output error
	}{
		{
			Name:   "Валидный пароль",
			Input:  "Amazing1!@#$%^&*()_+-=.,",
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
			Output: fmt.Errorf("пароль может состоять из латинских букв, цифр и специальных символов !@#$%%^&*()_+\\-="),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			err := ValidatePassword(tc.Input)
			require.Equal(t, tc.Output, err)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	t.Parallel()

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
		tc := tc
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
	t.Parallel()
	salt, hash, err := HashPassword("AmazingPassword1!")
	require.NoError(t, err)
	require.Equal(t, 8, len(salt))
	require.Equal(t, 32, len(hash))
}

func TestCheckPassword(t *testing.T) {
	t.Parallel()
	user := new(User)
	salt, hash, err := HashPassword("AmazingPassword1!")
	user.PasswordSalt = salt
	user.PasswordHash = hash
	require.NoError(t, err)
	require.Equal(t, user.CheckPassword("AmazingPassword1!"), true)
	require.Equal(t, user.CheckPassword("AmazingPassword!"), false)
}
