package http

import "github.com/gofiber/fiber/v3"

type route struct {
	f *fiber.App
}

func NewRoute(f *fiber.App) *route {
	return &route{f: f}
}

func (r *route) NewCustomerRoutes(handler *customerHandler) {
	api := r.f.Group("/api")

	api.Post("/v1/customer", handler.CreateCustomer)
	api.Get("/v1/customer", handler.FetchListCustomers)
	api.Get("/v1/customer/:customer_id", handler.FetchCustomerById)
	api.Put("/v1/customer/:customer_id", handler.UpdateCustomerById)
	api.Delete("/v1/customer/:customer_id", handler.DeleteCustomerById)
}

func (r *route) NewProductRoutes(handler *productHandler) {
	api := r.f.Group("/api")

	api.Post("/v1/product", handler.CreateProduct)
	api.Get("/v1/product", handler.FetchListProducts)
	api.Get("/v1/product/:product_id", handler.FetchProductById)
	api.Put("/v1/product/:product_id", handler.UpdateProductById)
	api.Delete("/v1/product/:product_id", handler.DeleteProductById)
}

func (r *route) NewOrderRoutes(handler *orderHandler) {
	api := r.f.Group("/api")

	api.Post("/v1/order", handler.CreateOrder)
	api.Get("/v1/order", handler.FetchListOrders)
	api.Get("/v1/order/:order_id", handler.FetchOrderById)
	api.Put("/v1/order/:order_id", handler.UpdateOrderById)
	api.Delete("/v1/order/:order_id", handler.DeleteOrderById)
}

func (r *route) NewRunNumberRoutes(handler *runNumberHandler) {
	api := r.f.Group("/api")

	api.Post("/v1/run_number", handler.CreateRunNumber)
}
