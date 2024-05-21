package use_case

import (
	"context"
	"fmt"
	"go-micro-services/src/product-service/config"
	pb "go-micro-services/src/product-service/delivery/grpc/pb/product"
	"go-micro-services/src/product-service/repository"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
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
		rollback := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    http.StatusInternalServerError,
			Message: "product-service UseCase, DetailProduct, begin failed, " + err.Error(),
			Data:    nil,
		}

		return result, rollback
	}
	productFound, err := productUseCase.ProductRepository.GetProductById(begin, id.Id)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    http.StatusBadRequest,
			Message: "product-service UseCase, DetailProduct is failed, GetProduct failed, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	if productFound == nil {
		rollback := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    http.StatusBadRequest,
			Message: "product-service UseCase, DetailProduct is failed, product is not found by id" + id.Id,
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &pb.ProductResponse{
		Code:    http.StatusOK,
		Message: "product-service UseCase, DetailProduct is succeed.",
		Data:    productFound,
	}

	return result, commit
}

func (productUseCase *ProductUseCase) UpdateProduct(ctx context.Context, request *pb.Update) (result *pb.ProductResponse, err error) {
	begin, err := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    http.StatusInternalServerError,
			Message: "product-service UseCase, UpdateProduct fail begin is failed," + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	foundProduct, err := productUseCase.ProductRepository.GetProductById(begin, request.Id)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    http.StatusBadRequest,
			Message: "product-service UseCase Update Product is failed, product is not found by id" + request.Id,
			Data:    nil,
		}
		return result, rollback
	}
	if foundProduct == nil {
		rollback := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    http.StatusBadRequest,
			Message: "product-service UseCase Update Product is failed, product is not found by id, " + request.Id,
			Data:    nil,
		}
		return result, rollback
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

	patchedProduct, err := productUseCase.ProductRepository.PatchOneById(begin, request.Id, foundProduct)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    http.StatusBadRequest,
			Message: "product-service UseCase, Query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &pb.ProductResponse{
		Code:    http.StatusOK,
		Message: "product-service UseCase Update Product is succes.",
		Data:    patchedProduct,
	}
	return result, commit
}

func (productUseCase *ProductUseCase) CreateProduct(ctx context.Context, request *pb.Create) (result *pb.ProductResponse, err error) {
	begin, err := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    http.StatusInternalServerError,
			Message: "ProductUseCase CreateProduct is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	if request.Name == "" || request.Price == 0 || request.Stock == 0 {
		rollback := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    http.StatusBadRequest,
			Message: "ProductUseCase CreateProduct is failed, Please input data correctly, data cannot be empty",
			Data:    nil,
		}
		return result, rollback
	}
	firstLetter := strings.ToUpper(string(request.Name))
	rand.Seed(time.Now().UnixNano())
	randomDigits := rand.Intn(900) + 100
	sku := fmt.Sprintf("%s%d", firstLetter, randomDigits)

	currentTime := null.NewTime(time.Now(), true)
	newproduct := &pb.Product{
		Id:         uuid.NewString(),
		Name:       request.Name,
		Sku:        sku,
		Price:      request.Price,
		Stock:      request.Stock,
		CategoryId: request.CategoryId,
		CreatedAt:  timestamppb.New(currentTime.Time),
		UpdatedAt:  timestamppb.New(currentTime.Time),
		DeletedAt:  timestamppb.New(currentTime.Time),
	}

	createdProduct, err := productUseCase.ProductRepository.CreateProduct(begin, newproduct)
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    http.StatusBadRequest,
			Message: "ProductUseCase CreateProduct is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &pb.ProductResponse{
		Code:    http.StatusCreated,
		Message: "ProductUseCase CreateProduct is success",
		Data:    createdProduct,
	}
	return result, commit
}
func (productUseCase *ProductUseCase) DeleteProduct(ctx context.Context, id *pb.ById) (result *pb.ProductResponse, err error) {
	begin, err := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    http.StatusInternalServerError,
			Message: "product-service UseCase, DeleteProduct is failed, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	deletedproduct, deletedproductErr := productUseCase.ProductRepository.DeleteOneById(begin, id.Id)
	if deletedproductErr != nil {
		rollback := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    http.StatusBadRequest,
			Message: "product-service UseCase, DeleteProduct is failed, " + deletedproductErr.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	if deletedproduct == nil {
		rollback := begin.Rollback()
		result = &pb.ProductResponse{
			Code:    http.StatusBadRequest,
			Message: "product-service UseCase, DeleteProduct is failed, product is not deleted by id, " + id.Id,
			Data:    nil,
		}
		return result, rollback
	}
	rollback := begin.Commit()
	result = &pb.ProductResponse{
		Code:    http.StatusOK,
		Message: "product-service UseCase DeleteProduct is succeed.",
		Data:    deletedproduct,
	}
	return result, rollback
}

func (productUseCase *ProductUseCase) ListProducts(context.Context, *pb.Empty) (result *pb.ProductResponseRepeated, err error) {
	begin, beginErr := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if beginErr != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("begin failed :%s", beginErr)
		result = &pb.ProductResponseRepeated{
			Code:    http.StatusInternalServerError,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}

	fetchproduct, fetchproductErr := productUseCase.ProductRepository.ListProducts(begin)
	if fetchproductErr != nil {
		rollback := begin.Rollback()
		errorMessage := fmt.Sprintf("product-service UseCase, ListProduct is failed, Getproduct failed : %s", fetchproductErr)
		result = &pb.ProductResponseRepeated{
			Code:    http.StatusBadRequest,
			Message: errorMessage,
			Data:    nil,
		}
		return result, rollback
	}

	if fetchproduct.Data == nil {
		rollback := begin.Rollback()
		result = &pb.ProductResponseRepeated{
			Code:    http.StatusBadRequest,
			Message: "product-service UseCase, ListProduct is failed, data product is empty ",
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &pb.ProductResponseRepeated{
		Code:    http.StatusOK,
		Message: "product-service UseCase, ListProduct is succeed.",
		Data:    fetchproduct.Data,
	}
	return result, commit
}
