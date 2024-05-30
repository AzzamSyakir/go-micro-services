package repository

import (
	"database/sql"
	"go-micro-services/grpc/pb"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	orderRepository := &OrderRepository{}
	return orderRepository
}
func DeserializeOrderRows(rows *sql.Rows) []*pb.Order {
	var foundOrders []*pb.Order
	for rows.Next() {
		foundOrder := &pb.Order{}
		var createdAt, updatedAt time.Time
		scanErr := rows.Scan(
			&foundOrder.Id,
			&foundOrder.UserId,
			&foundOrder.TotalPrice,
			&foundOrder.TotalPaid,
			&foundOrder.TotalReturn,
			&foundOrder.ReceiptCode,
			&createdAt,
			&updatedAt,
		)
		foundOrder.CreatedAt = timestamppb.New(createdAt)
		foundOrder.UpdatedAt = timestamppb.New(updatedAt)
		if scanErr != nil {
			panic(scanErr)
		}
		foundOrders = append(foundOrders, foundOrder)
	}
	return foundOrders
}
func DeserializeOrderProductRows(rows *sql.Rows) []*pb.OrderProduct {
	var foundOrderProducts []*pb.OrderProduct
	for rows.Next() {
		foundOrderProduct := &pb.OrderProduct{}
		var createdAt, updatedAt, deletedAt time.Time
		scanErr := rows.Scan(
			&foundOrderProduct.Id,
			&foundOrderProduct.OrderId,
			&foundOrderProduct.ProductId,
			&foundOrderProduct.TotalPrice,
			&foundOrderProduct.Qty,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		foundOrderProduct.CreatedAt = timestamppb.New(createdAt)
		foundOrderProduct.UpdatedAt = timestamppb.New(updatedAt)
		foundOrderProduct.DeletedAt = timestamppb.New(deletedAt)
		if scanErr != nil {
			panic(scanErr)
		}
		foundOrderProducts = append(foundOrderProducts, foundOrderProduct)
	}
	return foundOrderProducts
}

func (orderRepository *OrderRepository) ListOrders(begin *sql.Tx) (result *pb.OrderResponseRepeated, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = begin.Query(
		`SELECT id, user_id, total_price, total_paid, total_return, receipt_code, created_at, updated_at FROM "orders" `,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err

	}
	defer rows.Close()
	var orders []*pb.Order
	for rows.Next() {
		order := &pb.Order{}
		var createdAt, updatedAt, deletedAt time.Time
		scanErr := rows.Scan(
			&order.Id,
			&order.UserId,
			&order.TotalPrice,
			&order.TotalPaid,
			&order.TotalReturn,
			&order.ReceiptCode,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		order.CreatedAt = timestamppb.New(createdAt)
		order.UpdatedAt = timestamppb.New(updatedAt)
		if scanErr != nil {
			result = nil
			err = scanErr
			return result, err
		}
		orders = append(orders, order)
	}

	result = &pb.OrderResponseRepeated{
		Data: orders,
	}
	err = nil
	return result, err
}

func (orderRepository *OrderRepository) Order(begin *sql.Tx, orders *pb.Order) (result *pb.OrderResponse, err error) {
	rows, queryErr := begin.Query(
		`INSERT INTO "orders"(id, user_id, total_price, total_paid, total_return, receipt_code, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8)`,
		orders.Id,
		orders.UserId,
		orders.TotalPrice,
		orders.TotalPaid,
		orders.TotalReturn,
		orders.ReceiptCode,
		orders.CreatedAt.AsTime(),
		orders.UpdatedAt.AsTime(),
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

	ordersResponse := &pb.Order{
		Id:          orders.Id,
		UserId:      orders.UserId,
		TotalPrice:  orders.TotalPrice,
		TotalPaid:   orders.TotalPaid,
		TotalReturn: orders.TotalReturn,
		ReceiptCode: orders.ReceiptCode,
		CreatedAt:   orders.CreatedAt,
		UpdatedAt:   orders.UpdatedAt,
	}
	result = &pb.OrderResponse{
		Data: ordersResponse,
	}
	err = nil
	return result, err
}

func (orderRepository OrderRepository) DetailOrder(tx *sql.Tx, id string) (result *pb.Order, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = tx.Query(
		`SELECT id, user_id, total_price, total_paid, total_return, receipt_code,  created_at, updated_at FROM "orders" WHERE id=$1 LIMIT 1;`,
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

func (orderRepository *OrderRepository) OrderProducts(begin *sql.Tx, orderProducts *pb.OrderProduct) (result *pb.OrderProduct, err error) {
	rows, queryErr := begin.Query(
		`INSERT INTO "order_products"(id, order_id, product_id, total_price, qty, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7)`,
		orderProducts.Id,
		orderProducts.OrderId,
		orderProducts.ProductId,
		orderProducts.TotalPrice,
		orderProducts.Qty,
		orderProducts.CreatedAt.AsTime(),
		orderProducts.UpdatedAt.AsTime(),
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

	result = orderProducts
	err = nil
	return result, err
}

func (orderRepository *OrderRepository) GetOrderProductsByOrderId(tx *sql.Tx, order_id string) (result *pb.OrderProductResponse, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = tx.Query(
		`SELECT id, order_id, product_id,  total_price, qty,  created_at, updated_at FROM "order_products" WHERE order_id=$1;`,
		order_id,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err

	}
	defer rows.Close()
	var fetchOrderProducts []*pb.OrderProduct
	for rows.Next() {
		fetchOrderProduct := &pb.OrderProduct{}
		var createdAt, updatedAt time.Time
		scanErr := rows.Scan(
			&fetchOrderProduct.Id,
			&fetchOrderProduct.OrderId,
			&fetchOrderProduct.ProductId,
			&fetchOrderProduct.TotalPrice,
			&fetchOrderProduct.Qty,
			&createdAt,
			&updatedAt,
		)
		fetchOrderProduct.CreatedAt = timestamppb.New(createdAt)
		fetchOrderProduct.UpdatedAt = timestamppb.New(updatedAt)
		if scanErr != nil {
			result = nil
			err = scanErr
			return result, err
		}
		fetchOrderProducts = append(fetchOrderProducts, fetchOrderProduct)
	}

	result = &pb.OrderProductResponse{
		Data: fetchOrderProducts,
	}
	err = nil
	return result, err
}
