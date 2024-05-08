package auth

import (
	"context"
	auth "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/auth/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	authManager auth.AuthServiceClient
}

func NewGateway(connectAddr string) (*Gateway, error) {
	grpcConn, err := grpc.Dial(
		connectAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	authManager := auth.NewAuthServiceClient(grpcConn)

	_, err = authManager.Ping(context.Background(), &auth.Nothing{})
	if err != nil {
		return nil, err
	}

	return &Gateway{authManager: authManager}, nil
}

func (gate *Gateway) Logout(session string) error {
	_, err := gate.authManager.Logout(context.Background(), &auth.Session{Token: session})
	if err != nil {
		return err
	}
	return nil
}

func (gate *Gateway) LogoutAll(userID int) error {
	_, err := gate.authManager.LogoutAll(context.Background(), &auth.User{Id: uint64(userID)})
	if err != nil {
		return err
	}
	return nil
}

func (gate *Gateway) GetUserIDBySession(session string) (int, error) {
	user, err := gate.authManager.GetUserIDBySession(context.Background(), &auth.Session{Token: session})
	if err != nil {
		return 0, err
	}
	return int(user.Id), nil
}

func (gate *Gateway) CreateSession(userID int) (string, error) {
	session, err := gate.authManager.CreateSession(context.Background(), &auth.User{Id: uint64(userID)})
	if err != nil {
		return "", err
	}
	return session.Token, nil
}
