package service

import (
	"errors"
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

// Logout - выход из сессии
func (a AuthService) Logout(s string) error {
	if err := a.sessionRepo.DeleteSession(s); err != nil {
		return err
	}
	return nil
}

// LogoutAll - выход из всех сессий
func (a AuthService) LogoutAll(userID int) error {
	if err := a.sessionRepo.DeleteAllSessions(userID); err != nil {
		return err
	}
	return nil
}

// GetUserIDBySession - получение ID пользователя по сессии
func (a AuthService) GetUserIDBySession(session string) (int, error) {
	userID, err := a.sessionRepo.CheckSession(session)
	switch {
	case errors.Is(err, repository.ErrSessionNotFound):
		return 0, usecase.ErrSessionNotFound
	case err != nil:
		return 0, err
	default:
		return userID, nil
	}
}

// CreateSession - создание новой сессии
func (a AuthService) CreateSession(userID int) (string, error) {
	session, err := a.sessionRepo.NewSession(userID)
	if err != nil {
		return "", err
	}
	return session, nil
}
