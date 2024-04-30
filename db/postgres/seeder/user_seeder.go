package seeder

import (
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
		hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(user.Password.String), bcrypt.DefaultCost)
		if hashedPasswordErr != nil {
			panic(hashedPasswordErr)
		}
		password := null.NewString(string(hashedPassword), true)
		begin, beginErr := userSeeder.DatabaseConfig.CockroachdbDatabase.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			_, err = begin.Query(
				"INSERT INTO \"user\" (id, name, username, email, password, avatar_url, bio, is_verified, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);",
				user.Id,
				user.Name,
				user.Username,
				user.Email,
				password,
				user.AvatarUrl,
				user.Bio,
				user.IsVerified,
				user.CreatedAt,
				user.UpdatedAt,
				user.DeletedAt,
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
	}
}

func (userSeeder *UserSeeder) Down() {
	for _, user := range userSeeder.UserMock.Data {
		begin, beginErr := userSeeder.DatabaseConfig.CockroachdbDatabase.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			_, err = begin.Query(
				"DELETE FROM \"user\" WHERE id = $1 LIMIT 1;",
				user.Id,
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
	}
}
