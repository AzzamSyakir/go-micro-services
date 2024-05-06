package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-micro-services/src/auth-service/entity"
	model_request "go-micro-services/src/auth-service/model/request/controller"
	model_response "go-micro-services/src/auth-service/model/response"
	"net/http"
	"testing"

	"github.com/guregu/null"
	"github.com/stretchr/testify/assert"
)

type ProductWeb struct {
	Test *testing.T
	Path string
}

func NewProductWeb(test *testing.T) *ProductWeb {
	ProductWeb := &ProductWeb{
		Test: test,
		Path: "products",
	}
	return ProductWeb
}

func (ProductWeb *ProductWeb) Start() {
	ProductWeb.Test.Run("ProductWeb_GetProduct_Succeed", ProductWeb.FindOneById)
	ProductWeb.Test.Run("ProductWeb_DeleteProduct_Succeed", ProductWeb.DeleteOneById)
	ProductWeb.Test.Run("ProductWeb_UpdateProduct_Succeed", ProductWeb.PatchOneById)
	ProductWeb.Test.Run("ProductWeb_CreateProduct_Succeed", ProductWeb.CreateProduct)
	ProductWeb.Test.Run("ProductWeb_ListProduct_Succeed", ProductWeb.ListProduct)
}

func (ProductWeb *ProductWeb) FindOneById(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	testWeb.AllSeeder.Up()
	defer testWeb.AllSeeder.Down()

	selectedProductMock := testWeb.AllSeeder.Product.ProductMock.Data[0]

	url := fmt.Sprintf("%s/%s/%s", testWeb.Server.URL, ProductWeb.Path, selectedProductMock.Id.String)
	request, newRequestErr := http.NewRequest(http.MethodGet, url, http.NoBody)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
	request.Header.Set("authorization", "Bearer "+selectedSessionMock.AccessToken.String)
	response, doErr := http.DefaultClient.Do(request)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	if doErr != nil {
		t.Fatal(doErr)
	}

	bodyResponse := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, selectedProductMock.Id, bodyResponse.Data.Id)
	assert.Equal(t, selectedProductMock.Name, bodyResponse.Data.Name)
	assert.Equal(t, selectedProductMock.Sku, bodyResponse.Data.Sku)
	assert.Equal(t, selectedProductMock.Stock, bodyResponse.Data.Stock)
	assert.Equal(t, selectedProductMock.Price, bodyResponse.Data.Price)
	assert.Equal(t, selectedProductMock.CategoryId, bodyResponse.Data.CategoryId)
}

func (ProductWeb *ProductWeb) DeleteOneById(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	testWeb.AllSeeder.Up()
	defer testWeb.AllSeeder.Down()

	selectedProductMock := testWeb.AllSeeder.Product.ProductMock.Data[0]

	url := fmt.Sprintf("%s/%s/%s", testWeb.Server.URL, ProductWeb.Path, selectedProductMock.Id.String)
	request, newRequestErr := http.NewRequest(http.MethodDelete, url, http.NoBody)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
	request.Header.Set("authorization", "Bearer "+selectedSessionMock.AccessToken.String)
	response, doErr := http.DefaultClient.Do(request)
	if doErr != nil {
		t.Fatal(doErr)
	}

	bodyResponse := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, selectedProductMock.Id, bodyResponse.Data.Id)
	assert.Equal(t, selectedProductMock.Name, bodyResponse.Data.Name)
	assert.Equal(t, selectedProductMock.Sku, bodyResponse.Data.Sku)
	assert.Equal(t, selectedProductMock.Price, bodyResponse.Data.Price)
	assert.Equal(t, selectedProductMock.CategoryId, bodyResponse.Data.CategoryId)
}

func (ProductWeb *ProductWeb) PatchOneById(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	testWeb.AllSeeder.Up()
	defer testWeb.AllSeeder.Down()

	selectedProductMock := testWeb.AllSeeder.Product.ProductMock.Data[0]

	bodyRequest := &model_request.ProductPatchOneByIdRequest{}
	bodyRequest.Name = null.NewString(selectedProductMock.Name.String+"patched", true)
	bodyRequest.CategoryId = null.NewString(selectedProductMock.CategoryId.String+"patched", true)
	bodyRequest.Stock = null.NewInt(selectedProductMock.Stock.Int64, true)
	bodyRequest.Price = null.NewInt(selectedProductMock.Price.Int64, true)

	bodyRequestJsonByte, marshalErr := json.Marshal(bodyRequest)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}
	bodyRequestBuffer := bytes.NewBuffer(bodyRequestJsonByte)

	url := fmt.Sprintf("%s/%s/%s", testWeb.Server.URL, ProductWeb.Path, selectedProductMock.Id.String)
	request, newRequestErr := http.NewRequest(http.MethodPatch, url, bodyRequestBuffer)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
	request.Header.Set("authorization", "Bearer "+selectedSessionMock.AccessToken.String)
	response, doErr := http.DefaultClient.Do(request)
	if doErr != nil {
		t.Fatal(doErr)
	}

	bodyResponse := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, selectedProductMock.Id, bodyResponse.Data.Id)
	assert.Equal(t, bodyRequest.Name, bodyResponse.Data.Name)
	assert.Equal(t, bodyRequest.Stock, bodyResponse.Data.Stock)
	assert.Equal(t, bodyRequest.Price, bodyResponse.Data.Price)
	assert.Equal(t, bodyRequest.CategoryId, bodyResponse.Data.CategoryId)
}

func (productWeb *ProductWeb) CreateProduct(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	defer testWeb.AllSeeder.Down()

	mockAuth := testWeb.AllSeeder.Product.ProductMock.Data[0]

	bodyRequest := &model_request.CreateProduct{}
	bodyRequest.Name = null.NewString(mockAuth.Name.String, true)
	bodyRequest.Stock = null.NewInt(mockAuth.Stock.Int64, true)
	bodyRequest.CategoryId = null.NewString(mockAuth.CategoryId.String, true)
	bodyRequest.Price = null.NewInt(mockAuth.Price.Int64, true)

	bodyRequestJsonByte, marshalErr := json.Marshal(bodyRequest)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}
	bodyRequestBuffer := bytes.NewBuffer(bodyRequestJsonByte)

	url := fmt.Sprintf("%s/%s/register", testWeb.Server.URL, productWeb.Path)
	request, newRequestErr := http.NewRequest(http.MethodPost, url, bodyRequestBuffer)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}

	response, doErr := http.DefaultClient.Do(request)
	if doErr != nil {
		t.Fatal(doErr)
	}

	bodyResponse := &model_response.Response[*entity.Product]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}
	fmt.Println("bodyResponse", bodyResponse)
	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, mockAuth.Name.String, bodyResponse.Data.Name.String)
	assert.Equal(t, mockAuth.Sku.String, bodyResponse.Data.Sku.String)
	assert.Equal(t, mockAuth.Price.Int64, bodyResponse.Data.Price.Int64)
	assert.Equal(t, mockAuth.Price.Int64, bodyResponse.Data.Price.Int64)

	newProductMock := bodyResponse.Data
	testWeb.AllSeeder.Product.ProductMock.Data = append(testWeb.AllSeeder.Product.ProductMock.Data, newProductMock)
}

func (productWeb *ProductWeb) ListProduct(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	testWeb.AllSeeder.Up()
	defer testWeb.AllSeeder.Down()

	url := fmt.Sprintf("%s/%s", testWeb.Server.URL, productWeb.Path)
	request, newRequestErr := http.NewRequest(http.MethodGet, url, http.NoBody)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
	request.Header.Set("authorization", "Bearer "+selectedSessionMock.AccessToken.String)
	response, doErr := http.DefaultClient.Do(request)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	if doErr != nil {
		t.Fatal(doErr)
	}

	bodyResponse := &model_response.Response[[]*entity.Product]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, bodyResponse.Code, http.StatusOK)
}
