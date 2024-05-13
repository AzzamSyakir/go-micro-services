package use_case

import (
	"fmt"
	"go-micro-services/src/product-service/config"
	"go-micro-services/src/product-service/entity"
	model_request "go-micro-services/src/product-service/model/request/controller"
	model_response "go-micro-services/src/product-service/model/response"
	"go-micro-services/src/product-service/repository"
	"net/http"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/google/uuid"
	"github.com/guregu/null"
)

type CategoryUseCase struct {
	DatabaseConfig     *config.DatabaseConfig
	CategoryRepository *repository.CategoryRepository
}

func NewCategoryUseCase(
	databaseConfig *config.DatabaseConfig,
	categoryRepository *repository.CategoryRepository,

) *CategoryUseCase {
	categoryUseCase := &CategoryUseCase{
		DatabaseConfig:     databaseConfig,
		CategoryRepository: categoryRepository,
	}
	return categoryUseCase
}
func (categoryUseCase *CategoryUseCase) CreateCategory(request *model_request.CategoryRequest) (result *model_response.Response[*entity.Category], err error) {

	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusCreated,
			Message: "CategoryUseCase addCategory is failed, begin, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	currentTime := null.NewTime(time.Now(), true)
	newCategory := &entity.Category{
		Id:        null.NewString(uuid.NewString(), true),
		Name:      request.Name,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		DeletedAt: null.NewTime(time.Time{}, false),
	}

	createdCategory, err := categoryUseCase.CategoryRepository.CreateCategory(begin, newCategory)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusCreated,
			Message: "CategoryUseCase addCategory is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	commit := begin.Commit()
	result = &model_response.Response[*entity.Category]{
		Code:    http.StatusCreated,
		Message: "CategoryUseCase Register is succeed.",
		Data:    createdCategory,
	}
	return result, commit
}

func (categoryUseCase *CategoryUseCase) GetOneById(id string) (result *model_response.Response[*entity.Category]) {
	transaction, transactionErr := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if transactionErr != nil {
		errorMessage := fmt.Sprintf("transaction failed :%s", transactionErr)
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}

		return result
	}
	categoryFound, categoryFoundErr := categoryUseCase.CategoryRepository.GetOneById(transaction, id)
	if categoryFoundErr != nil {
		errorMessage := fmt.Sprintf("categoryUseCase GetOneById is failed, Getcategory failed : %s", categoryFoundErr)
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}

		return result
	}
	errorMessage := fmt.Sprintf("categoryUseCase GetOneById is failed, category is not found by id %s", id)
	if categoryFound == nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}

		return result
	}

	result = &model_response.Response[*entity.Category]{
		Code:    http.StatusOK,
		Message: "CategoryUseCase GetOneById is succeed.",
		Data:    categoryFound,
	}

	return result
}

func (categoryUseCase *CategoryUseCase) UpdateCategory(id string, request *model_request.CategoryRequest) (result *model_response.Response[*entity.Category]) {
	beginErr := crdb.Execute(func() (err error) {
		begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
		if err != nil {
			return err
		}

		foundCategory, err := categoryUseCase.CategoryRepository.GetOneById(begin, id)
		if err != nil {
			return err
		}
		if foundCategory == nil {
			err = begin.Rollback()
			result = &model_response.Response[*entity.Category]{
				Code:    http.StatusNotFound,
				Message: "CategoryUseCase Update Category is failed, category is not found by id.",
				Data:    nil,
			}
			return err
		}

		if request.Name.Valid {
			foundCategory.Name = request.Name
		}
		foundCategory.UpdatedAt = null.NewTime(time.Now(), true)

		patchedcategory, err := categoryUseCase.CategoryRepository.PatchOneById(begin, id, foundCategory)
		if err != nil {
			return err
		}

		err = begin.Commit()
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusOK,
			Message: "CategoryUseCase Update Category is succeed.",
			Data:    patchedcategory,
		}
		return err
	})

	if beginErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusInternalServerError,
			Message: "CategoryUseCase Update Category  is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}
	return result
}

func (categoryUseCase *CategoryUseCase) ListCategories() (result *model_response.Response[[]*entity.Category]) {
	transaction, transactionErr := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if transactionErr != nil {
		errorMessage := fmt.Sprintf("transaction failed :%s", transactionErr)
		result = &model_response.Response[[]*entity.Category]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}

		return result
	}

	listCategories, listCategoriesErr := categoryUseCase.CategoryRepository.ListCategories(transaction)
	if listCategoriesErr != nil {
		errorMessage := fmt.Sprintf("categoryUseCase ListCategory is failed, Get data category  failed : %s", listCategoriesErr)
		result = &model_response.Response[[]*entity.Category]{
			Code:    http.StatusNotFound,
			Message: errorMessage,
			Data:    nil,
		}

		return result
	}

	if listCategories.Data == nil {
		result = &model_response.Response[[]*entity.Category]{
			Code:    http.StatusNotFound,
			Message: "category UseCase ListCategories is failed, data category is empty ",
			Data:    nil,
		}

		return result
	}

	result = &model_response.Response[[]*entity.Category]{
		Code:    http.StatusOK,
		Message: "category UseCase Listcategories is succeed.",
		Data:    listCategories.Data,
	}

	return result
}

func (categoryUseCase *CategoryUseCase) DeleteCategory(id string) (result *model_response.Response[*entity.Category]) {
	beginErr := crdb.Execute(func() (err error) {
		begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
		if err != nil {
			return err
		}

		deletedcategory, deletedcategoryErr := categoryUseCase.CategoryRepository.DeleteOneById(begin, id)
		if deletedcategoryErr != nil {
			err = begin.Rollback()
			result = &model_response.Response[*entity.Category]{
				Code:    http.StatusNotFound,
				Message: "CategoryUseCase Deletecategory is failed, " + deletedcategoryErr.Error(),
				Data:    nil,
			}
			return err
		}
		if deletedcategory == nil {
			err = begin.Rollback()
			result = &model_response.Response[*entity.Category]{
				Code:    http.StatusNotFound,
				Message: "CategoryUseCase Deletecategory is failed, category is not deleted by id, " + id,
				Data:    nil,
			}
			return err
		}

		err = begin.Commit()
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusOK,
			Message: "CategoryUseCase Deletecategory is succeed.",
			Data:    deletedcategory,
		}
		return err
	})

	if beginErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusInternalServerError,
			Message: "CategoryUseCase Deletecategory is failed, " + beginErr.Error(),
			Data:    nil,
		}
	}

	return result
}
