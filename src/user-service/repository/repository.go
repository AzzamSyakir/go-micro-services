package repository

import (
	"database/sql"
	"go-micro-services/src/user-service/entity"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	userRepository := &UserRepository{}
	return userRepository
}
func DeserializeUserRows(rows *sql.Rows) []*entity.User {
	var foundUsers []*entity.User
	for rows.Next() {
		foundUser := &entity.User{}
		scanErr := rows.Scan(
			&foundUser.Id,
			&foundUser.Name,
			&foundUser.Balance,
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

func (userRepository *UserRepository) GetOneById(begin *sql.Tx, id string) (result *entity.User, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = begin.Query(
		`SELECT id, name, balance, created_at, updated_at, deleted_at FROM "users" WHERE id=$1 LIMIT 1;`,
		id,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err
	}
	defer rows.Close()

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

func (userRepository *UserRepository) PatchOneById(begin *sql.Tx, id string, toPatchUser *entity.User) (result *entity.User, err error) {
	rows, queryErr := begin.Query(
		`UPDATE "users" SET id=$1, name=$2,  balance=$3, created_at=$4, updated_at=$5, deleted_at=$6 WHERE id = $7 ;`,
		toPatchUser.Id,
		toPatchUser.Name,
		toPatchUser.Balance,
		toPatchUser.CreatedAt,
		toPatchUser.UpdatedAt,
		toPatchUser.DeletedAt,
		id,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return
	}
	defer rows.Close()

	result = toPatchUser
	err = nil
	return result, err
}
