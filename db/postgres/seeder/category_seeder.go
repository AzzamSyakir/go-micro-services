package seeder

import (
	"go-micro-services/src/auth-service/test/mock"
	"go-micro-services/src/product-service/config"

	"github.com/cockroachdb/cockroach-go/v2/crdb"
)

type CategorySeeder struct {
	DatabaseConfig *config.DatabaseConfig
	CategoryMock   *mock.CategoryMock
}

func NewCategorySeeder(
	databaseConfig *config.DatabaseConfig,
) *CategorySeeder {
	CategorySeeder := &CategorySeeder{
		DatabaseConfig: databaseConfig,
		CategoryMock:   mock.NewCategoryMock(),
	}
	return CategorySeeder
}

func (CategorySeeder *CategorySeeder) Up() {
	for _, category := range CategorySeeder.CategoryMock.Data {
		begin, beginErr := CategorySeeder.DatabaseConfig.ProductDB.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			_, err = begin.Query(
				"INSERT INTO categories (id, name, created_at, updated_at) VALUES ($1, $2, $3, $4;",
				category.Id,
				category.Name,
				category.CreatedAt,
				category.UpdatedAt,
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

func (CategorySeeder *CategorySeeder) Down() {
	for _, Category := range CategorySeeder.CategoryMock.Data {
		begin, beginErr := CategorySeeder.DatabaseConfig.ProductDB.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			_, err = begin.Query(
				"DELETE FROM Categories WHERE id = $1",
				Category.Id,
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
