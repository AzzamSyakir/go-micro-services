package repository

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	orderRepository := &OrderRepository{}
	return orderRepository

}
