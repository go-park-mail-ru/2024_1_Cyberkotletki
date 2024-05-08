package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/static"
	staticProto "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/delivery/grpc/static/proto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/postgres"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/service"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/connector"
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
	var conf config.StaticConfig
	defaults.SetDefaults(&conf)

	yamlData, err := yaml.Marshal(&conf)
	if err != nil {
		fmt.Printf("Ошибка при маршализации YAML: %v\n", err)
		return
	}
	file, err := os.Create("config_static.yaml")
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

func ParseParams() config.StaticConfig {
	yamlFile, err := os.ReadFile("config_static.yaml")
	if err != nil {
		log.Fatalf("Ошибка при чтении конфига сервера: %v", err)
	}
	var cfg config.StaticConfig
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Ошибка при парсинге конфига сервера: %v", err)
	}

	return cfg
}

func main() {
	genCfg := flag.Bool("generate-example-config", false, "Генерирует пример конфига, с которым умеет работать сервер")
	flag.Parse()
	if *genCfg {
		GenerateExampleConfig()
		return
	}
	logger := log.New("Static Microservice: ")
	params := ParseParams()
	logger.Printf("Параметры запуска сервера: %v \n", params)

	s3conn, err := connector.GetS3Connector(
		params.S3.Endpoint, params.S3.Region, params.S3.AccessKeyID, params.S3.SecretAccessKey,
	)
	if err != nil {
		logger.Fatal("Ошибка при подключении к S3: ", err)
	}
	db, err := connector.GetPostgresConnector(params.Postgres.ConnectURL)
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
