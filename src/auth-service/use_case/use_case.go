package use_case

import (
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
	AuthRepository *repository.UserRepository
}

func NewAuthUseCase(
	databaseConfig *config.DatabaseConfig,
	authRepository *repository.UserRepository,
) *AuthUseCase {
	authUseCase := &AuthUseCase{
		DatabaseConfig: databaseConfig,
		AuthRepository: authRepository,
	}
	return authUseCase
}

func (userUseCase *AuthUseCase) GetOneById(id string) (result *model_response.Response[*entity.Auth], err error) {
	transaction, transactionErr := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if transactionErr != nil {
		errorMessage := fmt.Sprintf("transaction failed :%s", transactionErr)
		result = &model_response.Response[*entity.Auth]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}
	GetOneById, GetOneByIdErr := userUseCase.AuthRepository.GetOneById(transaction, id)
	if GetOneByIdErr != nil {
		errorMessage := fmt.Sprintf("UserUseCase GetOneById is failed, GetUser failed : %s", GetOneByIdErr)
		result = &model_response.Response[*entity.Auth]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}
	if GetOneById == nil {
		errorMessage := fmt.Sprintf("Auth UseCase FindOneById is failed, Auth is not found by id %s", id)
		result = &model_response.Response[*entity.Auth]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}

	result = &model_response.Response[*entity.Auth]{
		Code:    http.StatusOK,
		Message: "Auth UseCase FindOneById is succeed.",
		Data:    GetOneById,
	}
	err = nil
	return result, err
}

func (userUseCase *AuthUseCase) UpdateBalance(id string, request *model_request.UserPatchOneByIdRequest) (result *model_response.Response[*entity.Auth]) {
	beginErr := crdb.Execute(func() (err error) {
		transaction, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
		if err != nil {
			return err
		}

		foundUser, err := userUseCase.AuthRepository.GetOneById(transaction, id)
		if err != nil {
			return err
		}
		if foundUser == nil {
			err = transaction.Rollback()
			result = &model_response.Response[*entity.Auth]{
				Code:    http.StatusNotFound,
				Message: "UserUserCase UpdateBalance is failed, Auth is not found by id.",
				Data:    nil,
			}
			return err
		}
		if request.Balance.Valid {
			foundUser.Balance = request.Balance
		} else {
			err = transaction.Rollback()
			result = &model_response.Response[*entity.Auth]{
				Code:    http.StatusNotFound,
				Message: "UserUserCase UpdateBalance is failed, balance is not provided ",
				Data:    nil,
			}
			return err
		}

		foundUser.UpdatedAt = null.NewTime(time.Now(), true)

		patchedUser, err := userUseCase.AuthRepository.PatchOneById(transaction, id, foundUser)
		if err != nil {
			return err
		}

		err = transaction.Commit()
		result = &model_response.Response[*entity.Auth]{
			Code:    http.StatusOK,
			Message: "UserUserCase UpdateBalance is succeed.",
			Data:    patchedUser,
		}
		return err
	})

	if beginErr != nil {
		result = &model_response.Response[*entity.Auth]{
			Code:    http.StatusInternalServerError,
			Message: "UserUserCase UpdateBalance  is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}

	return result
}

func (userUseCase *AuthUseCase) UpdateUser(id string, request *model_request.UserPatchOneByIdRequest) (result *model_response.Response[*entity.Auth]) {
	beginErr := crdb.Execute(func() (err error) {
		transaction, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
		if err != nil {
			return err
		}

		foundUser, err := userUseCase.AuthRepository.GetOneById(transaction, id)
		if err != nil {
			return err
		}
		if foundUser == nil {
			err = transaction.Rollback()
			result = &model_response.Response[*entity.Auth]{
				Code:    http.StatusNotFound,
				Message: "UserUserCase UpdateUser is failed, Auth is not found by id.",
				Data:    nil,
			}
			return err
		}
		if request.Name.Valid {
			foundUser.Name = request.Name
		}
		if request.Email.Valid {
			foundUser.Email = request.Email
		}
		if request.Password.Valid {
			foundUser.Password = request.Password
		}
		if request.Balance.Valid {
			foundUser.Balance = request.Balance
		}

		foundUser.UpdatedAt = null.NewTime(time.Now(), true)

		patchedUser, err := userUseCase.AuthRepository.PatchOneById(transaction, id, foundUser)
		if err != nil {
			return err
		}

		err = transaction.Commit()
		result = &model_response.Response[*entity.Auth]{
			Code:    http.StatusOK,
			Message: "UserUserCase UpdateUser is succeed.",
			Data:    patchedUser,
		}
		return err
	})

	if beginErr != nil {
		result = &model_response.Response[*entity.Auth]{
			Code:    http.StatusInternalServerError,
			Message: "UserUserCase UpdateUser  is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}

	return result
}

func (userUseCase *AuthUseCase) CreateUser(request *model_request.CreateUser) (result *model_response.Response[*entity.Auth]) {
	beginErr := crdb.Execute(func() (err error) {
		begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
		if err != nil {
			result = nil
			return err
		}

		hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(request.Password.String), bcrypt.DefaultCost)
		if hashedPasswordErr != nil {
			err = begin.Rollback()
			result = &model_response.Response[*entity.Auth]{
				Code:    http.StatusInternalServerError,
				Message: "UserUseCase Register is failed, password hashing is failed.",
				Data:    nil,
			}
			return err
		}

		currentTime := null.NewTime(time.Now(), true)
		newUser := &entity.Auth{
			Id:        null.NewString(uuid.NewString(), true),
			Name:      request.Name,
			Email:     request.Email,
			Password:  null.NewString(string(hashedPassword), true),
			Balance:   request.Balance,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: null.NewTime(time.Time{}, false),
		}

		createdUser, err := userUseCase.AuthRepository.CreateUser(begin, newUser)
		if err != nil {
			return err
		}

		err = begin.Commit()
		result = &model_response.Response[*entity.Auth]{
			Code:    http.StatusCreated,
			Message: "UserUseCase Register is succeed.",
			Data:    createdUser,
		}
		return err
	})

	if beginErr != nil {
		result = &model_response.Response[*entity.Auth]{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase Register  is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}

	return result
}

func (userUseCase *AuthUseCase) DeleteUser(id string) (result *model_response.Response[*entity.Auth]) {
	beginErr := crdb.Execute(func() (err error) {
		begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
		if err != nil {
			return err
		}

		deletedUser, deletedUserErr := userUseCase.AuthRepository.DeleteUser(begin, id)
		if deletedUserErr != nil {
			err = begin.Rollback()
			result = &model_response.Response[*entity.Auth]{
				Code:    http.StatusNotFound,
				Message: "UserUserCase DeleteUser is failed, " + deletedUserErr.Error(),
				Data:    nil,
			}
			return err
		}
		if deletedUser == nil {
			err = begin.Rollback()
			result = &model_response.Response[*entity.Auth]{
				Code:    http.StatusNotFound,
				Message: "UserUserCase DeleteUser is failed, auth is not deleted by id, " + id,
				Data:    nil,
			}
			return err
		}

		err = begin.Commit()
		result = &model_response.Response[*entity.Auth]{
			Code:    http.StatusOK,
			Message: "UserUserCase DeleteUser is succeed.",
			Data:    deletedUser,
		}
		return err
	})

	if beginErr != nil {
		result = &model_response.Response[*entity.Auth]{
			Code:    http.StatusInternalServerError,
			Message: "UserUserCase DeleteUser is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}

	return result
}

func (userUseCase *AuthUseCase) FetchUser() (result *model_response.Response[[]*entity.Auth], err error) {
	transaction, transactionErr := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if transactionErr != nil {
		errorMessage := fmt.Sprintf("transaction failed :%s", transactionErr)
		result = &model_response.Response[[]*entity.Auth]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}

	fetchUser, fetchUserErr := userUseCase.AuthRepository.FetchUser(transaction)
	if fetchUserErr != nil {
		errorMessage := fmt.Sprintf("UserUseCase fetchUser is failed, GetUser failed : %s", fetchUserErr)
		result = &model_response.Response[[]*entity.Auth]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}

	if fetchUser.Data == nil {
		result = &model_response.Response[[]*entity.Auth]{
			Code:    http.StatusNotFound,
			Message: "Auth UseCase FetchUser is failed, data Auth is empty ",
			Data:    nil,
		}
		err = nil
		return result, err
	}

	result = &model_response.Response[[]*entity.Auth]{
		Code:    http.StatusOK,
		Message: "Auth UseCase FetchUser is succeed.",
		Data:    fetchUser.Data,
	}
	err = nil
	return result, err
}
