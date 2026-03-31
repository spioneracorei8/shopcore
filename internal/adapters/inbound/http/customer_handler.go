package http

import (
	"net/http"
	"shopcore/internal/core/domain"
	"shopcore/internal/core/ports/inbound"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type customerHandler struct {
	customerUs inbound.CustomerUsecase
}

func NewCustomerHandlerImpl(customerUs inbound.CustomerUsecase) *customerHandler {
	return &customerHandler{
		customerUs: customerUs,
	}
}

func (h *customerHandler) CreateCustomer(c fiber.Ctx) error {
	var (
		ctx      = c.Context()
		customer domain.Customer
		validate = validator.New()
	)
	if err := c.Bind().JSON(&customer); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if err := validate.Struct(customer); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	err := h.customerUs.CreateCustomer(ctx, &customer)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.SendStatus(http.StatusCreated)
}

func (h *customerHandler) FetchListCustomers(c fiber.Ctx) error {
	var (
		ctx = c.Context()
	)

	customers, err := h.customerUs.FetchListCustomers(ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"customers": customers,
	})
}

func (h *customerHandler) FetchCustomerById(c fiber.Ctx) error {
	var (
		ctx = c.Context()
	)
	id, err := bson.ObjectIDFromHex(c.Params("customer_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	customer, err := h.customerUs.FetchCustomerById(ctx, &id)
	if err != nil {
		if strings.Contains(err.Error(), mongo.ErrNoDocuments.Error()) {
			return c.SendStatus(http.StatusNoContent)
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"customer": customer,
	})
}

func (h *customerHandler) UpdateCustomerById(c fiber.Ctx) error {
	var (
		ctx      = c.Context()
		req      domain.Customer
		validate = validator.New()
	)
	id, err := bson.ObjectIDFromHex(c.Params("customer_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if err := validate.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	customer, err := h.customerUs.UpdateCustomerById(ctx, &id, &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":  "successful",
		"customer": customer,
	})
}

func (h *customerHandler) DeleteCustomerById(c fiber.Ctx) error {
	var (
		ctx      = c.Context()
		customer domain.Customer
	)
	id, err := bson.ObjectIDFromHex(c.Params("customer_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	err = h.customerUs.DeleteCustomerById(ctx, &id, &customer)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}
