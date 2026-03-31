package app

import (
	"fmt"
	"shopcore/config"
	my_http "shopcore/internal/adapters/inbound/http"
	"shopcore/internal/adapters/outbound/mongodb"
	"shopcore/internal/core/services"

	"github.com/gofiber/fiber/v3"
)

type App struct {
	fiber *fiber.App
}

func New() *App {
	fiber := config.NewFiber()
	client := config.ConnectDatabase()

	//--------------------------
	// # REPOSITORY
	// --------------------------
	rnRepo := mongodb.NewRunNumberRepoImpl(client)
	customerRepo := mongodb.NewCustomerRepoImpl(client)
	productRepo := mongodb.NewProductRepoImpl(client)
	orderRepo := mongodb.NewOrderRepoImpl(client)

	//--------------------------
	// # USECASE
	// --------------------------
	rnUs := services.NewRunNumberUsecaseImpl(rnRepo)
	customerUs := services.NewCustomerUsecaseImpl(customerRepo)
	productUs := services.NewProductUsecaseImpl(productRepo)
	orderUs := services.NewOrderUsecaseImpl(rnUs, productUs, orderRepo)

	//--------------------------
	// # HANDLER
	// --------------------------
	rnHandler := my_http.NewRunNumberHandlerImpl(rnUs)
	customerHandler := my_http.NewCustomerHandlerImpl(customerUs)
	productHandler := my_http.NewProductHandlerImpl(productUs)
	orderHandler := my_http.NewOrderHandlerImpl(orderUs)

	route := my_http.NewRoute(fiber)
	route.NewCustomerRoutes(customerHandler)
	route.NewProductRoutes(productHandler)
	route.NewOrderRoutes(orderHandler)
	route.NewRunNumberRoutes(rnHandler)

	return &App{fiber: fiber}
}

func (a *App) Run() {
	a.fiber.Listen(fmt.Sprintf(":%s", config.APP_PORT))
}
