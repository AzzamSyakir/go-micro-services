package seeder

import (
	"go-micro-services/src/auth-service/test/mock"
	"go-micro-services/src/product-service/config"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
)

type ProductSeeder struct {
	DatabaseConfig *config.DatabaseConfig
	ProductMock    *mock.ProductMock
}

func NewProductSeeder(
	databaseConfig *config.DatabaseConfig,
	categorySeeder *CategorySeeder,
) *ProductSeeder {
	ProductSeeder := &ProductSeeder{
		DatabaseConfig: databaseConfig,
		ProductMock:    mock.NewProductMock(categorySeeder.CategoryMock),
	}
	return ProductSeeder
}

func (ProductSeeder *ProductSeeder) Up() {
	for _, Product := range ProductSeeder.ProductMock.Data {
		begin, beginErr := ProductSeeder.DatabaseConfig.ProductDB.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			_, err = begin.Query(
				"INSERT INTO Products (id, sku, name, stock, price, category_id created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);",
				Product.Id,
				Product.Sku,
				Product.Name,
				Product.Stock,
				Product.Price,
				Product.CategoryId,
				Product.CreatedAt,
				Product.UpdatedAt,
				Product.DeletedAt,
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

func (ProductSeeder *ProductSeeder) Down() {
	for _, Product := range ProductSeeder.ProductMock.Data {
		begin, beginErr := ProductSeeder.DatabaseConfig.ProductDB.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			_, err = begin.Query(
				"DELETE FROM Products WHERE id = $1",
				Product.Id,
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
