package http

import "go-micro-services/internal/use_case"

type OrderController struct {
	OrderUseCase *use_case.OrderUseCase
}

func NewOrderController(orderUseCase *use_case.OrderUseCase) *OrderController {
	orderControler := &OrderController{
		OrderUseCase: orderUseCase,
	}
	return orderControler
}
