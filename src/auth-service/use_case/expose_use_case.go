package use_case

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-micro-services/src/auth-service/client"
	"go-micro-services/src/auth-service/config"
	"go-micro-services/src/auth-service/delivery/grpc/pb"
	"go-micro-services/src/auth-service/entity"
	model_request "go-micro-services/src/auth-service/model/request/controller"
	model_response "go-micro-services/src/auth-service/model/response"
	"go-micro-services/src/auth-service/repository"
	"net/http"

	"github.com/guregu/null"
)

type ExposeUseCase struct {
	DatabaseConfig *config.DatabaseConfig
	AuthRepository *repository.AuthRepository
	Env            *config.EnvConfig
	userClient     *client.UserServiceClient
	productClient  *client.ProductServiceClient
	OrderClient    *client.OrderServiceClient
	CategoryClient *client.CategoryServiceClient
}

func NewExposeUseCase(
	databaseConfig *config.DatabaseConfig,
	authRepository *repository.AuthRepository,
	env *config.EnvConfig,
	initUserClient *client.UserServiceClient,
	initProductClient *client.ProductServiceClient,
	initOrderClient *client.OrderServiceClient,
	initCategoryClient *client.CategoryServiceClient,
) *ExposeUseCase {
	userUseCase := &ExposeUseCase{
		userClient:     initUserClient,
		productClient:  initProductClient,
		OrderClient:    initOrderClient,
		CategoryClient: initCategoryClient,
		DatabaseConfig: databaseConfig,
		AuthRepository: authRepository,
		Env:            env,
	}
	return userUseCase
}

// users
func (exposeUseCase *ExposeUseCase) ListUsers() (result *model_response.Response[[]*entity.User]) {
	ListUser, err := exposeUseCase.userClient.ListUsers()
	if err != nil {
		result = &model_response.Response[[]*entity.User]{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
		return result
	}
	var users []*entity.User
	for _, user := range ListUser.Data {
		userData := &entity.User{
			Id:        null.NewString(user.Id, true),
			Name:      null.NewString(user.Name, true),
			Email:     null.NewString(user.Email, true),
			Password:  null.NewString(user.Password, true),
			Balance:   null.NewInt(user.Balance, true),
			CreatedAt: null.NewTime(user.CreatedAt.AsTime(), true),
			UpdatedAt: null.NewTime(user.UpdatedAt.AsTime(), true),
			DeletedAt: null.NewTime(user.DeletedAt.AsTime(), true),
		}

		users = append(users, userData)
	}
	bodyResponseUser := &model_response.Response[[]*entity.User]{
		Code:    http.StatusOK,
		Message: ListUser.Message,
		Data:    users,
	}
	return bodyResponseUser
}
func (exposeUseCase *ExposeUseCase) CreateUser(request *model_request.RegisterRequest) (result *model_response.Response[*entity.User]) {
	req := &pb.CreateUserRequest{
		Name:     request.Name.String,
		Email:    request.Email.String,
		Password: request.Password.String,
		Balance:  request.Balance.Int64,
	}
	createUser, err := exposeUseCase.userClient.CreateUser(req)
	if err != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Data:    nil,
			Message: createUser.Message,
		}
		return
	}
	if createUser.Data == nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Data:    nil,
			Message: createUser.Message,
		}
		return
	}
	user := entity.User{
		Id:        null.NewString(createUser.Data.Id, true),
		Name:      null.NewString(createUser.Data.Name, true),
		Email:     null.NewString(createUser.Data.Email, true),
		Password:  null.NewString(createUser.Data.Password, true),
		Balance:   null.NewInt(createUser.Data.Balance, true),
		CreatedAt: null.NewTime(createUser.Data.CreatedAt.AsTime(), true),
		UpdatedAt: null.NewTime(createUser.Data.UpdatedAt.AsTime(), true),
	}
	bodyResponseUser := &model_response.Response[*entity.User]{
		Code:    http.StatusOK,
		Message: createUser.Message,
		Data:    &user,
	}
	return bodyResponseUser
}
func (exposeUseCase *ExposeUseCase) DeleteUser(id string) (result *model_response.Response[*entity.User]) {
	DeleteUser, err := exposeUseCase.userClient.DeleteUser(id)
	if err != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: DeleteUser.Message,
			Data:    nil,
		}
		return
	}
	if DeleteUser.Data == nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: DeleteUser.Message,
			Data:    nil,
		}
		return
	}
	user := entity.User{
		Id:        null.NewString(DeleteUser.Data.Id, true),
		Name:      null.NewString(DeleteUser.Data.Name, true),
		Email:     null.NewString(DeleteUser.Data.Email, true),
		Password:  null.NewString(DeleteUser.Data.Password, true),
		Balance:   null.NewInt(DeleteUser.Data.Balance, true),
		CreatedAt: null.NewTime(DeleteUser.Data.CreatedAt.AsTime(), true),
		UpdatedAt: null.NewTime(DeleteUser.Data.UpdatedAt.AsTime(), true),
	}
	bodyResponseUser := &model_response.Response[*entity.User]{
		Code:    http.StatusOK,
		Message: DeleteUser.Message,
		Data:    &user,
	}
	return bodyResponseUser
}
func (exposeUseCase *ExposeUseCase) UpdateUser(id string, request *model_request.UserPatchOneByIdRequest) (result *model_response.Response[*entity.User]) {
	req := &pb.UpdateUserRequest{}
	if id != "" {
		req.Id = id
	}
	if request.Name.Valid {
		req.Name = &request.Name.String
	}
	if request.Email.Valid {
		req.Email = &request.Email.String
	}
	if request.Password.Valid {
		req.Password = &request.Password.String
	}
	if request.Balance.Valid {
		req.Balance = &request.Balance.Int64
	}
	UpdateUser, err := exposeUseCase.userClient.UpdateUser(req)
	if err != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: UpdateUser.Message,
			Data:    nil,
		}
		return
	}
	if UpdateUser.Data == nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: UpdateUser.Message,
			Data:    nil,
		}
		return
	}
	user := entity.User{
		Id:        null.NewString(UpdateUser.Data.Id, true),
		Name:      null.NewString(UpdateUser.Data.Name, true),
		Email:     null.NewString(UpdateUser.Data.Email, true),
		Password:  null.NewString(UpdateUser.Data.Password, true),
		Balance:   null.NewInt(UpdateUser.Data.Balance, true),
		CreatedAt: null.NewTime(UpdateUser.Data.CreatedAt.AsTime(), true),
		UpdatedAt: null.NewTime(UpdateUser.Data.UpdatedAt.AsTime(), true),
	}
	bodyResponseUser := &model_response.Response[*entity.User]{
		Code:    http.StatusOK,
		Message: UpdateUser.Message,
		Data:    &user,
	}
	return bodyResponseUser
}
func (exposeUseCase *ExposeUseCase) DetailUser(id string) (result *model_response.Response[*entity.User]) {
	GetUser, err := exposeUseCase.userClient.GetUserById(id)
	if err != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: GetUser.Message,
			Data:    nil,
		}
		return
	}
	if GetUser.Data == nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: GetUser.Message,
			Data:    nil,
		}
		return
	}
	user := entity.User{
		Id:        null.NewString(GetUser.Data.Id, true),
		Name:      null.NewString(GetUser.Data.Name, true),
		Email:     null.NewString(GetUser.Data.Email, true),
		Password:  null.NewString(GetUser.Data.Password, true),
		Balance:   null.NewInt(GetUser.Data.Balance, true),
		CreatedAt: null.NewTime(GetUser.Data.CreatedAt.AsTime(), true),
		UpdatedAt: null.NewTime(GetUser.Data.UpdatedAt.AsTime(), true),
	}
	bodyResponseUser := &model_response.Response[*entity.User]{
		Code:    http.StatusOK,
		Message: GetUser.Message,
		Data:    &user,
	}
	return bodyResponseUser
}
func (exposeUseCase *ExposeUseCase) GetOneByEmail(email string) (result *model_response.Response[*entity.User]) {
	GetUser, err := exposeUseCase.userClient.GetUserByEmail(email)
	if err != nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: GetUser.Message,
			Data:    nil,
		}
		return
	}
	if GetUser.Data == nil {
		result = &model_response.Response[*entity.User]{
			Code:    http.StatusBadRequest,
			Message: GetUser.Message,
			Data:    nil,
		}
		return
	}
	user := entity.User{
		Id:        null.NewString(GetUser.Data.Id, true),
		Name:      null.NewString(GetUser.Data.Name, true),
		Email:     null.NewString(GetUser.Data.Email, true),
		Password:  null.NewString(GetUser.Data.Password, true),
		Balance:   null.NewInt(GetUser.Data.Balance, true),
		CreatedAt: null.NewTime(GetUser.Data.CreatedAt.AsTime(), true),
		UpdatedAt: null.NewTime(GetUser.Data.UpdatedAt.AsTime(), true),
	}
	bodyResponseUser := &model_response.Response[*entity.User]{
		Code:    http.StatusOK,
		Message: GetUser.Message,
		Data:    &user,
	}
	return bodyResponseUser
}

// product

func (exposeUseCase *ExposeUseCase) ListProducts() (result *model_response.Response[[]*entity.Product]) {
	ListProduct, err := exposeUseCase.productClient.ListProducts()
	if err != nil {
		result = &model_response.Response[[]*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
		return result
	}
	var products []*entity.Product
	for _, product := range ListProduct.Data {
		productData := &entity.Product{
			Id:         null.NewString(product.Id, true),
			Sku:        null.NewString(product.Sku, true),
			Name:       null.NewString(product.Name, true),
			Stock:      null.NewInt(product.Stock, true),
			Price:      null.NewInt(product.Price, true),
			CategoryId: null.NewString(product.CategoryId, true),
			CreatedAt:  null.NewTime(product.CreatedAt.AsTime(), true),
			UpdatedAt:  null.NewTime(product.UpdatedAt.AsTime(), true),
			DeletedAt:  null.NewTime(product.DeletedAt.AsTime(), true),
		}

		products = append(products, productData)
	}
	bodyResponseProduct := &model_response.Response[[]*entity.Product]{
		Code:    http.StatusOK,
		Message: ListProduct.Message,
		Data:    products,
	}
	return bodyResponseProduct
}
func (exposeUseCase *ExposeUseCase) CreateProduct(request *model_request.CreateProduct) (result *model_response.Response[*entity.Product]) {
	req := &pb.CreateProductRequest{
		Name:       request.Name.String,
		CategoryId: request.CategoryId.String,
		Price:      request.Price.Int64,
		Stock:      request.Stock.Int64,
	}
	createProduct, err := exposeUseCase.productClient.CreateProduct(req)
	if err != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Data:    nil,
			Message: createProduct.Message,
		}
		return
	}
	if createProduct.Data == nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Data:    nil,
			Message: createProduct.Message,
		}
		return
	}
	user := entity.Product{
		Id:         null.NewString(createProduct.Data.Id, true),
		Name:       null.NewString(createProduct.Data.Name, true),
		CategoryId: null.NewString(createProduct.Data.CategoryId, true),
		Price:      null.NewInt(createProduct.Data.Price, true),
		Stock:      null.NewInt(createProduct.Data.Stock, true),
		CreatedAt:  null.NewTime(createProduct.Data.CreatedAt.AsTime(), true),
		UpdatedAt:  null.NewTime(createProduct.Data.UpdatedAt.AsTime(), true),
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{
		Code:    http.StatusOK,
		Message: createProduct.Message,
		Data:    &user,
	}
	return bodyResponseProduct
}
func (exposeUseCase *ExposeUseCase) DeleteProduct(id string) (result *model_response.Response[*entity.Product]) {
	DeleteProduct, err := exposeUseCase.productClient.DeleteProduct(id)
	if err != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: DeleteProduct.Message,
			Data:    nil,
		}
		return
	}
	if DeleteProduct.Data == nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: DeleteProduct.Message,
			Data:    nil,
		}
		return
	}
	product := entity.Product{
		Id:         null.NewString(DeleteProduct.Data.Id, true),
		Name:       null.NewString(DeleteProduct.Data.Name, true),
		CategoryId: null.NewString(DeleteProduct.Data.CategoryId, true),
		Price:      null.NewInt(DeleteProduct.Data.Price, true),
		Stock:      null.NewInt(DeleteProduct.Data.Stock, true),
		CreatedAt:  null.NewTime(DeleteProduct.Data.CreatedAt.AsTime(), true),
		UpdatedAt:  null.NewTime(DeleteProduct.Data.UpdatedAt.AsTime(), true),
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{
		Code:    http.StatusOK,
		Message: DeleteProduct.Message,
		Data:    &product,
	}
	return bodyResponseProduct
}
func (exposeUseCase *ExposeUseCase) UpdateProduct(id string, request *model_request.ProductPatchOneByIdRequest) (result *model_response.Response[*entity.Product]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "products", id)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseProduct)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseProduct
}
func (exposeUseCase *ExposeUseCase) DetailProduct(id string) (result *model_response.Response[*entity.Product]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "products", id)

	newRequest, newRequestErr := http.NewRequest("GET", url, nil)

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseProduct := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseProduct)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Product]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseProduct
}

// category

func (exposeUseCase *ExposeUseCase) ListCategories() (result *model_response.Response[[]*entity.Category]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s", address, "categories")
	newRequest, newRequestErr := http.NewRequest("GET", url, nil)

	if newRequestErr != nil {
		result = &model_response.Response[[]*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[[]*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	Category := &model_response.Response[[]*entity.Category]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(Category)
	if decodeErr != nil {
		result = &model_response.Response[[]*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return Category
}
func (exposeUseCase *ExposeUseCase) CreateCategory(request *model_request.CategoryRequest) (result *model_response.Response[*entity.Category]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s", address, "categories")
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseCategory := &model_response.Response[*entity.Category]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseCategory)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseCategory
}
func (exposeUseCase *ExposeUseCase) DeleteCategory(id string) (result *model_response.Response[*entity.Category]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "categories", id)
	newRequest, newRequestErr := http.NewRequest("DELETE", url, nil)

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseCategory := &model_response.Response[*entity.Category]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseCategory)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseCategory
}
func (exposeUseCase *ExposeUseCase) UpdateCategory(id string, request *model_request.CategoryRequest) (result *model_response.Response[*entity.Category]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "categories", id)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseCategory := &model_response.Response[*entity.Category]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseCategory)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseCategory
}
func (exposeUseCase *ExposeUseCase) DetailCategory(id string) (result *model_response.Response[*entity.Category]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.ProductHost, exposeUseCase.Env.App.ProductPort)
	url := fmt.Sprintf("%s/%s/%s", address, "categories", id)
	newRequest, newRequestErr := http.NewRequest("GET", url, nil)

	if newRequestErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	bodyResponseCategory := &model_response.Response[*entity.Category]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(bodyResponseCategory)
	if decodeErr != nil {
		result = &model_response.Response[*entity.Category]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return bodyResponseCategory
}

// order

func (exposeUseCase *ExposeUseCase) Orders(tokenString string, request *model_request.OrderRequest) (result *model_response.Response[*model_response.OrderResponse]) {
	begin, err := exposeUseCase.DatabaseConfig.AuthDB.Connection.Begin()
	if err != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Data:    nil,
			Message: "AuthUseCase error, order is failed" + err.Error(),
		}
	}
	session, err := exposeUseCase.AuthRepository.FindOneByAccToken(begin, tokenString)
	if err != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Data:    nil,
			Message: "AuthUseCase error, order is failed" + err.Error(),
		}
	}
	userId := session.UserId
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.OrderHost, exposeUseCase.Env.App.OrderPort)
	url := fmt.Sprintf("%s/%s/%s", address, "orders", userId.String)
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}
	newRequest, newRequestErr := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))

	if newRequestErr != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	foundOrder := &model_response.Response[*model_response.OrderResponse]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(foundOrder)
	if decodeErr != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return foundOrder
}
func (exposeUseCase *ExposeUseCase) DetailOrder(id string) (result *model_response.Response[*model_response.OrderResponse]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.OrderHost, exposeUseCase.Env.App.OrderPort)
	url := fmt.Sprintf("%s/%s/%s", address, "orders", id)
	newRequest, newRequestErr := http.NewRequest("GET", url, nil)
	if newRequestErr != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	foundOrder := &model_response.Response[*model_response.OrderResponse]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(foundOrder)
	if decodeErr != nil {
		result = &model_response.Response[*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return foundOrder
}
func (exposeUseCase *ExposeUseCase) ListOrders() (result *model_response.Response[[]*model_response.OrderResponse]) {
	address := fmt.Sprintf("http://%s:%s", exposeUseCase.Env.App.OrderHost, exposeUseCase.Env.App.OrderPort)
	url := fmt.Sprintf("%s/%s", address, "orders")
	newRequest, newRequestErr := http.NewRequest("GET", url, nil)

	if newRequestErr != nil {
		result = &model_response.Response[[]*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: newRequestErr.Error(),
			Data:    nil,
		}
		return result
	}

	responseRequest, doErr := http.DefaultClient.Do(newRequest)
	if doErr != nil {
		result = &model_response.Response[[]*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: doErr.Error(),
			Data:    nil,
		}
		return result
	}
	foundOrders := &model_response.Response[[]*model_response.OrderResponse]{}
	decodeErr := json.NewDecoder(responseRequest.Body).Decode(foundOrders)
	if decodeErr != nil {
		result = &model_response.Response[[]*model_response.OrderResponse]{
			Code:    http.StatusBadRequest,
			Message: decodeErr.Error(),
			Data:    nil,
		}
	}
	return foundOrders
}
