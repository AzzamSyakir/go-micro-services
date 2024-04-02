package repository

import (
	"database/sql"
	"go-micro-services/internal/entity"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	orderRepository := &OrderRepository{}
	return orderRepository
}

func (orderRepository *OrderRepository) Order(begin *sql.Tx, orders *entity.Order) (result *entity.Order, err error) {
	_, queryErr := begin.Query(
		`INSERT INTO "orders"(id, user_id, name, total_price, total_paid, total_return, receipt_code, created_at, updated_at, deleted_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		orders.Id,
		orders.UserId,
		orders.Name,
		orders.TotalPrice,
		orders.TotalPaid,
		orders.TotalReturn,
		orders.ReceiptCode,
		orders.CreatedAt,
		orders.UpdatedAt,
		orders.DeletedAt,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return
	}

	result = orders
	err = nil
	return result, err
}

func (orderRepository *OrderRepository) OrderProducts(begin *sql.Tx, orderProducts entity.OrderProducts) (result *entity.OrderProducts, err error) {
	_, queryErr := begin.Query(
		`INSERT INTO "order_products"(id, order_id, product_id, total_price, qty, created_at, updated_at, deleted_at) values ($1, $2, $3, $4, $5, $6, $7, $8)`,
		orderProducts.Id,
		orderProducts.OrderId,
		orderProducts.ProductId,
		orderProducts.TotalPrice,
		orderProducts.Qty,
		orderProducts.CreatedAt,
		orderProducts.UpdatedAt,
		orderProducts.DeletedAt,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return
	}

	result = orderProducts
	err = nil
	return result, err
}
