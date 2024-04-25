package middleware

import (
	"go-micro-services/src/auth-service/config"
	"go-micro-services/src/auth-service/model"
	"go-micro-services/src/auth-service/model/response"
	"go-micro-services/src/auth-service/repository"
	"net/http"
	"strings"
	"time"
)

type AuthMiddleware struct {
	SessionRepository *repository.SessionRepository
	DatabaseConfig    *config.DatabaseConfig
}

func NewAuthMiddleware(sessionRepository repository.SessionRepository, databaseConfig *config.DatabaseConfig) *AuthMiddleware {
	return &AuthMiddleware{
		SessionRepository: &sessionRepository,
		DatabaseConfig:    databaseConfig,
	}
}

func (authMiddleware *AuthMiddleware) GetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, reader *http.Request) {
		token := reader.Header.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			result := &response.Response[any]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: Missing token",
				Data:    nil,
			}
			response.NewResponse(writer, result)
			return
		}

		transaction := reader.Context().Value("transaction").(*model.Transaction)

		foundSession, foundSessionErr := authMiddleware.SessionRepository.FindOneByAccToken(transaction.Tx, token)
		if foundSessionErr != nil {
			transaction.TxErr = foundSessionErr
			return
		}
		if foundSession == nil {
			rollbackErr := transaction.Tx.Rollback()
			if rollbackErr != nil {
				transaction.TxErr = rollbackErr
				return
			}
			result := &response.Response[any]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: Invalid Token",
				Data:    nil,
			}
			response.NewResponse(writer, result)
			return
		}
		if foundSession.AccessTokenExpiredAt.Time.Before(time.Now()) {
			result := &response.Response[any]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: Token expired",
				Data:    nil,
			}
			response.NewResponse(writer, result)
			return
		}
		next.ServeHTTP(writer, reader)
	})
}
