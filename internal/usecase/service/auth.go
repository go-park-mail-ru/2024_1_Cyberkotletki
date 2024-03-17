package service

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/DTO"
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

func (a AuthService) Register(regDTO DTO.Register) (string, error) {
	us := entity.NewUserEmpty()
	if err := entity.ValidateEmail(regDTO.Email); err != nil {
		return "", err
	}
	if err := entity.ValidatePassword(regDTO.Password); err != nil {
		return "", err
	}
	salt, hash := entity.HashPassword(regDTO.Password)
	us.Email = regDTO.Email
	us.PasswordHash = hash
	us.PasswordSalt = salt
	us, err := a.userRepo.AddUser(*us)
	if err != nil {
		return "", err
	} else {
		return a.sessionRepo.NewSession(us.Id), nil
	}
}

func (a AuthService) Login(login DTO.Login) (string, error) {
	us, err := a.userRepo.GetUserByEmail(login.Login)
	if err != nil {
		return "", err
	}
	if !us.CheckPassword(login.Password) {
		return "", entity.NewClientError("неверный пароль", entity.ErrBadRequest)
	}
	return a.sessionRepo.NewSession(us.Id), nil
}

func (a AuthService) IsAuth(s string) bool {
	return a.sessionRepo.CheckSession(s)
}

func (a AuthService) Logout(s string) error {
	if ok := a.sessionRepo.DeleteSession(s); ok {
		return nil
	}
	return entity.NewClientError("сессия недействительна")
}
