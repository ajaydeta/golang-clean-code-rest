package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (i *Handler) GetCustomers(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{})
}
