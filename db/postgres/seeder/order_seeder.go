package seeder

import (
	"go-micro-services/src/auth-service/test/mock"
	"go-micro-services/src/order-service/config"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
)

type OrderSeeder struct {
	DatabaseConfig *config.DatabaseConfig
	OrderMock      *mock.OrderMock
}

func NewOrderSeeder(
	databaseConfig *config.DatabaseConfig,
	userSeeder *UserSeeder,
) *OrderSeeder {
	OrderSeeder := &OrderSeeder{
		DatabaseConfig: databaseConfig,
		OrderMock:      mock.NewOrderMock(userSeeder.UserMock),
	}
	return OrderSeeder
}

func (OrderSeeder *OrderSeeder) Up() {
	for _, Order := range OrderSeeder.OrderMock.Data {
		begin, beginErr := OrderSeeder.DatabaseConfig.OrderDB.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			_, err = begin.Query(
				"INSERT INTO orders (id, user_id, total_price, total_paid, total_return, receipt_code, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);",
				Order.Id,
				Order.UserId,
				Order.TotalPrice,
				Order.TotalPaid,
				Order.TotalReturn,
				Order.ReceiptCode,
				Order.CreatedAt,
				Order.UpdatedAt,
				Order.DeletedAt,
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

func (OrderSeeder *OrderSeeder) Down() {
	for _, Order := range OrderSeeder.OrderMock.Data {
		begin, beginErr := OrderSeeder.DatabaseConfig.OrderDB.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			_, err = begin.Query(
				"DELETE FROM orders WHERE id = $1",
				Order.Id,
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
