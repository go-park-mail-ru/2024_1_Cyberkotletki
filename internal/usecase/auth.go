package usecase

import "errors"

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_auth.go
type Auth interface {
	// Logout - выход из сессии
	Logout(session string) error
	// LogoutAll - выход из всех сессий
	LogoutAll(userID int) error
	// GetUserIDBySession - получение ID пользователя по сессии.
	// Возвращает ErrSessionNotFound, если сессия не найдена
	GetUserIDBySession(session string) (int, error)
	// CreateSession - создание новой сессии
	CreateSession(userID int) (string, error)
}

var (
	ErrSessionNotFound = errors.New("сессия не найдена")
)
