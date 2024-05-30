package profanity

import (
	"context"
	"errors"
	"fmt"
	profanity "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/profanity/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Gateway struct {
	profanityManager profanity.ProfanityFilterClient
}

func NewGateway(connectAddr string) (*Gateway, error) {
	grpcConn, err := grpc.Dial(
		connectAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	profanityManager := profanity.NewProfanityFilterClient(grpcConn)

	_, err = profanityManager.Ping(context.Background(), &profanity.Nothing{})
	if err != nil {
		return nil, err
	}

	return &Gateway{profanityManager: profanityManager}, nil
}

func (gate *Gateway) FilterMessage(text string) (string, error) {
	filteredMessage, err := gate.profanityManager.FilterMessage(context.Background(), &profanity.Text{Text: text})
	if err != nil {
		return "", errors.Join(fmt.Errorf("ошибка при фильтрации сообщения"), err)
	}
	return filteredMessage.GetText(), nil
}
