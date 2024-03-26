package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_session.go
type Session interface {
	NewSession(id int) (string, error)
	CheckSession(session string) (bool, error)
	DeleteSession(session string) (bool, error)
}
