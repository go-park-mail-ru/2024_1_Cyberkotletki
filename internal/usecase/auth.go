package usecase

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/DTO"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mock_auth.go
type Auth interface {
	Register(DTO.Register) (string, error)
	Login(DTO.Login) (string, error)
	IsAuth(string) (bool, error)
	Logout(string) error
}
