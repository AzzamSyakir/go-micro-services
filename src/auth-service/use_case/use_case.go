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

	"github.com/cockroachdb/cockroach-go/v2/crdb"
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

func (authUseCase *AuthUseCase) Login(request *model_request.LoginRequest) (result *model_response.Response[*entity.Session]) {
	beginErr := crdb.Execute(func() (err error) {
		begin, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
		if err != nil {
			result = nil
			return err
		}

		foundUser := authUseCase.FindUserByEmail(request.Email.String)

		if foundUser.Errors == true {
			err = begin.Rollback()
			result = &model_response.Response[*entity.Session]{
				Code:    http.StatusNotFound,
				Message: foundUser.Message,
				Data:    nil,
			}
			return err
		}

		comparePasswordErr := bcrypt.CompareHashAndPassword([]byte(foundUser.Data.Password.String), []byte(request.Password.String))
		if comparePasswordErr != nil {
			err = begin.Rollback()
			result = &model_response.Response[*entity.Session]{
				Code:    http.StatusNotFound,
				Message: "AuthUseCase Login is failed, password is not match.",
				Data:    nil,
			}
			return err
		}

		accessToken := null.NewString(uuid.NewString(), true)
		refreshToken := null.NewString(uuid.NewString(), true)
		currentTime := null.NewTime(time.Now(), true)
		accessTokenExpiredAt := null.NewTime(currentTime.Time.Add(time.Minute*10), true)
		refreshTokenExpiredAt := null.NewTime(currentTime.Time.Add(time.Hour*24*2), true)

		foundSession, err := authUseCase.AuthRepository.FindOneByUserId(begin, foundUser.Data.Id.String)
		if err != nil {
			return err
		}

		if foundSession != nil {

			foundSession.AccessToken = accessToken
			foundSession.RefreshToken = refreshToken
			foundSession.AccessTokenExpiredAt = accessTokenExpiredAt
			foundSession.RefreshTokenExpiredAt = refreshTokenExpiredAt
			foundSession.UpdatedAt = currentTime
			patchedSession, err := authUseCase.AuthRepository.PatchOneById(begin, foundSession.Id.String, foundSession)
			if err != nil {
				return err
			}

			err = begin.Commit()
			result = &model_response.Response[*entity.Session]{
				Code:    http.StatusOK,
				Message: "AuthUseCase Login is succeed.",
				Data:    patchedSession,
			}
			return err
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
			return err
		}

		err = begin.Commit()
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusCreated,
			Message: "AuthUseCase Login is succeed.",
			Data:    createdSession,
		}
		return err
	})

	if beginErr != nil {
		result = &model_response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "AuthUseCase Login  is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}

	return result
}
func (authUseCase *AuthUseCase) FindUserByEmail(email string) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", authUseCase.Env.App.Host, authUseCase.Env.App.UserPort)
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
