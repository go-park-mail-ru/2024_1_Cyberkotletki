package config

type ServerMode string

const (
	DeployMode = "dev"
	TestMode   = "test"
	DevMode    = "prod"
)

type InitParams struct {
	Addr         string
	Mode         ServerMode
	GenSwagger   bool
	StaticFolder string
	CORS         string
}
