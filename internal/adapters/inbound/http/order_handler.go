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

type orderHandler struct {
	orderUs inbound.OrderUsecase
}

func NewOrderHandlerImpl(orderUs inbound.OrderUsecase) *orderHandler {
	return &orderHandler{
		orderUs: orderUs,
	}
}

func (h *orderHandler) CreateOrder(c fiber.Ctx) error {
	var (
		ctx      = c.Context()
		order    domain.Order
		validate = validator.New()
	)
	if err := c.Bind().JSON(&order); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if err := validate.Struct(order); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	err := h.orderUs.CreateOrder(ctx, &order)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.SendStatus(http.StatusCreated)
}

func (h *orderHandler) FetchListOrders(c fiber.Ctx) error {
	ctx := c.Context()
	orders, err := h.orderUs.FetchListOrders(ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if len(orders) == 0 {
		return c.SendStatus(http.StatusNoContent)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"orders": orders,
	})
}

func (h *orderHandler) FetchOrderById(c fiber.Ctx) error {
	ctx := c.Context()
	id, err := bson.ObjectIDFromHex(c.Params("order_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	order, err := h.orderUs.FetchOrderById(ctx, &id)
	if err != nil {
		if strings.Contains(err.Error(), mongo.ErrNoDocuments.Error()) {
			return c.SendStatus(http.StatusNoContent)
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"order": order,
	})
}

func (h *orderHandler) UpdateOrderById(c fiber.Ctx) error {
	var (
		ctx      = c.Context()
		order    domain.Order
		validate = validator.New()
	)
	id, err := bson.ObjectIDFromHex(c.Params("order_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	if err := c.Bind().JSON(&order); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if err := validate.Struct(order); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	response, err := h.orderUs.UpdateOrderById(ctx, &id, &order)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"order": response,
	})
}

func (h *orderHandler) DeleteOrderById(c fiber.Ctx) error {
	var (
		ctx   = c.Context()
		order domain.Order
	)
	id, err := bson.ObjectIDFromHex(c.Params("order_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if err := h.orderUs.DeleteOrderById(ctx, &id, &order); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}
