package middlewares

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"net/http"
)

type CORSConfig struct {
	AllowedOrigin string
}

func (corsConfig *CORSConfig) SetAllowedOriginsFromConfig(params config.InitParams) {
	// todo: чтение настроек cors из конфига
	switch params.Mode {
	case config.DevMode:
		corsConfig.AllowedOrigin = "0.0.0.0:8001"
	case config.DeployMode, config.TestMode:
		// todo
	}
}

func (corsConfig *CORSConfig) SetCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo добавить больше заголовков и настроить их через конфиг
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Origin", corsConfig.AllowedOrigin)
		w.Header().Set("Vary", "Origin")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
