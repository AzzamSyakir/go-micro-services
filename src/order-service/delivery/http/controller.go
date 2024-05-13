package http

import (
	"encoding/json"
	model_request "go-micro-services/src/order-service/model/request/controller"
	"go-micro-services/src/order-service/model/response"
	"go-micro-services/src/order-service/use_case"
	"net/http"

	"github.com/gorilla/mux"
)

type OrderController struct {
	OrderUseCase *use_case.OrderUseCase
}

func NewOrderController(orderUseCase *use_case.OrderUseCase) *OrderController {
	orderControler := &OrderController{
		OrderUseCase: orderUseCase,
	}
	return orderControler
}

func (orderController *OrderController) ListOrders(writer http.ResponseWriter, reader *http.Request) {
	result, _ := orderController.OrderUseCase.ListOrders()
	response.NewResponse(writer, result)
}
func (orderController *OrderController) DetailOrders(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	userId := vars["id"]
	result := orderController.OrderUseCase.DetailOrders(userId)
	response.NewResponse(writer, result)
}
func (orderController *OrderController) Orders(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	userId := vars["id"]
	request := &model_request.OrderRequest{}

	decodeErr := json.NewDecoder(reader.Body).Decode(request)
	if decodeErr != nil {
		http.Error(writer, "Failed to decode request body: "+decodeErr.Error(), http.StatusBadRequest)
		return
	}
	result, _ := orderController.OrderUseCase.Order(userId, request)
	response.NewResponse(writer, result)
}
