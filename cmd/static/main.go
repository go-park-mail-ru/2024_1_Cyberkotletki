package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/static"
	staticProto "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/static/proto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/connector"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"net"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.New("Static Microservice: ")
	params := config.ParseStaticServiceParams()
	logger.Printf("Параметры запуска сервера: %v \n", params)

	s3conn, err := connector.GetS3Connector(
		params.S3.Endpoint, params.S3.Region, params.S3.AccessKeyID, params.S3.SecretAccessKey,
	)
	if err != nil {
		logger.Fatal("Ошибка при подключении к S3: ", err)
	}
	db, err := connector.GetPostgresConnector(params.Postgres.GetConnectURL())
	if err != nil {
		logger.Fatal("Ошибка при подключении к Postgres: ", err)
	}
	staticRepository := postgres.NewStaticRepository(db, s3conn, params.S3.BucketName, params.MaxFileSize)
	staticUseCase := service.NewStaticService(staticRepository)
	staticService := static.NewGrpc(staticUseCase)
	server := grpc.NewServer()
	staticProto.RegisterStaticServiceServer(server, staticService)
	addr := fmt.Sprintf("%s:%d", params.IP, params.Port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Невозможно прослушать порт:", err)
	}
	logger.Info("Слушаем grpc по адресу", addr)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	defer stop()
	go func() {
		err := server.Serve(lis)
		if err != nil {
			logger.Error(err)
		}
	}()
	<-ctx.Done()
	server.GracefulStop()
}
