package request

import (
	"go-micro-services/grpc/pb"

	"github.com/guregu/null"
)

type LoginRequest struct {
	Email    null.String `json:"email"`
	Password null.String `json:"password"`
}

type UserPatchOneByIdRequest struct {
	Name     null.String `json:"name"`
	Email    null.String `json:"email"`
	Password null.String `json:"password"`
	Balance  null.Int    `json:"balance"`
}
type RegisterRequest struct {
	Name     null.String `json:"name"`
	Balance  null.Int    `json:"balance"`
	Email    null.String `json:"email"`
	Password null.String `json:"password"`
}
type ProductPatchOneByIdRequest struct {
	Name       null.String `json:"name"`
	Stock      null.Int    `json:"stock"`
	Price      null.Int    `json:"price"`
	CategoryId null.String `json:"category_id"`
}
type CreateProduct struct {
	CategoryId null.String `json:"category_id"`
	Name       null.String `json:"name"`
	Stock      null.Int    `json:"stock"`
	Price      null.Int    `json:"price"`
}
type CategoryRequest struct {
	Name null.String `json:"name"`
}
type OrderRequest struct {
	TotalPaid null.Int                  `json:"total_paid"`
	Products  []*pb.OrderProductRequest `json:"products"`
}
