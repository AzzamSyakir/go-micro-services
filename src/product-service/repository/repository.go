package repository

import (
	"database/sql"
	pb "go-micro-services/src/product-service/delivery/grpc/pb/product"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	productRepository := &ProductRepository{}
	return productRepository
}
func (productRepository *ProductRepository) CreateProduct(begin *sql.Tx, toCreateproduct *pb.Product) (result *pb.Product, err error) {
	_, queryErr := begin.Query(
		`INSERT INTO "products" (id, sku, name, stock, price, category_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		toCreateproduct.Id,
		toCreateproduct.Sku,
		toCreateproduct.Name,
		toCreateproduct.Stock,
		toCreateproduct.Price,
		toCreateproduct.CategoryId,
		toCreateproduct.CreatedAt.AsTime(),
		toCreateproduct.UpdatedAt.AsTime(),
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return
	}

	result = toCreateproduct
	err = nil
	return result, err
}

func DeserializeProductRows(rows *sql.Rows) []*pb.Product {
	var foundProducts []*pb.Product
	for rows.Next() {
		foundProduct := &pb.Product{}
		var createdAt, updatedAt *time.Time
		scanErr := rows.Scan(
			&foundProduct.Id,
			&foundProduct.Sku,
			&foundProduct.Name,
			&foundProduct.Stock,
			&foundProduct.Price,
			&foundProduct.CategoryId,
			&createdAt,
			&updatedAt,
		)
		foundProduct.CreatedAt = timestamppb.New(*createdAt)
		foundProduct.UpdatedAt = timestamppb.New(*updatedAt)
		if scanErr != nil {
			panic(scanErr)
		}
		foundProducts = append(foundProducts, foundProduct)
	}
	return foundProducts
}

func (productRepository ProductRepository) GetProductById(tx *sql.Tx, id string) (result *pb.Product, err error) {
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

func (productRepository *ProductRepository) PatchOneById(begin *sql.Tx, id string, toPatchProduct *pb.Product) (result *pb.Product, err error) {
	rows, queryErr := begin.Query(
		`UPDATE "products" SET name=$1,  stock=$2, price=$3, updated_at=$4 WHERE id = $5 ;`,
		toPatchProduct.Name,
		toPatchProduct.Stock,
		toPatchProduct.Price,
		toPatchProduct.UpdatedAt.AsTime(),
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

func (productRepository *ProductRepository) ListProducts(begin *sql.Tx) (result *pb.ProductResponseRepeated, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = begin.Query(
		`SELECT id, name, sku, stock, price, category_id, created_at, updated_at, deleted_at FROM "products" `,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err

	}
	defer rows.Close()
	var products []*pb.Product
	for rows.Next() {
		product := &pb.Product{}
		var createdAt, updatedAt time.Time
		scanErr := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Sku,
			&product.Stock,
			&product.Price,
			&product.CategoryId,
			&createdAt,
			&updatedAt,
		)
		product.CreatedAt = timestamppb.New(createdAt)
		product.UpdatedAt = timestamppb.New(updatedAt)
		if scanErr != nil {
			result = nil
			err = scanErr
			return result, err
		}
		products = append(products, product)
	}

	result = &pb.ProductResponseRepeated{
		Data: products,
	}
	err = nil
	return result, err
}

func (productRepository *ProductRepository) DeleteOneById(begin *sql.Tx, id string) (result *pb.Product, err error) {
	rows, queryErr := begin.Query(
		`DELETE FROM "products" WHERE id=$1 RETURNING id, name, sku, stock, price, category_id, created_at, updated_at`,
		id,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return
	}
	foundproducts := DeserializeProductRows(rows)
	if len(foundproducts) == 0 {
		result = nil
		err = nil
		return result, err
	}

	result = foundproducts[0]
	err = nil
	return result, err
}
