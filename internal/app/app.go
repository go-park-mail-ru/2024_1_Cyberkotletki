package app

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

type InitParams struct {
	Addr string
	Mode config.ServerMode
}

// todo: сделать RunParams и как-то логически разделить его с InitParams

func Init(logger *log.Logger, params InitParams) *http.Server {
	// todo: ко 2 рк надо будет сделать более продвинутый инит

	router := mux.NewRouter()
	// todo: сделать нормальный роутер :)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return &http.Server{
		Addr:         params.Addr,
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

func Run(server *http.Server, params InitParams) error {
	if params.Mode != config.DeployMode {
		return server.ListenAndServe()
	} else {
		// todo
		// server.ListenAndServeTLS()
		return nil
	}
}

func Shutdown(server *http.Server, logger *log.Logger, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	logger.Println("Завершение работы сервера...")

	// todo: занести время для выполнения graceful shutdown в конфиг
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Не удалось завершить работу http-сервера: %v\n", err)
	}
	close(done)
}
