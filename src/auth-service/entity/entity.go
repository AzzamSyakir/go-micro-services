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
	Name      null.String `json:"name"`
	Email     null.String `json:"email"`
	Password  null.String `json:"password"`
	Balance   null.Int    `json:"balance"`
	CreatedAt null.Time   `json:"created_at"`
	UpdatedAt null.Time   `json:"updated_at"`
	DeletedAt null.Time   `json:"deleted_at"`
}
