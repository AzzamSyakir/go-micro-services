package config

import (
	"os"
)

type AppEnv struct {
	Host         string
	UserHost     string
	ProductHost  string
	OrderHost    string
	AuthHttpPort string
	OrderPort    string
	UserPort     string
	ProductPort  string
	AuthGrpcPort string
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
	AuthDB *PostgresEnv
}

func NewEnvConfig() *EnvConfig {
	envConfig := &EnvConfig{
		App: &AppEnv{
			Host:         os.Getenv("AUTH_HOST"),
			UserHost:     os.Getenv("USER_HOST"),
			ProductHost:  os.Getenv("PRODUCT_HOST"),
			OrderHost:    os.Getenv("ORDER_HOST"),
			AuthHttpPort: os.Getenv("AUTH_PORT"),
			AuthGrpcPort: os.Getenv("AUTH_GRPC_PORT"),
			OrderPort:    os.Getenv("ORDER_PORT"),
			UserPort:     os.Getenv("USER_PORT"),
			ProductPort:  os.Getenv("PRODUCT_PORT"),
		},
		AuthDB: &PostgresEnv{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_AUTH_PORT"),
			Auth:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: "auth_db",
		},
	}
	return envConfig
}
