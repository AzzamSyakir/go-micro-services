package config

import (
	"database/sql"
	"fmt"
	"os"
)

type DatabaseConfig struct {
	PostgresDatabase *PostgresDatabase
}

type PostgresDatabase struct {
	Connection *sql.DB
}

func InitUserDb() *PostgresDatabase {

	sqlInfo := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
	)

	connection, err := sql.Open("mysql", sqlInfo)
	if err != nil {
		panic(err)
	}
	PostgresDatabase := &PostgresDatabase{
		Connection: connection,
	}
	return PostgresDatabase
}
func InitProductDb() *PostgresDatabase {

	sqlInfo := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_PRODUCT"),
	)

	connection, err := sql.Open("mysql", sqlInfo)
	if err != nil {
		panic(err)
	}

	PostgresDatabase := &PostgresDatabase{
		Connection: connection,
	}
	return PostgresDatabase
}
