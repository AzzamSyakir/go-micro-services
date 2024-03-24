package repository

import (
	"database/sql"
	"go-micro-services/internal/entity"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	userRepository := &UserRepository{}
	return userRepository
}
func (userRepository *UserRepository) GetOneById(begin *sql.Tx, id string) (result *entity.User, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = begin.Query(
		`SELECT id, name, saldo, created_at, updated_at, deleted_at FROM "users" WHERE id=$1 LIMIT 1;`,
		id,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err
	}

	foundUsers := DeserializeUserRows(rows)
	if len(foundUsers) == 0 {
		result = nil
		err = nil
		return result, err
	}

	result = foundUsers[0]
	err = nil
	return result, err
}
func DeserializeUserRows(rows *sql.Rows) []*entity.User {
	var foundUsers []*entity.User
	for rows.Next() {
		foundUser := &entity.User{}
		scanErr := rows.Scan(
			&foundUser.Id,
			&foundUser.Name,
			&foundUser.Saldo,
			&foundUser.CreatedAt,
			&foundUser.UpdatedAt,
			&foundUser.DeletedAt,
		)
		if scanErr != nil {
			panic(scanErr)
		}
		foundUsers = append(foundUsers, foundUser)
	}
	return foundUsers
}
