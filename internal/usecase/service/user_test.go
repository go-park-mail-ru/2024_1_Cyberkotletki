package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"testing"
)

func TestUserService_Register(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name              string
		Input             *dto.Register
		ExpectedID        int
		ExpectedErr       error
		SetupUserRepoMock func(repo *mockrepo.MockUser)
	}{
		{
			Name:        "Успешная регистрация",
			Input:       &dto.Register{Email: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  1,
			ExpectedErr: nil,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().AddUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.User{ID: 1}, nil)
			},
		},
		{
			Name:              "Некорректный email",
			Input:             &dto.Register{Email: "email", Password: "AmazingPassword123!"},
			ExpectedID:        -1,
			ExpectedErr:       nil,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {},
		},
		{
			Name:              "Некорректный пароль",
			Input:             &dto.Register{Email: "email@email.com", Password: "pass"},
			ExpectedID:        -1,
			ExpectedErr:       nil,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {},
		},
		{
			Name:        "Ошибка добавления пользователя",
			Input:       &dto.Register{Email: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  -1,
			ExpectedErr: nil,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().AddUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			userService := NewUserService(mockUserRepo, nil)
			tc.SetupUserRepoMock(mockUserRepo)
			id, err := userService.Register(tc.Input)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedID, id)
		})
	}
}

func TestUserService_Login(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name              string
		Input             *dto.Login
		ExpectedID        int
		ExpectedErr       error
		SetupUserRepoMock func(repo *mockrepo.MockUser)
	}{
		{
			Name:        "Успешный вход",
			Input:       &dto.Login{Login: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  1,
			ExpectedErr: nil,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				salt, hash, _ := entity.HashPassword("AmazingPassword123!")
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
			},
		},
		{
			Name:        "Неверный пароль",
			Input:       &dto.Login{Login: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  -1,
			ExpectedErr: errors.New("error"),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				salt, hash, _ := entity.HashPassword("BadPassword1!")
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
			},
		},
		{
			Name:        "Пользователь не найден",
			Input:       &dto.Login{Login: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  -1,
			ExpectedErr: errors.New("error"),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, errors.New("error"))
			},
		},
		{
			Name:        "Ошибка получения пользователя",
			Input:       &dto.Login{Login: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  -1,
			ExpectedErr: errors.New("error"),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, errors.New("error"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			userService := NewUserService(mockUserRepo, nil)
			tc.SetupUserRepoMock(mockUserRepo)
			id, err := userService.Login(tc.Input)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedID, id)
		})
	}
}

func TestUserService_UpdateAvatar(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name              string
		UserID            int
		Upload            io.Reader
		ExpectedErr       error
		SetupUserRepoMock func(repo *mockrepo.MockUser)
	}{
		{
			Name:        "Успешное обновление аватара",
			UserID:      1,
			Upload:      nil,
			ExpectedErr: nil,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().UpdateUser(gomock.Any()).Return(nil)
			},
		},
		{
			Name:        "Ошибка обновления аватара",
			UserID:      1,
			Upload:      nil,
			ExpectedErr: errors.New("error"),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().UpdateUser(gomock.Any()).Return(errors.New("error"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			userService := NewUserService(mockUserRepo, nil)
			tc.SetupUserRepoMock(mockUserRepo)
			err := userService.UpdateAvatar(tc.UserID, tc.Upload)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestUserService_UpdateInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name              string
		UserID            int
		Update            *dto.UserUpdate
		ExpectedErr       error
		SetupUserRepoMock func(repo *mockrepo.MockUser)
	}{
		{
			Name:        "Успешное обновление информации",
			UserID:      1,
			Update:      &dto.UserUpdate{Name: "name", Email: "email@email.com"},
			ExpectedErr: nil,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().UpdateUser(gomock.Any()).Return(nil)
			},
		},
		{
			Name:        "Ошибка обновления информации",
			UserID:      1,
			Update:      &dto.UserUpdate{Name: "name", Email: "email@email.com"},
			ExpectedErr: errors.New("error"),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().UpdateUser(gomock.Any()).Return(errors.New("error"))
			},
		},
		{
			Name:              "Некорректный email",
			UserID:            1,
			Update:            &dto.UserUpdate{Name: "name", Email: "email"},
			ExpectedErr:       errors.New("error"),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {},
		},
		{
			Name:              "Некорректное имя",
			UserID:            1,
			Update:            &dto.UserUpdate{Name: "VeryVeryVeryVeryVeryVeryVeryVeryLongName", Email: "email@email.com"},
			ExpectedErr:       errors.New("error"),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			userService := NewUserService(mockUserRepo, nil)
			tc.SetupUserRepoMock(mockUserRepo)
			err := userService.UpdateInfo(tc.UserID, tc.Update)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name              string
		UserID            int
		ExpectedProfile   *dto.UserProfile
		ExpectedErr       error
		SetupUserRepoMock func(userRepo *mockrepo.MockUser, reviewRepo *mockrepo.MockReview, staticRepo *mockrepo.MockStatic)
	}{
		{
			Name:            "Успешное получение пользователя",
			UserID:          1,
			ExpectedProfile: &dto.UserProfile{ID: 1, Name: "name", Email: "email", Avatar: "avatar", Rating: 0},
			ExpectedErr:     nil,
			SetupUserRepoMock: func(userRepo *mockrepo.MockUser, reviewRepo *mockrepo.MockReview, staticRepo *mockrepo.MockStatic) {
				userRepo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1, Name: "name", Email: "email", AvatarUploadID: 1}, nil)
				staticRepo.EXPECT().GetStatic(gomock.Any()).Return("avatar", nil)
			},
		},
		{
			Name:            "Ошибка получения пользователя",
			UserID:          1,
			ExpectedProfile: nil,
			ExpectedErr:     errors.New("error"),
			SetupUserRepoMock: func(userRepo *mockrepo.MockUser, reviewRepo *mockrepo.MockReview, staticRepo *mockrepo.MockStatic) {
				userRepo.EXPECT().GetUserByID(gomock.Any()).Return(nil, errors.New("error"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			mockReviewRepo := mockrepo.NewMockReview(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			userService := NewUserService(mockUserRepo, nil)
			tc.SetupUserRepoMock(mockUserRepo, mockReviewRepo, mockStaticRepo)
			profile, err := userService.GetUser(tc.UserID)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedProfile, profile)
		})
	}
}

func TestUserService_UpdatePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name              string
		UserID            int
		Update            *dto.UpdatePassword
		ExpectedErr       error
		SetupUserRepoMock func(repo *mockrepo.MockUser)
	}{
		{
			Name:        "Успешное обновление пароля",
			UserID:      1,
			Update:      &dto.UpdatePassword{NewPassword: "AmazingPassword123!", OldPassword: "AmazingPassword123!"},
			ExpectedErr: nil,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				salt, hash, _ := entity.HashPassword("AmazingPassword123!")
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
				repo.EXPECT().UpdateUser(gomock.Any()).Return(nil)
			},
		},
		{
			Name:        "Ошибка обновления пароля",
			UserID:      1,
			Update:      &dto.UpdatePassword{NewPassword: "AmazingPassword123!", OldPassword: "AmazingPassword123!"},
			ExpectedErr: errors.New("error"),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				salt, hash, _ := entity.HashPassword("AmazingPassword123!")
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
				repo.EXPECT().UpdateUser(gomock.Any()).Return(errors.New("error"))
			},
		},
		{
			Name:              "Некорректный пароль",
			UserID:            1,
			Update:            &dto.UpdatePassword{NewPassword: "pass", OldPassword: "AmazingPassword123!"},
			ExpectedErr:       errors.New("error"),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {},
		},
		{
			Name:        "Неверный старый пароль",
			UserID:      1,
			Update:      &dto.UpdatePassword{OldPassword: "AmazingPassword123!", NewPassword: "AmazingPassword123!"},
			ExpectedErr: errors.New("error"),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				salt, hash, _ := entity.HashPassword("AmazingPassword1!")
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
			},
		},
		{
			Name:        "Ошибка получения пользователя",
			UserID:      1,
			Update:      &dto.UpdatePassword{OldPassword: "AmazingPassword123!", NewPassword: "AmazingPassword123!"},
			ExpectedErr: errors.New("error"),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(nil, errors.New("error"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockUserRepo := mockrepo.NewMockUser(ctrl)
			userService := NewUserService(mockUserRepo, nil)
			tc.SetupUserRepoMock(mockUserRepo)
			err := userService.UpdatePassword(tc.UserID, tc.Update)
			require.Equal(t, tc.ExpectedErr, err)
		})
	}
}
