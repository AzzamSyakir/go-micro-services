package mock

import (
	"go-micro-services/src/auth-service/entity"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type ProductMock struct {
	CategoryMock *CategoryMock
	Data         []*entity.Product
}

func NewProductMock(
	categoryMock *CategoryMock,
) *ProductMock {
	currentTime := time.Now()
	currentTimeRfc3339 := currentTime.Format(time.RFC3339)
	currentTimeFromRfc3339, parseErr := time.Parse(time.RFC3339, currentTimeRfc3339)
	if parseErr != nil {
		panic(parseErr)
	}
	productMock := &ProductMock{
		Data: []*entity.Product{
			{
				Id:         null.NewString(uuid.NewString(), true),
				Sku:        null.NewString(uuid.NewString(), true),
				Name:       null.NewString(uuid.NewString(), true),
				Stock:      null.NewInt(0, true),
				Price:      null.NewInt(0, true),
				CategoryId: null.NewString(categoryMock.Data[0].Id.String, true),
				CreatedAt:  null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt:  null.NewTime(currentTimeFromRfc3339.Add(time.Duration(time.Duration.Seconds(0))), true),
				DeletedAt:  null.NewTime(time.Time{}, false),
			},
			{
				Id:         null.NewString(uuid.NewString(), true),
				Sku:        null.NewString(uuid.NewString(), true),
				Name:       null.NewString(uuid.NewString(), true),
				Stock:      null.NewInt(0, true),
				Price:      null.NewInt(0, true),
				CategoryId: null.NewString(categoryMock.Data[1].Id.String, true),
				CreatedAt:  null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt:  null.NewTime(currentTimeFromRfc3339.Add(time.Duration(time.Duration.Seconds(0))), true),
				DeletedAt:  null.NewTime(time.Time{}, false),
			},
		},
	}
	return productMock
}
