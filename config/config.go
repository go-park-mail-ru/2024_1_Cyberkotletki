package config

import (
	"fmt"
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
	Addr     string `yaml:"addr"     default:"host.docker.internal:6379"`
	Password string `yaml:"password" default:""`
	DB       int    `yaml:"db"       default:"0"`
}

type PostgresDatabase struct {
	// default исключительно для примера
	// nolint
	ConnectURL string `yaml:"connect_url" default:"postgres://kinoskop_admin:admin_secret_password@localhost:5432/kinoskop?sslmode=disable"`
}

type Config struct {
	HTTP struct {
		CORSAllowedOrigins string `yaml:"cors_allowed_origins" default:"http://localhost:8000"`
		SecureCookies      bool   `yaml:"secure_cookies"       default:"false"`
		Server             Server `yaml:"server"`
	} `yaml:"http"`
	Microservices struct {
		Auth struct {
			Addr                 string `yaml:"auth_addr"          default:"localhost:8081"`
			HTTPSessionAliveTime int    `yaml:"session_alive_time" default:"86400"`
		} `yaml:"auth_service"`
		Static struct {
			Addr        string `yaml:"static_addr"   default:"localhost:8082"`
			MaxFileSize int    `yaml:"max_file_size" default:"10485760"`
		} `yaml:"static_service"`
	} `yaml:"microservices"`
	Postgres PostgresDatabase `yaml:"postgres"`
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
		AccessKeyID     string `yaml:"access_key_id"`
		SecretAccessKey string `yaml:"secret_access_key"`
		Region          string `yaml:"region"            default:"ru-msk"`
		Endpoint        string `yaml:"endpoint"          default:"https://hb.vkcs.cloud"`
		BucketName      string `yaml:"bucket_name"       default:"kinoskop_dev"`
	} `yaml:"s3"`
	Postgres PostgresDatabase `yaml:"postgres"`
}

func (cfg *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", cfg.HTTP.Server.IP, cfg.HTTP.Server.Port)
}
