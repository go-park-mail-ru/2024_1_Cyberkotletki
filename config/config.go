package config

import (
	"fmt"
)

type Server struct {
	IP                      string `yaml:"ip"                        default:"0.0.0.0"`
	Port                    int    `yaml:"port"                      default:"8080"`
	WriteTimeout            int    `yaml:"write_timeout"             default:"15"`
	ReadTimeout             int    `yaml:"read_timeout"              default:"15"`
	ReadHeaderTimeout       int    `yaml:"read_header_timeout"       default:"15"`
	GracefulShutdownTimeout int    `yaml:"graceful_shutdown_timeout" default:"60"`
}

type RedisDatabase struct {
	Addr     string `yaml:"addr" default:"redis:6379"`
	Password string `yaml:"-"`
	DB       int    `yaml:"db"   default:"0"`
}

type PostgresDatabase struct {
	IP   string `yaml:"ip"   default:"postgres"`
	Port int    `yaml:"port" default:"5432"`
	User string `yaml:"-"`
	Pass string `yaml:"-"`
}

type Config struct {
	HTTP struct {
		CORSAllowedOrigins string `yaml:"cors_allowed_origins" default:"http://localhost:8000"`
		SecureCookies      bool   `yaml:"secure_cookies"       default:"false"`
		Server             Server `yaml:"server"`
	} `yaml:"http"`
	Microservices struct {
		Auth struct {
			Addr                 string `yaml:"auth_addr"          default:"auth:8081"`
			HTTPSessionAliveTime int    `yaml:"session_alive_time" default:"86400"`
		} `yaml:"auth_service"`
		Static struct {
			Addr        string `yaml:"static_addr"   default:"static:8082"`
			MaxFileSize int    `yaml:"max_file_size" default:"10485760"`
		} `yaml:"static_service"`
		ProfanityFilter struct {
			Addr string `yaml:"profanity_filter_addr" default:"profanity:8050"`
		} `yaml:"profanity_filter_service"`
	} `yaml:"microservices"`
	ContentSecretKey string           `yaml:"-"`
	Postgres         PostgresDatabase `yaml:"postgres"`
}

type AuthConfig struct {
	IP               string        `yaml:"ip"                 default:"0.0.0.0"`
	Port             int           `yaml:"port"               default:"8081"`
	SessionAliveTime int           `yaml:"session_alive_time" default:"86400"`
	Redis            RedisDatabase `yaml:"redis"`
}

type StaticConfig struct {
	IP          string `yaml:"ip"            default:"0.0.0.0"`
	Port        int    `yaml:"port"          default:"8082"`
	MaxFileSize int    `yaml:"max_file_size" default:"10485760"`
	S3          struct {
		AccessKeyID     string `yaml:"-"`
		SecretAccessKey string `yaml:"-"`
		Region          string `yaml:"region"      default:"ru-msk"`
		Endpoint        string `yaml:"endpoint"    default:"https://hb.vkcs.cloud"`
		BucketName      string `yaml:"bucket_name" default:"kinoskop_dev"`
	} `yaml:"s3"`
	Postgres PostgresDatabase `yaml:"postgres"`
}

func (cfg *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", cfg.HTTP.Server.IP, cfg.HTTP.Server.Port)
}

func (cfg *PostgresDatabase) GetConnectURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/kinoskop?sslmode=disable",
		cfg.User, cfg.Pass, cfg.IP, cfg.Port)
}
