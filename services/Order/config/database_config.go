package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	OrderDB *PostgresDatabase
}

type PostgresDatabase struct {
	Connection *sql.DB
}

func NewDBConfig(envConfig *EnvConfig) *DatabaseConfig {
	databaseConfig := &DatabaseConfig{
		OrderDB: NewOrderDB(envConfig),
	}
	return databaseConfig
}

func NewOrderDB(envConfig *EnvConfig) *PostgresDatabase {
	var url string
	if envConfig.OrderDB.Password == "" {
		url = fmt.Sprintf(
			"postgresql://%s@%s:%s/%s",
			envConfig.OrderDB.User,
			envConfig.OrderDB.Host,
			envConfig.OrderDB.Port,
			envConfig.OrderDB.Database,
		)
	} else {
		url = fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			envConfig.OrderDB.User,
			envConfig.OrderDB.Password,
			envConfig.OrderDB.Host,
			envConfig.OrderDB.Port,
			envConfig.OrderDB.Database,
		)
	}

	connection, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}

	orderDB := &PostgresDatabase{
		Connection: connection,
	}
	return orderDB
}
