package static

import (
	"bytes"
	"context"
	static "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/static/proto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"strings"
)

const (
	bufSize = 1024
)

type Gateway struct {
	staticManager static.StaticServiceClient
}

func NewGateway(connectAddr string) (*Gateway, error) {
	grpcConn, err := grpc.Dial(
		connectAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	staticManager := static.NewStaticServiceClient(grpcConn)

	_, err = staticManager.Ping(context.Background(), &static.Nothing{})
	if err != nil {
		return nil, err
	}

	return &Gateway{staticManager: staticManager}, nil
}

func (gate *Gateway) GetStatic(staticID int) (string, error) {
	staticFile, err := gate.staticManager.GetStatic(context.Background(), &static.Static{Id: uint64(staticID)})
	if err != nil {
		if strings.Contains(err.Error(), repository.ErrStaticNotFound.Error()) {
			return "", usecase.ErrStaticNotFound
		}
		return "", err
	}
	switch staticFile.Error {
	case "ErrStaticNotFound":
		return "", usecase.ErrStaticNotFound
	default:
		return staticFile.Uri, nil
	}
}

func (gate *Gateway) GetStaticFile(staticURI string) (io.ReadSeeker, error) {
	stream, err := gate.staticManager.GetStaticFile(context.Background(), &static.Static{Uri: staticURI})
	if err != nil {
		if strings.Contains(err.Error(), repository.ErrStaticNotFound.Error()) {
			return nil, usecase.ErrStaticNotFound
		}
		return nil, err
	}

	var buffer []byte
	for {
		chunk, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			if strings.Contains(err.Error(), repository.ErrStaticNotFound.Error()) {
				return nil, usecase.ErrStaticNotFound
			}
			return nil, err
		}
		buffer = append(buffer, chunk.Chunk...)
	}

	err = stream.CloseSend()
	if err != nil {
		return nil, err
	}

	return io.ReadSeeker(bytes.NewReader(buffer)), nil
}

func (gate *Gateway) UploadAvatar(reader io.ReadSeeker) (int, error) {
	stream, err := gate.staticManager.UploadAvatar(context.Background())
	if err != nil {
		return -1, err
	}

	buffer := make([]byte, bufSize)

	for {
		bytesRead, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return -1, err
		}

		err = stream.Send(&static.StaticUpload{
			Chunk: buffer[:bytesRead],
		})
		if err != nil {
			return -1, err
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		return -1, err
	}
	switch response.Error {
	case "ErrStaticTooBigFile":
		return -1, usecase.ErrStaticTooBigFile
	case "ErrStaticNotImage":
		return -1, usecase.ErrStaticNotImage
	case "ErrStaticImageDimensions":
		return -1, usecase.ErrStaticImageDimensions
	default:
		return int(response.Id), nil
	}
}
