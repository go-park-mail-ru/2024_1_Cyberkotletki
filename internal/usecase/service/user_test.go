package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	mock_usecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
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
			ExpectedErr:       usecase.UserIncorrectDataError{Err: errors.New("невалидная почта")},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {},
		},
		{
			Name:              "Некорректный пароль",
			Input:             &dto.Register{Email: "email@email.com", Password: "pass"},
			ExpectedID:        -1,
			ExpectedErr:       usecase.UserIncorrectDataError{Err: errors.New("пароль должен содержать не менее 8 символов")},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {},
		},
		{
			Name:        "Пользователь уже существует",
			Input:       &dto.Register{Email: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  -1,
			ExpectedErr: usecase.ErrUserAlreadyExists,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().AddUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, repository.ErrUserAlreadyExists)
			},
		},
		{
			Name:        "Ошибка добавления пользователя",
			Input:       &dto.Register{Email: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  -1,
			ExpectedErr: usecase.UserIncorrectDataError{Err: repository.ErrUserIncorrectData},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().AddUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, repository.ErrUserIncorrectData)
			},
		},
		{
			Name:        "Внутренняя ошибка",
			Input:       &dto.Register{Email: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  -1,
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при регистрации пользователя"), errors.New("123")),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().AddUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("123"))
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
			ExpectedErr: usecase.UserIncorrectDataError{Err: errors.New("неверный пароль")},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				salt, hash, _ := entity.HashPassword("BadPassword1!")
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
			},
		},
		{
			Name:        "Пользователь не найден",
			Input:       &dto.Login{Login: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  -1,
			ExpectedErr: usecase.ErrUserNotFound,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, repository.ErrUserNotFound)
			},
		},
		{
			Name:        "Ошибка получения пользователя",
			Input:       &dto.Login{Login: "email@email.com", Password: "AmazingPassword123!"},
			ExpectedID:  -1,
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при поиске пользователя"), errors.New("error")),
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
		SetupStaticUCMock func(staticUC *mock_usecase.MockStatic)
	}{
		{
			Name:        "Успешное обновление аватара",
			UserID:      1,
			Upload:      nil,
			ExpectedErr: nil,
			SetupStaticUCMock: func(staticUC *mock_usecase.MockStatic) {
				staticUC.EXPECT().UploadAvatar(gomock.Any()).Return(1, nil)
			},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1}, nil)
				repo.EXPECT().UpdateUser(gomock.Any()).Return(nil)
			},
		},
		{
			Name:              "Пользователь не найден",
			UserID:            1,
			Upload:            nil,
			ExpectedErr:       usecase.ErrUserNotFound,
			SetupStaticUCMock: func(staticUC *mock_usecase.MockStatic) {},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(nil, repository.ErrUserNotFound)
			},
		},
		{
			Name:              "Ошибка получения пользователя",
			UserID:            1,
			Upload:            nil,
			ExpectedErr:       entity.UsecaseWrap(errors.New("ошибка при поиске пользователя"), errors.New("error")),
			SetupStaticUCMock: func(staticUC *mock_usecase.MockStatic) {},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(nil, errors.New("error"))
			},
		},
		{
			Name:        "Слишком большой файл",
			UserID:      1,
			Upload:      nil,
			ExpectedErr: usecase.UserIncorrectDataError{Err: usecase.ErrStaticTooBigFile},
			SetupStaticUCMock: func(staticUC *mock_usecase.MockStatic) {
				staticUC.EXPECT().UploadAvatar(gomock.Any()).Return(0, usecase.ErrStaticTooBigFile)
			},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1}, nil)
			},
		},
		{
			Name:        "Файл не является изображением",
			UserID:      1,
			Upload:      nil,
			ExpectedErr: usecase.UserIncorrectDataError{Err: usecase.ErrStaticNotImage},
			SetupStaticUCMock: func(staticUC *mock_usecase.MockStatic) {
				staticUC.EXPECT().UploadAvatar(gomock.Any()).Return(0, usecase.ErrStaticNotImage)
			},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1}, nil)
			},
		},
		{
			Name:        "Некорректные размеры изображения",
			UserID:      1,
			Upload:      nil,
			ExpectedErr: usecase.UserIncorrectDataError{Err: usecase.ErrStaticImageDimensions},
			SetupStaticUCMock: func(staticUC *mock_usecase.MockStatic) {
				staticUC.EXPECT().UploadAvatar(gomock.Any()).Return(0, usecase.ErrStaticImageDimensions)
			},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1}, nil)
			},
		},
		{
			Name:        "Ошибка загрузки аватара",
			UserID:      1,
			Upload:      nil,
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при загрузке аватара"), errors.New("error")),
			SetupStaticUCMock: func(staticUC *mock_usecase.MockStatic) {
				staticUC.EXPECT().UploadAvatar(gomock.Any()).Return(0, errors.New("error"))
			},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1}, nil)
			},
		},
		{
			Name:        "Ошибка обновления пользователя",
			UserID:      1,
			Upload:      nil,
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при обновлении пользователя"), errors.New("error")),
			SetupStaticUCMock: func(staticUC *mock_usecase.MockStatic) {
				staticUC.EXPECT().UploadAvatar(gomock.Any()).Return(1, nil)
			},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1}, nil)
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
			staticService := mock_usecase.NewMockStatic(ctrl)
			userService := NewUserService(mockUserRepo, staticService)
			tc.SetupStaticUCMock(staticService)
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
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1}, nil)
				repo.EXPECT().UpdateUser(gomock.Any()).Return(nil)
			},
		},
		{
			Name:        "Пользователь не найден",
			UserID:      1,
			Update:      &dto.UserUpdate{Name: "name", Email: "email@email.com"},
			ExpectedErr: usecase.ErrUserNotFound,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(nil, repository.ErrUserNotFound)
			},
		},
		{
			Name:        "Ошибка получения пользователя",
			UserID:      1,
			Update:      &dto.UserUpdate{Name: "name", Email: "email@email.com"},
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при поиске пользователя"), errors.New("error")),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(nil, errors.New("error"))
			},
		},
		{
			Name:        "Некорректный email",
			UserID:      1,
			Update:      &dto.UserUpdate{Name: "name", Email: "email"},
			ExpectedErr: usecase.UserIncorrectDataError{Err: errors.New("невалидная почта")},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1}, nil)
			},
		},
		{
			Name:        "Некорректное имя",
			UserID:      1,
			Update:      &dto.UserUpdate{Name: "VeryVeryVeryVeryVeryVeryVeryVeryLongName", Email: "email@email.com"},
			ExpectedErr: usecase.UserIncorrectDataError{Err: errors.New("имя не может быть длиннее 30 символов")},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1}, nil)
			},
		},
		{
			Name:        "Ошибка обновления пользователя",
			UserID:      1,
			Update:      &dto.UserUpdate{Name: "name", Email: "email@email.com"},
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при обновлении пользователя"), errors.New("error")),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1}, nil)
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
		SetupUserRepoMock func(userRepo *mockrepo.MockUser, reviewRepo *mockrepo.MockReview, staticUC *mock_usecase.MockStatic)
	}{
		{
			Name:            "Успешное получение пользователя",
			UserID:          1,
			ExpectedProfile: &dto.UserProfile{ID: 1, Name: "name", Email: "email", Avatar: "avatar", Rating: 0},
			ExpectedErr:     nil,
			SetupUserRepoMock: func(userRepo *mockrepo.MockUser, reviewRepo *mockrepo.MockReview, staticUC *mock_usecase.MockStatic) {
				userRepo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1, Name: "name", Email: "email", AvatarUploadID: 1}, nil)
				staticUC.EXPECT().GetStatic(gomock.Any()).Return("avatar", nil)
			},
		},
		{
			Name:            "Пользователь не найден",
			UserID:          1,
			ExpectedProfile: nil,
			ExpectedErr:     usecase.ErrUserNotFound,
			SetupUserRepoMock: func(userRepo *mockrepo.MockUser, reviewRepo *mockrepo.MockReview, staticUC *mock_usecase.MockStatic) {
				userRepo.EXPECT().GetUserByID(gomock.Any()).Return(nil, repository.ErrUserNotFound)
			},
		},
		{
			Name:            "Ошибка получения пользователя",
			UserID:          1,
			ExpectedProfile: nil,
			ExpectedErr:     entity.UsecaseWrap(errors.New("ошибка при поиске пользователя"), errors.New("error")),
			SetupUserRepoMock: func(userRepo *mockrepo.MockUser, reviewRepo *mockrepo.MockReview, staticUC *mock_usecase.MockStatic) {
				userRepo.EXPECT().GetUserByID(gomock.Any()).Return(nil, errors.New("error"))
			},
		},
		{
			Name:            "Аватар не найден",
			UserID:          1,
			ExpectedProfile: &dto.UserProfile{ID: 1, Name: "name", Email: "email", Avatar: "", Rating: 0},
			ExpectedErr:     nil,
			SetupUserRepoMock: func(userRepo *mockrepo.MockUser, reviewRepo *mockrepo.MockReview, staticUC *mock_usecase.MockStatic) {
				userRepo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1, Name: "name", Email: "email", AvatarUploadID: 1}, nil)
				staticUC.EXPECT().GetStatic(gomock.Any()).Return("", usecase.ErrStaticNotFound)
			},
		},
		{
			Name:            "Ошибка получения аватара",
			UserID:          1,
			ExpectedProfile: nil,
			ExpectedErr:     entity.UsecaseWrap(errors.New("ошибка при получении аватара"), errors.New("error")),
			SetupUserRepoMock: func(userRepo *mockrepo.MockUser, reviewRepo *mockrepo.MockReview, staticUC *mock_usecase.MockStatic) {
				userRepo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1, Name: "name", Email: "email", AvatarUploadID: 1}, nil)
				staticUC.EXPECT().GetStatic(gomock.Any()).Return("", errors.New("error"))
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
			mockStaticUC := mock_usecase.NewMockStatic(ctrl)
			userService := NewUserService(mockUserRepo, mockStaticUC)
			tc.SetupUserRepoMock(mockUserRepo, mockReviewRepo, mockStaticUC)
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
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при обновлении пользователя"), errors.New("error")),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				salt, hash, _ := entity.HashPassword("AmazingPassword123!")
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
				repo.EXPECT().UpdateUser(gomock.Any()).Return(errors.New("error"))
			},
		},
		{
			Name:        "Некорректный пароль",
			UserID:      1,
			Update:      &dto.UpdatePassword{NewPassword: "pass", OldPassword: "AmazingPassword123!"},
			ExpectedErr: usecase.UserIncorrectDataError{Err: errors.New("пароль должен содержать не менее 8 символов")},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				salt, hash, _ := entity.HashPassword("AmazingPassword123!")
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
			},
		},
		{
			Name:        "Неверный старый пароль",
			UserID:      1,
			Update:      &dto.UpdatePassword{OldPassword: "AmazingPassword123!", NewPassword: "AmazingPassword123!"},
			ExpectedErr: usecase.UserIncorrectDataError{Err: errors.New("неверный пароль")},
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				salt, hash, _ := entity.HashPassword("AmazingPassword1!")
				repo.EXPECT().GetUserByID(gomock.Any()).Return(&entity.User{ID: 1, PasswordHash: hash, PasswordSalt: salt}, nil)
			},
		},
		{
			Name:        "Ошибка получения пользователя",
			UserID:      1,
			Update:      &dto.UpdatePassword{OldPassword: "AmazingPassword123!", NewPassword: "AmazingPassword123!"},
			ExpectedErr: entity.UsecaseWrap(errors.New("ошибка при поиске пользователя"), errors.New("error")),
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(nil, errors.New("error"))
			},
		},
		{
			Name:        "Пользователь не найден",
			UserID:      1,
			Update:      &dto.UpdatePassword{OldPassword: "AmazingPassword123!", NewPassword: "AmazingPassword123!"},
			ExpectedErr: usecase.ErrUserNotFound,
			SetupUserRepoMock: func(repo *mockrepo.MockUser) {
				repo.EXPECT().GetUserByID(gomock.Any()).Return(nil, repository.ErrUserNotFound)
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
