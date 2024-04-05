package service

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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
				repo.EXPECT().AddUser(gomock.Any()).Return(&entity.User{ID: 1}, nil)
			},
		},
		{
			Name:              "Некорректный email",
			Input:             &dto.Register{Email: "email", Password: "AmazingPassword123!"},
			ExpectedID:        -1,
			ExpectedErr:       entity.NewClientError("невалидная почта", entity.ErrBadRequest),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {},
		},
		{
			Name:              "Некорректный пароль",
			Input:             &dto.Register{Email: "email@email.com", Password: "pass"},
			ExpectedID:        -1,
			ExpectedErr:       entity.NewClientError("пароль должен содержать не менее 8 символов", entity.ErrBadRequest),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {},
		},
		{
			Name:        "Ошибка добавления пользователя",
			Input:       &dto.Register{Email: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  -1,
			ExpectedErr: entity.NewClientError("ошибка добавления пользователя", entity.ErrInternal),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().AddUser(gomock.Any()).Return(nil, entity.NewClientError("ошибка добавления пользователя", entity.ErrInternal))
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
				repo.EXPECT().GetUser(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
			},
		},
		{
			Name:        "Неверный пароль",
			Input:       &dto.Login{Login: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  -1,
			ExpectedErr: entity.NewClientError("неверный пароль", entity.ErrForbidden),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				salt, hash, _ := entity.HashPassword("BadPassword1!")
				repo.EXPECT().GetUser(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
			},
		},
		{
			Name:        "Пользователь не найден",
			Input:       &dto.Login{Login: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  -1,
			ExpectedErr: entity.NewClientError("пользователь не найден", entity.ErrNotFound),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUser(gomock.Any()).Return(nil, entity.NewClientError("пользователь не найден", entity.ErrNotFound))
			},
		},
		{
			Name:        "Ошибка получения пользователя",
			Input:       &dto.Login{Login: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  -1,
			ExpectedErr: entity.NewClientError("ошибка получения пользователя", entity.ErrInternal),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUser(gomock.Any()).Return(nil, entity.NewClientError("ошибка получения пользователя", entity.ErrInternal))
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
		UploadID          int
		ExpectedErr       error
		SetupUserRepoMock func(repo *mockrepo.MockUser)
	}{
		{
			Name:        "Успешное обновление аватара",
			UserID:      1,
			UploadID:    1,
			ExpectedErr: nil,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			Name:        "Ошибка обновления аватара",
			UserID:      1,
			UploadID:    1,
			ExpectedErr: entity.NewClientError("ошибка обновления пользователя", entity.ErrInternal),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(entity.NewClientError("ошибка обновления пользователя", entity.ErrInternal))
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
			err := userService.UpdateAvatar(tc.UserID, tc.UploadID)
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
				repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			Name:        "Ошибка обновления информации",
			UserID:      1,
			Update:      &dto.UserUpdate{Name: "name", Email: "email@email.com"},
			ExpectedErr: entity.NewClientError("ошибка обновления пользователя", entity.ErrInternal),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(entity.NewClientError("ошибка обновления пользователя", entity.ErrInternal))
			},
		},
		{
			Name:              "Некорректный email",
			UserID:            1,
			Update:            &dto.UserUpdate{Name: "name", Email: "email"},
			ExpectedErr:       entity.NewClientError("невалидная почта", entity.ErrBadRequest),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {},
		},
		{
			Name:              "Некорректное имя",
			UserID:            1,
			Update:            &dto.UserUpdate{Name: "VeryVeryVeryVeryVeryVeryVeryVeryLongName", Email: "email@email.com"},
			ExpectedErr:       entity.NewClientError("имя не может быть длиннее 30 символов", entity.ErrBadRequest),
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
		SetupUserRepoMock func(userRepo *mockrepo.MockUser, reviewRepo *mockrepo.MockReview)
	}{
		{
			Name:            "Успешное получение пользователя",
			UserID:          1,
			ExpectedProfile: &dto.UserProfile{Name: "name", Email: "email", Avatar: 1},
			ExpectedErr:     nil,
			SetupUserRepoMock: func(userRepo *mockrepo.MockUser, reviewRepo *mockrepo.MockReview) {
				reviewRepo.EXPECT().GetAuthorRating(gomock.Any()).Return(0, nil)
				userRepo.EXPECT().GetUser(gomock.Any()).Return(&entity.User{Name: "name", Email: "email", AvatarUploadID: 1}, nil)
			},
		},
		{
			Name:            "Ошибка получения пользователя",
			UserID:          1,
			ExpectedProfile: nil,
			ExpectedErr:     entity.NewClientError("ошибка получения пользователя", entity.ErrInternal),
			SetupUserRepoMock: func(userRepo *mockrepo.MockUser, reviewRepo *mockrepo.MockReview) {
				userRepo.EXPECT().GetUser(gomock.Any()).Return(nil, entity.NewClientError("ошибка получения пользователя", entity.ErrInternal))
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
			userService := NewUserService(mockUserRepo, mockReviewRepo)
			tc.SetupUserRepoMock(mockUserRepo, mockReviewRepo)
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
				repo.EXPECT().GetUser(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
				repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(nil)
			},
		},
		{
			Name:        "Ошибка обновления пароля",
			UserID:      1,
			Update:      &dto.UpdatePassword{NewPassword: "AmazingPassword123!", OldPassword: "AmazingPassword123!"},
			ExpectedErr: entity.NewClientError("ошибка обновления пользователя", entity.ErrInternal),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				salt, hash, _ := entity.HashPassword("AmazingPassword123!")
				repo.EXPECT().GetUser(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
				repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(entity.NewClientError("ошибка обновления пользователя", entity.ErrInternal))
			},
		},
		{
			Name:              "Некорректный пароль",
			UserID:            1,
			Update:            &dto.UpdatePassword{NewPassword: "pass", OldPassword: "AmazingPassword123!"},
			ExpectedErr:       entity.NewClientError("пароль должен содержать не менее 8 символов", entity.ErrBadRequest),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {},
		},
		{
			Name:        "Неверный старый пароль",
			UserID:      1,
			Update:      &dto.UpdatePassword{OldPassword: "AmazingPassword123!", NewPassword: "AmazingPassword123!"},
			ExpectedErr: entity.NewClientError("неверный пароль", entity.ErrBadRequest),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				salt, hash, _ := entity.HashPassword("AmazingPassword1!")
				repo.EXPECT().GetUser(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
			},
		},
		{
			Name:        "Ошибка получения пользователя",
			UserID:      1,
			Update:      &dto.UpdatePassword{OldPassword: "AmazingPassword123!", NewPassword: "AmazingPassword123!"},
			ExpectedErr: entity.NewClientError("ошибка получения пользователя", entity.ErrInternal),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUser(gomock.Any()).Return(nil, entity.NewClientError("ошибка получения пользователя", entity.ErrInternal))
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
