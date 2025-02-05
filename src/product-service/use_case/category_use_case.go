package use_case

import (
	"context"
	"fmt"
	"go-micro-services/grpc/pb"
	"go-micro-services/src/product-service/config"
	"go-micro-services/src/product-service/repository"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CategoryUseCase struct {
	pb.UnimplementedCategoryServiceServer
	DatabaseConfig     *config.DatabaseConfig
	CategoryRepository *repository.CategoryRepository
}

func NewCategoryUseCase(
	databaseConfig *config.DatabaseConfig,
	categoryRepository *repository.CategoryRepository,

) *CategoryUseCase {
	categoryUseCase := &CategoryUseCase{
		UnimplementedCategoryServiceServer: pb.UnimplementedCategoryServiceServer{},
		DatabaseConfig:                     databaseConfig,
		CategoryRepository:                 categoryRepository,
	}
	return categoryUseCase
}
func (categoryUseCase *CategoryUseCase) GetCategoryById(ctx context.Context, id *pb.ById) (result *pb.CategoryResponse, err error) {
	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to retrieve category: Unable to start database transaction. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	categoryFound, err := categoryUseCase.CategoryRepository.GetProductById(begin, id.Id)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to retrieve category: Database query error while fetching category with ID %s. Error: %v. Rollback status: %v", id.Id, err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if categoryFound == nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("Category not found: No category exists with ID %s. Rollback status: %v", id.Id, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if commitErr := begin.Commit(); commitErr != nil {
		return &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize the transaction while retrieving category with ID %s. Error: %v", id.Id, commitErr),
			Data:    nil,
		}, commitErr
	}

	return &pb.CategoryResponse{
		Code:    int64(codes.OK),
		Message: fmt.Sprintf("Category retrieved successfully. Category ID: %s", id.Id),
		Data:    categoryFound,
	}, nil
}
func (categoryUseCase *CategoryUseCase) UpdateCategory(ctx context.Context, request *pb.UpdateCategoryRequest) (result *pb.CategoryResponse, err error) {
	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to update category: Unable to start database transaction. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	foundCategory, err := categoryUseCase.CategoryRepository.GetProductById(begin, request.Id)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to update category: Database query error while retrieving category with ID %s. Error: %v. Rollback status: %v", request.Id, err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if foundCategory == nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("Failed to update category: No category found with ID %s. Rollback status: %v", request.Id, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if request.Name != nil {
		foundCategory.Name = request.GetName()
	}
	foundCategory.UpdatedAt = timestamppb.New(time.Now())

	patchedCategory, err := categoryUseCase.CategoryRepository.PatchOneById(begin, request.Id, foundCategory)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to update category: Error updating category in the database. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if commitErr := begin.Commit(); commitErr != nil {
		return &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize the transaction while updating category with ID %s. Error: %v", request.Id, commitErr),
			Data:    nil,
		}, commitErr
	}

	return &pb.CategoryResponse{
		Code:    int64(codes.OK),
		Message: fmt.Sprintf("Category updated successfully. Category ID: %s", request.Id),
		Data:    patchedCategory,
	}, nil
}
func (categoryUseCase *CategoryUseCase) CreateCategory(ctx context.Context, request *pb.CreateCategoryRequest) (result *pb.CategoryResponse, err error) {
	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to create category: Unable to start database transaction. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if request.Name == "" {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.InvalidArgument),
			Message: fmt.Sprintf("Failed to create category: Category name is required. Rollback status: %v", rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	currentTime := time.Now()
	newCategory := &pb.Category{
		Id:        uuid.NewString(),
		Name:      request.Name,
		CreatedAt: timestamppb.New(currentTime),
		UpdatedAt: timestamppb.New(currentTime),
	}

	createdCategory, err := categoryUseCase.CategoryRepository.CreateCategory(begin, newCategory)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to create category: Database insertion error. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if commitErr := begin.Commit(); commitErr != nil {
		return &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize the transaction while creating category. Error: %v", commitErr),
			Data:    nil,
		}, commitErr
	}

	return &pb.CategoryResponse{
		Code:    int64(codes.OK),
		Message: fmt.Sprintf("Category created successfully. Category ID: %s, Name: %s", createdCategory.Id, createdCategory.Name),
		Data:    createdCategory,
	}, nil
}
func (categoryUseCase *CategoryUseCase) DeleteCategory(ctx context.Context, id *pb.ById) (result *pb.CategoryResponse, err error) {
	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to delete category: Unable to start database transaction. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if id.Id == "" {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.InvalidArgument),
			Message: fmt.Sprintf("Failed to delete category: Category ID is required. Rollback status: %v", rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	deletedCategory, err := categoryUseCase.CategoryRepository.DeleteOneById(begin, id.Id)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to delete category: Database deletion error. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if deletedCategory == nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponse{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("Failed to delete category: No category found with ID %s. Rollback status: %v", id.Id, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if commitErr := begin.Commit(); commitErr != nil {
		return &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize the transaction while deleting category. Error: %v", commitErr),
			Data:    nil,
		}, commitErr
	}

	return &pb.CategoryResponse{
		Code:    int64(codes.OK),
		Message: fmt.Sprintf("Category deleted successfully. Category ID: %s", deletedCategory.Id),
		Data:    deletedCategory,
	}, nil
}
func (categoryUseCase *CategoryUseCase) ListCategories(ctx context.Context, _ *pb.Empty) (result *pb.CategoryResponseRepeated, err error) {
	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponseRepeated{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to retrieve categories: Unable to start database transaction. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	listCategories, err := categoryUseCase.CategoryRepository.ListCategories(begin)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponseRepeated{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to retrieve categories: Database query error. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if listCategories == nil || len(listCategories.Data) == 0 {
		rollbackErr := begin.Rollback()
		return &pb.CategoryResponseRepeated{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("No categories found in the database. Rollback status: %v", rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if commitErr := begin.Commit(); commitErr != nil {
		return &pb.CategoryResponseRepeated{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize the transaction while retrieving categories. Error: %v", commitErr),
			Data:    nil,
		}, commitErr
	}

	return &pb.CategoryResponseRepeated{
		Code:    int64(codes.OK),
		Message: "Successfully retrieved category list.",
		Data:    listCategories.Data,
	}, nil
}
