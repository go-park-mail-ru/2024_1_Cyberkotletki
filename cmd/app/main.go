package main

import (
	"errors"
	_ "github.com/go-park-mail-ru/2024_1_Cyberkotletki/docs"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/app"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// @title API Киноскопа
// @version 1.0
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
