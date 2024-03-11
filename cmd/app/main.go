package main

import (
	"errors"
	"fmt"
	_ "github.com/go-park-mail-ru/2024_1_Cyberkotletki/docs"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/app"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"time"
)

func ParseParams(logger *log.Logger) config.InitParams {
	if err := godotenv.Load(); err != nil {
		logger.Fatal("Не удалось загрузить .env файл. Убедитесь, что он расположен в корне проекта")
	}
	listenAddress := os.Getenv("LISTEN_ADDRESS")
	listenPort := os.Getenv("LISTEN_PORT")
	serverMode := os.Getenv("SERVER_MODE")
	switch serverMode {
	case "PRODUCTION":
		serverMode = "prod"
	case "TEST":
		serverMode = "test"
	default:
		serverMode = "dev"
	}
	genSwagger := os.Getenv("GEN_SWAGGER") == "true"
	staticFolder := os.Getenv("STATIC_DIR")
	staticDefaultFolder, _ := os.Getwd()
	if staticFolder == "" {
		staticFolder = filepath.Join(staticDefaultFolder, "assets", "examples", "static")
	}
	cors := os.Getenv("CORS")
	sessionAliveTime, err := strconv.ParseInt(os.Getenv("SESSION_ALIVE_TIME"), 10, 64)
	if err != nil {
		logger.Fatal("Параметр SESSION_ALIVE_TIME должен быть валидным натуральным числом")
	}
	cookiesSecure := os.Getenv("COOKIES_SECURE") == "true"
	gracefulShutdownTime, err := strconv.ParseInt(os.Getenv("GRACEFUL_SHUTDOWN_TIME"), 10, 64)
	if err != nil {
		logger.Fatal("Параметр GRACEFUL_SHUTDOWN_TIME должен быть валидным натуральным числом")
	}
	return config.InitParams{
		Addr:                 fmt.Sprintf("%s:%s", listenAddress, listenPort),
		Mode:                 config.ServerMode(serverMode),
		GenSwagger:           genSwagger,
		StaticFolder:         staticFolder,
		CORS:                 cors,
		SessionAliveTime:     time.Duration(sessionAliveTime) * time.Second,
		CookiesSecure:        cookiesSecure,
		GracefulShutdownTime: time.Duration(gracefulShutdownTime) * time.Second,
	}
}

// @title API Киноскопа
// @version 1.0
// @BasePath  /api
func main() {
	logger := log.New(os.Stdout, "server: ", log.LstdFlags)
	params := ParseParams(logger)
	logger.Printf("Параметры запуска сервера: %v \n", params)

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := app.Init(logger, params)
	go app.Shutdown(server, logger, quit, done, params)

	logger.Println("Сервер запущен по адресу", params.Addr)
	if err := app.Run(server, params); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("Сервер перестал обрабатывать запросы по адресу %s: %v\n", params.Addr, err)
	}

	<-done
	logger.Println("Сервер завершил свою работу")
}
