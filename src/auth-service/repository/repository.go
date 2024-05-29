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

func DeserializeSessionRows(rows *sql.Rows) []*entity.Session {
	var foundSessions []*entity.Session
	for rows.Next() {
		foundSession := &entity.Session{}
		scanErr := rows.Scan(
			&foundSession.Id,
			&foundSession.UserId,
			&foundSession.AccessToken,
			&foundSession.RefreshToken,
			&foundSession.AccessTokenExpiredAt,
			&foundSession.RefreshTokenExpiredAt,
			&foundSession.CreatedAt,
			&foundSession.UpdatedAt,
		)
		if scanErr != nil {
			panic(scanErr)
		}
		foundSessions = append(foundSessions, foundSession)
	}
	return foundSessions
}
func (sessionRepository *AuthRepository) CreateSession(begin *sql.Tx, toCreateSession *entity.Session) (result *entity.Session, err error) {
	_, queryErr := begin.Query(
		`INSERT INTO sessions (id, user_id, access_token, refresh_token, access_token_expired_at, refresh_token_expired_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`,
		toCreateSession.Id,
		toCreateSession.UserId,
		toCreateSession.AccessToken,
		toCreateSession.RefreshToken,
		toCreateSession.AccessTokenExpiredAt,
		toCreateSession.RefreshTokenExpiredAt,
		toCreateSession.CreatedAt,
		toCreateSession.UpdatedAt,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err
	}

	result = toCreateSession
	err = nil
	return result, err
}

func (sessionRepository *AuthRepository) FindOneByAccToken(begin *sql.Tx, accessToken string) (result *entity.Session, err error) {
	rows, queryErr := begin.Query(
		`SELECT id, user_id, access_token, refresh_token, access_token_expired_at, refresh_token_expired_at, created_at, updated_at FROM sessions WHERE access_token=$1 LIMIT 1;`,
		accessToken,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err
	}

	foundSessions := DeserializeSessionRows(rows)
	if len(foundSessions) == 0 {
		result = nil
		err = nil
		return result, err
	}

	result = foundSessions[0]
	err = nil
	return result, err
}

func (sessionRepository *AuthRepository) GetOneByUserId(begin *sql.Tx, userId string) (result *entity.Session, err error) {
	rows, queryErr := begin.Query(
		`SELECT id, user_id, access_token, refresh_token, access_token_expired_at, refresh_token_expired_at, created_at, updated_at FROM sessions WHERE user_id=$1 LIMIT 1;`,
		userId,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err
	}

	foundSessions := DeserializeSessionRows(rows)
	if len(foundSessions) == 0 {
		result = nil
		err = nil
		return result, err
	}

	result = foundSessions[0]
	err = nil
	return result, err
}

func (sessionRepository *AuthRepository) FindOneByRefToken(begin *sql.Tx, refreshToken string) (result *entity.Session, err error) {
	rows, queryErr := begin.Query(
		`SELECT id, user_id, access_token, refresh_token, access_token_expired_at, refresh_token_expired_at, created_at, updated_at FROM sessions WHERE refresh_token=$1 LIMIT 1;`,
		refreshToken,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err
	}

	foundSessions := DeserializeSessionRows(rows)
	if len(foundSessions) == 0 {
		result = nil
		err = nil
		return result, err
	}

	result = foundSessions[0]
	err = nil
	return result, err
}

func (sessionRepository *AuthRepository) PatchOneById(begin *sql.Tx, id string, toPatchSession *entity.Session) (result *entity.Session, err error) {
	_, queryErr := begin.Query(
		`UPDATE sessions SET id=$1, user_id=$2, access_token=$3, refresh_token=$4, access_token_expired_at=$5, refresh_token_expired_at=$6, created_at=$7, updated_at=$8 WHERE id=$9;`,
		toPatchSession.Id,
		toPatchSession.UserId,
		toPatchSession.AccessToken,
		toPatchSession.RefreshToken,
		toPatchSession.AccessTokenExpiredAt,
		toPatchSession.RefreshTokenExpiredAt,
		toPatchSession.CreatedAt,
		toPatchSession.UpdatedAt,
		id,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err
	}

	result = toPatchSession
	err = nil
	return result, err
}
func (sessionRepository *AuthRepository) DeleteOneById(begin *sql.Tx, id string) (result *entity.Session, err error) {
	rows, queryErr := begin.Query(
		`DELETE FROM sessions WHERE id=$1  RETURNING id, user_id, access_token, refresh_token, access_token_expired_at, refresh_token_expired_at, created_at, updated_at;`,
		id,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return
	}

	foundSessions := DeserializeSessionRows(rows)
	if len(foundSessions) == 0 {
		result = nil
		err = nil
		return result, err
	}
	result = foundSessions[0]
	err = nil
	return result, err
}
