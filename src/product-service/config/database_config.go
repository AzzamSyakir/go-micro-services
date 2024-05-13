package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	ProductDB *PostgresDatabase
}

type PostgresDatabase struct {
	Connection *sql.DB
}

func NewDBConfig(envConfig *EnvConfig) *DatabaseConfig {
	databaseConfig := &DatabaseConfig{
		ProductDB: NewProductDB(envConfig),
	}
	return databaseConfig
}

func NewProductDB(envConfig *EnvConfig) *PostgresDatabase {
	var url string
	if envConfig.ProductDB.Password == "" {
		url = fmt.Sprintf(
			"postgresql://%s@%s:%s/%s",
			envConfig.ProductDB.User,
			envConfig.ProductDB.Host,
			envConfig.ProductDB.Port,
			envConfig.ProductDB.Database,
		)
	} else {
		url = fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			envConfig.ProductDB.User,
			envConfig.ProductDB.Password,
			envConfig.ProductDB.Host,
			envConfig.ProductDB.Port,
			envConfig.ProductDB.Database,
		)
	}

	connection, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	connection.SetConnMaxLifetime(300)
	connection.SetMaxIdleConns(10)
	connection.SetMaxOpenConns(10)
	productDB := &PostgresDatabase{
		Connection: connection,
	}
	return productDB
}
