package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_session.go
type Session interface {
	NewSession(id int) (string, error)
	CheckSession(session string) (int, error)
	DeleteAllSessions(userID int) error
	DeleteSession(session string) error
}
