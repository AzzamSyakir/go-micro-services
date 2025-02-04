package use_case

import (
	"context"
	"fmt"
	"go-micro-services/grpc/pb"
	"go-micro-services/src/auth-service/config"
	"go-micro-services/src/auth-service/delivery/grpc/client"
	"go-micro-services/src/auth-service/entity"
	model_request "go-micro-services/src/auth-service/model/request/controller"
	model_response "go-micro-services/src/auth-service/model/response"
	"go-micro-services/src/auth-service/repository"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	pb.UnimplementedAuthServiceServer
	DatabaseConfig *config.DatabaseConfig
	AuthRepository *repository.AuthRepository
	Env            *config.EnvConfig
	userClient     *client.UserServiceClient
}

func NewAuthUseCase(
	databaseConfig *config.DatabaseConfig,
	authRepository *repository.AuthRepository,
	env *config.EnvConfig,
	initUserClient *client.UserServiceClient,
) *AuthUseCase {
	authUseCase := &AuthUseCase{
		UnimplementedAuthServiceServer: pb.UnimplementedAuthServiceServer{},
		userClient:                     initUserClient,
		DatabaseConfig:                 databaseConfig,
		AuthRepository:                 authRepository,
		Env:                            env,
	}
	return authUseCase
}

func (authUseCase *AuthUseCase) Login(request *model_request.LoginRequest) (result *model_response.Response[*entity.Session]) {
	tx, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		tx.Rollback()
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Login failed: Unable to start database transaction. " + err.Error(),
			Data:    nil,
		}
	}

	foundUser, err := authUseCase.userClient.GetUserByEmail(request.Email.String)
	if err != nil {
		tx.Rollback()
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Login failed: Error retrieving user data. " + err.Error(),
			Data:    nil,
		}
	}

	if foundUser.Data == nil {
		tx.Rollback()
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusNotFound,
			Message: "Login failed: User not found.",
			Data:    nil,
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Data.Password), []byte(request.Password.String))
	if err != nil {
		tx.Rollback()
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusUnauthorized,
			Message: "Login failed: Incorrect password.",
			Data:    nil,
		}
	}

	accessToken := null.NewString(uuid.NewString(), true)
	refreshToken := null.NewString(uuid.NewString(), true)
	currentTime := null.NewTime(time.Now(), true)
	accessTokenExpiredAt := null.NewTime(currentTime.Time.Add(time.Minute*10), true)
	refreshTokenExpiredAt := null.NewTime(currentTime.Time.Add(time.Hour*24*2), true)

	foundSession, err := authUseCase.AuthRepository.GetOneByUserId(tx, foundUser.Data.Id)
	if err != nil {
		tx.Rollback()
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Login failed: Error retrieving session from database. " + err.Error(),
			Data:    nil,
		}
	}

	if foundSession != nil {
		foundSession.AccessToken = accessToken
		foundSession.RefreshToken = refreshToken
		foundSession.AccessTokenExpiredAt = accessTokenExpiredAt
		foundSession.RefreshTokenExpiredAt = refreshTokenExpiredAt
		foundSession.UpdatedAt = currentTime

		patchedSession, err := authUseCase.AuthRepository.PatchOneById(tx, foundSession.Id.String, foundSession)
		if err != nil {
			tx.Rollback()
			return &model_response.Response[*entity.Session]{
				Code:    http.StatusInternalServerError,
				Message: "Login failed: Unable to update session. " + err.Error(),
				Data:    nil,
			}
		}

		if commitErr := tx.Commit(); commitErr != nil {
			return &model_response.Response[*entity.Session]{
				Code:    http.StatusInternalServerError,
				Message: "Login failed: Unable to commit transaction. " + commitErr.Error(),
				Data:    nil,
			}
		}

		return &model_response.Response[*entity.Session]{
			Code:    http.StatusOK,
			Message: "Login successful.",
			Data:    patchedSession,
		}
	}

	newSession := &entity.Session{
		Id:                    null.NewString(uuid.NewString(), true),
		UserId:                null.NewString(foundUser.Data.Id, true),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiredAt:  accessTokenExpiredAt,
		RefreshTokenExpiredAt: refreshTokenExpiredAt,
		CreatedAt:             currentTime,
		UpdatedAt:             currentTime,
	}

	createdSession, err := authUseCase.AuthRepository.CreateSession(tx, newSession)
	if err != nil {
		tx.Rollback()
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Login failed: Unable to create new session. " + err.Error(),
			Data:    nil,
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Login failed: Unable to finalize transaction. " + commitErr.Error(),
			Data:    nil,
		}
	}

	return &model_response.Response[*entity.Session]{
		Code:    http.StatusOK,
		Message: "Login successful.",
		Data:    createdSession,
	}
}
func (authUseCase *AuthUseCase) Logout(accessToken string) (result *model_response.Response[*entity.Session]) {
	tx, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		tx.Rollback()
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Logout failed: Unable to start database transaction. " + err.Error(),
			Data:    nil,
		}
	}

	if accessToken == "" {
		tx.Rollback()
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "Logout failed: Access token is required.",
			Data:    nil,
		}
	}

	foundSession, err := authUseCase.AuthRepository.FindOneByAccToken(tx, accessToken)
	if err != nil {
		tx.Rollback()
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusUnauthorized,
			Message: "Logout failed: Invalid access token. " + err.Error(),
			Data:    nil,
		}
	}

	if foundSession == nil {
		tx.Rollback()
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusNotFound,
			Message: "Logout failed: No active session found for the provided access token.",
			Data:    nil,
		}
	}

	deletedSession, err := authUseCase.AuthRepository.DeleteOneById(tx, foundSession.Id.String)
	if err != nil {
		tx.Rollback()
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Logout failed: Unable to delete session from database. " + err.Error(),
			Data:    nil,
		}
	}

	if deletedSession == nil {
		tx.Rollback()
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Logout failed: Session deletion was unsuccessful.",
			Data:    nil,
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Logout failed: Unable to finalize transaction. " + commitErr.Error(),
			Data:    nil,
		}
	}

	return &model_response.Response[*entity.Session]{
		Code:    http.StatusOK,
		Message: "Logout successful.",
		Data:    deletedSession,
	}
}
func (authUseCase *AuthUseCase) LogoutWithUserId(context context.Context, id *pb.ByUserId) (empty *pb.Empty, err error) {
	tx, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		tx.Rollback()
		return &pb.Empty{}, fmt.Errorf("Logout failed: Unable to start database transaction. %v", err)
	}

	if id.Id == "" {
		tx.Rollback()
		return &pb.Empty{}, fmt.Errorf("Logout failed: User ID is required")
	}

	foundSession, err := authUseCase.AuthRepository.GetOneByUserId(tx, id.Id)
	if err != nil {
		tx.Rollback()
		return &pb.Empty{}, fmt.Errorf("Logout failed: Error retrieving session from database. %v", err)
	}
	if foundSession == nil {
		tx.Rollback()
		return &pb.Empty{}, fmt.Errorf("Logout failed: No active session found for the given user ID")
	}

	_, err = authUseCase.AuthRepository.DeleteOneByUserId(tx, foundSession.UserId.String)
	if err != nil {
		tx.Rollback()
		return &pb.Empty{}, fmt.Errorf("Logout failed: Unable to delete session. %v", err)
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		return &pb.Empty{}, fmt.Errorf("Logout failed: Error committing transaction. %v", commitErr)
	}

	return &pb.Empty{}, nil
}
func (authUseCase *AuthUseCase) GetNewAccessToken(refreshToken string) (result *model_response.Response[*entity.Session]) {
	tx, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		tx.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate new access token: Unable to start database transaction. " + err.Error(),
			Data:    nil,
		}
		return result
	}

	if refreshToken == "" {
		tx.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "Failed to generate new access token: Refresh token is required.",
			Data:    nil,
		}
		return result
	}

	foundSession, err := authUseCase.AuthRepository.FindOneByRefToken(tx, refreshToken)
	if err != nil {
		tx.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate new access token: Error retrieving session from database. " + err.Error(),
			Data:    nil,
		}
		return result
	}

	if foundSession == nil {
		tx.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusNotFound,
			Message: "Failed to generate new access token: No session found for the provided refresh token.",
			Data:    nil,
		}
		return result
	}

	if foundSession.RefreshTokenExpiredAt.Time.Before(time.Now()) {
		tx.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusUnauthorized,
			Message: "Failed to generate new access token: Refresh token has expired.",
			Data:    nil,
		}
		return result
	}

	foundSession.AccessToken = null.NewString(uuid.NewString(), true)
	foundSession.UpdatedAt = null.NewTime(time.Now(), true)
	patchedSession, err := authUseCase.AuthRepository.PatchOneById(tx, foundSession.Id.String, foundSession)
	if err != nil {
		tx.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate new access token: Unable to update session in database. " + err.Error(),
			Data:    nil,
		}
		return result
	}

	commit := tx.Commit()
	if commit != nil {
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate new access token: Error committing transaction. " + commit.Error(),
			Data:    nil,
		}
		tx.Commit()
		return result
	}

	result = &model_response.Response[*entity.Session]{
		Code:    http.StatusOK,
		Message: "New access token generated successfully.",
		Data:    patchedSession,
	}
	return result
}
func (authUseCase *AuthUseCase) ListSessions() (result *model_response.Response[[]*entity.Session]) {
	tx, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		tx.Rollback()
		return &model_response.Response[[]*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Error: Unable to start database transaction. " + err.Error(),
			Data:    nil,
		}
	}

	sessions, err := authUseCase.AuthRepository.ListSession(tx)
	if err != nil {
		tx.Rollback()
		return &model_response.Response[[]*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Error: Failed to retrieve session data from the database. " + err.Error(),
			Data:    nil,
		}
	}

	if len(sessions) == 0 {
		tx.Rollback()
		return &model_response.Response[[]*entity.Session]{
			Code:    http.StatusNotFound,
			Message: "No active sessions found in the database.",
			Data:    nil,
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		tx.Rollback()
		return &model_response.Response[[]*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Error: Failed to finalize the database transaction. " + commitErr.Error(),
			Data:    nil,
		}
	}

	return &model_response.Response[[]*entity.Session]{
		Code:    http.StatusOK,
		Message: "Successfully retrieved all active sessions.",
		Data:    sessions,
	}
}
