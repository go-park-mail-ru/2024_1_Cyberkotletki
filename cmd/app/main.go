package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	_ "github.com/go-park-mail-ru/2024_1_Cyberkotletki/docs"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/app"
	_ "github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v3"
)

func GenerateExampleConfig() {
	// Создание экземпляра структуры с использованием значений по умолчанию из тегов default
	var conf config.Config
	defaults.SetDefaults(&conf)

	yamlData, err := yaml.Marshal(&conf)
	if err != nil {
		fmt.Printf("Ошибка при маршализации YAML: %v\n", err)
		return
	}
	file, err := os.Create("config.yaml")
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

func ParseParams() config.Config {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка при чтении конфига сервера: %v", err)
	}
	var cfg config.Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Ошибка при парсинге конфига сервера: %v", err)
	}

	return cfg
}

// @title API Киноскопа
// @version 1.0
// @Description сервис Киноскоп (аналог кинопоиска)
// @BasePath  /api
// @securityDefinitions.apikey _csrf
// @in header
// @name x-csrf
func main() {
	genCfg := flag.Bool("generate-example-config", false, "Генерирует пример конфига, с которым умеет работать сервер")
	flag.Parse()
	if *genCfg {
		GenerateExampleConfig()
		return
	}

	logger := log.New("server: ")
	params := ParseParams()
	logger.Printf("Параметры запуска сервера: %v \n", params)

	echoServer := app.Init(logger, params)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	go app.Run(echoServer, params)

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(params.HTTP.Server.GracefulShutdownTimeout)*time.Second,
	)
	defer cancel()
	app.Shutdown(ctx, echoServer)
}
