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
	orderWeb.Test.Run("OrderWeb_GetOrder_Succeed", orderWeb.FindOneById)
	orderWeb.Test.Run("OrderWeb_Order_Succeed", orderWeb.Order)
	// orderWeb.Test.Run("OrderWeb_ListOrder_Succeed", orderWeb.ListOrder)
}

func (orderWeb *OrderWeb) FindOneById(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	testWeb.AllSeeder.Up()
	defer testWeb.AllSeeder.Down()

	selectedOrderMock := testWeb.AllSeeder.Order.OrderMock.Data[0]

	url := fmt.Sprintf("%s/%s/%s", testWeb.Server.URL, orderWeb.Path, selectedOrderMock.Id.String)
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
	bodyResponse := &model_response.Response[*entity.Order]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
}

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

	bodyResponse := &model_response.Response[*model_response.OrderResponse]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.NotEqual(t, nil, bodyResponse.Data)
	newOrderMock := &entity.Order{
		Id:          bodyResponse.Data.Id,
		UserId:      bodyResponse.Data.UserId,
		ReceiptCode: bodyResponse.Data.ReceiptCode,
		TotalPrice:  bodyResponse.Data.TotalPrice,
		TotalPaid:   bodyRequest.TotalPaid,
		TotalReturn: bodyResponse.Data.TotalReturn,
		CreatedAt:   bodyResponse.Data.CreatedAt,
		UpdatedAt:   bodyResponse.Data.UpdatedAt,
	}
	newOrderProductMock := bodyResponse.Data.Products
	testWeb.AllSeeder.Order.OrderMock.Data = append(testWeb.AllSeeder.Order.OrderMock.Data, newOrderMock)
	testWeb.AllSeeder.OrderProduct.OrderProductMock.Data = newOrderProductMock
}

func (orderWeb *OrderWeb) ListOrder(t *testing.T) {
	t.Parallel()

	testWeb := GetTestWeb()
	testWeb.AllSeeder.Up()
	defer testWeb.AllSeeder.Down()

	url := fmt.Sprintf("%s/%s", testWeb.Server.URL, orderWeb.Path)
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

	bodyResponse := &model_response.Response[[]*entity.Order]{}
	decodeErr := json.NewDecoder(response.Body).Decode(bodyResponse)
	if decodeErr != nil {
		t.Fatal(decodeErr)
	}
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "application/json", response.Header.Get("Content-Type"))
	assert.Equal(t, bodyResponse.Code, http.StatusOK)
}
