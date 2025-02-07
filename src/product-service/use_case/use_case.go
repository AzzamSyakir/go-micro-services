package use_case

import (
	"context"
	"fmt"
	"go-micro-services/grpc/pb"
	"go-micro-services/src/product-service/config"
	"go-micro-services/src/product-service/repository"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductUseCase struct {
	pb.UnimplementedProductServiceServer
	DatabaseConfig    *config.DatabaseConfig
	ProductRepository *repository.ProductRepository
}

func NewProductUseCase(
	databaseConfig *config.DatabaseConfig,
	productRepository *repository.ProductRepository,

) *ProductUseCase {
	productUseCase := &ProductUseCase{
		UnimplementedProductServiceServer: pb.UnimplementedProductServiceServer{},
		DatabaseConfig:                    databaseConfig,
		ProductRepository:                 productRepository,
	}
	return productUseCase
}
func (productUseCase *ProductUseCase) GetProductById(ctx context.Context, id *pb.ById) (result *pb.ProductResponse, err error) {
	begin, err := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to retrieve product details: Unable to start database transaction. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	productFound, err := productUseCase.ProductRepository.GetProductById(begin, id.Id)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to retrieve product details: Database query error. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if productFound == nil {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("Failed to retrieve product details: No product found with ID %s. Rollback status: %v", id.Id, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if commitErr := begin.Commit(); commitErr != nil {
		return &pb.ProductResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize transaction while retrieving product details. Error: %v", commitErr),
			Data:    nil,
		}, commitErr
	}

	return &pb.ProductResponse{
		Code:    int64(codes.OK),
		Message: fmt.Sprintf("Product details retrieved successfully for ID %s.", id.Id),
		Data:    productFound,
	}, nil
}
func (productUseCase *ProductUseCase) UpdateProduct(ctx context.Context, request *pb.UpdateProductRequest) (result *pb.ProductResponse, err error) {
	begin, err := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to update product: Unable to start database transaction. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	foundProduct, err := productUseCase.ProductRepository.GetProductById(begin, request.Id)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to update product: Database query error while retrieving product details. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if foundProduct == nil {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("Failed to update product: No product found with ID %s. Rollback status: %v", request.Id, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}
	if request.CategoryId != nil {
		rollbackErr := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    int64(codes.InvalidArgument),
			Message: "Update failed. Category Id is a must.",
			Data:    nil,
		}
		return result, rollbackErr
	}
	if request.Name == nil || request.Stock == nil || request.Price == nil {
		rollbackErr := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    int64(codes.InvalidArgument),
			Message: "Update failed. At least one field (Name, Stock, Price) must be provided for update.",
			Data:    nil,
		}
		return result, rollbackErr
	}
	if *request.Price <= 0 {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.InvalidArgument),
			Message: "Failed to update product: Price must be greater than zero.",
			Data:    nil,
		}, rollbackErr
	}
	if request.Name != nil {
		foundProduct.Name = *request.Name
	}
	if request.Stock != nil {
		foundProduct.Stock = *request.Stock
	}
	if request.Price != nil {
		foundProduct.Price = *request.Price
	}
	if request.CategoryId != nil {
		foundProduct.CategoryId = *request.CategoryId
	}
	foundProduct.UpdatedAt = timestamppb.Now()

	updatedProduct, err := productUseCase.ProductRepository.PatchOneById(begin, request.Id, foundProduct)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to update product: Database update operation failed. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if commitErr := begin.Commit(); commitErr != nil {
		return &pb.ProductResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize transaction while updating product. Error: %v", commitErr),
			Data:    nil,
		}, commitErr
	}

	return &pb.ProductResponse{
		Code:    int64(codes.OK),
		Message: fmt.Sprintf("Product updated successfully. ID: %s", request.Id),
		Data:    updatedProduct,
	}, nil
}
func (productUseCase *ProductUseCase) CreateProduct(ctx context.Context, request *pb.CreateProductRequest) (result *pb.ProductResponse, err error) {
	begin, err := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to create product: Unable to start database transaction. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}
	if request.Name == "" {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.InvalidArgument),
			Message: "Failed to create product: Product name cannot be empty.",
			Data:    nil,
		}, rollbackErr
	}
	if request.Price <= 0 {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.InvalidArgument),
			Message: "Failed to create product: Price must be greater than zero.",
			Data:    nil,
		}, rollbackErr
	}
	if request.Stock < 0 {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.InvalidArgument),
			Message: "Failed to create product: Stock cannot be negative.",
			Data:    nil,
		}, rollbackErr
	}
	if request.CategoryId == "" {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.InvalidArgument),
			Message: "Failed to create product: Category ID cannot be empty.",
			Data:    nil,
		}, rollbackErr
	}

	firstLetter := strings.ToUpper(string(request.Name[0]))
	rand.Seed(time.Now().UnixNano())
	randomDigits := rand.Intn(900) + 100
	sku := fmt.Sprintf("%s%d", firstLetter, randomDigits)

	currentTime := null.NewTime(time.Now(), true)
	newProduct := &pb.Product{
		Id:         uuid.NewString(),
		Name:       request.Name,
		Sku:        sku,
		Price:      request.Price,
		Stock:      request.Stock,
		CategoryId: request.CategoryId,
		CreatedAt:  timestamppb.New(currentTime.Time),
		UpdatedAt:  timestamppb.New(currentTime.Time),
	}

	createdProduct, err := productUseCase.ProductRepository.CreateProduct(begin, newProduct)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to create product: Database insert operation failed. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if commitErr := begin.Commit(); commitErr != nil {
		return &pb.ProductResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize transaction while creating product. Error: %v", commitErr),
			Data:    nil,
		}, commitErr
	}

	return &pb.ProductResponse{
		Code:    http.StatusCreated,
		Message: fmt.Sprintf("Product successfully created. ID: %s, SKU: %s", createdProduct.Id, createdProduct.Sku),
		Data:    createdProduct,
	}, nil
}
func (productUseCase *ProductUseCase) DeleteProduct(ctx context.Context, id *pb.ById) (result *pb.ProductResponse, err error) {
	begin, err := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to delete product: Unable to start database transaction. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if id.Id == "" {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.InvalidArgument),
			Message: "Failed to delete product: Product ID cannot be empty.",
			Data:    nil,
		}, rollbackErr
	}

	deletedProduct, err := productUseCase.ProductRepository.DeleteOneById(begin, id.Id)
	if err != nil {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to delete product: Database query error. Error: %v. Rollback status: %v", err, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if deletedProduct == nil {
		rollbackErr := begin.Rollback()
		return &pb.ProductResponse{
			Code:    int64(codes.NotFound),
			Message: fmt.Sprintf("Failed to delete product: No product found with ID %s. Rollback status: %v", id.Id, rollbackErr),
			Data:    nil,
		}, rollbackErr
	}

	if commitErr := begin.Commit(); commitErr != nil {
		return &pb.ProductResponse{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize transaction while deleting product. Error: %v", commitErr),
			Data:    nil,
		}, commitErr
	}

	return &pb.ProductResponse{
		Code:    int64(codes.OK),
		Message: fmt.Sprintf("Product with ID %s has been successfully deleted.", id.Id),
		Data:    deletedProduct,
	}, nil
}
func (productUseCase *ProductUseCase) ListProducts(ctx context.Context, request *pb.Empty) (result *pb.ProductResponseRepeated, err error) {
	begin, beginErr := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if beginErr != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("Failed to start database transaction. Error: %s. Rollback status: %v", beginErr, rollback)
		result = &pb.ProductResponseRepeated{
			Code:    int64(codes.Internal),
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}

	fetchProduct, fetchProductErr := productUseCase.ProductRepository.ListProducts(begin)
	if fetchProductErr != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("Failed to fetch product list from the database. Error: %s. Rollback status: %v", fetchProductErr, rollback)
		result = &pb.ProductResponseRepeated{
			Code:    int64(codes.Canceled),
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}

	if fetchProduct.Data == nil {
		rollback := begin.Rollback()
		result = &pb.ProductResponseRepeated{
			Code:    int64(codes.Canceled),
			Message: "Product list is empty. No products found in the database.",
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	if commitErr := commit; commitErr != nil {
		result = &pb.ProductResponseRepeated{
			Code:    int64(codes.Internal),
			Message: fmt.Sprintf("Failed to finalize transaction while fetching product list. Error: %v", commitErr),
			Data:    nil,
		}
		return result, commitErr
	}

	result = &pb.ProductResponseRepeated{
		Code:    int64(codes.OK),
		Message: "Product list fetched successfully.",
		Data:    fetchProduct.Data,
	}
	return result, nil
}
