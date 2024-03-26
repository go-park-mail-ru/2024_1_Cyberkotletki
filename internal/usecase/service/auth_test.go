package service

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
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
		Input                *dto.Register
		ExpectedErr          error
		SetupUserRepoMock    func(repo *mockrepo.MockUser, regDTO *dto.Register)
		SetupSessionRepoMock func(repo *mockrepo.MockSession)
	}{
		{
			Name: "Успешная регистрация",
			Input: &dto.Register{
				Email:    "test@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: nil,
			SetupUserRepoMock: func(repo *mockrepo.MockUser, regDTO *dto.Register) {
				repo.EXPECT().AddUser(gomock.Cond(func(u any) bool { return u.(*entity.User).Email == regDTO.Email })).Return(&entity.User{}, nil)
			},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().NewSession(gomock.Any()).Return("session", nil).AnyTimes()
			},
		},
		{
			Name: "Пользователь уже существует",
			Input: &dto.Register{
				Email:    "existing@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: entity.ErrAlreadyExists,
			SetupUserRepoMock: func(repo *mockrepo.MockUser, regDTO *dto.Register) {
				repo.EXPECT().AddUser(gomock.Cond(func(u any) bool { return u.(*entity.User).Email == regDTO.Email })).Return(nil, entity.ErrAlreadyExists)
			},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
		{
			Name: "Пароль содержит недопустимые символы",
			Input: &dto.Register{
				Email:    "existing2@example.com",
				Password: "Албурабызлык",
			},
			ExpectedErr:          fmt.Errorf("пароль должен состоять из латинских букв, цифр и специальных символов !@#$%%^&*"),
			SetupUserRepoMock:    func(repo *mockrepo.MockUser, regDTO *dto.Register) {},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
		{
			Name: "Невалидная почта",
			Input: &dto.Register{
				Email:    "почта@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr:          fmt.Errorf("невалидная почта"),
			SetupUserRepoMock:    func(repo *mockrepo.MockUser, regDTO *dto.Register) {},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			tc.SetupUserRepoMock(mockUserRepo, tc.Input)
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
		Input                *dto.Login
		ExpectedErr          error
		SetupUserRepoMock    func(repo *mockrepo.MockUser, logDTO *dto.Login)
		SetupSessionRepoMock func(repo *mockrepo.MockSession)
	}{
		{
			Name: "Успешная авторизация",
			Input: &dto.Login{
				Login:    "test@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: nil,
			SetupUserRepoMock: func(repo *mockrepo.MockUser, logDTO *dto.Login) {
				salt, hashed, _ := entity.HashPassword("AmazingPassword1!")
				repo.EXPECT().GetUserByEmail(gomock.Cond(func(u any) bool { return u == logDTO.Login })).Return(&entity.User{Email: "test@example.com", PasswordSalt: salt, PasswordHash: hashed}, nil)
			},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().NewSession(gomock.Any()).Return("session", nil)
			},
		},
		{
			Name: "Несуществующий пользователь",
			Input: &dto.Login{
				Login:    "test123@example.com",
				Password: "AmazingPassword1!",
			},
			ExpectedErr: entity.ErrNotFound,
			SetupUserRepoMock: func(repo *mockrepo.MockUser, logDTO *dto.Login) {
				repo.EXPECT().GetUserByEmail(gomock.Cond(func(u any) bool { return u == logDTO.Login })).Return(nil, entity.ErrNotFound)
			},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
		{
			Name: "Неверный пароль",
			Input: &dto.Login{
				Login:    "test@example.com",
				Password: "BadPassword1!",
			},
			ExpectedErr: fmt.Errorf("неверный пароль"),
			SetupUserRepoMock: func(repo *mockrepo.MockUser, logDTO *dto.Login) {
				repo.EXPECT().GetUserByEmail(gomock.Cond(func(u any) bool { return u == logDTO.Login })).Return(&entity.User{}, nil)
			},
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			tc.SetupUserRepoMock(mockUserRepo, tc.Input)
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
			Input:    "session1",
			Expected: true,
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().CheckSession("session1").Return(true, nil)
			},
		},
		{
			Name:     "Несуществующая сессия",
			Input:    "session2",
			Expected: false,
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().CheckSession("session2").Return(false, nil)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			tc.SetupSessionRepoMock(mockSessionRepo)
			isAuth, _ := authService.IsAuth(tc.Input)
			require.Equal(t, tc.Expected, isAuth)
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
			Input:       "session1",
			ExpectedErr: nil,
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().DeleteSession("session1").Return(true, nil)
			},
		},
		{
			Name:        "Несуществующая сессия",
			Input:       "session2",
			ExpectedErr: fmt.Errorf("не удалось удалить сессию"),
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().DeleteSession("session2").Return(false, fmt.Errorf("не удалось удалить сессию"))
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
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
