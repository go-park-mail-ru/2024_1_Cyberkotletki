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

func Init(logger *log.Logger, params config.InitParams) *http.Server {

	// REST API
	router := mux.NewRouter()
	rest.RegisterRoutes(router, params)

	// Swagger
	if params.GenSwagger {
		cmd := exec.Command("swag", "init", "--dir", "cmd/app,internal/transport/rest", "--parseDependency")
		if out, err := cmd.Output(); err != nil {
			logger.Fatalf("Не удалось сгенерировать документацию сваггер по причине: %s", out)
		} else {
			logger.Printf("Логи swagger кодогена:\n%s", out)
		}
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
		// для secure-соединения лучше использовать nginx
		// return server.ListenAndServeTLS()
		return nil
	}
}

func Shutdown(server *http.Server, logger *log.Logger, quit <-chan os.Signal, done chan<- bool, params config.InitParams) {
	<-quit
	logger.Println("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), params.GracefulShutdownTime)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Не удалось завершить работу http-сервера: %v\n", err)
	}
	close(done)
}
