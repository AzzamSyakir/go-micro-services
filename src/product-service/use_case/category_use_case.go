package use_case

import (
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
			Message: "CategoryUseCase AddCategory is failed, begin fail, " + err.Error(),
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
			Message: "CategoryUseCase AddCategory is failed, query to db fail, " + err.Error(),
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

func (categoryUseCase *CategoryUseCase) GetOneById(id string) (result *model_response.Response[*entity.Category], err error) {
	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusNotFound,
			Message: "CategoryUseCase GetCategory is failed, begin fail, " + err.Error(),
			Data:    nil,
		}

		return result, rollback
	}
	categoryFound, err := categoryUseCase.CategoryRepository.GetOneById(begin, id)
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusNotFound,
			Message: "CategoryUseCase GetCategory is failed, query to db fail, " + err.Error(),
			Data:    nil,
		}

		return result, rollback
	}
	rollback := begin.Rollback()
	if categoryFound == nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusNotFound,
			Message: "CategoryUseCase GetCategory is failed, category not found by id, " + id,
			Data:    nil,
		}

		return result, rollback
	}
	commit := begin.Commit()
	result = &model_response.Response[*entity.Category]{
		Code:    http.StatusOK,
		Message: "CategoryUseCase GetOneById is succeed.",
		Data:    categoryFound,
	}

	return result, commit
}

func (categoryUseCase *CategoryUseCase) UpdateCategory(id string, request *model_request.CategoryRequest) (result *model_response.Response[*entity.Category], err error) {

	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusCreated,
			Message: "CategoryUseCase UpdateCategory is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	foundCategory, err := categoryUseCase.CategoryRepository.GetOneById(begin, id)
	if err != nil {
		begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
		if err != nil {
			rollback := begin.Rollback()
			result = &model_response.Response[*entity.Category]{
				Code:    http.StatusCreated,
				Message: "CategoryUseCase UpdateCategory is failed, query to db fail, " + err.Error(),
				Data:    nil,
			}
			return result, rollback
		}
	}
	if foundCategory == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusNotFound,
			Message: "CategoryUseCase Update Category is failed, category is not found by id, " + id,
			Data:    nil,
		}
		return result, rollback
	}

	if request.Name.Valid {
		foundCategory.Name = request.Name
	}
	foundCategory.UpdatedAt = null.NewTime(time.Now(), true)

	patchedcategory, err := categoryUseCase.CategoryRepository.PatchOneById(begin, id, foundCategory)
	if err != nil {
		begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
		if err != nil {
			rollback := begin.Rollback()
			result = &model_response.Response[*entity.Category]{
				Code:    http.StatusCreated,
				Message: "CategoryUseCase UpdateCategory is failed, query to db fail, " + err.Error(),
				Data:    nil,
			}
			return result, rollback
		}
	}

	commit := begin.Commit()
	result = &model_response.Response[*entity.Category]{
		Code:    http.StatusOK,
		Message: "CategoryUseCase UpdateCategory is succeed.",
		Data:    patchedcategory,
	}
	return result, commit
}

func (categoryUseCase *CategoryUseCase) ListCategories() (result *model_response.Response[[]*entity.Category], err error) {
	begin, err := categoryUseCase.DatabaseConfig.ProductDB.Connection.Begin()
	if err != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[[]*entity.Category]{
			Code:    http.StatusCreated,
			Message: "CategoryUseCase ListCategory is failed, begin fail, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	listCategories, listCategoriesErr := categoryUseCase.CategoryRepository.ListCategories(begin)
	if listCategoriesErr != nil {
		rollback := begin.Rollback()
		result = &model_response.Response[[]*entity.Category]{
			Code:    http.StatusCreated,
			Message: "CategoryUseCase ListCategory is failed, Query to db, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}

	if listCategories.Data == nil {
		rollback := begin.Rollback()
		result = &model_response.Response[[]*entity.Category]{
			Code:    http.StatusCreated,
			Message: "CategoryUseCase UpdateCategory is failed, Category is empty, " + err.Error(),
			Data:    nil,
		}
		return result, rollback
	}
	commit := begin.Commit()
	result = &model_response.Response[[]*entity.Category]{
		Code:    http.StatusCreated,
		Message: "CategoryUseCase ListCategory is Succed, ",
		Data:    listCategories.Data,
	}
	return result, commit
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
