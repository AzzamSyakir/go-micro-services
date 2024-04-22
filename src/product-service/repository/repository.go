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
func DeserializeProductRows(rows *sql.Rows) []*entity.product-service {
	var foundProducts []*entity.product-service
	for rows.Next() {
		foundProduct := &entity.product-service{}
		scanErr := rows.Scan(
			&foundProduct.Id,
			&foundProduct.Sku,
			&foundProduct.Name,
			&foundProduct.Stock,
			&foundProduct.Price,
			&foundProduct.CategoryId,
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
func (productRepository ProductRepository) GetOneById(tx *sql.Tx, id string) (result *entity.product-service, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = tx.Query(
		`SELECT id, sku, name,  stock, price, category_id,  created_at, updated_at, deleted_at FROM "products" WHERE id=$1 LIMIT 1;`,
		id,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err
	}
	defer rows.Close()

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
func (productRepository *ProductRepository) PatchOneById(begin *sql.Tx, id string, toPatchProduct *entity.product-service) (result *entity.product-service, err error) {
	rows, queryErr := begin.Query(
		`UPDATE "products" SET name=$1,  stock=$2, price=$3, updated_at=$4 WHERE id = $5 ;`,
		toPatchProduct.Name,
		toPatchProduct.Stock,
		toPatchProduct.Price,
		toPatchProduct.UpdatedAt,
		id,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return
	}
	defer rows.Close()

	result = toPatchProduct
	err = nil
	return result, err
}
