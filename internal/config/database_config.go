package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	UserDB    *PostgresDatabase
	ProductDB *PostgresDatabase
	OrderDB   *PostgresDatabase
}

type PostgresDatabase struct {
	Connection *sql.DB
}

func NewProductDBConfig(envConfig *EnvConfig) *DatabaseConfig {
	databaseConfig := &DatabaseConfig{
		ProductDB: NewProductDB(envConfig),
	}
	return databaseConfig
}
func NewUserDBConfig(envConfig *EnvConfig) *DatabaseConfig {
	databaseConfig := &DatabaseConfig{
		UserDB: NewUserDB(envConfig),
	}
	return databaseConfig
}
func NewOrderDBConfig(envConfig *EnvConfig) *DatabaseConfig {
	databaseConfig := &DatabaseConfig{
		OrderDB: NewOrderDB(envConfig),
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

	userDB := &PostgresDatabase{
		Connection: connection,
	}
	return userDB
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

	productDB := &PostgresDatabase{
		Connection: connection,
	}
	return productDB
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
