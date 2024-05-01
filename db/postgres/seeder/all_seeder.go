package seeder

import (
	"fmt"
)

type AllSeeder struct {
	User         *UserSeeder
	Session      *SessionSeeder
	Product      *ProductSeeder
	Order        *OrderSeeder
	OrderProduct *OrderProductSeeder
	Category     *CategorySeeder
}

func NewAllSeeder(
	user *UserSeeder,
	session *SessionSeeder,
	product *ProductSeeder,
	order *OrderSeeder,
	orderProduct *OrderProductSeeder,
	category *CategorySeeder,
) *AllSeeder {
	allSeeder := &AllSeeder{
		User:         user,
		Session:      session,
		Product:      product,
		Order:        order,
		OrderProduct: orderProduct,
		Category:     category,
	}
	return allSeeder
}

func (allSeeder *AllSeeder) Up() {
	fmt.Println("Seeder up started.")
	allSeeder.User.Up()
	allSeeder.Session.Up()
	allSeeder.Category.Up()
	allSeeder.Product.Up()
	allSeeder.Order.Up()
	allSeeder.OrderProduct.Up()
	fmt.Println("Seeder up finished.")
}

func (allSeeder *AllSeeder) Down() {
	fmt.Println("Seeder down started.")
	allSeeder.User.Down()
	allSeeder.Session.Down()
	allSeeder.Product.Down()
	allSeeder.OrderProduct.Down()
	allSeeder.Order.Down()
	allSeeder.Category.Down()
	fmt.Println("Seeder down finished.")
}
