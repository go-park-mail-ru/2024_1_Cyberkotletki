package config

type ServerMode int

const (
	DeployMode = iota
	TestMode
	DevMode
)

type InitParams struct {
	Addr       string
	Mode       ServerMode
	GenSwagger bool
}

// todo переименовать файл?
