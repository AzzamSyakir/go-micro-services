package http

import (
	"encoding/json"
	"github.com/gorilla/mux"
	model_request "go-micro-services/internal/model/request/controller"
	"go-micro-services/internal/model/response"
	"go-micro-services/internal/use_case"
	"net/http"
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

func (orderController *OrderController) Orders(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	userId := vars["id"]
	request := &model_request.OrderRequest{}

	decodeErr := json.NewDecoder(reader.Body).Decode(request)

	if decodeErr != nil {
		http.Error(writer, "Invalid request body", http.StatusBadRequest)
	}

	result := orderController.OrderUseCase.Order(userId, request)
	response.NewResponse(writer, result)

}
