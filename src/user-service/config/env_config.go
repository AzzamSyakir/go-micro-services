package config

import (
	"os"
)

type AppEnv struct {
	Host         string
	Port         string
	AuthHost     string
	AuthGrpcPort string
}

type PostgresEnv struct {
	Host     string
	Port     string
	User     string
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
			Host:         os.Getenv("AUTH_HOST"),
			Port:         os.Getenv("USER_PORT"),
			AuthHost:     os.Getenv("AUTH_HOST"),
			AuthGrpcPort: os.Getenv("AUTH_GRPC_PORT"),
		},
		UserDB: &PostgresEnv{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_USER_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: "user_db",
		},
	}
	return envConfig
}
