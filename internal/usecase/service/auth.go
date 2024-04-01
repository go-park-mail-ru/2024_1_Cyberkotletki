package service

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
)

type AuthService struct {
	sessionRepo repository.Session
}

func NewAuthService(sessionRepo repository.Session) usecase.Auth {
	return &AuthService{
		sessionRepo: sessionRepo,
	}
}

func (a AuthService) Logout(s string) error {
	if err := a.sessionRepo.DeleteSession(s); err != nil {
		return err
	}
	return nil
}

func (a AuthService) LogoutAll(userId int) error {
	if err := a.sessionRepo.DeleteAllSessions(userId); err != nil {
		return err
	}
	return nil
}

func (a AuthService) GetUserIDBySession(session string) (int, error) {
	userId, err := a.sessionRepo.CheckSession(session)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func (a AuthService) CreateSession(userId int) (string, error) {
	session, err := a.sessionRepo.NewSession(userId)
	if err != nil {
		return "", err
	}
	return session, nil
}
