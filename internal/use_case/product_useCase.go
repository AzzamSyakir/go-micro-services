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
	productRepository *repository.ProductRepository
}

func NewProductUseCase(productRepository *repository.ProductRepository) *ProductUseCase {
	productUseCase := &ProductUseCase{
		productRepository: productRepository,
	}
	return productUseCase
}

func (productUseCase *ProductUseCase) GetOneById(ctx context.Context, id string) (result *response.Response[*entity.User], err error) {
	transaction := ctx.Value("transaction").(*model.Transaction)

	foundUser, foundUserErr := productUseCase.productRepository.GetOneById(transaction.Tx, id)
	if foundUserErr != nil {
		transaction.TxErr = foundUserErr
		result = nil
		err = foundUserErr
		return result, err
	}
	if foundUser == nil {
		rollbackErr := transaction.Tx.Rollback()
		if rollbackErr != nil {
			transaction.TxErr = rollbackErr
			result = nil
			err = rollbackErr
			return result, err
		}
		result = &response.Response[*entity.User]{
			Code:    http.StatusNotFound,
			Message: "UserUserCase FindOneById is failed, user is not found by id.",
			Data:    nil,
		}
		err = nil
		return result, err
	}

	result = &response.Response[*entity.User]{
		Code:    http.StatusOK,
		Message: "UserUserCase FindOneById is succeed.",
		Data:    foundUser,
	}
	err = nil
	return result, err
}
