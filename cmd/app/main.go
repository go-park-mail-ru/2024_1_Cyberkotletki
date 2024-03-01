package main

import (
	"errors"
	"flag"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/app"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var (
	listenAddress string
	listenPort    string
	serverMode    config.ServerMode
)

func main() {
	flag.StringVar(&listenAddress, "listen-address", "0.0.0.0", "server listen address")
	flag.StringVar(&listenPort, "listen-port", ":8000", "server listen port")
	flag.IntVar((*int)(&serverMode), "server-mode", 2, "0 = deploy\n1 = test\n2 = dev")
	flag.Parse()
	params := config.InitParams{
		Addr: listenAddress + listenPort,
		Mode: serverMode,
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
