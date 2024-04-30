package mock

import (
	"go-micro-services/src/auth-service/entity"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type OrderMock struct {
	UserMock *UserMock
	Data     []*entity.Order
}

func NewOrderMock(
	userMock *UserMock,
) *OrderMock {
	currentTime := time.Now()
	currentTimeRfc3339 := currentTime.Format(time.RFC3339)
	currentTimeFromRfc3339, parseErr := time.Parse(time.RFC3339, currentTimeRfc3339)
	if parseErr != nil {
		panic(parseErr)
	}
	orderMock := &OrderMock{
		Data: []*entity.Order{
			{
				Id:          null.NewString(uuid.NewString(), true),
				UserId:      null.NewString(userMock.Data[0].Id.String, true),
				ReceiptCode: null.NewString(uuid.NewString(), true),
				TotalPrice:  null.NewInt(0, true),
				TotalPaid:   null.NewInt(0, true),
				TotalReturn: null.NewInt(0, true),
				CreatedAt:   null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt:   null.NewTime(currentTimeFromRfc3339.Add(time.Duration(time.Duration.Seconds(0))), true),
				DeletedAt:   null.NewTime(time.Time{}, false),
			},
			{
				Id:          null.NewString(uuid.NewString(), true),
				UserId:      null.NewString(userMock.Data[1].Id.String, true),
				ReceiptCode: null.NewString(uuid.NewString(), true),
				TotalPrice:  null.NewInt(0, true),
				TotalPaid:   null.NewInt(0, true),
				TotalReturn: null.NewInt(0, true),
				CreatedAt:   null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt:   null.NewTime(currentTimeFromRfc3339.Add(time.Duration(time.Duration.Seconds(0))), true),
				DeletedAt:   null.NewTime(time.Time{}, false),
			},
		},
	}
	return orderMock
}
