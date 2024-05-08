package entity

import (
	"github.com/guregu/null"
)

type Session struct {
	Id                    null.String `json:"id"`
	UserId                null.String `json:"user_id"`
	AccessToken           null.String `json:"access_token"`
	RefreshToken          null.String `json:"refresh_token"`
	AccessTokenExpiredAt  null.Time   `json:"access_token_expired_at"`
	RefreshTokenExpiredAt null.Time   `json:"refresh_token_expired_at"`
	UpdatedAt             null.Time   `json:"updated_at"`
	CreatedAt             null.Time   `json:"created_at"`
	DeletedAt             null.Time   `json:"deleted_at"`
}

type User struct {
	Id        null.String `json:"id"`
	Name      null.String `json:"name" `
	Email     null.String `json:"email"`
	Password  null.String `json:"password"`
	Balance   null.Int    `json:"balance"`
	CreatedAt null.Time   `json:"created_at"`
	UpdatedAt null.Time   `json:"updated_at"`
	DeletedAt null.Time   `json:"deleted_at"`
}
type Product struct {
	Id         null.String `json:"id"`
	Sku        null.String `json:"sku"`
	Name       null.String `json:"name"`
	Stock      null.Int    `json:"stock"`
	Price      null.Int    `json:"price"`
	CategoryId null.String `json:"category_id"`
	CreatedAt  null.Time   `json:"created_at"`
	UpdatedAt  null.Time   `json:"updated_at"`
	DeletedAt  null.Time   `json:"deleted_at"`
}

type Category struct {
	Id        null.String `json:"id"`
	Name      null.String `json:"name"`
	CreatedAt null.Time   `json:"created_at"`
	UpdatedAt null.Time   `json:"updated_at"`
	DeletedAt null.Time   `json:"deleted_at"`
}

type OrderProducts struct {
	Id         null.String `json:"id"`
	OrderId    null.String `json:"user_id"`
	ProductId  null.String `json:"product_id"`
	TotalPrice null.Int    `json:"total_price"`
	Qty        null.Int    `json:"qty"`
	CreatedAt  null.Time   `json:"created_at"`
	UpdatedAt  null.Time   `json:"updated_at"`
	DeletedAt  null.Time   `json:"deleted_at"`
}
type Order struct {
	Id          null.String `json:"id"`
	UserId      null.String `json:"user_id"`
	ReceiptCode null.String `json:"receipt_code"`
	TotalPrice  null.Int    `json:"total_price"`
	TotalPaid   null.Int    `json:"total_paid"`
	TotalReturn null.Int    `json:"total_return"`
	CreatedAt   null.Time   `json:"created_at"`
	UpdatedAt   null.Time   `json:"updated_at"`
	DeletedAt   null.Time   `json:"deleted_at"`
}
