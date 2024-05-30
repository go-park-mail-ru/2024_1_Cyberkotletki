package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/auth"
	authProto "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/auth/proto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/redis"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/connector"
	"github.com/labstack/gommon/log"
	"google.golang.org/grpc"
	"net"
	"os/signal"
	"syscall"
)

func main() {
	logger := log.New("Auth Microservice: ")
	params := config.ParseAuthServiceParams()
	logger.Printf("Параметры запуска сервера: %v \n", params)

	redisConn, err := connector.GetRedisConnector(params.Redis.Addr, params.Redis.Password, params.Redis.DB)
	if err != nil {
		logger.Fatal("Ошибка при подключении к Redis: ", err)
	}
	authRepository := redis.NewSessionRepository(redisConn, params.SessionAliveTime)
	authUseCase := service.NewAuthService(authRepository)
	authService := auth.NewGrpc(authUseCase)
	server := grpc.NewServer()
	authProto.RegisterAuthServiceServer(server, authService)
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
