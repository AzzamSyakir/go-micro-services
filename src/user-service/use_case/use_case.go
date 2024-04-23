package use_case

import (
	"fmt"
	"go-micro-services/src/user-service/config"
	"go-micro-services/src/user-service/entity"
	model_request "go-micro-services/src/user-service/model/request/controller"
	model_response "go-micro-services/src/user-service/model/response"
	"go-micro-services/src/user-service/repository"
	"net/http"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
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

func (userUseCase *UserUseCase) GetOneById(id string) (result *model_response.Response[*entity.User], err error) {
	transaction, transactionErr := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if transactionErr != nil {
		errorMessage := fmt.Sprintf("transaction failed :%s", transactionErr)
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}
	GetOneById, GetOneByIdErr := userUseCase.UserRepository.GetOneById(transaction, id)
	if GetOneByIdErr != nil {
		errorMessage := fmt.Sprintf("UserUseCase GetOneById is failed, GetUser failed : %s", GetOneByIdErr)
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}
	if GetOneById == nil {
		errorMessage := fmt.Sprintf("User UseCase FindOneById is failed, User is not found by id %s", id)
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}

	result = &model_response.Response[*entity.User]{
		Code:    http.StatusOK,
		Message: "User UseCase FindOneById is succeed.",
		Data:    GetOneById,
	}
	err = nil
	return result, err
}

func (userUseCase *UserUseCase) PatchOneByIdFromRequest(id string, request *model_request.UserPatchOneByIdRequest) (result *model_response.Response[*entity.User]) {
	beginErr := crdb.Execute(func() (err error) {
		transaction, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
		if err != nil {
			return err
		}

		foundUser, err := userUseCase.UserRepository.GetOneById(transaction, id)
		if err != nil {
			return err
		}
		if foundUser == nil {
			err = transaction.Rollback()
			result = &model_response.Response[*entity.User]{
				Code:    http.StatusNotFound,
				Message: "UserUserCase PatchOneByIdFromRequest is failed, User is not found by id.",
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

		patchedUser, err := userUseCase.UserRepository.PatchOneById(transaction, id, foundUser)
		if err != nil {
			return err
		}

		err = transaction.Commit()
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusOK,
			Message: "UserUserCase PatchOneByIdFromRequest is succeed.",
			Data:    patchedUser,
		}
		return err
	})

	if beginErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusInternalServerError,
			Message: "UserUserCase PatchOneByIdFromRequest  is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}

	return result
}

func (userUseCase *UserUseCase) CreateUser(request *model_request.CreateUser) (result *model_response.Response[*entity.User]) {
	beginErr := crdb.Execute(func() (err error) {
		begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
		if err != nil {
			result = nil
			return err
		}

		hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(request.Password.String), bcrypt.DefaultCost)
		if hashedPasswordErr != nil {
			err = begin.Rollback()
			result = &model_response.Response[*entity.User]{
				Code:    http.StatusInternalServerError,
				Message: "UserUseCase Register is failed, password hashing is failed.",
				Data:    nil,
			}
			return err
		}

		currentTime := null.NewTime(time.Now(), true)
		newUser := &entity.User{
			Id:        null.NewString(uuid.NewString(), true),
			Name:      request.Name,
			Email:     request.Email,
			Password:  null.NewString(string(hashedPassword), true),
			Balance:   request.Balance,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			DeletedAt: null.NewTime(time.Time{}, false),
		}

		createdUser, err := userUseCase.UserRepository.CreateUser(begin, newUser)
		if err != nil {
			return err
		}

		err = begin.Commit()
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusCreated,
			Message: "UserUseCase Register is succeed.",
			Data:    createdUser,
		}
		return err
	})

	if beginErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase Register  is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}

	return result
}

func (userUseCase *UserUseCase) DeleteUser(id string) (result *model_response.Response[*entity.User]) {
	beginErr := crdb.Execute(func() (err error) {
		begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
		if err != nil {
			return err
		}

		deletedUser, deletedUserErr := userUseCase.UserRepository.DeleteUser(begin, id)
		if deletedUserErr != nil {
			err = begin.Rollback()
			result = &model_response.Response[*entity.User]{
				Code:    http.StatusNotFound,
				Message: "UserUserCase DeleteOneById is failed, " + deletedUserErr.Error(),
				Data:    nil,
			}
			return err
		}
		if deletedUser == nil {
			err = begin.Rollback()
			result = &model_response.Response[*entity.User]{
				Code:    http.StatusNotFound,
				Message: "UserUserCase DeleteOneById is failed, user is not deleted by id, " + id,
				Data:    nil,
			}
			return err
		}

		err = begin.Commit()
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusOK,
			Message: "UserUserCase DeleteOneById is succeed.",
			Data:    deletedUser,
		}
		return err
	})

	if beginErr != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusInternalServerError,
			Message: "UserUserCase DeleteOneById is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}

	return result
}
func (userUseCase *UserUseCase) FetchUser() (result *model_response.Response[[]*entity.User], err error) {
	transaction, transactionErr := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if transactionErr != nil {
		errorMessage := fmt.Sprintf("transaction failed :%s", transactionErr)
		result = &model_response.Response[[]*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}

	fetchUser, fetchUserErr := userUseCase.UserRepository.FetchUser(transaction)
	if fetchUserErr != nil {
		errorMessage := fmt.Sprintf("UserUseCase fetchUser is failed, GetUser failed : %s", fetchUserErr)
		result = &model_response.Response[[]*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}

	if fetchUser.Data == nil {
		result = &model_response.Response[[]*entity.User]{
			Code:    http.StatusNotFound,
			Message: "User UseCase FetchUser is failed, data User is empty ",
			Data:    nil,
		}
		err = nil
		return result, err
	}

	result = &model_response.Response[[]*entity.User]{
		Code:    http.StatusOK,
		Message: "User UseCase FetchUser is succeed.",
		Data:    fetchUser.Data,
	}
	err = nil
	return result, err
}
