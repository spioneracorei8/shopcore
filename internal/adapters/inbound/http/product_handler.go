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

type productHandler struct {
	productUs inbound.ProductUsecase
}

func NewProductHandlerImpl(productUs inbound.ProductUsecase) *productHandler {
	return &productHandler{
		productUs: productUs,
	}
}

func (h *productHandler) CreateProduct(c fiber.Ctx) error {
	var (
		ctx      = c.Context()
		product  domain.Product
		validate = validator.New()
	)
	if err := c.Bind().JSON(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if err := validate.Struct(product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	err := h.productUs.CreateProduct(ctx, &product)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.SendStatus(http.StatusCreated)
}

func (h *productHandler) FetchListProducts(c fiber.Ctx) error {
	ctx := c.Context()
	products, err := h.productUs.FetchListProducts(ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if len(products) == 0 {
		return c.SendStatus(http.StatusNoContent)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"products": products,
	})
}

func (h *productHandler) FetchProductById(c fiber.Ctx) error {
	ctx := c.Context()
	id, err := bson.ObjectIDFromHex(c.Params("product_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	product, err := h.productUs.FetchProductById(ctx, &id)
	if err != nil {
		if strings.Contains(err.Error(), mongo.ErrNoDocuments.Error()) {
			return c.SendStatus(http.StatusNoContent)
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"product": product,
	})
}

func (h *productHandler) UpdateProductById(c fiber.Ctx) error {
	var (
		ctx      = c.Context()
		product  domain.Product
		validate = validator.New()
	)
	id, err := bson.ObjectIDFromHex(c.Params("product_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}
	if err := c.Bind().JSON(&product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if err := validate.Struct(product); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	response, err := h.productUs.UpdateProductById(ctx, &id, &product)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"product": response,
	})
}

func (h *productHandler) DeleteProductById(c fiber.Ctx) error {
	var (
		ctx     = c.Context()
		product domain.Product
	)
	id, err := bson.ObjectIDFromHex(c.Params("product_id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if err := h.productUs.DeleteProductById(ctx, &id, &product); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.SendStatus(http.StatusNoContent)
}
