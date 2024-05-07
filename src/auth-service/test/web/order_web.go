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

type OrderWeb struct {
	Test *testing.T
	Path string
}

func NewOrderWeb(test *testing.T) *OrderWeb {
	orderWeb := &OrderWeb{
		Test: test,
		Path: "orders",
	}
	return orderWeb
}

func (orderWeb *OrderWeb) Start() {
	// orderWeb.Test.Run("CategoryWeb_GetCategory_Succeed", orderWeb.FindOneById)
	// orderWeb.Test.Run("CategoryWeb_DeleteCategory_Succeed", orderWeb.DeleteOneById)
	// orderWeb.Test.Run("CategoryWeb_UpdateCategory_Succeed", orderWeb.PatchOneById)
	orderWeb.Test.Run("OrderWeb_Order_Succeed", orderWeb.Order)
	// orderWeb.Test.Run("CategoryWeb_ListCategory_Succeed", orderWeb.ListOrder)
}

// func (orderWeb *OrderWeb) FindOneById(t *testing.T) {
// 	t.Parallel()

// 	testWeb := GetTestWeb()
// 	testWeb.AllSeeder.Up()
// 	defer testWeb.AllSeeder.Down()

// 	selectedCategoryMock := testWeb.AllSeeder.Category.CategoryMock.Data[0]

// 	url := fmt.Sprintf("%s/%s/%s", testWeb.Server.URL, orderWeb.Path, selectedCategoryMock.Id.String)
// 	request, newRequestErr := http.NewRequest(http.MethodGet, url, http.NoBody)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}
// 	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
// 	request.Header.Set("authorization", "Bearer "+selectedSessionMock.AccessToken.String)
// 	response, doErr := http.DefaultClient.Do(request)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}
// 	if doErr != nil {
// 		t.Fatal(doErr)
// 	}
// 	bodyResponse := &model_response.Response[*entity.Category]{}
// 	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
// 	if decodeErr != nil {
// 		t.Fatal(decodeErr)
// 	}
// 	assert.Equal(t, http.StatusOK, response.StatusCode)
// 	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
// }

// func (orderWeb *OrderWeb) DeleteOneById(t *testing.T) {
// 	t.Parallel()

// 	testWeb := GetTestWeb()
// 	testWeb.AllSeeder.Session.Up()
// 	testWeb.AllSeeder.Category.Up()
// 	defer testWeb.AllSeeder.Down()

// 	selectedCategoryMock := testWeb.AllSeeder.Category.CategoryMock.Data[0]

// 	url := fmt.Sprintf("%s/%s/%s", testWeb.Server.URL, orderWeb.Path, selectedCategoryMock.Id.String)
// 	request, newRequestErr := http.NewRequest(http.MethodDelete, url, http.NoBody)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}
// 	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
// 	request.Header.Set("Authorization", "Bearer "+selectedSessionMock.AccessToken.String)
// 	response, doErr := http.DefaultClient.Do(request)
// 	if doErr != nil {
// 		t.Fatal(doErr)
// 	}

// 	bodyResponse := &model_response.Response[*entity.Category]{}
// 	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
// 	if decodeErr != nil {
// 		t.Fatal(decodeErr)
// 	}
// 	assert.Equal(t, http.StatusOK, response.StatusCode)
// 	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
// }

// func (orderWeb *OrderWeb) PatchOneById(t *testing.T) {
// 	t.Parallel()

// 	testWeb := GetTestWeb()
// 	testWeb.AllSeeder.Up()
// 	defer testWeb.AllSeeder.Down()

// 	selectedCategoryMock := testWeb.AllSeeder.Category.CategoryMock.Data[0]

// 	bodyRequest := &model_request.CategoryRequest{}
// 	bodyRequest.Name = null.NewString(selectedCategoryMock.Name.String+"patched", true)

// 	bodyRequestJsonByte, marshalErr := json.Marshal(bodyRequest)
// 	if marshalErr != nil {
// 		t.Fatal(marshalErr)
// 	}
// 	bodyRequestBuffer := bytes.NewBuffer(bodyRequestJsonByte)

// 	url := fmt.Sprintf("%s/%s/%s", testWeb.Server.URL, orderWeb.Path, selectedCategoryMock.Id.String)
// 	request, newRequestErr := http.NewRequest(http.MethodPatch, url, bodyRequestBuffer)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}
// 	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
// 	request.Header.Set("authorization", "Bearer "+selectedSessionMock.AccessToken.String)
// 	response, doErr := http.DefaultClient.Do(request)
// 	if doErr != nil {
// 		t.Fatal(doErr)
// 	}

// 	bodyResponse := &model_response.Response[*entity.Category]{}
// 	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
// 	if decodeErr != nil {
// 		t.Fatal(decodeErr)
// 	}

// 	assert.Equal(t, http.StatusOK, response.StatusCode)
// 	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
// 	assert.Equal(t, selectedCategoryMock.Id, bodyResponse.Data.Id)
// 	assert.Equal(t, bodyRequest.Name, bodyResponse.Data.Name)
// }

func (orderWeb *OrderWeb) Order(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	defer testWeb.AllSeeder.Down()
	testWeb.AllSeeder.User.Up()
	testWeb.AllSeeder.Session.Up()
	testWeb.AllSeeder.Category.Up()
	testWeb.AllSeeder.Product.Up()

	orderMock := testWeb.AllSeeder.Order.OrderMock.Data[0]
	bodyRequest := &model_request.OrderRequest{}
	var productRequest []model_request.OrderProducts
	for _, productMock := range testWeb.AllSeeder.Product.ProductMock.Data {
		productID := productMock.Id.String
		qty := 1
		productRequest = append(productRequest, model_request.OrderProducts{
			ProductId: null.NewString(productID, true),
			Qty:       null.NewInt(int64(qty), true),
		})
	}
	bodyRequest.Products = productRequest

	bodyRequest.TotalPaid = null.NewInt(orderMock.TotalPaid.Int64, true)

	bodyRequestJsonByte, marshalErr := json.Marshal(bodyRequest)
	if marshalErr != nil {
		t.Fatal(marshalErr)
	}
	bodyRequestBuffer := bytes.NewBuffer(bodyRequestJsonByte)

	url := fmt.Sprintf("%s/%s", testWeb.Server.URL, orderWeb.Path)
	request, newRequestErr := http.NewRequest(http.MethodPost, url, bodyRequestBuffer)
	if newRequestErr != nil {
		t.Fatal(newRequestErr)
	}
	selectedSessionOrder := testWeb.AllSeeder.Session.SessionMock.Data[0]
	request.Header.Set("authorization", "Bearer "+selectedSessionOrder.AccessToken.String)
	response, doErr := http.DefaultClient.Do(request)
	if doErr != nil {
		t.Fatal(doErr)
	}

	bodyResponse := &model_response.Response[*entity.Category]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))

	newCategoryMock := bodyResponse.Data
	testWeb.AllSeeder.Category.CategoryMock.Data = append(testWeb.AllSeeder.Category.CategoryMock.Data, newCategoryMock)
}

// func (orderWeb *OrderWeb) ListOrder(t *testing.T) {
// 	t.Parallel()

// 	testWeb := GetTestWeb()
// 	testWeb.AllSeeder.Up()
// 	defer testWeb.AllSeeder.Down()

// 	url := fmt.Sprintf("%s/%s", testWeb.Server.URL, orderWeb.Path)
// 	request, newRequestErr := http.NewRequest(http.MethodGet, url, http.NoBody)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}
// 	selectedSessionMock := testWeb.AllSeeder.Session.SessionMock.Data[0]
// 	request.Header.Set("authorization", "Bearer "+selectedSessionMock.AccessToken.String)
// 	response, doErr := http.DefaultClient.Do(request)
// 	if newRequestErr != nil {
// 		t.Fatal(newRequestErr)
// 	}
// 	if doErr != nil {
// 		t.Fatal(doErr)
// 	}

// 	bodyResponse := &model_response.Response[[]*entity.Category]{}
// 	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
// 	if decodeErr != nil {
// 		t.Fatal(decodeErr)
// 	}
// 	assert.Equal(t, http.StatusOK, response.StatusCode)
// 	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
// 	assert.Equal(t, bodyResponse.Code, http.StatusOK)
// }
