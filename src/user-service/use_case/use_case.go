package use_case

import (
	"fmt"
	"go-micro-services/src/user-service/config"
	"go-micro-services/src/user-service/delivery/grpc/pb"
	"go-micro-services/src/user-service/repository"
	"net/http"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserUseCase struct {
	pb.UnimplementedUserServiceServer
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

func (userUseCase *UserUseCase) GetOneById(id string) (result *pb.UserResponse, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase GetOneById is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	GetOneById, GetOneByIdErr := userUseCase.UserRepository.GetOneById(begin, id)
	if GetOneByIdErr != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("UserUseCase GetOneById is failed, GetUser failed : %s", GetOneByIdErr)
		result = &pb.UserResponse{
			Code:    http.StatusBadRequest,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	if GetOneById == nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("User UseCase FindOneById is failed, User is not found by id %s", id)
		result = &pb.UserResponse{
			Code:    http.StatusBadRequest,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &pb.UserResponse{
		Code:    http.StatusOK,
		Message: "User UseCase FindOneById is succeed.",
		Data:    GetOneById,
	}
	return result, commit
}
func (userUseCase *UserUseCase) GetOneByEmail(email string) (result *pb.UserResponse, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase GetOneByEmail is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	GetOneByEmail, GetOneByEmailErr := userUseCase.UserRepository.GetOneByEmail(begin, email)
	if GetOneByEmailErr != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("UserUseCase GetOneByEmail is failed, GetUser failed : %s", GetOneByEmailErr)
		result = &pb.UserResponse{
			Code:    http.StatusBadRequest,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	if GetOneByEmail == nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("User UseCase FindOneByemail is failed, User is not found by email %s", email)
		result = &pb.UserResponse{
			Code:    http.StatusBadRequest,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &pb.UserResponse{
		Code:    http.StatusOK,
		Message: "User UseCase FindOneById is succeed.",
		Data:    GetOneByEmail,
	}
	return result, commit
}
func (userUseCase *UserUseCase) UpdateUser(userId string, request *pb.Update) (result *pb.UserResponse, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase UpdateUser is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	foundUser, err := userUseCase.UserRepository.GetOneById(begin, userId)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    http.StatusBadRequest,
			Message: "UserUseCase UpdateUser is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	if foundUser == nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    http.StatusBadRequest,
			Message: "UserUserCase UpdateUser is failed, User is not found by id " + userId,
			Data:    nil,
		}
		return result, rollback
	}
	if request.Name != nil {
		foundUser.Name.Value = *request.Name
	}
	if request.Email != nil {
		foundUser.Email.Value = *request.Email
	}
	if request.Balance != nil {
		foundUser.Balance.Value = *request.Balance
	}
	if request.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
		if err != nil {
			rollback := begin.Rollback()
			result = &pb.UserResponse{
				Code:    http.StatusBadRequest,
				Message: "UserUseCase UpdateUser is failed, password hashing is failed, " + err.Error(),
				Data:    nil,
			}
			return result, rollback
		}

		foundUser.Password.Value = string(hashedPassword)
	}
	if request.Balance != nil {
		foundUser.Balance.Value = *request.Balance
	}
	time := time.Now()
	foundUser.UpdatedAt = timestamppb.New(time)

	patchedUser, err := userUseCase.UserRepository.PatchOneById(begin, userId, foundUser)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase UpdateUser is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &pb.UserResponse{
		Code:    http.StatusOK,
		Message: "UserUserCase UpdateUser is succeed.",
		Data:    patchedUser,
	}
	return result, commit
}
func (userUseCase *UserUseCase) CreateUser(request *pb.Create) (result *pb.UserResponse, err error) {

	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase Register is failed, begin fail," + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(request.Password.Value), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		err = begin.Rollback()
		result = &pb.UserResponse{
			Code:    http.StatusBadRequest,
			Message: "UserUseCase Register is failed, password hashing is failed.",
			Data:    nil,
		}
		return result, err
	}

	currentTime := null.NewTime(time.Now(), true)
	newUser := &pb.User{
		Id:        uuid.NewString(),
		Name:      request.Name,
		Email:     request.Email,
		Password:  &wrappers.StringValue{Value: string(hashedPassword)},
		Balance:   request.Balance,
		CreatedAt: timestamppb.New(currentTime.Time),
		UpdatedAt: timestamppb.New(currentTime.Time),
		DeletedAt: &timestamppb.Timestamp{},
	}

	createdUser, err := userUseCase.UserRepository.CreateUser(begin, newUser)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    http.StatusInternalServerError,
			Message: "UserUseCase Register is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &pb.UserResponse{
		Code:    http.StatusCreated,
		Message: "UserUseCase Register is succeed.",
		Data:    createdUser,
	}
	return result, commit
}
func (userUseCase *UserUseCase) DeleteUser(id string) (result *pb.UserResponse, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		return result, err
	}

	deletedUser, deletedUserErr := userUseCase.UserRepository.DeleteUser(begin, id)
	if deletedUserErr != nil {
		err = begin.Rollback()
		result = &pb.UserResponse{
			Code:    http.StatusInternalServerError,
			Message: "UserUserCase DeleteUser is failed, " + deletedUserErr.Error(),
			Data:    nil,
		}
		return result, err
	}
	if deletedUser == nil {
		err = begin.Rollback()
		result = &pb.UserResponse{
			Code:    http.StatusBadRequest,
			Message: "UserUserCase DeleteUser is failed, user is not deleted by id, " + id,
			Data:    nil,
		}
		return result, err
	}

	err = begin.Commit()
	result = &pb.UserResponse{
		Code:    http.StatusOK,
		Message: "UserUserCase DeleteUser is succeed.",
		Data:    deletedUser,
	}
	return result, err
}
func (userUseCase *UserUseCase) ListUsers() (result *pb.UserResponseRepeated, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("begin failed :%s", err)
		result = &pb.UserResponseRepeated{
			Code:    http.StatusInternalServerError,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}

	ListUser, err := userUseCase.UserRepository.ListUser(begin)
	if err != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("UserUseCase ListUser is failed, query failed : %s", err)
		result = &pb.UserResponseRepeated{
			Code:    http.StatusInternalServerError,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}

	if ListUser.Data == nil {
		rollback := begin.Rollback()
		result = &pb.UserResponseRepeated{
			Code:    http.StatusNotFound,
			Message: "User UseCase ListUser is failed, data User is empty ",
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &pb.UserResponseRepeated{
		Code:    http.StatusOK,
		Message: "User UseCase ListUser is succeed.",
		Data:    ListUser.Data,
	}
	return result, commit
}
