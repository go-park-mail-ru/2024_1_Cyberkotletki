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
)

// @title API Киноскопа
// @version 1.0
func main() {
	logger := log.New(os.Stdout, "server: ", log.LstdFlags)

	if err := godotenv.Load(); err != nil {
		log.Fatal("Не удалось загрузить .env файл. Убедитесь, что он расположен в корне проекта")
	}
	staticDefaultFolder, _ := os.Getwd()
	listenAddress := os.Getenv("LISTEN_ADDRESS")
	listenPort := os.Getenv("LISTEN_PORT")
	serverMode, err := strconv.ParseInt(os.Getenv("SERVER_MODE"), 10, 32)
	if err != nil || serverMode > 2 || serverMode < 0 {
		log.Fatal("Неправильное использование параметра SERVER_MODE")
	}
	var genSwagger bool
	if swagger := os.Getenv("GEN_SWAGGER"); swagger == "true" {
		genSwagger = true
	} else {
		genSwagger = false
	}
	staticFolder := os.Getenv("STATIC_DIR")
	if staticFolder == "" {
		staticFolder = filepath.Join(staticDefaultFolder, "assets", "examples", "static")
	}
	cors := os.Getenv("CORS")
	params := config.InitParams{
		Addr:         fmt.Sprintf("%s:%s", listenAddress, listenPort),
		Mode:         config.ServerMode(serverMode),
		GenSwagger:   genSwagger,
		StaticFolder: staticFolder,
		CORS:         cors,
	}

	logger.Printf("Параметры запуска сервера: %v \n", params)

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := app.Init(logger, params)
	go app.Shutdown(server, logger, quit, done)

	logger.Println("Сервер запущен по адресу", params.Addr)
	if err := app.Run(server, params); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("Сервер перестал обрабатывать запросы по адресу %s: %v\n", params.Addr, err)
	}

	<-done
	logger.Println("Сервер завершил свою работу")
}
