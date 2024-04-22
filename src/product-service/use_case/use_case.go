package use_case

import (
	"fmt"
	"go-micro-services/src/Order/entity"
	"go-micro-services/src/product-service/config"
	model_request "go-micro-services/src/product-service/model/request/controller"
	"go-micro-services/src/product-service/model/response"
	"go-micro-services/src/product-service/repository"
	"net/http"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
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
func (productUseCase *ProductUseCase) GetOneById(id string) (result *response.Response[*entity.product], err error) {
	transaction, transactionErr := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if transactionErr != nil {
		errorMessage := fmt.Sprintf("transaction failed :%s", transactionErr)
		result = &response.Response[*entity.product-service]{
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
		result = &response.Response[*entity.product-service]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}
	errorMessage := fmt.Sprintf("productUseCase FindOneById is failed, product is not found by id %s", id)
	if productFound == nil {
		result = &response.Response[*entity.product-service]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}

	result = &response.Response[*entity.product-service]{
		Code:    http.StatusOK,
		Message: "product-service UseCase FindOneById is succeed.",
		Data:    productFound,
	}
	err = nil
	return result, err
}
func (productUseCase *ProductUseCase) PatchOneByIdFromRequest(id string, request *model_request.ProductPatchOneByIdRequest) (result *response.Response[*entity.Product]) {
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
			result = &response.Response[*entity.product-service]{
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
		result = &response.Response[*entity.product-service]{
			Code:    http.StatusOK,
			Message: "ProductProductCase PatchOneByIdFromRequest is succeed.",
			Data:    patchedProduct,
		}
		return err
	})

	if beginErr != nil {
		result = &response.Response[*entity.product-service]{
			Code:    http.StatusInternalServerError,
			Message: "ProductProductCase PatchOneByIdFromRequest  is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}
	return result
}
