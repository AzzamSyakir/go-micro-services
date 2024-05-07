package seeder

import (
	"go-micro-services/src/auth-service/test/mock"
	"go-micro-services/src/order-service/config"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
)

type OrderProductSeeder struct {
	DatabaseConfig   *config.DatabaseConfig
	OrderProductMock *mock.OrderProductMock
}

func NewOrderProductSeeder(
	databaseConfig *config.DatabaseConfig,
	orderSeeder *OrderSeeder,
	productSeeder *ProductSeeder,
) *OrderProductSeeder {
	OrderProductSeeder := &OrderProductSeeder{
		DatabaseConfig:   databaseConfig,
		OrderProductMock: mock.NewOrderProductMock(orderSeeder.OrderMock, productSeeder.ProductMock),
	}
	return OrderProductSeeder
}

func (OrderProductSeeder *OrderProductSeeder) Up() {
	for _, OrderProduct := range OrderProductSeeder.OrderProductMock.Data {
		begin, beginErr := OrderProductSeeder.DatabaseConfig.OrderDB.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			_, err = begin.Query(
				"INSERT INTO order_products (id, order_id, product_id, total_price, qty, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);",
				OrderProduct.Id,
				OrderProduct.OrderId,
				OrderProduct.ProductId,
				OrderProduct.TotalPrice,
				OrderProduct.Qty,
				OrderProduct.CreatedAt,
				OrderProduct.UpdatedAt,
				OrderProduct.DeletedAt,
			)
			return err
		})
		if queryErr != nil {
			panic(queryErr)
		}
		commitErr := crdb.Execute(func() (err error) {
			err = begin.Commit()
			return err
		})
		if commitErr != nil {
			panic(commitErr)
		}
	}
}

func (OrderProductSeeder *OrderProductSeeder) Down() {
	for _, OrderProduct := range OrderProductSeeder.OrderProductMock.Data {
		begin, beginErr := OrderProductSeeder.DatabaseConfig.OrderDB.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}
		queryErr := crdb.Execute(func() (err error) {
			_, err = begin.Query(
				"DELETE FROM order_products WHERE id = $1",
				OrderProduct.Id,
			)
			return err
		})
		if queryErr != nil {
			panic(queryErr)
		}
		commitErr := crdb.Execute(func() (err error) {
			err = begin.Commit()
			return err
		})
		if commitErr != nil {
			panic(commitErr)
		}
	}
}
