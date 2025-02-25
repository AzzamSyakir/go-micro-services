package config

import "os"

type AppEnv struct {
	Host        string
	UserHost    string
	ProductHost string
	OrderPort   string
	UserPort    string
	ProductPort string
}

type PostgresEnv struct {
	Host      string
	OrderPort string
	User      string
	Password  string
	Database  string
}

type EnvConfig struct {
	App     *AppEnv
	OrderDB *PostgresEnv
}

func NewEnvConfig() *EnvConfig {
	envConfig := &EnvConfig{
		App: &AppEnv{
			Host:        os.Getenv("AUTH_HOST"),
			UserHost:    os.Getenv("USER_HOST"),
			ProductHost: os.Getenv("PRODUCT_HOST"),
			OrderPort:   os.Getenv("ORDER_PORT"),
			UserPort:    os.Getenv("USER_PORT"),
			ProductPort: os.Getenv("PRODUCT_PORT"),
		},
		OrderDB: &PostgresEnv{
			Host:      os.Getenv("POSTGRES_HOST"),
			OrderPort: os.Getenv("POSTGRES_ORDER_PORT"),
			User:      os.Getenv("POSTGRES_USER"),
			Password:  os.Getenv("POSTGRES_PASSWORD"),
			Database:  "order_db",
		},
	}
	return envConfig
}
