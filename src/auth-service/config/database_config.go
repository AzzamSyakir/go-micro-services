package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	AuthDB *PostgresDatabase
}

type PostgresDatabase struct {
	Connection *sql.DB
}

func NewAuthDBConfig(envConfig *EnvConfig) *DatabaseConfig {
	databaseConfig := &DatabaseConfig{
		AuthDB: NewAuthDB(envConfig),
	}
	return databaseConfig
}

func NewAuthDB(envConfig *EnvConfig) *PostgresDatabase {
	var url string
	if envConfig.AuthDB.Password == "" {
		url = fmt.Sprintf(
			"postgresql://%s@%s:%s/%s",
			envConfig.AuthDB.Auth,
			envConfig.AuthDB.Host,
			envConfig.AuthDB.Port,
			envConfig.AuthDB.Database,
		)
	} else {
		url = fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			envConfig.AuthDB.Auth,
			envConfig.AuthDB.Password,
			envConfig.AuthDB.Host,
			envConfig.AuthDB.Port,
			envConfig.AuthDB.Database,
		)
	}

	connection, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	connection.SetConnMaxIdleTime(10 * time.Second)
	connection.SetConnMaxLifetime(30 * time.Second)
	connection.SetMaxOpenConns(500)
	userDB := &PostgresDatabase{
		Connection: connection,
	}
	return userDB
}
