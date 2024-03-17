package service

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/DTO"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestAuth_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockrepo.NewMockUser(ctrl)
	mockSessionRepo := mockrepo.NewMockSession(ctrl)
	authService := AuthService{
		userRepo:    mockUserRepo,
		sessionRepo: mockSessionRepo,
	}

	testCases := []struct {
		Name                 string
		Input                DTO.Register
		ExpectedErr          error
		SetupUserRepoMock    func(repo *mockrepo.MockUser)
		SetupSessionRepoMock func(repo *mockrepo.MockSession)
	}{
		{
			Name: "Успешная регистрация",
			Input: DTO.Register{
				Email:    "test@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: nil,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().AddUser(gomock.Any()).Return(&entity.User{}, nil)
			},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().NewSession(gomock.Any()).Return("session").AnyTimes()
			},
		},
		{
			Name: "Пользователь уже существует",
			Input: DTO.Register{
				Email:    "existing@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: entity.ErrBadRequest,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().AddUser(gomock.Any()).Return(nil, entity.ErrBadRequest).AnyTimes()
			},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
		{
			Name: "Пароль содержит недопустимые символы",
			Input: DTO.Register{
				Email:    "existing@example.com",
				Password: "Албурабызлык",
			},
			ExpectedErr:          fmt.Errorf("пароль должен состоять из латинских букв, цифр и специальных символов !@#$%%^&*"),
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
		{
			Name: "Пароль не содержит заглавной буквы",
			Input: DTO.Register{
				Email:    "existing@example.com",
				Password: "amazingpassword1!",
			},
			ExpectedErr:          fmt.Errorf("пароль должен содержать как минимум одну заглавную букву"),
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
		{
			Name: "Пароль не содержит строчной буквы",
			Input: DTO.Register{
				Email:    "existing@example.com",
				Password: "AMAZINGPASSWORD1!",
			},
			ExpectedErr:          fmt.Errorf("пароль должен содержать как минимум одну строчную букву"),
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
		{
			Name: "Пароль не содержит спецсимвола",
			Input: DTO.Register{
				Email:    "existing@example.com",
				Password: "AmazingPassword1",
			},
			ExpectedErr:          fmt.Errorf("пароль должен содержать как минимум один из специальных символов !@#$%%^&*"),
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
		{
			Name: "Пароль не содержит цифры",
			Input: DTO.Register{
				Email:    "existing@example.com",
				Password: "AmazingPassword!",
			},
			ExpectedErr:          fmt.Errorf("пароль должен содержать как минимум одну цифру"),
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
		{
			Name: "Короткий пароль",
			Input: DTO.Register{
				Email:    "existing@example.com",
				Password: "Short1!",
			},
			ExpectedErr:          fmt.Errorf("пароль должен содержать не менее 8 символов"),
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
		{
			Name: "Слишком длинный пароль",
			Input: DTO.Register{
				Email:    "existing@example.com",
				Password: "ItIsAVeryLongPasswordAndIThinkItIsWrong1!",
			},
			ExpectedErr:          fmt.Errorf("пароль должен содержать не более 32 символов"),
			SetupUserRepoMock:    func(repo *mockrepo.MockUser) {},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.SetupUserRepoMock(mockUserRepo)
			tc.SetupSessionRepoMock(mockSessionRepo)

			_, err := authService.Register(tc.Input)
			if tc.ExpectedErr != nil {
				require.EqualError(t, err, tc.ExpectedErr.Error())
			}
		})
	}
}

func TestAuth_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockrepo.NewMockUser(ctrl)
	mockSessionRepo := mockrepo.NewMockSession(ctrl)
	authService := AuthService{
		userRepo:    mockUserRepo,
		sessionRepo: mockSessionRepo,
	}

	testCases := []struct {
		Name                 string
		Input                DTO.Login
		ExpectedErr          error
		SetupUserRepoMock    func(repo *mockrepo.MockUser)
		SetupSessionRepoMock func(repo *mockrepo.MockSession)
	}{
		{
			Name: "Успешная авторизация",
			Input: DTO.Login{
				Login:    "test@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: nil,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				salt, hashed := entity.HashPassword("AmazingPassword1!")
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&entity.User{Email: "test@example.com", PasswordSalt: salt, PasswordHash: hashed}, nil)
			},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().NewSession(gomock.Any()).Return("session")
			},
		},
		{
			Name: "Несуществующий пользователь",
			Input: DTO.Login{
				Login:    "test123@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: entity.ErrNotFound,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, entity.ErrNotFound)
			},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
		{
			Name: "Неверный пароль",
			Input: DTO.Login{
				Login:    "test@example.com",
				Password: "BadPassword1!",
			},
			ExpectedErr: fmt.Errorf("неверный пароль"),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&entity.User{}, nil)
			},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.SetupUserRepoMock(mockUserRepo)
			tc.SetupSessionRepoMock(mockSessionRepo)

			_, err := authService.Login(tc.Input)
			if tc.ExpectedErr != nil {
				require.EqualError(t, err, tc.ExpectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAuth_IsAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockrepo.NewMockUser(ctrl)
	mockSessionRepo := mockrepo.NewMockSession(ctrl)
	authService := AuthService{
		userRepo:    mockUserRepo,
		sessionRepo: mockSessionRepo,
	}

	testCases := []struct {
		Name                 string
		Input                string
		Expected             bool
		SetupSessionRepoMock func(repo *mockrepo.MockSession)
	}{
		{
			Name:     "Существующая сессия",
			Input:    "session",
			Expected: true,
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().CheckSession("session").Return(true)
			},
		},
		{
			Name:     "Несуществующая сессия",
			Input:    "session",
			Expected: false,
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().CheckSession("session").Return(false)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.SetupSessionRepoMock(mockSessionRepo)
			require.Equal(t, tc.Expected, authService.IsAuth(tc.Input))
		})
	}
}

func TestAuth_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockrepo.NewMockUser(ctrl)
	mockSessionRepo := mockrepo.NewMockSession(ctrl)
	authService := AuthService{
		userRepo:    mockUserRepo,
		sessionRepo: mockSessionRepo,
	}

	testCases := []struct {
		Name                 string
		Input                string
		ExpectedErr          error
		SetupSessionRepoMock func(repo *mockrepo.MockSession)
	}{
		{
			Name:        "Существующая сессия",
			Input:       "session",
			ExpectedErr: nil,
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().DeleteSession("session").Return(true)
			},
		},
		{
			Name:        "Несуществующая сессия",
			Input:       "session",
			ExpectedErr: fmt.Errorf("сессия недействительна"),
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().DeleteSession("session").Return(false)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.SetupSessionRepoMock(mockSessionRepo)
			err := authService.Logout(tc.Input)
			if tc.ExpectedErr != nil {
				require.EqualError(t, err, tc.ExpectedErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
