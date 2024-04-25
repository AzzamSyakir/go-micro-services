package config

import (
	"os"
)

type AppEnv struct {
	Host string
	Port string
}

type PostgresEnv struct {
	Host     string
	Port     string
	Auth     string
	Password string
	Database string
}

type EnvConfig struct {
	App    *AppEnv
	UserDB *PostgresEnv
}

func NewEnvConfig() *EnvConfig {
	envConfig := &EnvConfig{
		App: &AppEnv{
			Host: os.Getenv("GATEWAY_APP_HOST"),
			Port: os.Getenv("AUTH_SERVICES_PORT"),
		},
		UserDB: &PostgresEnv{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_AUTH_PORT"),
			Auth:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: "order_db",
		},
	}
	return envConfig
}
