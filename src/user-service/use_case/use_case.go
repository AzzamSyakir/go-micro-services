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
	begin, beginErr := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if beginErr != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("begin failed :%s", beginErr)
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	GetOneById, GetOneByIdErr := userUseCase.UserRepository.GetOneById(begin, id)
	if GetOneByIdErr != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("UserUseCase GetOneById is failed, GetUser failed : %s", GetOneByIdErr)
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	if GetOneById == nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("User UseCase FindOneById is failed, User is not found by id %s", id)
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &model_response.Response[*entity.User]{
		Code:    http.StatusOK,
		Message: "User UseCase FindOneById is succeed.",
		Data:    GetOneById,
	}
	return result, commit
}

func (userUseCase *UserUseCase) GetOneByEmail(email string) (result *model_response.Response[*entity.User], err error) {
	begin, beginErr := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if beginErr != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("begin failed :%s", beginErr)
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	GetOneByEmail, GetOneByEmailErr := userUseCase.UserRepository.GetOneByEmail(begin, email)
	if GetOneByEmailErr != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("UserUseCase GetOneByEmail is failed, GetUser failed : %s", GetOneByEmailErr)
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	if GetOneByEmail == nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("User UseCase FindOneByemail is failed, User is not found by email %s", email)
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &model_response.Response[*entity.User]{
		Code:    http.StatusOK,
		Message: "User UseCase FindOneById is succeed.",
		Data:    GetOneByEmail,
	}
	return result, commit
}

func (userUseCase *UserUseCase) UpdateUser(userId string, request *model_request.UserPatchOneByIdRequest) (result *model_response.Response[*entity.User], err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase UpdateUser is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	foundUser, err := userUseCase.UserRepository.GetOneById(begin, userId)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase UpdateUser is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	if foundUser == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: "UserUserCase UpdateUser is failed, User is not found by id " + userId,
			Data:    nil,
		}
		return result, rollback
	}
	if request.Name.Valid {
		foundUser.Name = request.Name
	}
	if request.Email.Valid {
		foundUser.Email = request.Email
	}
	if request.Balance.Valid {
		foundUser.Balance = request.Balance
	}
	if request.Password.Valid {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password.String), bcrypt.DefaultCost)
		if err != nil {
			rollback := begin.Rollback()
			result = &model_response.Response[*entity.User]{
				Code:    http.StatusInternalServerError,
				Message: "UserUseCase UpdateUser is failed, password hashing is failed, " + err.Error(),
				Data:    nil,
			}
			return result, rollback
		}

		foundUser.Password = null.NewString(string(hashedPassword), true)
	}
	if request.Balance.Valid {
		foundUser.Balance = request.Balance
	}

	foundUser.UpdatedAt = null.NewTime(time.Now(), true)

	patchedUser, err := userUseCase.UserRepository.PatchOneById(begin, userId, foundUser)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase UpdateUser is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &model_response.Response[*entity.User]{
		Code:    http.StatusOK,
		Message: "UserUserCase UpdateUser is succeed.",
		Data:    patchedUser,
	}
	return result, commit
}

func (userUseCase *UserUseCase) CreateUser(request *model_request.CreateUser) (result *model_response.Response[*entity.User], err error) {

	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase Register is failed, begin fail," + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(request.Password.String), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		err = begin.Rollback()
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase Register is failed, password hashing is failed.",
			Data:    nil,
		}
		return result, err
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
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase Register is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &model_response.Response[*entity.User]{
		Code:    http.StatusCreated,
		Message: "UserUseCase Register is succeed.",
		Data:    createdUser,
	}
	return result, commit
}
func (userUseCase *UserUseCase) DeleteUser(id string) (result *model_response.Response[*entity.User], err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		return result, err
	}

	deletedUser, deletedUserErr := userUseCase.UserRepository.DeleteUser(begin, id)
	if deletedUserErr != nil {
		err = begin.Rollback()
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: "UserUserCase DeleteUser is failed, " + deletedUserErr.Error(),
			Data:    nil,
		}
		return result, err
	}
	if deletedUser == nil {
		err = begin.Rollback()
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: "UserUserCase DeleteUser is failed, user is not deleted by id, " + id,
			Data:    nil,
		}
		return result, err
	}

	err = begin.Commit()
	result = &model_response.Response[*entity.User]{
		Code:    http.StatusOK,
		Message: "UserUserCase DeleteUser is succeed.",
		Data:    deletedUser,
	}
	return result, err
}

func (userUseCase *UserUseCase) ListUser() (result *model_response.Response[[]*entity.User], err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("begin failed :%s", err)
		result = &model_response.Response[[]*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}

	ListUser, err := userUseCase.UserRepository.ListUser(begin)
	if err != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("UserUseCase ListUser is failed, query failed : %s", err)
		result = &model_response.Response[[]*entity.User]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}

	if ListUser.Data == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[[]*entity.User]{
			Code:    http.StatusNotFound,
			Message: "User UseCase ListUser is failed, data User is empty ",
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &model_response.Response[[]*entity.User]{
		Code:    http.StatusOK,
		Message: "User UseCase ListUser is succeed.",
		Data:    ListUser.Data,
	}
	return result, commit
}
