package app

import (
	"context"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/user"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/transport/rest"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// todo: сделать RunParams и как-то логически разделить его с InitParams

func Init(logger *log.Logger, params config.InitParams) *http.Server {
	// todo: ко 2 рк надо будет сделать более продвинутый инит

	// REST API
	router := mux.NewRouter()
	rest.RegisterRoutes(router, params)

	// Swagger
	cmd := exec.Command("swag", "init", "--dir", "cmd/app,internal/transport/rest", "--parseDependency")
	if out, err := cmd.Output(); err != nil {
		logger.Fatal("Не удалось сгенерировать документацию сваггер по причине: ", err)
	} else {
		logger.Printf("Логи swagger кодогена:\n%s", out)
	}

	// DB
	user.UsersDatabase.InitDB()
	content.FilmsDatabase.InitDB()

	return &http.Server{
		Addr:         params.Addr,
		Handler:      router,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

func Run(server *http.Server, params config.InitParams) error {
	if params.Mode == config.DevMode {
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
