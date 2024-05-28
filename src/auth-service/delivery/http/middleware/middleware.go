package middleware

import (
	"go-micro-services/src/auth-service/config"
	"go-micro-services/src/auth-service/model/response"
	"go-micro-services/src/auth-service/repository"
	"net/http"
	"strings"
	"time"

	"github.com/guregu/null"
)

type AuthMiddleware struct {
	SessionRepository *repository.AuthRepository
	DatabaseConfig    *config.DatabaseConfig
}

func NewAuthMiddleware(sessionRepository repository.AuthRepository, databaseConfig *config.DatabaseConfig) *AuthMiddleware {
	return &AuthMiddleware{
		SessionRepository: &sessionRepository,
		DatabaseConfig:    databaseConfig,
	}
}

func (authMiddleware *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.Replace(token, "Bearer ", "", 1)
		if token == "" {
			result := &response.Response[interface{}]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: Missing token",
			}
			response.NewResponse(w, result)
			return
		}

		begin, err := authMiddleware.DatabaseConfig.AuthDB.Connection.Begin()
		if err != nil {
			begin.Rollback()
			result := &response.Response[interface{}]{
				Code:    http.StatusInternalServerError,
				Message: "transaction error",
			}
			response.NewResponse(w, result)
			return
		}

		session, err := authMiddleware.SessionRepository.FindOneByAccToken(begin, token)
		if err != nil {
			begin.Rollback()
			result := &response.Response[interface{}]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: token not found",
			}
			response.NewResponse(w, result)
			return
		}
		if session == nil {
			begin.Rollback()
			result := &response.Response[interface{}]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: Invalid Token",
			}
			response.NewResponse(w, result)
			return
		}
		if session.AccessTokenExpiredAt == null.NewTime(time.Now(), true) {
			begin.Rollback()
			result := &response.Response[interface{}]{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: Token expired",
			}
			response.NewResponse(w, result)
			return
		}
		begin.Commit()
		next.ServeHTTP(w, r)
	})
}
