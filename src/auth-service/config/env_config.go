package config

import (
	"os"
)

type AppEnv struct {
	Host        string
	AuthPort    string
	OrderPort   string
	UserPort    string
	ProductPort string
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
			Host:        os.Getenv("GATEWAY_APP_HOST"),
			AuthPort:    os.Getenv("AUTH_SERVICES_PORT"),
			OrderPort:   os.Getenv("ORDER_SERVICES_PORT"),
			UserPort:    os.Getenv("USER_SERVICES_PORT"),
			ProductPort: os.Getenv("PRODUCT_SERVICES_PORT"),
		},
		AuthDB: &PostgresEnv{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_AUTH_PORT"),
			Auth:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: "order_db",
		},
	}
	return envConfig
}
