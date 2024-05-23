package dto

import (
	"github.com/gofiber/fiber/v2"
	"synapsis-challenge/internal/core/domain"
	"synapsis-challenge/shared"
)

func GetBaseFilter(c *fiber.Ctx) domain.Filter {
	return domain.Filter{
		Limit:  shared.StringToInt64(c.Query("limit"), 10),
		Offset: shared.StringToInt64(c.Query("offset"), 0),
		Search: c.Query("q"),
		Sort:   c.Query("sort"),
		Order:  c.Query("order"),
	}
}
