package static

import (
	"bytes"
	"context"
	"errors"
	staticProto "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/static/proto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"io"
	"net/url"
)

type Grpc struct {
	staticProto.UnimplementedStaticServiceServer
	staticUC usecase.Static
}

func NewGrpc(staticUC usecase.Static) *Grpc {
	return &Grpc{staticUC: staticUC}
}

func (service *Grpc) GetStatic(_ context.Context, static *staticProto.Static) (*staticProto.Static, error) {
	uri, err := service.staticUC.GetStatic(int(static.GetId()))
	if err != nil {
		return nil, err
	}
	return &staticProto.Static{Uri: uri}, nil
}

func (service *Grpc) UploadAvatar(stream staticProto.StaticService_UploadAvatarServer) error {
	var bytesAvatar []byte

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		bytesAvatar = append(bytesAvatar, chunk.Chunk...)
	}

	reader := bytes.NewReader(bytesAvatar)
	staticID, err := service.staticUC.UploadAvatar(reader)
	switch {
	case errors.Is(err, usecase.ErrStaticTooBigFile):
		return stream.SendAndClose(&staticProto.Static{Error: "ErrStaticTooBigFile"})
	case errors.Is(err, usecase.ErrStaticNotImage):
		return stream.SendAndClose(&staticProto.Static{Error: "ErrStaticNotImage"})
	case errors.Is(err, usecase.ErrStaticImageDimensions):
		return stream.SendAndClose(&staticProto.Static{Error: "ErrStaticImageDimensions"})
	case err != nil:
		return err
	default:
		return stream.SendAndClose(&staticProto.Static{Id: uint64(staticID)})
	}
}

func (service *Grpc) GetStaticFile(
	static *staticProto.Static,
	stream staticProto.StaticService_GetStaticFileServer,
) error {
	uri, err := url.QueryUnescape(static.GetUri())
	if err != nil {
		return err
	}
	file, err := service.staticUC.GetStaticFile(uri)
	if err != nil {
		return err
	}
	buffer := make([]byte, 1024)
	for {
		bytesCount, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = stream.Send(&staticProto.StaticUpload{Chunk: buffer[:bytesCount]})
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *Grpc) Ping(context.Context, *staticProto.Nothing) (*staticProto.Nothing, error) {
	return &staticProto.Nothing{}, nil
}
