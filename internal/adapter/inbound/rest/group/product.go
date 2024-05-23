package group

import (
	"github.com/gofiber/fiber/v2"
	"synapsis-challenge/internal/adapter/inbound/rest/handler"
)

func NewProductRequest(app *fiber.App, handler *handler.Handler) {
	customer := app.Group("/product")
	customer.Get("/", handler.ListProduct)
	customer.Get("/count", handler.GetCountProduct)
	customer.Get("/:id", handler.GetProduct)
}
