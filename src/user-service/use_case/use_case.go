package use_case

import (
	"context"
	"fmt"
	"go-micro-services/src/user-service/config"
	"go-micro-services/src/user-service/delivery/grpc/pb"
	"go-micro-services/src/user-service/repository"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
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
	return &UserUseCase{
		UnimplementedUserServiceServer: pb.UnimplementedUserServiceServer{},
		DatabaseConfig:                 databaseConfig,
		UserRepository:                 userRepository,
	}
}

func (userUseCase *UserUseCase) GetUserById(context context.Context, id *pb.ById) (result *pb.UserResponse, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: "UserUseCase GetUserById is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	GetUserById, GetUserByIdErr := userUseCase.UserRepository.GetUserById(begin, id.Id)
	if GetUserByIdErr != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("UserUseCase GetUserById is failed, GetUserById failed : %s", GetUserByIdErr)
		result = &pb.UserResponse{
			Code:    int64(codes.Canceled),
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	if GetUserById == nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("User UseCase GetOneById is failed, User is not found by id %s", id)
		result = &pb.UserResponse{
			Code:    int64(codes.Canceled),
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &pb.UserResponse{
		Code:    int64(codes.OK),
		Message: "User UseCase GetOneById is succeed.",
		Data:    GetUserById,
	}
	return result, commit
}
func (userUseCase *UserUseCase) GetUserByEmail(context context.Context, email *pb.ByEmail) (result *pb.UserResponse, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: "UserUseCase GetUserByEmail is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	GetUserByEmail, GetUserByEmailErr := userUseCase.UserRepository.GetUserByEmail(begin, email.Email)
	if GetUserByEmailErr != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("UserUseCase GetUserByEmail is failed, GetUserById failed : %s", GetUserByEmailErr)
		result = &pb.UserResponse{
			Code:    int64(codes.Canceled),
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	if GetUserByEmail == nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("User UseCase FindOneByemail is failed, User is not found by email %s", email)
		result = &pb.UserResponse{
			Code:    int64(codes.Canceled),
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &pb.UserResponse{
		Code:    int64(codes.OK),
		Message: "User UseCase GetOneById is succeed.",
		Data:    GetUserByEmail,
	}
	return result, commit
}
func (userUseCase *UserUseCase) UpdateUser(context context.Context, request *pb.Update) (result *pb.UserResponse, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: "UserUseCase UpdateUser is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	foundUser, err := userUseCase.UserRepository.GetUserById(begin, request.Id)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Canceled),
			Message: "UserUseCase UpdateUser is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	if foundUser == nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Canceled),
			Message: "UserUserCase UpdateUser is failed, User is not found by id " + request.Id,
			Data:    nil,
		}
		return result, rollback
	}
	if request.Name != nil {
		foundUser.Name = *request.Name
	}
	if request.Email != nil {
		foundUser.Email = *request.Email
	}
	if request.Balance != nil {
		foundUser.Balance = *request.Balance
	}
	if request.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
		if err != nil {
			rollback := begin.Rollback()
			result = &pb.UserResponse{
				Code:    int64(codes.Canceled),
				Message: "UserUseCase UpdateUser is failed, password hashing is failed, " + err.Error(),
				Data:    nil,
			}
			return result, rollback
		}

		foundUser.Password = string(hashedPassword)
	}
	if request.Balance != nil {
		foundUser.Balance = *request.Balance
	}
	time := time.Now()
	foundUser.UpdatedAt = timestamppb.New(time)
	patchedUser, err := userUseCase.UserRepository.PatchOneById(begin, request.Id, foundUser)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: "UserUseCase UpdateUser is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &pb.UserResponse{
		Code:    int64(codes.OK),
		Message: "UserUserCase UpdateUser is succeed.",
		Data:    patchedUser,
	}
	return result, commit
}
func (userUseCase *UserUseCase) CreateUser(context context.Context, request *pb.Create) (result *pb.UserResponse, err error) {

	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: "UserUseCase Register is failed, begin fail," + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		err = begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Canceled),
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
		Password:  string(hashedPassword),
		Balance:   request.Balance,
		CreatedAt: timestamppb.New(currentTime.Time),
		UpdatedAt: timestamppb.New(currentTime.Time),
		DeletedAt: &timestamppb.Timestamp{},
	}

	createdUser, err := userUseCase.UserRepository.CreateUser(begin, newUser)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: "UserUseCase Register is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &pb.UserResponse{
		Code:    int64(codes.OK),
		Message: "UserUseCase Register is succeed.",
		Data:    createdUser,
	}
	return result, commit
}
func (userUseCase *UserUseCase) DeleteUser(context context.Context, id *pb.ById) (result *pb.UserResponse, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		return result, err
	}

	deletedUser, deletedUserErr := userUseCase.UserRepository.DeleteUser(begin, id.Id)
	if deletedUserErr != nil {
		err = begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: "UserUserCase DeleteUser is failed, " + deletedUserErr.Error(),
			Data:    nil,
		}
		return result, err
	}
	if deletedUser == nil {
		err = begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Canceled),
			Message: "UserUserCase DeleteUser is failed, user is not deleted by id, " + id.Id,
			Data:    nil,
		}
		return result, err
	}

	err = begin.Commit()
	result = &pb.UserResponse{
		Code:    int64(codes.OK),
		Message: "UserUserCase DeleteUser is succeed.",
		Data:    deletedUser,
	}
	return result, err
}
func (userUseCase *UserUseCase) ListUsers(ctx context.Context, empty *pb.Empty) (result *pb.UserResponseRepeated, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("begin failed :%s", err)
		result = &pb.UserResponseRepeated{
			Code:    int64(codes.Internal),
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
			Code:    int64(codes.Internal),
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}

	if ListUser.Data == nil {
		rollback := begin.Rollback()
		result = &pb.UserResponseRepeated{
			Code:    int64(codes.Canceled),
			Message: "User UseCase ListUser is failed, data User is empty ",
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &pb.UserResponseRepeated{
		Code:    int64(codes.OK),
		Message: "User UseCase ListUser is succeed.",
		Data:    ListUser.Data,
	}
	return result, commit
}
