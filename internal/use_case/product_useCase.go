package use_case

import (
	"context"
	"go-micro-services/internal/config"
	"go-micro-services/internal/entity"
	"go-micro-services/internal/model"
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
func (productUseCase *ProductUseCase) GetOneById(ctx context.Context, id string) (result *response.Response[*entity.Product], err error) {
	transaction := ctx.Value("transaction").(*model.Transaction)

	productFound, productFoundErr := productUseCase.ProductRepository.GetOneById(transaction.Tx, id)
	if productFoundErr != nil {
		transaction.TxErr = productFoundErr
		result = nil
		err = productFoundErr
		return result, err
	}
	if productFound == nil {
		rollbackErr := transaction.Tx.Rollback()
		if rollbackErr != nil {
			transaction.TxErr = rollbackErr
			result = nil
			err = rollbackErr
			return result, err
		}
		result = &response.Response[*entity.Product]{
			Code:    http.StatusNotFound,
			Message: "ProductUseCase FindOneById is failed, product is not found by id.",
			Data:    nil,
		}
		err = nil
		return result, err
	}

	result = &response.Response[*entity.Product]{
		Code:    http.StatusOK,
		Message: "ProductUseCase FindOneById is succeed.",
		Data:    productFound,
	}
	err = nil
	return result, err
}
