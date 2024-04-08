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

func (a AuthService) LogoutAll(userID int) error {
	if err := a.sessionRepo.DeleteAllSessions(userID); err != nil {
		return err
	}
	return nil
}

func (a AuthService) GetUserIDBySession(session string) (int, error) {
	userID, err := a.sessionRepo.CheckSession(session)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (a AuthService) CreateSession(userID int) (string, error) {
	session, err := a.sessionRepo.NewSession(userID)
	if err != nil {
		return "", err
	}
	return session, nil
}
