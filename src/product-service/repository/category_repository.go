package repository

import (
	"database/sql"
	pb "go-micro-services/src/product-service/delivery/grpc/pb/category"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	categoryRepository := &CategoryRepository{}
	return categoryRepository
}
func (CategoryRepository *CategoryRepository) CreateCategory(begin *sql.Tx, toCreateCategory *pb.Category) (result *pb.Category, err error) {
	_, queryErr := begin.Query(
		`INSERT INTO "categories" (id, name, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5);`,
		toCreateCategory.Id,
		toCreateCategory.Name,
		toCreateCategory.CreatedAt.AsTime(),
		toCreateCategory.UpdatedAt.AsTime(),
		toCreateCategory.DeletedAt.AsTime(),
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return
	}

	result = toCreateCategory
	err = nil
	return result, err
}

func DeserializeCategoryRows(rows *sql.Rows) []*pb.Category {
	var foundCategories []*pb.Category
	for rows.Next() {
		foundCategory := &pb.Category{}
		var createdAt, updatedAt, deletedAt time.Time
		scanErr := rows.Scan(
			&foundCategory.Id,
			&foundCategory.Name,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		if scanErr != nil {
			panic(scanErr)
		}
		foundCategory.CreatedAt = timestamppb.New(createdAt)
		foundCategory.UpdatedAt = timestamppb.New(updatedAt)
		foundCategory.DeletedAt = timestamppb.New(deletedAt)
		foundCategories = append(foundCategories, foundCategory)
	}
	return foundCategories
}

func (categoryRepository CategoryRepository) GetProductById(tx *sql.Tx, id string) (result *pb.Category, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = tx.Query(
		`SELECT id, name, created_at, updated_at, deleted_at FROM "categories" WHERE id=$1 LIMIT 1;`,
		id,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err
	}
	defer rows.Close()

	foundCategory := DeserializeCategoryRows(rows)
	if len(foundCategory) == 0 {
		result = nil
		err = nil
		return result, err
	}

	result = foundCategory[0]
	err = nil
	return result, err
}

func (categoryRepository *CategoryRepository) PatchOneById(begin *sql.Tx, id string, toPatchCategory *pb.Category) (result *pb.Category, err error) {
	rows, queryErr := begin.Query(
		`UPDATE "categories" SET name=$1, updated_at=$2 WHERE id = $3 ;`,
		toPatchCategory.Name,
		toPatchCategory.UpdatedAt.AsTime(),
		id,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return
	}
	defer rows.Close()

	result = toPatchCategory
	err = nil
	return result, err
}

func (categoryRepository *CategoryRepository) ListCategories(begin *sql.Tx) (result *pb.CategoryResponseRepeated, err error) {
	var rows *sql.Rows
	var queryErr error
	rows, queryErr = begin.Query(
		`SELECT id, name, created_at, updated_at, deleted_at FROM "categories" `,
	)

	if queryErr != nil {
		result = nil
		err = queryErr
		return result, err

	}
	defer rows.Close()
	var categories []*pb.Category
	for rows.Next() {
		category := &pb.Category{}
		var createdAt, updatedAt, deletedAt time.Time
		scanErr := rows.Scan(
			&category.Id,
			&category.Name,
			&createdAt,
			&updatedAt,
			&deletedAt,
		)
		if scanErr != nil {
			panic(scanErr)
		}
		category.CreatedAt = timestamppb.New(createdAt)
		category.UpdatedAt = timestamppb.New(updatedAt)
		category.DeletedAt = timestamppb.New(deletedAt)

		categories = append(categories, category)
	}

	result = &pb.CategoryResponseRepeated{
		Data: categories,
	}
	err = nil
	return result, err
}

func (CategoryRepository *CategoryRepository) DeleteOneById(begin *sql.Tx, id string) (result *pb.Category, err error) {
	rows, queryErr := begin.Query(
		`DELETE FROM "categories" WHERE id=$1 RETURNING id, name, created_at, updated_at, deleted_at`,
		id,
	)
	if queryErr != nil {
		result = nil
		err = queryErr
		return
	}
	foundCategory := DeserializeCategoryRows(rows)
	if len(foundCategory) == 0 {
		result = nil
		err = nil
		return result, err
	}

	result = foundCategory[0]
	err = nil
	return result, err
}
