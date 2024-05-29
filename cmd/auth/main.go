package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/auth"
	authProto "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/auth/proto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/redis"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/connector"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/mcuadros/go-defaults"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func GenerateExampleConfig() {
	// Создание экземпляра структуры с использованием значений по умолчанию из тегов default
	var conf config.AuthConfig
	defaults.SetDefaults(&conf)

	yamlData, err := yaml.Marshal(&conf)
	if err != nil {
		fmt.Printf("Ошибка при маршализации YAML: %v\n", err)
		return
	}
	file, err := os.Create("config_auth.yaml")
	if err != nil {
		fmt.Printf("Ошибка при создании файла: %v\n", err)
		return
	}
	_, err = file.Write(yamlData)
	if err != nil {
		fmt.Printf("Ошибка при записи в файл: %v\n", err)
		return
	}
	err = file.Close()
	if err != nil {
		fmt.Printf("Ошибка при записи в файл: %v\n", err)
		return
	}

	fmt.Println("Конфигурационный файл успешно создан.")
}

func ParseParams() config.AuthConfig {
	var cfg config.AuthConfig
	// читаем конфиг
	yamlFile, err := os.ReadFile("config_auth.yaml")
	if err != nil {
		log.Fatalf("Ошибка при чтении конфига сервера: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Ошибка при парсинге конфига сервера: %v", err)
	}
	// читаем переменные окружения
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки .env файла")
	}
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	return cfg
}

func main() {
	genCfg := flag.Bool("generate-example-config", false, "Генерирует пример конфига, с которым умеет работать сервер")
	flag.Parse()
	if *genCfg {
		GenerateExampleConfig()
		return
	}
	logger := log.New("Auth Microservice: ")
	params := ParseParams()
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
