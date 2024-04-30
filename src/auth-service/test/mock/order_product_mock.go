package mock

import (
	"go-micro-services/src/auth-service/entity"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type OrderProductsMock struct {
	OrderMock   *OrderMock
	ProductMock *ProductMock
	Data        []*entity.OrderProducts
}

func NewOrderProductsMock(
	orderMock *OrderMock,
	productMock *ProductMock,
) *OrderProductsMock {
	currentTime := time.Now()
	currentTimeRfc3339 := currentTime.Format(time.RFC3339)
	currentTimeFromRfc3339, parseErr := time.Parse(time.RFC3339, currentTimeRfc3339)
	if parseErr != nil {
		panic(parseErr)
	}
	orderProductMock := &OrderProductsMock{
		Data: []*entity.OrderProducts{
			{
				Id:         null.NewString(uuid.NewString(), true),
				OrderId:    null.NewString(orderMock.Data[0].Id.String, true),
				ProductId:  null.NewString(productMock.Data[0].Id.String, true),
				TotalPrice: null.NewInt(0, true),
				Qty:        null.NewInt(0, true),
				CreatedAt:  null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt:  null.NewTime(currentTimeFromRfc3339.Add(time.Duration(time.Duration.Seconds(0))), true),
				DeletedAt:  null.NewTime(time.Time{}, false),
			},
			{
				Id:         null.NewString(uuid.NewString(), true),
				OrderId:    null.NewString(orderMock.Data[1].Id.String, true),
				ProductId:  null.NewString(productMock.Data[1].Id.String, true),
				TotalPrice: null.NewInt(0, true),
				Qty:        null.NewInt(0, true),
				CreatedAt:  null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt:  null.NewTime(currentTimeFromRfc3339.Add(time.Duration(time.Duration.Seconds(0))), true),
				DeletedAt:  null.NewTime(time.Time{}, false),
			},
		},
	}
	return orderProductMock
}
