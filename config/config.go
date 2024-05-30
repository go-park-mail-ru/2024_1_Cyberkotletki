package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"github.com/mcuadros/go-defaults"
	"gopkg.in/yaml.v3"
	"os"
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

func GenerateExampleConfigs() {
	type cfg struct {
		config   any
		filename string
	}

	configs := []cfg{
		{Config{}, "config.yaml"},
		{StaticConfig{}, "config_static.yaml"},
		{AuthConfig{}, "config_auth.yaml"},
	}
	for _, conf := range configs {
		defaults.SetDefaults(conf.config)
		yamlData, err := yaml.Marshal(conf.config)
		if err != nil {
			fmt.Printf("Ошибка при маршализации YAML: %v\n", err)
			return
		}
		file, err := os.Create("config.yaml")
		if err != nil {
			fmt.Printf("Ошибка при создании файла: %v\n", err)
			return
		}
		_, err = file.Write(yamlData)
		if err != nil {
			fmt.Printf("Ошибка при записи в файл: %v\n", err)
			return
		}
		err = file.Close()
		if err != nil {
			fmt.Printf("Ошибка при записи в файл: %v\n", err)
			return
		}

		fmt.Println("Конфигурационный файл успешно создан.")
	}
}

func ParseCoreServiceParams() Config {
	var cfg Config
	// читаем конфиг
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка при чтении конфига сервера: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Ошибка при парсинге конфига сервера: %v", err)
	}
	// читаем переменные окружения
	err = godotenv.Load()
	if err != nil {
		fmt.Println(".env файл отсутствует, переменные будут загружены из переменных окружения")
	}
	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Pass = os.Getenv("POSTGRES_PASSWORD")
	cfg.ContentSecretKey = os.Getenv("CONTENT_SECRET_KEY")
	return cfg
}

func ParseAuthServiceParams() AuthConfig {
	var cfg AuthConfig
	// читаем конфиг
	yamlFile, err := os.ReadFile("config_auth.yaml")
	if err != nil {
		log.Fatalf("Ошибка при чтении конфига сервера: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Ошибка при парсинге конфига сервера: %v", err)
	}
	// читаем переменные окружения
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки .env файла")
	}
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	return cfg
}

func ParseStaticServiceParams() StaticConfig {
	var cfg StaticConfig
	// читаем конфиг
	yamlFile, err := os.ReadFile("config_static.yaml")
	if err != nil {
		log.Fatalf("Ошибка при чтении конфига сервера: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Ошибка при парсинге конфига сервера: %v", err)
	}
	// читаем переменные окружения
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Ошибка загрузки .env файла")
	}
	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Pass = os.Getenv("POSTGRES_PASSWORD")
	cfg.S3.AccessKeyID = os.Getenv("S3_ACCESS_KEY_ID")
	cfg.S3.SecretAccessKey = os.Getenv("S3_SECRET_ACCESS_KEY")
	return cfg
}
