package middlewares

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"net/http"
)

type CORSConfig struct {
	AllowedOrigin string
}

func (corsConfig *CORSConfig) SetAllowedOriginsFromConfig(params config.InitParams) {
	corsConfig.AllowedOrigin = params.CORS
}

func (corsConfig *CORSConfig) SetCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Max-Age", "86400")
		w.Header().Set("Access-Control-Allow-Origin", corsConfig.AllowedOrigin)
		w.Header().Set("Vary", "Origin")
		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)
	})
}
