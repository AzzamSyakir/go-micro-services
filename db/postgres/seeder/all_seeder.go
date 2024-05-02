package seeder

import (
	"fmt"
)

type AllSeeder struct {
	User    *UserSeeder
	Session *SessionSeeder
}

func NewAllSeeder(
	user *UserSeeder,
	session *SessionSeeder) *AllSeeder {
	allSeeder := &AllSeeder{
		User:    user,
		Session: session,
	}
	return allSeeder
}

func (allSeeder *AllSeeder) Up() {
	fmt.Println("Seeder up started.")
	allSeeder.User.Up()
	allSeeder.Session.Up()
	fmt.Println("Seeder up finished.")
}

func (allSeeder *AllSeeder) Down() {
	fmt.Println("Seeder down started.")
	allSeeder.Session.Down()
	allSeeder.User.Down()
	fmt.Println("Seeder down finished.")
}
