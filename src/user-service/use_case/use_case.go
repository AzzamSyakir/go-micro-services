package use_case

import (
	"context"
	"fmt"
	"go-micro-services/grpc/pb"
	"go-micro-services/src/user-service/config"
	"go-micro-services/src/user-service/delivery/grpc/client"
	"go-micro-services/src/user-service/repository"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserUseCase struct {
	AuthClient *client.AuthServiceClient
	pb.UnimplementedUserServiceServer
	DatabaseConfig *config.DatabaseConfig
	UserRepository *repository.UserRepository
}

func NewUserUseCase(
	authClient *client.AuthServiceClient,
	databaseConfig *config.DatabaseConfig,
	userRepository *repository.UserRepository,
) *UserUseCase {
	return &UserUseCase{
		AuthClient:                     authClient,
		UnimplementedUserServiceServer: pb.UnimplementedUserServiceServer{},
		DatabaseConfig:                 databaseConfig,
		UserRepository:                 userRepository,
	}
}

func (userUseCase *UserUseCase) GetUserById(ctx context.Context, id *pb.ById) (result *pb.UserResponse, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to start database transaction. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	user, userErr := userUseCase.UserRepository.GetUserById(begin, id.Id)
	if userErr != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Canceled),
			Message: fmt.Sprintf("Failed to retrieve user by ID. Database query error: %v. Rollback status: %v", userErr, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	if user == nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("User not found for ID %s. Rollback status: %v", id.Id, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	commitErr := begin.Commit()
	if commitErr != nil {
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize transaction while retrieving user. Error: %v", commitErr),
			Data:    nil,
		}
		return result, commitErr
	}

	result = &pb.UserResponse{
		Code:    int64(codes.OK),
		Message: "User retrieved successfully.",
		Data:    user,
	}
	return result, nil
}
func (userUseCase *UserUseCase) GetUserByEmail(ctx context.Context, email *pb.ByEmail) (result *pb.UserResponse, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to start database transaction for retrieving user by email. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	user, userErr := userUseCase.UserRepository.GetUserByEmail(begin, email.Email)
	if userErr != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Canceled),
			Message: fmt.Sprintf("Failed to retrieve user by email. Database query error: %v. Rollback status: %v", userErr, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	if user == nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("User not found for email %s. Rollback status: %v", email.Email, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	commitErr := begin.Commit()
	if commitErr != nil {
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize transaction while retrieving user by email. Error: %v", commitErr),
			Data:    nil,
		}
		return result, commitErr
	}

	result = &pb.UserResponse{
		Code:    int64(codes.OK),
		Message: "User retrieved successfully by email.",
		Data:    user,
	}
	return result, nil
}
func (userUseCase *UserUseCase) UpdateUser(ctx context.Context, request *pb.UpdateUserRequest) (result *pb.UserResponse, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if request.Name == nil && request.Email == nil && request.Password == nil && request.Balance == nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.InvalidArgument),
			Message: "Update failed. At least one field (Name, Email, Password, or Balance) must be provided for update.",
			Data:    nil,
		}
		return result, rollbackErr
	}
	if err != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to start transaction for updating user. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	existingUser, err := userUseCase.UserRepository.GetUserById(begin, request.Id)
	if err != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to query user by ID during update process. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}
	if existingUser == nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("User with ID %s not found for update. Rollback status: %v", request.Id, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	if request.Name != nil {
		existingUser.Name = *request.Name
	}
	if request.Email != nil {
		existingUser.Email = *request.Email
	}
	if request.Balance != nil {
		existingUser.Balance = *request.Balance
	}
	if request.Password != nil {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(*request.Password), bcrypt.DefaultCost)
		if hashErr != nil {
			rollbackErr := begin.Rollback()
			result = &pb.UserResponse{
				Code:    int64(codes.Internal),
				Message: fmt.Sprintf("Password hashing failed during user update. Error: %v. Rollback status: %v", hashErr, rollbackErr),
				Data:    nil,
			}
			return result, rollbackErr
		}
		existingUser.Password = string(hashedPassword)
	}

	existingUser.UpdatedAt = timestamppb.New(time.Now())

	updatedUser, err := userUseCase.UserRepository.PatchOneById(begin, request.Id, existingUser)
	if err != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to update user data in database. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	commitErr := begin.Commit()
	if commitErr != nil {
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to commit transaction after user update. Error: %v", commitErr),
			Data:    nil,
		}
		return result, commitErr
	}

	result = &pb.UserResponse{
		Code:    int64(codes.OK),
		Message: "User updated successfully.",
		Data:    updatedUser,
	}
	return result, nil
}
func (userUseCase *UserUseCase) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (result *pb.UserResponse, err error) {

	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to start transaction for user registration. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}
	if request.Name == "" || request.Email == "" || request.Password == "" {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.InvalidArgument),
			Message: "Registration failed. Name, Email, and Password are required and cannot be empty.",
			Data:    nil,
		}
		return result, rollbackErr
	}
	hashedPassword, hashedPasswordErr := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if hashedPasswordErr != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Password hashing failed during user registration. Error: %v. Rollback status: %v", hashedPasswordErr, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
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
	}

	createdUser, err := userUseCase.UserRepository.CreateUser(begin, newUser)
	if err != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to insert new user into database. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	commitErr := begin.Commit()
	if commitErr != nil {
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to commit transaction after user creation. Error: %v", commitErr),
			Data:    nil,
		}
		return result, commitErr
	}

	result = &pb.UserResponse{
		Code:    int64(codes.OK),
		Message: "User successfully registered.",
		Data:    createdUser,
	}
	return result, nil
}
func (userUseCase *UserUseCase) DeleteUser(ctx context.Context, id *pb.ById) (result *pb.UserResponse, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to start transaction for user deletion. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	userId := &pb.ByUserId{Id: id.Id}
	userUseCase.AuthClient.LogoutWithUserId(userId)

	deletedUser, deletedUserErr := userUseCase.UserRepository.DeleteUser(begin, id.Id)
	if deletedUserErr != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to delete user from database. Error: %v. Rollback status: %v", deletedUserErr, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	if deletedUser == nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponse{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("User with ID %s not found. Deletion failed. Rollback status: %v", id.Id, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	commitErr := begin.Commit()
	if commitErr != nil {
		result = &pb.UserResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to commit transaction after user deletion. Error: %v", commitErr),
			Data:    nil,
		}
		return result, commitErr
	}

	result = &pb.UserResponse{
		Code:    int64(codes.OK),
		Message: fmt.Sprintf("User with ID %s has been successfully deleted.", id.Id),
		Data:    deletedUser,
	}
	return result, nil
}
func (userUseCase *UserUseCase) ListUsers(ctx context.Context, empty *pb.Empty) (result *pb.UserResponseRepeated, err error) {
	begin, err := userUseCase.DatabaseConfig.UserDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponseRepeated{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to start database transaction for retrieving users. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	users, queryErr := userUseCase.UserRepository.ListUser(begin)
	if queryErr != nil {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponseRepeated{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to retrieve users from database. Query error: %v. Rollback status: %v", queryErr, rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	if len(users.Data) == 0 {
		rollbackErr := begin.Rollback()
		result = &pb.UserResponseRepeated{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("No users found in the system. Rollback status: %v", rollbackErr),
			Data:    nil,
		}
		return result, rollbackErr
	}

	commitErr := begin.Commit()
	if commitErr != nil {
		result = &pb.UserResponseRepeated{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to commit transaction after retrieving users. Error: %v", commitErr),
			Data:    nil,
		}
		return result, commitErr
	}

	result = &pb.UserResponseRepeated{
		Code:    int64(codes.OK),
		Message: fmt.Sprintf("Successfully retrieved %d users.", len(users.Data)),
		Data:    users.Data,
	}
	return result, nil
}
