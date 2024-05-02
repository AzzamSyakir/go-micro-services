package web

import (
	seeder "go-micro-services/db/postgres/seeder"
	auth_container "go-micro-services/src/auth-service/container"
	user_container "go-micro-services/src/user-service/container"
	"net/http/httptest"
)

type TestWeb struct {
	Server        *httptest.Server
	AllSeeder     *seeder.AllSeeder
	UserContainer *user_container.WebContainer
	AuthContainer *auth_container.WebContainer
}

func NewTestWeb() *TestWeb {
	userWebContainer := user_container.NewWebContainer()
	authWebContainer := auth_container.NewWebContainer()

	server := httptest.NewServer(authWebContainer.Route.Router)

	testWeb := &TestWeb{
		Server:        server,
		UserContainer: userWebContainer,
		AuthContainer: authWebContainer,
	}

	return testWeb
}

func (web *TestWeb) GetAllSeeder() *seeder.AllSeeder {
	userSeeder := seeder.NewUserSeeder(web.UserContainer.UserDB)
	sessionSeeder := seeder.NewSessionSeeder(web.AuthContainer.AuthDB, userSeeder)
	seederConfig := seeder.NewAllSeeder(
		userSeeder,
		sessionSeeder,
	)
	return seederConfig
}

func GetTestWeb() *TestWeb {
	testWeb := NewTestWeb()
	testWeb.AllSeeder = testWeb.GetAllSeeder()
	return testWeb
}
