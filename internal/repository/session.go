package repository

import "errors"

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_session.go
type Session interface {
	// NewSession создает новую сессию для пользователя и возвращает ее ID
	NewSession(id int) (string, error)
	// CheckSession проверяет, что сессия существует и возвращает ID пользователя
	// Если сессия не существует, возвращает ошибку ErrSessionNotFound
	CheckSession(session string) (int, error)
	// DeleteAllSessions удаляет все сессии пользователя
	DeleteAllSessions(userID int) error
	// DeleteSession удаляет сессию
	DeleteSession(session string) error
}

var (
	ErrSessionNotFound = errors.New("пользователь с такой сессией не найден")
)
