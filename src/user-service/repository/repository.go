package repository

import (
	"database/sql"
	"fmt"
	"go-micro-services/src/user-service/entity"
	model_response "go-micro-services/src/user-service/model/response"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	userRepository := &UserRepository{}
	return userRepository
}

func (userRepository *UserRepository) CreateUser(begin *sql.Tx, toCreateUser *entity.User) (result *entity.User, err error) {
	_, queryErr := begin.Query(
		`INSERT INTO "users" (id, name, email, password, balance, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		toCreateUser.Id,
		toCreateUser.Name,
		toCreateUser.Email,
		toCreateUser.Password,
		toCreateUser.Balance,
		toCreateUser.CreatedAt,
		toCreateUser.UpdatedAt,
		toCreateUser.DeletedAt,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return
	}

	result = toCreateUser
	err = nil
	return result, err

}

func (userRepository *UserRepository) FetchUser(begin *sql.Tx) (result *model_response.Response[[]*entity.User], err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = begin.Query(
		`SELECT id, name, email, password, balance, created_at, updated_at, deleted_at FROM "users" `,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err

	}
	defer rows.Close()
	var fetchUsers []*entity.User
	for rows.Next() {
		fetchUser := &entity.User{}
		scanErr := rows.Scan(
			&fetchUser.Id,
			&fetchUser.Name,
			&fetchUser.Email,
			&fetchUser.Password,
			&fetchUser.Balance,
			&fetchUser.CreatedAt,
			&fetchUser.UpdatedAt,
			&fetchUser.DeletedAt,
		)
		if scanErr != nil {
			result = nil
			err = scanErr
			return result, err
		}
		fetchUsers = append(fetchUsers, fetchUser)
	}

	result = &model_response.Response[[]*entity.User]{
		Data: fetchUsers,
	}
	fmt.Println("rows fetchUser", rows)
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
			&foundUser.Email,
			&foundUser.Password,
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
		`SELECT id, name, email, password, balance, created_at, updated_at, deleted_at FROM "users" WHERE id=$1 LIMIT 1;`,
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
		`UPDATE "users" SET id=$1, name=$2, email=$3, password=$4, balance=$5, created_at=$6, updated_at=$7, deleted_at=$8 WHERE id = $9 ;`,
		toPatchUser.Id,
		toPatchUser.Name,
		toPatchUser.Email,
		toPatchUser.Password,
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

func (userRepository *UserRepository) DeleteUser(begin *sql.Tx, id string) (result *entity.User, err error) {
	rows, queryErr := begin.Query(
		`DELETE FROM "users" WHERE id=$1 RETURNING id, name,  email, password, balance, created_at, updated_at, deleted_at`,
		id,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return
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
