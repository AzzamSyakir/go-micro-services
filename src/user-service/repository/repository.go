package repository

import (
	"database/sql"
	"go-micro-services/src/user-service/delivery/grpc/pb"

	"github.com/guregu/null"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	userRepository := &UserRepository{}
	return userRepository
}
func DeserializeUserRows(rows *sql.Rows) []*pb.User {
	var foundUsers []*pb.User
	for rows.Next() {
		foundUser := &pb.User{}
		var createdAt, updatedAt, deletedAt null.Time
		scanErr := rows.Scan(
			&foundUser.Id,
			&foundUser.Name,
			&foundUser.Email,
			&foundUser.Password,
			&foundUser.Balance,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		foundUser.CreatedAt = timestamppb.New(createdAt.Time)
		foundUser.UpdatedAt = timestamppb.New(updatedAt.Time)
		foundUser.DeletedAt = timestamppb.New(deletedAt.Time)
		if scanErr != nil {
			panic(scanErr)
		}
		foundUsers = append(foundUsers, foundUser)
	}
	return foundUsers
}

func (userRepository *UserRepository) CreateUser(begin *sql.Tx, toCreateUser *pb.User) (result *pb.User, err error) {
	_, queryErr := begin.Query(
		`INSERT INTO "users" (id, name, email, password, balance, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		toCreateUser.Id,
		toCreateUser.Name,
		toCreateUser.Email,
		toCreateUser.Password,
		toCreateUser.Balance,
		toCreateUser.CreatedAt.AsTime(),
		toCreateUser.UpdatedAt.AsTime(),
		toCreateUser.DeletedAt.AsTime(),
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

func (userRepository *UserRepository) ListUser(begin *sql.Tx) (result *pb.UserResponseRepeated, err error) {
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
	var ListUsers []*pb.User
	var createdAt, updatedAt, deletedAt null.Time
	for rows.Next() {
		ListUser := &pb.User{}
		scanErr := rows.Scan(
			&ListUser.Id,
			&ListUser.Name,
			&ListUser.Email,
			&ListUser.Password,
			&ListUser.Balance,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		ListUser.CreatedAt = timestamppb.New(createdAt.Time)
		ListUser.UpdatedAt = timestamppb.New(updatedAt.Time)
		ListUser.DeletedAt = timestamppb.New(deletedAt.Time)
		if scanErr != nil {
			result = nil
			err = scanErr
			return result, err
		}
		ListUsers = append(ListUsers, ListUser)
	}

	result = &pb.UserResponseRepeated{
		Data: ListUsers,
	}
	err = nil
	return result, err
}

func (userRepository *UserRepository) GetUserById(begin *sql.Tx, id string) (result *pb.User, err error) {
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

func (userRepository *UserRepository) GetUserByEmail(begin *sql.Tx, email string) (result *pb.User, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = begin.Query(
		`SELECT id, name, email, password, balance, created_at, updated_at, deleted_at FROM "users" WHERE email=$1 LIMIT 1;`,
		email,
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

func (userRepository *UserRepository) PatchOneById(begin *sql.Tx, id string, toPatchUser *pb.User) (result *pb.User, err error) {
	rows, queryErr := begin.Query(
		`UPDATE "users" SET id=$1, name=$2, email=$3, password=$4, balance=$5, created_at=$6, updated_at=$7, deleted_at=$8 WHERE id = $9 ;`,
		toPatchUser.Id,
		toPatchUser.Name,
		toPatchUser.Email,
		toPatchUser.Password,
		toPatchUser.Balance,
		toPatchUser.CreatedAt.AsTime(),
		toPatchUser.UpdatedAt.AsTime(),
		toPatchUser.DeletedAt.AsTime(),
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

func (userRepository *UserRepository) DeleteUser(begin *sql.Tx, id string) (result *pb.User, err error) {
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
