package config

import "fmt"

type Server struct {
	IP                      string `yaml:"ip"                        default:"localhost"`
	Port                    int    `yaml:"port"                      default:"8080"`
	WriteTimeout            int    `yaml:"write_timeout"             default:"15"`
	ReadTimeout             int    `yaml:"read_timeout"              default:"15"`
	ReadHeaderTimeout       int    `yaml:"read_header_timeout"       default:"15"`
	GracefulShutdownTimeout int    `yaml:"graceful_shutdown_timeout" default:"60"`
}

type Config struct {
	HTTP struct {
		StaticFolder       string `yaml:"static_folder"        default:"assets/examples/static"`
		CORSAllowedOrigins string `yaml:"cors_allowed_origins" default:"http://localhost:8000/"`
		SecureCookies      bool   `yaml:"secure_cookies"       default:"false"`
		Server             Server `yaml:"server"`
	} `yaml:"http"`
	Auth struct {
		SessionAliveTime int `yaml:"session_alive_time" default:"86400"`
		Redis            struct {
			Addr     string `yaml:"addr"     default:"localhost:6379"`
			Password string `yaml:"password" default:""`
			DB       int    `yaml:"db"       default:"0"`
		} `yaml:"redis"`
	} `yaml:"auth"`
}

func (cfg *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", cfg.HTTP.Server.IP, cfg.HTTP.Server.Port)
}
