package http

import (
	"net/http"
	"shopcore/internal/core/domain"
	"shopcore/internal/core/ports/inbound"

	"github.com/gofiber/fiber/v3"
)

type runNumberHandler struct {
	runNumberUs inbound.RunNumberUsecase
}

func NewRunNumberHandlerImpl(runNumberUs inbound.RunNumberUsecase) *runNumberHandler {
	return &runNumberHandler{
		runNumberUs: runNumberUs,
	}
}

func (h *runNumberHandler) CreateRunNumber(c fiber.Ctx) error {
	var (
		ctx = c.Context()
		rn  domain.RunNumber
	)
	if err := c.Bind().JSON(&rn); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	err := h.runNumberUs.CreateRunNumber(ctx, &rn)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.SendStatus(http.StatusCreated)
}
