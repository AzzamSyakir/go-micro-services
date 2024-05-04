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
	"golang.org/x/crypto/bcrypt"
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
	ProductWeb.Test.Run("ProductWeb_FindOneById_Succeed", ProductWeb.FindOneById)
	ProductWeb.Test.Run("ProductWeb_FindOneByEmail_Succeed", ProductWeb.FindOneByEmail)
	ProductWeb.Test.Run("ProductWeb_ProductPatchOneByIdRequest_Succeed", ProductWeb.PatchOneById)
	ProductWeb.Test.Run("ProductWeb_DeleteOneById_Succeed", ProductWeb.DeleteOneById)
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
	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
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
	assert.Equal(t, selectedProductMock.Email, bodyResponse.Data.Email)
	assert.Equal(t, selectedProductMock.Balance, bodyResponse.Data.Balance)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(selectedProductMock.Password.String)))
}

func (ProductWeb *ProductWeb) FindOneByEmail(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	testWeb.AllSeeder.Up()
	defer testWeb.AllSeeder.Down()

	selectedProductMock := testWeb.AllSeeder.Product.ProductMock.Data[0]

	url := fmt.Sprintf("%s/%s/email/%s", testWeb.Server.URL, ProductWeb.Path, selectedProductMock.Email.String)
	request, newRequestErr := http.NewRequest(http.MethodGet, url, http.NoBody)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
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
	assert.Equal(t, selectedProductMock.Email, bodyResponse.Data.Email)
	assert.Equal(t, selectedProductMock.Balance, bodyResponse.Data.Balance)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(selectedProductMock.Password.String)))
}

func (ProductWeb *ProductWeb) PatchOneById(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	testWeb.AllSeeder.Up()
	defer testWeb.AllSeeder.Down()

	selectedProductMock := testWeb.AllSeeder.Product.ProductMock.Data[0]

	bodyRequest := &model_request.ProductPatchOneByIdRequest{}
	bodyRequest.Name = null.NewString(selectedProductMock.Name.String+"patched", true)
	bodyRequest.Email = null.NewString(selectedProductMock.Email.String+"patched", true)
	bodyRequest.Balance = null.NewInt(selectedProductMock.Balance.Int64, true)
	bodyRequest.Password = null.NewString(selectedProductMock.Password.String+"patched", true)

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
	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
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
	assert.Equal(t, bodyRequest.Email, bodyResponse.Data.Email)
	assert.Equal(t, bodyRequest.Balance, bodyResponse.Data.Balance)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(bodyRequest.Password.String)))
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
	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
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
	assert.Equal(t, selectedProductMock.Email, bodyResponse.Data.Email)
	assert.Equal(t, selectedProductMock.Balance, bodyResponse.Data.Balance)
	assert.NoError(t, bcrypt.CompareHashAndPassword([]byte(bodyResponse.Data.Password.String), []byte(selectedProductMock.Password.String)))
}
