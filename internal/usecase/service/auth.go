package service

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
)

type AuthService struct {
	userRepo    repository.User
	sessionRepo repository.Session
}

func NewAuthService(userRepo repository.User, sessionRepo repository.Session) usecase.Auth {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (a AuthService) Register(regDTO *dto.Register) (string, error) {
	user := entity.NewUserEmpty()
	if err := entity.ValidateEmail(regDTO.Email); err != nil {
		return "", err
	}
	if err := entity.ValidatePassword(regDTO.Password); err != nil {
		return "", err
	}
	salt, hash, err := entity.HashPassword(regDTO.Password)
	if err != nil {
		return "", err
	}
	user.Email = regDTO.Email
	user.PasswordHash = hash
	user.PasswordSalt = salt
	user, err = a.userRepo.AddUser(user)
	if err != nil {
		return "", err
	}
	session, err := a.sessionRepo.NewSession(user.ID)
	if err != nil {
		return "", err
	}
	return session, nil
}

func (a AuthService) Login(login *dto.Login) (string, error) {
	user, err := a.userRepo.GetUserByEmail(login.Login)
	if err != nil {
		return "", err
	}
	if !user.CheckPassword(login.Password) {
		return "", entity.NewClientError("неверный пароль", entity.ErrForbidden)
	}
	session, err := a.sessionRepo.NewSession(user.ID)
	if err != nil {
		return "", err
	}
	return session, nil
}

func (a AuthService) IsAuth(s string) (bool, error) {
	isAuth, err := a.sessionRepo.CheckSession(s)
	if err != nil {
		return false, err
	}
	return isAuth, nil
}

func (a AuthService) Logout(s string) error {
	if _, err := a.sessionRepo.DeleteSession(s); err != nil {
		return err
	}
	return nil
}
