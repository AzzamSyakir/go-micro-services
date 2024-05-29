package seeder

import (
	"database/sql"
	"go-micro-services/src/auth-service/test/mock"
	"go-micro-services/src/user-service/config"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
)

type UserSeeder struct {
	DatabaseConfig *config.DatabaseConfig
	UserMock       *mock.UserMock
}

func NewUserSeeder(
	databaseConfig *config.DatabaseConfig,
) *UserSeeder {
	userSeeder := &UserSeeder{
		DatabaseConfig: databaseConfig,
		UserMock:       mock.NewUserMock(),
	}
	return userSeeder
}

func (userSeeder *UserSeeder) Up() {
	for _, user := range userSeeder.UserMock.Data {
		var rows *sql.Rows
		hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(user.Password.String), bcrypt.DefaultCost)
		if hashedPasswordErr != nil {
			panic(hashedPasswordErr)
		}
		password := null.NewString(string(hashedPassword), true)
		begin, beginErr := userSeeder.DatabaseConfig.UserDB.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			rows, err = begin.Query(
				"INSERT INTO users (id, name, balance, email, password, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $5);",
				user.Id,
				user.Name,
				user.Balance,
				user.Email,
				password,
				user.CreatedAt,
				user.UpdatedAt,
			)
			return err
		})
		if queryErr != nil {
			panic(queryErr)
		}

		commitErr := crdb.Execute(func() (err error) {
			err = begin.Commit()
			return err
		})
		if commitErr != nil {
			panic(commitErr)
		}
		defer rows.Close()
	}
}

func (userSeeder *UserSeeder) Down() {
	for _, user := range userSeeder.UserMock.Data {
		var query *sql.Rows

		begin, beginErr := userSeeder.DatabaseConfig.UserDB.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			query, err = begin.Query(
				"DELETE FROM users WHERE id = $1",
				user.Id,
			)
			return err
		})
		if queryErr != nil {
			panic(queryErr)
		}
		defer query.Close()
		commitErr := crdb.Execute(func() (err error) {
			err = begin.Commit()
			return err
		})
		if commitErr != nil {
			panic(commitErr)
		}
	}
}
