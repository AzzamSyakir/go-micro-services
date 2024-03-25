package use_case

import (
	"context"
	"fmt"
	"github.com/guregu/null"
	"go-micro-services/internal/config"
	"go-micro-services/internal/entity"
	model_request "go-micro-services/internal/model/request/controller"
	"go-micro-services/internal/model/response"
	"go-micro-services/internal/repository"
	"net/http"
	"time"
)

type UserUseCase struct {
	DatabaseConfig *config.DatabaseConfig
	UserRepository *repository.UserRepository
}

func NewUserUseCase(
	databaseConfig *config.DatabaseConfig,
	userRepository *repository.UserRepository,

) *UserUseCase {
	userUseCase := &UserUseCase{
		DatabaseConfig: databaseConfig,
		UserRepository: userRepository,
	}
	return userUseCase
}
func (userUseCase *UserUseCase) GetOneById(id string) (result *response.Response[*entity.User], err error) {
	transaction, transactionErr := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if transactionErr != nil {
		errorMessage := fmt.Sprintf("transaction failed :", transactionErr)
		result = &response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}

	foundUser, foundUserErr := userUseCase.UserRepository.GetOneById(transaction, id)
	if foundUserErr != nil {
		result = nil
		err = foundUserErr
		return result, err
	}
	if foundUser == nil {
		errorMessage := fmt.Sprintf("User UseCase FindOneById is failed, user is not found by id.", id)
		result = &response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}

	result = &response.Response[*entity.User]{
		Code:    http.StatusOK,
		Message: "User UseCase FindOneById is succeed.",
		Data:    foundUser,
	}
	err = nil
	return result, err
}
func (userUseCase *UserUseCase) PatchOneByIdFromRequest(ctx context.Context, id string, request *model_request.UserPatchOneByIdRequest) (result *response.Response[*entity.User], err error) {
	transaction, transactionErr := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if transactionErr != nil {
		errorMessage := fmt.Sprintf("transaction failed :", transactionErr)
		result = &response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}

	foundUser, foundUserErr := userUseCase.UserRepository.GetOneById(transaction, id)
	if foundUserErr != nil {
		result = nil
		err = foundUserErr
		return result, err
	}
	if foundUser == nil {
		errorMessage := fmt.Sprintf("User UseCase PatchOneById is failed, user is not found by id : %s", id)
		result = &response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}

	if request.Name.Valid {
		foundUser.Name = request.Name
	}
	if request.Saldo.Valid {
		foundUser.Saldo = request.Saldo
	}

	foundUser.UpdatedAt = null.NewTime(time.Now(), true)

	patchedUser, patchedUserErr := userUseCase.UserRepository.PatchOneById(transaction, id, foundUser)
	if patchedUserErr != nil {
		result = nil
		err = patchedUserErr
		return result, err
	}

	result = &response.Response[*entity.User]{
		Code:    http.StatusOK,
		Message: "UserUserCase PatchOneByIdFromRequest is succeed.",
		Data:    patchedUser,
	}
	err = nil
	return result, err
}
