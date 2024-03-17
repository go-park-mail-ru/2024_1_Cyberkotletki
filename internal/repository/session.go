package repository

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_session.go
type Session interface {
	NewSession(id int) string
	CheckSession(session string) bool
	DeleteSession(session string) bool
}
