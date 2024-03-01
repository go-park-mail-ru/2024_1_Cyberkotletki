package config

type ServerMode int

const (
	DeployMode = iota
	TestMode
	DevMode
)

type InitParams struct {
	Addr string
	Mode ServerMode
}

// todo переименовать файл?
