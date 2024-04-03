package config

import (
	"fmt"
	"github.com/mcuadros/go-defaults"
)

type Server struct {
	IP                      string `yaml:"ip"                        default:"localhost"`
	Port                    int    `yaml:"port"                      default:"8080"`
	WriteTimeout            int    `yaml:"write_timeout"             default:"15"`
	ReadTimeout             int    `yaml:"read_timeout"              default:"15"`
	ReadHeaderTimeout       int    `yaml:"read_header_timeout"       default:"15"`
	GracefulShutdownTimeout int    `yaml:"graceful_shutdown_timeout" default:"60"`
}

type RedisDatabase struct {
	Addr     string `yaml:"addr"     default:"localhost:6379"`
	Password string `yaml:"password" default:""`
	DB       int    `yaml:"db"       default:"0"`
}

type PostgresDatabase struct {
	// nolint
	ConnectURL string `yaml:"connect_url" default:"postgres://kinoskop_admin:admin_secret_password@localhost:5432/kinoskop?sslmode=disable"`
}

type Config struct {
	HTTP struct {
		StaticFolder       string `yaml:"static_folder"        default:"assets/examples/static"`
		CORSAllowedOrigins string `yaml:"cors_allowed_origins" default:"http://localhost:8000/"`
		SecureCookies      bool   `yaml:"secure_cookies"       default:"false"`
		Server             Server `yaml:"server"`
	} `yaml:"http"`
	Auth struct {
		SessionAliveTime int           `yaml:"session_alive_time" default:"86400"`
		Redis            RedisDatabase `yaml:"redis"`
	} `yaml:"auth_service"`
	User struct {
		Postgres PostgresDatabase `yaml:"postgres"`
	} `yaml:"user_service"`
	Static struct {
		MaxFileSize int              `yaml:"max_file_size" default:"10485760"`
		Path        string           `yaml:"path"          default:"assets/examples/static"`
		Postgres    PostgresDatabase `yaml:"postgres"`
	} `yaml:"static_service"`
}

func (cfg *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", cfg.HTTP.Server.IP, cfg.HTTP.Server.Port)
}

func NewConfigWithDefaults() *Config {
	var conf Config
	defaults.SetDefaults(&conf)
	return &conf
}
