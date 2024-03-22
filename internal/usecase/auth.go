package usecase

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_auth.go
type Auth interface {
	Register(*dto.Register) (string, error)
	Login(*dto.Login) (string, error)
	IsAuth(string) (bool, error)
	Logout(string) error
}
