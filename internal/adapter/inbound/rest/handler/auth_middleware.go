package handler

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
)

func (h *Handler) VerifyAuth(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	svc := h.GetServiceRegistry().GetCustomerService()

	if len(token) == 0 {
		return unauthorizedResponse(c)
	}

	token = strings.Replace(token, "Bearer ", "", -1)

	customerId, err := svc.VerifyToken(token)
	if err != nil {
		log.Println(err.Error())
		return unauthorizedResponse(c)
	}

	c.Context().SetUserValue("customerId", customerId)

	return c.Next()
}

func unauthorizedResponse(c *fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "unauthorized",
	})
}
