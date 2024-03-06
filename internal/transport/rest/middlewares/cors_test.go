package middlewares

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/config"
	"testing"
)

func Test_SetAllowedOriginsFromConfig(t *testing.T) {
	tests := []struct {
		name   string
		params config.InitParams
		want   string
	}{
		{
			name: "Test with CORS set to '*'",
			params: config.InitParams{
				CORS: "*",
			},
			want: "*",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corsConfig := &CORSConfig{}
			corsConfig.SetAllowedOriginsFromConfig(tt.params)
			if corsConfig.AllowedOrigin != tt.want {
				t.Errorf("SetAllowedOriginsFromConfig() = %v, want %v", corsConfig.AllowedOrigin, tt.want)
			}
		})
	}
}
