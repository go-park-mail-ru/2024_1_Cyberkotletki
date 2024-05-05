package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"io"
)

type UserService struct {
	userRepo repository.User
	staticUC usecase.Static
}

func NewUserService(userRepo repository.User, staticUC usecase.Static) usecase.User {
	return &UserService{
		userRepo: userRepo,
		staticUC: staticUC,
	}
}

func (u *UserService) Register(regDTO *dto.Register) (int, error) {
	if err := entity.ValidateEmail(regDTO.Email); err != nil {
		return -1, usecase.UserIncorrectDataError{Err: err}
	}
	if err := entity.ValidatePassword(regDTO.Password); err != nil {
		return -1, usecase.UserIncorrectDataError{Err: err}
	}
	salt, hash, err := entity.HashPassword(regDTO.Password)
	if err != nil {
		return -1, entity.UsecaseWrap(errors.New("ошибка при хешировании пароля"), err)
	}
	user, err := u.userRepo.AddUser(regDTO.Email, hash, salt)
	switch {
	case errors.Is(err, repository.ErrUserAlreadyExists):
		return -1, usecase.ErrUserAlreadyExists
	case errors.Is(err, repository.ErrUserIncorrectData):
		return -1, usecase.UserIncorrectDataError{Err: err}
	case err != nil:
		return -1, entity.UsecaseWrap(errors.New("ошибка при регистрации пользователя"), err)
	}
	return user.ID, nil
}

func (u *UserService) Login(login *dto.Login) (int, error) {
	user, err := u.userRepo.GetUserByEmail(login.Login)
	if errors.Is(err, repository.ErrUserNotFound) {
		return -1, usecase.ErrUserNotFound
	}
	if err != nil {
		return -1, entity.UsecaseWrap(errors.New("ошибка при поиске пользователя"), err)
	}
	if !user.CheckPassword(login.Password) {
		return -1, usecase.UserIncorrectDataError{Err: errors.New("неверный пароль")}
	}
	return user.ID, nil
}

func (u *UserService) UpdateAvatar(userID int, reader io.Reader) error {
	user, err := u.userRepo.GetUserByID(userID)
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		return usecase.ErrUserNotFound
	case err != nil:
		return entity.UsecaseWrap(errors.New("ошибка при поиске пользователя"), err)
	}
	uploadID, err := u.staticUC.UploadAvatar(reader)
	switch {
	case errors.Is(err, usecase.ErrStaticTooBigFile):
		return usecase.UserIncorrectDataError{Err: err}
	case errors.Is(err, usecase.ErrStaticNotImage):
		return usecase.UserIncorrectDataError{Err: err}
	case errors.Is(err, usecase.ErrStaticImageDimensions):
		return usecase.UserIncorrectDataError{Err: err}
	case err != nil:
		return entity.UsecaseWrap(errors.New("ошибка при загрузке аватара"), err)
	}
	user.AvatarUploadID = uploadID
	err = u.userRepo.UpdateUser(user)
	if err != nil {
		return entity.UsecaseWrap(errors.New("ошибка при обновлении пользователя"), err)
	}
	return nil
}

func (u *UserService) UpdateInfo(userID int, update *dto.UserUpdate) error {
	user, err := u.userRepo.GetUserByID(userID)
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		return usecase.ErrUserNotFound
	case err != nil:
		return entity.UsecaseWrap(errors.New("ошибка при поиске пользователя"), err)
	}
	if err = entity.ValidateEmail(update.Email); err != nil {
		return usecase.UserIncorrectDataError{Err: err}
	}
	user.Email = update.Email
	if err = entity.ValidateName(update.Name); err != nil {
		return usecase.UserIncorrectDataError{Err: err}
	}
	user.Name = update.Name
	err = u.userRepo.UpdateUser(user)
	if err != nil {
		return entity.UsecaseWrap(errors.New("ошибка при обновлении пользователя"), err)
	}
	return nil
}

func (u *UserService) GetUser(userID int) (*dto.UserProfile, error) {
	user, err := u.userRepo.GetUserByID(userID)
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		return nil, usecase.ErrUserNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при поиске пользователя"), err)
	}
	var avatar string
	if user.AvatarUploadID > 0 {
		avatar, err = u.staticUC.GetStatic(user.AvatarUploadID)
		switch {
		case errors.Is(err, usecase.ErrStaticNotFound):
			avatar = ""
		case err != nil:
			return nil, entity.UsecaseWrap(errors.New("ошибка при получении аватара"), err)
		}
	}
	return &dto.UserProfile{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Rating: user.Rating,
		Avatar: avatar,
	}, nil
}

func (u *UserService) UpdatePassword(userID int, update *dto.UpdatePassword) error {
	user, err := u.userRepo.GetUserByID(userID)
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		return usecase.ErrUserNotFound
	case err != nil:
		return entity.UsecaseWrap(errors.New("ошибка при поиске пользователя"), err)
	}
	if err = entity.ValidatePassword(update.NewPassword); err != nil {
		return usecase.UserIncorrectDataError{Err: err}
	}
	if !user.CheckPassword(update.OldPassword) {
		return usecase.UserIncorrectDataError{Err: errors.New("неверный пароль")}
	}
	salt, hash, err := entity.HashPassword(update.NewPassword)
	if err != nil {
		return entity.UsecaseWrap(errors.New("ошибка при хешировании пароля"), err)
	}
	user.PasswordHash = hash
	user.PasswordSalt = salt
	err = u.userRepo.UpdateUser(user)
	if err != nil {
		return entity.UsecaseWrap(errors.New("ошибка при обновлении пользователя"), err)
	}
	return nil
}
