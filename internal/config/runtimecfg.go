package config

import "time"

type ServerMode string

const (
	DeployMode = "prod"
	TestMode   = "test"
	DevMode    = "dev"
)

type InitParams struct {
	// Addr - адрес, на котором работает сервер
	Addr string
	// Mode - режим работы. prod/test/dev
	Mode ServerMode
	// GenSwagger - генерировать ли swagger документацию к api
	GenSwagger bool
	// StaticFolder - путь к папке со статикой (если она не обсуживается через nginx)
	StaticFolder string
	// CORS - адрес в формате {protocol}://{address} для cors заголовка в api-запросах
	CORS string
	// SessionAliveTime - время жизни сессии в секундах
	SessionAliveTime time.Duration
	// CookiesSecure - флаг Secure у Cookies
	CookiesSecure bool
	// GracefulShutdownTime - время в секундах для выполнения graceful shutdown, прежде чем сервер насильно выключится
	GracefulShutdownTime time.Duration
}
