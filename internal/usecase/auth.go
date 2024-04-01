package usecase

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_auth.go
type Auth interface {
	Logout(session string) error
	LogoutAll(userID int) error
	GetUserIDBySession(session string) (int, error)
	CreateSession(userID int) (string, error)
}
