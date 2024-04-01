package service

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
)

type UserService struct {
	userRepo repository.User
}

func NewUserService(userRepo repository.User) usecase.User {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) Register(regDTO *dto.Register) (int, error) {
	user := entity.NewUserEmpty()
	if err := entity.ValidateEmail(regDTO.Email); err != nil {
		return -1, err
	}
	if err := entity.ValidatePassword(regDTO.Password); err != nil {
		return -1, err
	}
	salt, hash, err := entity.HashPassword(regDTO.Password)
	if err != nil {
		return -1, err
	}
	user.Email = regDTO.Email
	user.PasswordHash = hash
	user.PasswordSalt = salt
	user, err = u.userRepo.AddUser(user)
	if err != nil {
		return -1, err
	}
	return user.ID, nil
}

func (u *UserService) Login(login *dto.Login) (int, error) {
	user, err := u.userRepo.GetUser(map[string]interface{}{"email": login.Login})
	if err != nil {
		return -1, err
	}
	if !user.CheckPassword(login.Password) {
		return -1, entity.NewClientError("неверный пароль", entity.ErrForbidden)
	}
	return user.ID, nil
}

func (u *UserService) UpdateAvatar(userID, uploadID int) error {
	err := u.userRepo.UpdateUser(
		map[string]interface{}{"id": userID},
		map[string]interface{}{"avatar_upload_id": uploadID},
	)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) UpdateInfo(userID int, update *dto.UserUpdate) error {
	if err := entity.ValidateEmail(update.Email); err != nil {
		return err
	}
	if err := entity.ValidateName(update.Name); err != nil {
		return err
	}
	err := u.userRepo.UpdateUser(
		map[string]interface{}{"id": userID},
		map[string]interface{}{"name": update.Name, "email": update.Email},
	)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) GetUser(userID int) (*dto.UserProfile, error) {
	user, err := u.userRepo.GetUser(map[string]interface{}{"id": userID})
	if err != nil {
		return nil, err
	}
	return &dto.UserProfile{
		Name:   user.Name,
		Email:  user.Email,
		Avatar: user.AvatarUploadID,
	}, nil
}

func (u *UserService) UpdatePassword(userID int, update *dto.UpdatePassword) error {
	if err := entity.ValidatePassword(update.NewPassword); err != nil {
		return err
	}
	user, err := u.userRepo.GetUser(map[string]interface{}{"id": userID})
	if err != nil {
		return err
	}
	if !user.CheckPassword(update.OldPassword) {
		return entity.NewClientError("неверный пароль", entity.ErrBadRequest)
	}
	salt, hash, err := entity.HashPassword(update.NewPassword)
	if err != nil {
		return err
	}
	err = u.userRepo.UpdateUser(
		map[string]interface{}{"id": userID},
		map[string]interface{}{"password_hashed": hash, "salt_password": salt},
	)
	if err != nil {
		return err
	}
	return nil
}
