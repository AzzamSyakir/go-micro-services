package use_case

import (
	"fmt"
	"go-micro-services/src/product-service/config"
	"go-micro-services/src/product-service/entity"
	model_request "go-micro-services/src/product-service/model/request/controller"
	model_response "go-micro-services/src/product-service/model/response"
	"go-micro-services/src/product-service/repository"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type ProductUseCase struct {
	DatabaseConfig    *config.DatabaseConfig
	ProductRepository *repository.ProductRepository
}

func NewProductUseCase(
	databaseConfig *config.DatabaseConfig,
	productRepository *repository.ProductRepository,

) *ProductUseCase {
	productUseCase := &ProductUseCase{
		DatabaseConfig:    databaseConfig,
		ProductRepository: productRepository,
	}
	return productUseCase
}
func (productUseCase *ProductUseCase) Createproduct(request *model_request.CreateProduct) (result *model_response.Response[*entity.Product]) {
	beginErr := crdb.Execute(func() (err error) {
		begin, err := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
		if err != nil {
			result = nil
			return err
		}
		if request.Name.String == "" || request.Price.Int64 == 0 || request.Stock.Int64 == 0 {
			err = begin.Rollback()
			result = &model_response.Response[*entity.Product]{
				Code:    http.StatusBadRequest,
				Message: "Please input data correctly, data cannot be empty",
				Data:    nil,
			}
			return err
		}
		firstLetter := strings.ToUpper(string(request.Name.String[0]))
		rand.Seed(time.Now().UnixNano())
		randomDigits := rand.Intn(900) + 100
		sku := fmt.Sprintf("%s%d", firstLetter, randomDigits)

		currentTime := null.NewTime(time.Now(), true)
		newproduct := &entity.Product{
			Id:         null.NewString(uuid.NewString(), true),
			Name:       request.Name,
			Sku:        null.NewString(sku, true),
			Price:      request.Price,
			Stock:      request.Stock,
			CategoryId: request.CategoryId,
			CreatedAt:  currentTime,
			UpdatedAt:  currentTime,
			DeletedAt:  null.NewTime(time.Time{}, false),
		}

		createdProduct, err := productUseCase.ProductRepository.CreateProduct(begin, newproduct)
		if err != nil {
			return err
		}

		err = begin.Commit()
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusCreated,
			Message: "productUseCase addProduct is succeed.",
			Data:    createdProduct,
		}
		return err
	})

	if beginErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusInternalServerError,
			Message: "productUseCase addProduct  is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}
	return result
}
func (productUseCase *ProductUseCase) GetOneById(id string) (result *model_response.Response[*entity.Product], err error) {
	transaction, transactionErr := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if transactionErr != nil {
		errorMessage := fmt.Sprintf("transaction failed :%s", transactionErr)
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}
	productFound, productFoundErr := productUseCase.ProductRepository.GetOneById(transaction, id)
	if productFoundErr != nil {
		errorMessage := fmt.Sprintf("ProductUseCase GetOneById is failed, GetProduct failed : %s", productFoundErr)
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}
	errorMessage := fmt.Sprintf("productUseCase FindOneById is failed, product is not found by id %s", id)
	if productFound == nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}

	result = &model_response.Response[*entity.Product]{
		Code:    http.StatusOK,
		Message: "product-service UseCase FindOneById is succeed.",
		Data:    productFound,
	}
	err = nil
	return result, err
}
func (productUseCase *ProductUseCase) PatchOneByIdFromRequest(id string, request *model_request.ProductPatchOneByIdRequest) (result *model_response.Response[*entity.Product]) {
	beginErr := crdb.Execute(func() (err error) {
		begin, err := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
		if err != nil {
			return err
		}

		foundProduct, err := productUseCase.ProductRepository.GetOneById(begin, id)
		if err != nil {
			return err
		}
		if foundProduct == nil {
			err = begin.Rollback()
			result = &model_response.Response[*entity.Product]{
				Code:    http.StatusNotFound,
				Message: "ProductProductCase PatchOneByIdFromRequest is failed, product is not found by id.",
				Data:    nil,
			}
			return err
		}

		if request.Stock.Valid {
			foundProduct.Stock = request.Stock
		}
		foundProduct.UpdatedAt = null.NewTime(time.Now(), true)

		patchedProduct, err := productUseCase.ProductRepository.PatchOneById(begin, id, foundProduct)
		if err != nil {
			return err
		}

		err = begin.Commit()
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusOK,
			Message: "ProductProductCase PatchOneByIdFromRequest is succeed.",
			Data:    patchedProduct,
		}
		return err
	})

	if beginErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusInternalServerError,
			Message: "ProductProductCase PatchOneByIdFromRequest  is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}
	return result
}
