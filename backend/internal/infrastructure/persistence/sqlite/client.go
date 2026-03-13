package sqlite

import "example.com/m/internal/infrastructure/env"

type Config struct {
	DBPath string
}

func NewConfig() *Config {
	path := env.GetString("DB_PATH", ":memory:")
	return &Config{
		DBPath: path,
	}
}