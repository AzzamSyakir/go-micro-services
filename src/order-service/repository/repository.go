package repository

import (
	"database/sql"
	"fmt"
	"go-micro-services/src/order-service/entity"
	model_response "go-micro-services/src/order-service/model/response"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	orderRepository := &OrderRepository{}
	return orderRepository
}
func DeserializeOrderRows(rows *sql.Rows) []*entity.Order {
	var foundOrders []*entity.Order
	for rows.Next() {
		foundOrder := &entity.Order{}
		scanErr := rows.Scan(
			&foundOrder.Id,
			&foundOrder.UserId,
			&foundOrder.TotalPrice,
			&foundOrder.TotalPaid,
			&foundOrder.TotalReturn,
			&foundOrder.ReceiptCode,
			&foundOrder.CreatedAt,
			&foundOrder.UpdatedAt,
			&foundOrder.DeletedAt,
		)
		if scanErr != nil {
			panic(scanErr)
		}
		foundOrders = append(foundOrders, foundOrder)
	}
	return foundOrders
}
func (orderRepository *OrderRepository) ListOrders(begin *sql.Tx) (result *model_response.Response[[]*entity.Order], err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = begin.Query(
		`SELECT id, user_id, total_price, total_paid, total_return, receipt_code, created_at, updated_at, deleted_at FROM "orders" `,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err

	}
	defer rows.Close()
	var orders []*entity.Order
	for rows.Next() {
		order := &entity.Order{}
		scanErr := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.TotalPrice,
			&order.TotalPaid,
			&order.TotalReturn,
			&order.ReceiptCode,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.DeletedAt,
		)
		if scanErr != nil {
			result = nil
			err = scanErr
			return result, err
		}
		orders = append(orders, order)
	}

	result = &model_response.Response[[]*entity.Order]{
		Data: orders,
	}
	err = nil
	return result, err
}
func (orderRepository *OrderRepository) Order(begin *sql.Tx, orders *entity.Order) (result *model_response.Response[*model_response.OrderResponse], err error) {
	rows, queryErr := begin.Query(
		`INSERT INTO "orders"(id, user_id, total_price, total_paid, total_return, receipt_code, created_at, updated_at, deleted_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		orders.Id,
		orders.UserId,
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
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	ordersResponse := &model_response.OrderResponse{
		Id:          orders.Id,
		UserId:      orders.UserId,
		TotalPrice:  orders.TotalPrice,
		TotalPaid:   orders.TotalPaid,
		TotalReturn: orders.TotalReturn,
		ReceiptCode: orders.ReceiptCode,
		CreatedAt:   orders.CreatedAt,
		UpdatedAt:   orders.UpdatedAt,
	}
	result = &model_response.Response[*model_response.OrderResponse]{
		Data: ordersResponse,
	}
	err = nil
	return result, err
}
func (orderRepository *OrderRepository) OrderProducts(begin *sql.Tx, orderProducts *entity.OrderProducts) (result *entity.OrderProducts, err error) {
	rows, queryErr := begin.Query(
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
		fmt.Println("queryErr : ", queryErr)
		result = nil
		err = queryErr
		return
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	result = orderProducts
	err = nil
	return result, err
}
func (orderRepository OrderRepository) GetOneById(tx *sql.Tx, id string) (result *entity.Order, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = tx.Query(
		`SELECT id, user_id, total_price,  total_paid, total_return, receipt_code,  created_at, updated_at, deleted_at FROM "orders" WHERE id=$1 LIMIT 1;`,
		id,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err
	}
	defer rows.Close()

	foundOrder := DeserializeOrderRows(rows)
	if len(foundOrder) == 0 {
		result = nil
		err = nil
		return result, err
	}

	result = foundOrder[0]
	err = nil
	return result, err
}
