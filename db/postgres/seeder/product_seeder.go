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

func NewProductMock(
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
				"INSERT INTO \"products\" (id, sku ,name, stock, price, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5,);",
				Product.Id,
				Product.Sku,
				Product.Name,
				Product.Stock,
				Product.Price,
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
				"DELETE FROM \"products\" WHERE id = $1 LIMIT 1;",
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
