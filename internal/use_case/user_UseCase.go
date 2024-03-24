package use_case

import (
	"fmt"
	"go-micro-services/internal/config"
	"go-micro-services/internal/entity"
	"go-micro-services/internal/model/response"
	"go-micro-services/internal/repository"
	"net/http"
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
		result = &response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: "User UseCase FindOneById is failed, user is not found by id.",
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
