package use_case

import (
	"fmt"
	"go-micro-services/internal/config"
	"go-micro-services/internal/entity"
	"go-micro-services/internal/model/response"
	"go-micro-services/internal/repository"
	"net/http"
)

type ProductUseCase struct {
	DatabaseConfig    *config.DatabaseConfig
	ProductRepository *repository.ProductRepository
}

func NewProductUseCase(
	databaseConfig *config.DatabaseConfig,
	userRepository *repository.ProductRepository,

) *ProductUseCase {
	userUseCase := &ProductUseCase{
		DatabaseConfig:    databaseConfig,
		ProductRepository: userRepository,
	}
	return userUseCase
}
func (productUseCase *ProductUseCase) GetOneById(id string) (result *response.Response[*entity.Product], err error) {
	transaction, transactionErr := productUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if transactionErr != nil {
		errorMessage := fmt.Sprintf("transaction failed :", transactionErr)
		result = &response.Response[*entity.Product]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}
		err = nil
		return result, err
	}
	productFound, productFoundErr := productUseCase.ProductRepository.GetOneById(transaction, id)
	if productFoundErr != nil {
		result = nil
		err = productFoundErr
		return result, err
	}
	if productFound == nil {
		result = &response.Response[*entity.Product]{
			Code:    http.StatusNotFound,
			Message: "Product UseCase FindOneById is failed, product is not found by id.",
			Data:    nil,
		}
		err = nil
		return result, err
	}

	result = &response.Response[*entity.Product]{
		Code:    http.StatusOK,
		Message: "Product UseCase FindOneById is succeed.",
		Data:    productFound,
	}
	err = nil
	return result, err
}
