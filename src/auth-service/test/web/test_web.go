package web

import (
	seeder "go-micro-services/db/postgres/seeder"
	auth_container "go-micro-services/src/auth-service/container"
	order_container "go-micro-services/src/order-service/container"
	product_container "go-micro-services/src/product-service/container"
	user_container "go-micro-services/src/user-service/container"
	"net/http/httptest"
)

type TestWeb struct {
	Server           *httptest.Server
	AllSeeder        *seeder.AllSeeder
	UserContainer    *user_container.WebContainer
	OrderContainer   *order_container.WebContainer
	ProductContainer *product_container.WebContainer
	AuthContainer    *auth_container.WebContainer
}

func NewTestWeb() *TestWeb {
	userWebContainer := user_container.NewWebContainer()
	productWebContainer := product_container.NewWebContainer()
	authWebContainer := auth_container.NewWebContainer()
	orderWebContainer := order_container.NewWebContainer()

	server := httptest.NewServer(authWebContainer.Route.Router)

	testWeb := &TestWeb{
		Server:           server,
		UserContainer:    userWebContainer,
		AuthContainer:    authWebContainer,
		ProductContainer: productWebContainer,
		OrderContainer:   orderWebContainer,
	}

	return testWeb
}

func (web *TestWeb) GetAllSeeder() *seeder.AllSeeder {
	userSeeder := seeder.NewUserSeeder(web.UserContainer.UserDB)
	categorySeeder := seeder.NewCategorySeeder(web.ProductContainer.ProductDB)
	productSeeder := seeder.NewProductSeeder(web.ProductContainer.ProductDB, categorySeeder)
	sessionSeeder := seeder.NewSessionSeeder(web.AuthContainer.AuthDB, userSeeder)
	orderSeeder := seeder.NewOrderSeeder(web.OrderContainer.OrderDB, userSeeder)
	seederConfig := seeder.NewAllSeeder(
		userSeeder,
		sessionSeeder,
		categorySeeder,
		productSeeder,
		orderSeeder,
	)
	return seederConfig
}

func GetTestWeb() *TestWeb {
	testWeb := NewTestWeb()
	testWeb.AllSeeder = testWeb.GetAllSeeder()
	return testWeb
}
