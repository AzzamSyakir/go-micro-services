package use_case

import (
	"encoding/json"
	"fmt"
	"go-micro-services/src/auth-service/config"
	"go-micro-services/src/auth-service/entity"
	model_request "go-micro-services/src/auth-service/model/request/controller"
	"go-micro-services/src/auth-service/model/response"
	"go-micro-services/src/auth-service/repository"
	model_response "go-micro-services/src/order-service/model/response"
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
) *AuthUseCase {
	authUseCase := &AuthUseCase{
		DatabaseConfig: databaseConfig,
		AuthRepository: authRepository,
	}
	return authUseCase
}

func (authUseCase *AuthUseCase) Login(request *model_request.LoginRequest) (result *response.Response[*entity.Session]) {
	beginErr := crdb.Execute(func() (err error) {
		begin, err := authUseCase.DatabaseConfig.AuthDB.Connection.Begin()
		if err != nil {
			result = nil
			return err
		}

		foundUser := authUseCase.FindOneByEmail(request.Email.String)

		if foundUser == nil {
			err = begin.Rollback()
			result = &response.Response[*entity.Session]{
				Code:    http.StatusNotFound,
				Message: "AuthUseCase Login is failed, user is not found by email.",
				Data:    nil,
			}
			return err
		}

		comparePasswordErr := bcrypt.CompareHashAndPassword([]byte(foundUser.Data.Password.String), []byte(request.Password.String))
		if comparePasswordErr != nil {
			err = begin.Rollback()
			result = &response.Response[*entity.Session]{
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

		foundSession, err := authUseCase.FindOneByUserId(foundUser.Data.Id.String)
		if err != nil {
			return err
		}

		if foundSession != nil {
			foundSession.AccessToken = accessToken
			foundSession.RefreshToken = refreshToken
			foundSession.AccessTokenExpiredAt = accessTokenExpiredAt
			foundSession.RefreshTokenExpiredAt = refreshTokenExpiredAt
			foundSession.UpdatedAt = currentTime
			patchedSession, err := authUseCase.PatchOneById(foundSession.Id.String, foundSession)
			if err != nil {
				return err
			}

			err = begin.Commit()
			result = &response.Response[*entity.Session]{
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
		result = &response.Response[*entity.Session]{
			Code:    http.StatusCreated,
			Message: "AuthUseCase Login is succeed.",
			Data:    createdSession,
		}
		return err
	})

	if beginErr != nil {
		result = &response.Response[*entity.Session]{
			Code:    http.StatusInternalServerError,
			Message: "AuthUseCase Login  is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}

	return result
}
func (authUseCase *AuthUseCase) FindOneByEmail(email string) (result *model_response.Response[*entity.User]) {
	address := fmt.Sprintf("http://%s:%s", authUseCase.Env.App.Host, authUseCase.Env.App.UserPort)
	url := fmt.Sprintf("%s/%s/%s/%s", address, "users", email)
	newRequest, newRequestErr := http.NewRequest("GET", url, nil)

	if newRequestErr != nil {
		err = nil
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: "OrderUseCase failed, UpdateBalance user is failed," + newRequestErr.Error(),
			Data:    nil,
		}
		return result, err
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: "AuthUseCase failed, GetUser by email user is failed," + doErr.Error(),
			Data:    nil,
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
		}
	}
	return bodyResponseUser
}
func (authUseCase *AuthUseCase) FindOneByUserId(id string) (result *entity.Session, err error) {
	return
}
func (authUseCase *AuthUseCase) PatchOneById(id string, toPatchSession *entity.Session) (result *entity.Session, err error) {
	return
}
