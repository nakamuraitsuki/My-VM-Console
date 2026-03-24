package sshca

import "example.com/m/internal/infrastructure/env"

type Config struct {
	privateKeyPath string
	jumpUser       string
}

func NewConfig() *Config {
	return &Config{
		privateKeyPath: env.GetString("CA_PRIVATE_KEY_PATH", ""),
		jumpUser:       env.GetString("CA_JUMP_USER", "jumpuser"),
	}
}
