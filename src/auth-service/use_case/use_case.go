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

func (authUseCase *AuthUseCase) Login(request *model_request.LoginRequest) (result *model_response.Response[*entity.Session], err error) {
	begin, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Login failed: Unable to start database transaction. " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	foundUser, err := authUseCase.userClient.GetUserByEmail(request.Email.String)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "Login failed: " + foundUser.Message,
			Data:    nil,
		}
		return result, rollback
	}
	if foundUser.Data == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "Login failed: User not found.",
			Data:    nil,
		}
		return result, rollback
	}

	comparePasswordErr := bcrypt.CompareHashAndPassword([]byte(foundUser.Data.Password), []byte(request.Password.String))
	if comparePasswordErr != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusUnauthorized,
			Message: "Login failed: Incorrect password.",
			Data:    nil,
		}
		return result, rollback
	}

	accessToken := null.NewString(uuid.NewString(), true)
	refreshToken := null.NewString(uuid.NewString(), true)
	currentTime := null.NewTime(time.Now(), true)
	accessTokenExpiredAt := null.NewTime(currentTime.Time.Add(time.Minute*10), true)
	refreshTokenExpiredAt := null.NewTime(currentTime.Time.Add(time.Hour*24*2), true)

	foundSession, err := authUseCase.AuthRepository.GetOneByUserId(begin, foundUser.Data.Id)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Login failed: Error retrieving session from database. " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	if foundSession != nil {
		foundSession.AccessToken = accessToken
		foundSession.RefreshToken = refreshToken
		foundSession.AccessTokenExpiredAt = accessTokenExpiredAt
		foundSession.RefreshTokenExpiredAt = refreshTokenExpiredAt
		foundSession.UpdatedAt = currentTime

		patchedSession, err := authUseCase.AuthRepository.PatchOneById(begin, foundSession.Id.String, foundSession)
		if err != nil {
			rollback := begin.Rollback()
			result = &model_response.Response[*entity.Session]{
				Code:    http.StatusInternalServerError,
				Message: "Login failed: Unable to update session. " + err.Error(),
				Data:    nil,
			}
			return result, rollback
		}

		commit := begin.Commit()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusOK,
			Message: "Login successful.",
			Data:    patchedSession,
		}
		return result, commit
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

	createdSession, err := authUseCase.AuthRepository.CreateSession(begin, newSession)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Login failed: Unable to create new session. " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &model_response.Response[*entity.Session]{
		Code:    http.StatusOK,
		Message: "Login successful.",
		Data:    createdSession,
	}
	return result, commit
}
func (authUseCase *AuthUseCase) Logout(accessToken string) (result *model_response.Response[*entity.Session], err error) {
	begin, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Logout failed: Unable to start database transaction. " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	foundSession, err := authUseCase.AuthRepository.FindOneByAccToken(begin, accessToken)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusUnauthorized,
			Message: "Logout failed: Invalid access token. " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	if foundSession == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusNotFound,
			Message: "Logout failed: Session not found for the provided access token.",
			Data:    nil,
		}
		return result, rollback
	}

	deletedSession, err := authUseCase.AuthRepository.DeleteOneById(begin, foundSession.Id.String)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Logout failed: Unable to delete session from database. " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	if deletedSession == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Logout failed: Failed to remove session.",
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &model_response.Response[*entity.Session]{
		Code:    http.StatusOK,
		Message: "Logout successful.",
		Data:    deletedSession,
	}
	return result, commit
}
func (authUseCase *AuthUseCase) LogoutWithUserId(context context.Context, id *pb.ByUserId) (empty *pb.Empty, err error) {
	begin, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		begin.Rollback()
		return &pb.Empty{}, fmt.Errorf("Logout failed: Unable to start database transaction. %v", err)
	}

	foundSession, err := authUseCase.AuthRepository.GetOneByUserId(begin, id.Id)
	if err != nil {
		begin.Rollback()
		return &pb.Empty{}, fmt.Errorf("Logout failed: Error retrieving session from database. %v", err)
	}
	if foundSession == nil {
		begin.Rollback()
		return &pb.Empty{}, fmt.Errorf("Logout failed: No active session found for the given user ID")
	}

	_, err = authUseCase.AuthRepository.DeleteOneByUserId(begin, foundSession.UserId.String)
	if err != nil {
		begin.Rollback()
		return &pb.Empty{}, fmt.Errorf("Logout failed: Unable to delete session. %v", err)
	}

	commitErr := begin.Commit()
	if commitErr != nil {
		return &pb.Empty{}, fmt.Errorf("Logout failed: Error committing transaction. %v", commitErr)
	}

	return &pb.Empty{}, nil
}
func (authUseCase *AuthUseCase) GetNewAccessToken(refreshToken string) (result *model_response.Response[*entity.Session], err error) {
	begin, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate new access token: Unable to start database transaction. " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	foundSession, err := authUseCase.AuthRepository.FindOneByRefToken(begin, refreshToken)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate new access token: Error retrieving session from database. " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	if foundSession == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusNotFound,
			Message: "Failed to generate new access token: No session found for the provided refresh token.",
			Data:    nil,
		}
		return result, rollback
	}

	if foundSession.RefreshTokenExpiredAt.Time.Before(time.Now()) {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusUnauthorized,
			Message: "Failed to generate new access token: Refresh token has expired.",
			Data:    nil,
		}
		return result, rollback
	}

	foundSession.AccessToken = null.NewString(uuid.NewString(), true)
	foundSession.UpdatedAt = null.NewTime(time.Now(), true)
	patchedSession, err := authUseCase.AuthRepository.PatchOneById(begin, foundSession.Id.String, foundSession)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate new access token: Unable to update session in database. " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &model_response.Response[*entity.Session]{
		Code:    http.StatusOK,
		Message: "New access token generated successfully.",
		Data:    patchedSession,
	}
	return result, commit
}
func (authUseCase *AuthUseCase) ListSessions() (result *model_response.Response[[]*entity.Session]) {
	tx, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		tx.Rollback()
		return &model_response.Response[[]*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Error: Unable to start database transaction.",
			Data:    nil,
		}
	}

	sessions, err := authUseCase.AuthRepository.ListSession(tx)
	if err != nil {
		tx.Rollback()
		return &model_response.Response[[]*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Error: Failed to retrieve session data from the database.",
			Data:    nil,
		}
	}

	if len(sessions) == 0 {
		tx.Rollback()
		return &model_response.Response[[]*entity.Session]{
			Code:    http.StatusNotFound,
			Message: "No active sessions found.",
			Data:    nil,
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return &model_response.Response[[]*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "Error: Failed to finalize the database transaction.",
			Data:    nil,
		}
	}

	return &model_response.Response[[]*entity.Session]{
		Code:    http.StatusOK,
		Message: "Successfully retrieved all active sessions.",
		Data:    sessions,
	}
}
