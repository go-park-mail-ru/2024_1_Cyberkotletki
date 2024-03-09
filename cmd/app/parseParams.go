package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"github.com/joho/godotenv"
	"log"
	"os"
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
