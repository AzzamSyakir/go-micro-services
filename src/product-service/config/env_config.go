package config

import "os"

type AppEnv struct {
	Host string
	Port string
}

type PostgresEnv struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type EnvConfig struct {
	App       *AppEnv
	ProductDB *PostgresEnv
}

func NewEnvConfig() *EnvConfig {
	envConfig := &EnvConfig{
		App: &AppEnv{
			Host: os.Getenv("GATEWAY_APP_HOST"),
			Port: os.Getenv("PRODUCT_SERVICES_PORT"),
		},
		ProductDB: &PostgresEnv{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PRODUCT_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: "product_db",
		},
	}
	return envConfig
}
