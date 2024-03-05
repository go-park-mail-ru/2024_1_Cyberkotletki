package main

import (
	"errors"
	"flag"
	"fmt"
	_ "github.com/go-park-mail-ru/2024_1_Cyberkotletki/docs"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/app"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
)

var (
	listenAddress string
	listenPort    string
	serverMode    config.ServerMode
	genSwagger    bool
	staticFolder  string
	cors          string
)

// @title API Киноскопа
// @version 1.0
func main() {
	staticDefaultFolder, _ := os.Getwd()
	flag.StringVar(&listenAddress, "listen-address", "localhost", "Адрес сервера")
	flag.StringVar(&listenPort, "listen-port", "8000", "Порт сервера")
	flag.IntVar((*int)(&serverMode), "server-mode", 2, "0 = deploy\n1 = test\n2 = dev")
	flag.BoolVar(&genSwagger, "generate-swagger", true, "true = сгенерировать swagger")
	flag.StringVar(&staticFolder, "static-folder", filepath.Join(staticDefaultFolder, "assets", "examples", "static"), "путь до папки со статикой")
	flag.StringVar(&staticFolder, cors, "http://localhost", "Параметр Access-Control-Allow-Origin")
	flag.Parse()
	params := config.InitParams{
		Addr:         fmt.Sprintf("%s:%s", listenAddress, listenPort),
		Mode:         serverMode,
		GenSwagger:   genSwagger,
		StaticFolder: staticFolder,
		CORS:         cors,
	}

	logger := log.New(os.Stdout, "server: ", log.LstdFlags)

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := app.Init(logger, params)
	go app.Shutdown(server, logger, quit, done)

	logger.Println("Сервер запущен по адресу", params.Addr)
	if err := app.Run(server, params); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("Невозможно обрабатывать запросы по адресу %s: %v\n", params.Addr, err)
	}

	<-done
	logger.Println("Сервер завершил свою работу")
}
