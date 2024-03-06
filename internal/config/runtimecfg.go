package config

type ServerMode string

const (
	DeployMode = "prod"
	TestMode   = "test"
	DevMode    = "dev"
)

type InitParams struct {
	Addr         string
	Mode         ServerMode
	GenSwagger   bool
	StaticFolder string
	CORS         string
}
