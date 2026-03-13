package incus

import "example.com/m/internal/infrastructure/env"

type Config struct {
	SocketPath string
}

func NewConfig() *Config {
	return &Config{
		SocketPath: env.GetString("INCUS_SOCKET_PATH", ""),
	}
}