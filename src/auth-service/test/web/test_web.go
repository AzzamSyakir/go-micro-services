package web

import (
	seeder "go-micro-services/db/postgres/seeder"
	AuthContainer "go-micro-services/src/auth-service/container"
	OrderContainer "go-micro-services/src/order-service/container"
	productContainer "go-micro-services/src/product-service/container"
	userContainer "go-micro-services/src/user-service/container"
	"net/http/httptest"
)

type TestWeb struct {
	AllSeeder        *seeder.AllSeeder
	AuthServer       *httptest.Server
	AuthContainer    *AuthContainer.WebContainer
	UserServer       *httptest.Server
	UserContainer    *userContainer.WebContainer
	ProductServer    *httptest.Server
	ProductContainer *productContainer.WebContainer
	OrderServer      *httptest.Server
	OrderContainer   *OrderContainer.WebContainer
}

func NewTestWeb() *TestWeb {
	authWebContainer := AuthContainer.NewWebContainer()
	userWebContainer := userContainer.NewWebContainer()
	productWebContainer := productContainer.NewWebContainer()
	orderWebContainer := OrderContainer.NewWebContainer()

	authServer := httptest.NewServer(authWebContainer.Route.Router)
	userServer := httptest.NewServer(userWebContainer.Route.Router)
	productServer := httptest.NewServer(productWebContainer.Route.Router)
	orderServer := httptest.NewServer(orderWebContainer.Route.Router)

	testWeb := &TestWeb{
		AuthServer:       authServer,
		AuthContainer:    authWebContainer,
		UserServer:       userServer,
		UserContainer:    userWebContainer,
		ProductServer:    productServer,
		ProductContainer: productWebContainer,
		OrderServer:      orderServer,
		OrderContainer:   orderWebContainer,
	}

	return testWeb
}

func (web *TestWeb) GetAllSeeder() *seeder.AllSeeder {
	userSeeder := seeder.NewUserSeeder(web.UserContainer.UserDB)
	sessionSeeder := seeder.NewSessionSeeder(web.AuthContainer.AuthDB, userSeeder)
	categorySeeder := seeder.NewCategorSeeder(web.ProductContainer.ProductDB)
	productSeeder := seeder.NewProductSeeder(web.ProductContainer.ProductDB, categorySeeder)
	orderSeeder := seeder.NewOrderSeeder(web.OrderContainer.OrderDB, userSeeder)
	orderProductSeeder := seeder.NewOrderProductSeeder(web.OrderContainer.OrderDB, orderSeeder, productSeeder)
	seederConfig := seeder.NewAllSeeder(
		userSeeder,
		sessionSeeder,
		productSeeder,
		orderSeeder,
		orderProductSeeder,
		categorySeeder,
	)
	return seederConfig
}

func GetTestWeb() *TestWeb {
	testWeb := NewTestWeb()
	testWeb.AllSeeder = testWeb.GetAllSeeder()
	return testWeb
}
