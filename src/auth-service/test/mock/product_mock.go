package mock

import (
	"go-micro-services/src/auth-service/entity"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
)

type ProductMock struct {
	Data         []*entity.Product
	CategoryMock *CategoryMock
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
	ProductMock := &ProductMock{
		CategoryMock: categoryMock,
		Data: []*entity.Product{
			{
				Id:         null.NewString(uuid.NewString(), true),
				Sku:        null.NewString(uuid.NewString(), true),
				Name:       null.NewString(uuid.NewString(), true),
				Stock:      null.NewInt(100000, true),
				Price:      null.NewInt(100000, true),
				CategoryId: null.NewString(categoryMock.Data[0].Id.String, true),
				CreatedAt:  null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt:  null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				DeletedAt:  null.NewTime(time.Time{}, false),
			},
			{
				Id:         null.NewString(uuid.NewString(), true),
				Sku:        null.NewString(uuid.NewString(), true),
				Name:       null.NewString(uuid.NewString(), true),
				Stock:      null.NewInt(100000, true),
				Price:      null.NewInt(100000, true),
				CategoryId: null.NewString(categoryMock.Data[0].Id.String, true),
				CreatedAt:  null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				UpdatedAt:  null.NewTime(currentTimeFromRfc3339.Add(0*time.Second), true),
				DeletedAt:  null.NewTime(time.Time{}, false),
			},
		},
	}
	return ProductMock
}
