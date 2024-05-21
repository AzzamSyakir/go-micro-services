package use_case

import (
	"context"
	"go-micro-services/src/product-service/config"
	pb "go-micro-services/src/product-service/delivery/grpc/pb/category"
	"go-micro-services/src/product-service/repository"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/google/uuid"
	"github.com/guregu/null"
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
func (categoryUseCase *CategoryUseCase) GetCategoryById(ctx context.Context, id *pb.ByIdCategory) (result *pb.CategoryResponse, err error) {
	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: "CategoryUseCase GetCategory is failed, begin fail, " + err.Error(),
			Data:    nil,
		}

		return result, rollback
	}
	categoryFound, err := categoryUseCase.CategoryRepository.GetProductById(begin, id.Id)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.CategoryResponse{
			Code:    int64(codes.Canceled),
			Message: "CategoryUseCase GetCategory is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}

		return result, rollback
	}
	rollback := begin.Rollback()
	if categoryFound == nil {
		result = &pb.CategoryResponse{
			Code:    int64(codes.Canceled),
			Message: "CategoryUseCase GetCategory is failed, category not found by id, " + id.Id,
			Data:    nil,
		}

		return result, rollback
	}
	commit := begin.Commit()
	result = &pb.CategoryResponse{
		Code:    int64(codes.OK),
		Message: "CategoryUseCase GetProductById is succeed.",
		Data:    categoryFound,
	}

	return result, commit
}

func (categoryUseCase *CategoryUseCase) UpdateCategory(ctx context.Context, request *pb.RequestUpdate) (result *pb.CategoryResponse, err error) {

	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: "CategoryUseCase UpdateCategory is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	foundCategory, err := categoryUseCase.CategoryRepository.GetProductById(begin, request.Id)
	if err != nil {
		begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
		if err != nil {
			rollback := begin.Rollback()
			result = &pb.CategoryResponse{
				Code:    int64(codes.Canceled),
				Message: "CategoryUseCase UpdateCategory is failed, query to db fail, " + err.Error(),
				Data:    nil,
			}
			return result, rollback
		}
	}
	if foundCategory == nil {
		rollback := begin.Rollback()
		result = &pb.CategoryResponse{
			Code:    int64(codes.Canceled),
			Message: "CategoryUseCase Update Category is failed, category is not found by id, " + request.Id,
			Data:    nil,
		}
		return result, rollback
	}

	if request.Name != nil {
		foundCategory.Name = request.GetName()
	}
	foundCategory.UpdatedAt = timestamppb.New(time.Now())

	patchedcategory, err := categoryUseCase.CategoryRepository.PatchOneById(begin, request.Id, foundCategory)
	if err != nil {
		begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
		if err != nil {
			rollback := begin.Rollback()
			result = &pb.CategoryResponse{
				Code:    int64(codes.OK),
				Message: "CategoryUseCase UpdateCategory is failed, query to db fail, " + err.Error(),
				Data:    nil,
			}
			return result, rollback
		}
	}

	commit := begin.Commit()
	result = &pb.CategoryResponse{
		Code:    int64(codes.OK),
		Message: "CategoryUseCase UpdateCategory is succeed.",
		Data:    patchedcategory,
	}
	return result, commit
}

func (categoryUseCase *CategoryUseCase) CreateCategory(ctx context.Context, request *pb.RequestCreate) (result *pb.CategoryResponse, err error) {

	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: "CategoryUseCase AddCategory is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	currentTime := null.NewTime(time.Now(), true)
	newCategory := &pb.Category{
		Id:        uuid.NewString(),
		Name:      request.Name,
		CreatedAt: timestamppb.New(currentTime.Time),
		UpdatedAt: timestamppb.New(currentTime.Time),
		DeletedAt: timestamppb.New(time.Time{}),
	}

	createdCategory, err := categoryUseCase.CategoryRepository.CreateCategory(begin, newCategory)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.CategoryResponse{
			Code:    int64(codes.Canceled),
			Message: "CategoryUseCase AddCategory is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &pb.CategoryResponse{
		Code:    int64(codes.Canceled),
		Message: "CategoryUseCase Register is succeed.",
		Data:    createdCategory,
	}
	return result, commit
}

func (categoryUseCase *CategoryUseCase) DeleteCategory(ctx context.Context, id *pb.ByIdCategory) (result *pb.CategoryResponse, err error) {
	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.CategoryResponse{
			Code:    int64(codes.Internal),
			Message: "CategoryUseCase DeleteCategory is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	deletedcategory, err := categoryUseCase.CategoryRepository.DeleteOneById(begin, id.Id)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.CategoryResponse{
			Code:    int64(codes.Canceled),
			Message: "CategoryUseCase DeleteCategory is failed, Query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	if deletedcategory == nil {
		rollback := begin.Rollback()
		result = &pb.CategoryResponse{
			Code:    int64(codes.Canceled),
			Message: "CategoryUseCase DeleteCategory is failed, category is not deleted by id , " + id.Id,
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &pb.CategoryResponse{
		Code:    int64(codes.OK),
		Message: "CategoryUseCase DeleteCategory is succed.",
		Data:    deletedcategory,
	}
	return result, commit
}

func (categoryUseCase *CategoryUseCase) ListCategories() (result *pb.CategoryResponseRepeated, err error) {
	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.CategoryResponseRepeated{
			Code:    int64(codes.Internal),
			Message: "CategoryUseCase ListCategory is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	listCategories, err := categoryUseCase.CategoryRepository.ListCategories(begin)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.CategoryResponseRepeated{
			Code:    int64(codes.Canceled),
			Message: "CategoryUseCase ListCategory is failed, Query to db, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	if listCategories.Data == nil {
		rollback := begin.Rollback()
		result = &pb.CategoryResponseRepeated{
			Code:    int64(codes.Canceled),
			Message: "CategoryUseCase UpdateCategory is failed, Category is empty, ",
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &pb.CategoryResponseRepeated{
		Code:    int64(codes.OK),
		Message: "CategoryUseCase ListCategory is Succed, ",
		Data:    listCategories.Data,
	}
	return result, commit
}
