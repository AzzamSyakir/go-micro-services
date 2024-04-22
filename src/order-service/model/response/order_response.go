package response

import (
	"go-micro-services/src/order-service/entity"

	"github.com/guregu/null"
)

type OrderResponse struct {
	Id          null.String            `json:"id"`
	UserId      null.String            `json:"user_id"`
	Name        null.String            `json:"name"`
	ReceiptCode null.String            `json:"receipt_code"`
	TotalPrice  null.Int               `json:"total_price"`
	TotalPaid   null.Int               `json:"total_paid"`
	TotalReturn null.Int               `json:"total_return"`
	CreatedAt   null.Time              `json:"created_at"`
	UpdatedAt   null.Time              `json:"updated_at"`
	Products    []entity.OrderProducts `json:"products"`
}
