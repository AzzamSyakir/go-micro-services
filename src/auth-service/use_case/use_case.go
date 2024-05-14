package use_case

import (
	"encoding/json"
	"fmt"
	"go-micro-services/src/auth-service/config"
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
	DatabaseConfig *config.DatabaseConfig
	AuthRepository *repository.AuthRepository
	Env            *config.EnvConfig
}

func NewAuthUseCase(
	databaseConfig *config.DatabaseConfig,
	authRepository *repository.AuthRepository,
	env *config.EnvConfig,
) *AuthUseCase {
	authUseCase := &AuthUseCase{
		DatabaseConfig: databaseConfig,
		AuthRepository: authRepository,
		Env:            env,
	}
	return authUseCase
}

func (authUseCase *AuthUseCase) Login(request *model_request.LoginRequest) (result *model_response.Response[*entity.Session], err error) {
	begin, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase Login failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	foundUser := authUseCase.FindUserByEmail(request.Email.String)
	if foundUser.Data == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase Login fail, GetUser failed, " + foundUser.Message,
			Data:    nil,
		}
		return result, rollback
	}

	comparePasswordErr := bcrypt.CompareHashAndPassword([]byte(foundUser.Data.Password.String), []byte(request.Password.String))
	if comparePasswordErr != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusNotFound,
			Message: "AuthUseCase Login is failed, password is not match.",
			Data:    nil,
		}
		return result, rollback
	}

	accessToken := null.NewString(uuid.NewString(), true)
	refreshToken := null.NewString(uuid.NewString(), true)
	currentTime := null.NewTime(time.Now(), true)
	accessTokenExpiredAt := null.NewTime(currentTime.Time.Add(time.Minute*10), true)
	refreshTokenExpiredAt := null.NewTime(currentTime.Time.Add(time.Hour*24*2), true)

	foundSession, err := authUseCase.AuthRepository.GetOneByUserId(begin, foundUser.Data.Id.String)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase Login failed, query to db fail, " + err.Error(),
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
				Code:    http.StatusBadRequest,
				Message: "AuthUseCase Login failed, query to db fail, " + err.Error(),
				Data:    nil,
			}
			return result, rollback
		}

		commit := begin.Commit()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase Login is succeed",
			Data:    patchedSession,
		}
		return result, commit
	}

	newSession := &entity.Session{
		Id:                    null.NewString(uuid.NewString(), true),
		UserId:                foundUser.Data.Id,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiredAt:  accessTokenExpiredAt,
		RefreshTokenExpiredAt: refreshTokenExpiredAt,
		CreatedAt:             currentTime,
		UpdatedAt:             currentTime,
		DeletedAt:             null.NewTime(time.Time{}, false),
	}

	createdSession, err := authUseCase.AuthRepository.CreateSession(begin, newSession)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase Login failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &model_response.Response[*entity.Session]{
		Code:    http.StatusBadRequest,
		Message: "AuthUseCase Login is succeed",
		Data:    createdSession,
	}
	return result, commit
}

func (authUseCase *AuthUseCase) Logout(accessToken string) (result *model_response.Response[*entity.Session], err error) {
	begin, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase Logout failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	foundSession, err := authUseCase.AuthRepository.FindOneByAccToken(begin, accessToken)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase Logout failed, Invalid token, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	if foundSession == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase Logout is failed, session is not found by access token.",
			Data:    nil,
		}
		return result, rollback
	}
	deletedSession, err := authUseCase.AuthRepository.DeleteOneById(begin, foundSession.Id.String)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase Logout failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	if deletedSession == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase Logout failed, delete session failed",
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &model_response.Response[*entity.Session]{
		Code:    http.StatusOK,
		Message: "AuthUseCase Logout is succeed.",
		Data:    deletedSession,
	}
	return result, commit
}

func (authUseCase *AuthUseCase) GetNewAccessToken(refreshToken string) (result *model_response.Response[*entity.Session], err error) {
	begin, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase GetNewAccesToken failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	foundSession, err := authUseCase.AuthRepository.FindOneByRefToken(begin, refreshToken)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase GetNewAccesToken failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	if foundSession == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase GetNewAccesToken  failed, session is not found by refresh token.",
			Data:    nil,
		}
		return result, rollback
	}

	if foundSession.RefreshTokenExpiredAt.Time.Before(time.Now()) {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusNotFound,
			Message: "AuthUseCase GetNewAccessToken is failed, refresh token is expired.",
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
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase GetNewAccesToken  failed, query to db fail," + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &model_response.Response[*entity.Session]{
		Code:    http.StatusOK,
		Message: "AuthUseCase GetNewAccessToken is succeed.",
		Data:    patchedSession,
	}
	return result, commit

}

func (authUseCase *AuthUseCase) FindUserByEmail(email string) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", authUseCase.Env.App.UserHost, authUseCase.Env.App.UserPort)
	url := fmt.Sprintf("%s/%s/%s/%s", address, "users", "email", email)
	newRequest, newRequestErr := http.NewRequest("GET", url, nil)
	if newRequestErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase failed, GetUser by email user is failed," + newRequestErr.Error(),
			Data:    nil,
			Errors:  true,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase failed, GetUser by email user is failed," + doErr.Error(),
			Data:    nil,
			Errors:  true,
		}
		return result
	}
	bodyResponseUser := &model_response.Response[*entity.User]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseUser)
	if decodeErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase fail, GetUser by email user is failed," + decodeErr.Error(),
			Data:    nil,
			Errors:  true,
		}
		return result
	}
	return bodyResponseUser
}
