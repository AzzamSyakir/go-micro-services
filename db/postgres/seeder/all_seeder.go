package seeder

import (
	"fmt"
)

type AllSeeder struct {
	User     *UserSeeder
	Session  *SessionSeeder
	Category *CategorySeeder
	Product  *ProductSeeder
	Order    *OrderSeeder
}

func NewAllSeeder(
	user *UserSeeder,
	session *SessionSeeder,
	category *CategorySeeder,
	product *ProductSeeder,
	order *OrderSeeder,
) *AllSeeder {
	allSeeder := &AllSeeder{
		User:     user,
		Session:  session,
		Category: category,
		Product:  product,
		Order:    order,
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
	fmt.Println("Seeder up finished.")
}

func (allSeeder *AllSeeder) Down() {
	fmt.Println("Seeder down started.")
	allSeeder.Session.Down()
	allSeeder.User.Down()
	allSeeder.Product.Down()
	allSeeder.Category.Down()
	allSeeder.Order.Down()
	fmt.Println("Seeder down finished.")
}
