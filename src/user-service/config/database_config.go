package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	UserDB *PostgresDatabase
}

type PostgresDatabase struct {
	Connection *sql.DB
}

func NewUserDBConfig(envConfig *EnvConfig) *DatabaseConfig {
	databaseConfig := &DatabaseConfig{
		UserDB: NewUserDB(envConfig),
	}
	return databaseConfig
}

func NewUserDB(envConfig *EnvConfig) *PostgresDatabase {
	var url string
	if envConfig.UserDB.Password == "" {
		url = fmt.Sprintf(
			"postgresql://%s@%s:%s/%s",
			envConfig.UserDB.User,
			envConfig.UserDB.Host,
			envConfig.UserDB.Port,
			envConfig.UserDB.Database,
		)
	} else {
		url = fmt.Sprintf(
			"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
			envConfig.UserDB.User,
			envConfig.UserDB.Password,
			envConfig.UserDB.Host,
			envConfig.UserDB.Port,
			envConfig.UserDB.Database,
		)
	}

	connection, err := sql.Open("postgres", url)
	if err != nil {
		panic(err)
	}
	connection.SetConnMaxLifetime(300 * time.Second)
	connection.SetMaxIdleConns(10)
	connection.SetMaxOpenConns(10)

	userDB := &PostgresDatabase{
		Connection: connection,
	}
	return userDB
}
