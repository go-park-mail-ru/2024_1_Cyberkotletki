package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	_ "github.com/go-park-mail-ru/2024_1_Cyberkotletki/docs"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/app"
	_ "github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v3"
	"os"
	"os/signal"
	"time"
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

func ParseParams(logger *log.Logger) config.Config {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка при чтении конфига сервера: %v", err)
	}
	var cfg config.Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Ошибка при парсинге конфига сервера: %v", err)
	}

	/*
		при подключении БД все данные надо будет хранить в .env
		if err = godotenv.Load(".env"); err != nil {
			logger.Fatal("Не удалось загрузить .env файл. Убедитесь, что он расположен в корне проекта")
		}
	*/

	return cfg
}

// @title API Киноскопа
// @version 1.0
// @BasePath  /api
func main() {
	genCfg := flag.Bool("generate-example-config", false, "Генерирует пример конфига, с которым умеет работать сервер")
	flag.Parse()
	if *genCfg {
		GenerateExampleConfig()
		return
	}

	logger := log.New("server: ")
	params := ParseParams(logger)
	logger.Printf("Параметры запуска сервера: %v \n", params)

	echoServer := app.Init(logger, params)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer stop()
	go app.Run(echoServer, params)

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(params.HTTP.Server.GracefulShutdownTimeout)*time.Second)
	defer cancel()
	app.Shutdown(echoServer, ctx)
}
