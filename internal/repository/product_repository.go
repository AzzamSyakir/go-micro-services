package repository

import (
	"database/sql"
	"go-micro-services/internal/entity"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	productRepository := &ProductRepository{}
	return productRepository
}
func (productRepository ProductRepository) GetOneById(tx *sql.Tx, id string) (result *entity.Product, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = tx.Query(
		`SELECT id, name, price, stock, created_at, updated_at, deleted_at FROM "products" WHERE id=$1 LIMIT 1;`,
		id,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err
	}

	foundProducts := DeserializeProductRows(rows)
	if len(foundProducts) == 0 {
		result = nil
		err = nil
		return result, err
	}

	result = foundProducts[0]
	err = nil
	return result, err
}
func DeserializeProductRows(rows *sql.Rows) []*entity.Product {
	var foundProducts []*entity.Product
	for rows.Next() {
		foundProduct := &entity.Product{}
		scanErr := rows.Scan(
			&foundProduct.Id,
			&foundProduct.Name,
			&foundProduct.Price,
			&foundProduct.Stock,
			&foundProduct.CreatedAt,
			&foundProduct.UpdatedAt,
			&foundProduct.DeletedAt,
		)
		if scanErr != nil {
			panic(scanErr)
		}
		foundProducts = append(foundProducts, foundProduct)
	}
	return foundProducts
}
