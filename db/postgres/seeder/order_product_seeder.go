package seeder

import (
	"go-micro-services/src/auth-service/test/mock"
	"go-micro-services/src/order-service/config"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
)

type OrderProductSeeder struct {
	DatabaseConfig   *config.DatabaseConfig
	OrderProductMock *mock.OrderProductsMock
}

func NewOrderProductMock(
	databaseConfig *config.DatabaseConfig,
	orderSeeder *OrderSeeder,
	productSeeder *ProductSeeder,
) *OrderProductSeeder {
	OrderProductSeeder := &OrderProductSeeder{
		DatabaseConfig:   databaseConfig,
		OrderProductMock: mock.NewOrderProductsMock(orderSeeder.OrderMock, productSeeder.ProductMock),
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
				"INSERT INTO \"order_products\" (id, order_id, product_id, total_price, qty ,created_at, updated_at deleted_at) VALUES ($1, $2, $3, $4, $5,);",
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

func (orderProductSeeder *OrderProductSeeder) Down() {
	for _, category := range orderProductSeeder.OrderProductMock.Data {
		begin, beginErr := orderProductSeeder.DatabaseConfig.OrderDB.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			_, err = begin.Query(
				"DELETE FROM \"categories\" WHERE id = $1 LIMIT 1;",
				category.Id,
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
