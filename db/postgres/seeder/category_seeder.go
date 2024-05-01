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

func NewCategorSeeder(
	databaseConfig *config.DatabaseConfig,
) *CategorySeeder {
	categorySeeder := &CategorySeeder{
		DatabaseConfig: databaseConfig,
		CategoryMock:   mock.NewCategoryMock(),
	}
	return categorySeeder
}

func (categorySeeder *CategorySeeder) Up() {
	for _, category := range categorySeeder.CategoryMock.Data {
		begin, beginErr := categorySeeder.DatabaseConfig.ProductDB.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			_, err = begin.Query(
				"INSERT INTO \"categories\" (id, name, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5);",
				category.Id,
				category.Name,
				category.CreatedAt,
				category.UpdatedAt,
				category.DeletedAt,
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

func (categorySeeder *CategorySeeder) Down() {
	for _, category := range categorySeeder.CategoryMock.Data {
		begin, beginErr := categorySeeder.DatabaseConfig.ProductDB.Connection.Begin()
		if beginErr != nil {
			panic(beginErr)
		}

		queryErr := crdb.Execute(func() (err error) {
			_, err = begin.Query(
				"DELETE FROM \"categories\" WHERE id = $1;",
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
