package mock

import (
	"go-micro-services/src/auth-service/entity"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type OrderProductMock struct {
	Data        []*entity.OrderProducts
	orderMock   *OrderMock
	productMock *ProductMock
}

func NewOrderProductMock(
	orderMock *OrderMock,
	productMock *ProductMock,
) *OrderProductMock {
	currentTime := time.Now()
	currentTimeRfc3339 := currentTime.Format(time.RFC3339)
	currentTimeFromRfc3339, parseErr := time.Parse(time.RFC3339, currentTimeRfc3339)
	if parseErr != nil {
		panic(parseErr)
	}
	OrderProductMock := &OrderProductMock{
		orderMock:   orderMock,
		productMock: productMock,
		Data: []*entity.OrderProducts{
			{
				Id:         null.NewString(uuid.NewString(), true),
				OrderId:    null.NewString(orderMock.Data[0].Id.String, true),
				ProductId:  null.NewString(productMock.Data[0].Id.String, true),
				TotalPrice: null.NewInt(100000, true),
				Qty:        null.NewInt(200000, true),
				CreatedAt:  null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt:  null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
			},
			{
				Id:         null.NewString(uuid.NewString(), true),
				OrderId:    null.NewString(orderMock.Data[1].Id.String, true),
				ProductId:  null.NewString(productMock.Data[1].Id.String, true),
				TotalPrice: null.NewInt(100000, true),
				Qty:        null.NewInt(200000, true),
				CreatedAt:  null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt:  null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
			},
		},
	}
	return OrderProductMock
}
