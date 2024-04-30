package seeder

import (
	"fmt"
)

type AllSeeder struct {
	User    *UserSeeder
	Session *SessionSeeder
	Post    *PostSeeder
}

func NewAllSeeder(
	user *UserSeeder,
	session *SessionSeeder,
	post *PostSeeder,
) *AllSeeder {
	allSeeder := &AllSeeder{
		User:    user,
		Session: session,
		Post:    post,
	}
	return allSeeder
}

func (allSeeder *AllSeeder) Up() {
	fmt.Println("Seeder up started.")
	allSeeder.User.Up()
	allSeeder.Session.Up()
	allSeeder.Post.Up()
	fmt.Println("Seeder up finished.")
}

func (allSeeder *AllSeeder) Down() {
	fmt.Println("Seeder down started.")
	allSeeder.Post.Down()
	allSeeder.Session.Down()
	allSeeder.User.Down()
	fmt.Println("Seeder down finished.")
}
