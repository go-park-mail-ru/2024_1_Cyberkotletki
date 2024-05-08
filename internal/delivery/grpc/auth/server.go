package auth

import (
	"context"
	authProto "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/auth/proto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
)

type Grpc struct {
	authProto.UnimplementedAuthServiceServer
	authUC usecase.Auth
}

func NewGrpc(authUC usecase.Auth) *Grpc {
	return &Grpc{authUC: authUC}
}

func (service *Grpc) Logout(_ context.Context, session *authProto.Session) (*authProto.Nothing, error) {
	err := service.authUC.Logout(session.GetToken())
	if err != nil {
		return nil, err
	}
	return &authProto.Nothing{}, nil
}

func (service *Grpc) LogoutAll(_ context.Context, userID *authProto.User) (*authProto.Nothing, error) {
	err := service.authUC.LogoutAll(int(userID.GetId()))
	if err != nil {
		return nil, err
	}
	return &authProto.Nothing{}, nil
}

func (service *Grpc) GetUserIDBySession(_ context.Context, session *authProto.Session) (*authProto.User, error) {
	userID, err := service.authUC.GetUserIDBySession(session.GetToken())
	if err != nil {
		return nil, err
	}
	return &authProto.User{Id: uint64(userID)}, nil
}

func (service *Grpc) CreateSession(_ context.Context, userID *authProto.User) (*authProto.Session, error) {
	session, err := service.authUC.CreateSession(int(userID.GetId()))
	if err != nil {
		return nil, err
	}
	return &authProto.Session{Token: session}, nil
}

func (service *Grpc) Ping(_ context.Context, _ *authProto.Nothing) (*authProto.Nothing, error) {
	return &authProto.Nothing{}, nil
}
