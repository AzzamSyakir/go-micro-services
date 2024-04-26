package repository

import (
	"database/sql"
	"go-micro-services/src/auth-service/entity"
)

type AuthRepository struct {
}

func NewAuthRepository() *AuthRepository {
	AuthRepository := &AuthRepository{}
	return AuthRepository
}

func (AuthRepository *AuthRepository) CreateSession(begin *sql.Tx, toCreateSession *entity.Session) (result *entity.Session, err error) {
	_, queryErr := begin.Query(
		`INSERT INTO "users" (id, name, email, password, balance, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		toCreateSession.Id,
		toCreateSession.Name,
		toCreateSession.Email,
		toCreateSession.Password,
		toCreateSession.Balance,
		toCreateSession.CreatedAt,
		toCreateSession.UpdatedAt,
		toCreateSession.DeletedAt,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return
	}

	result = toCreateSession
	err = nil
	return result, err
}
