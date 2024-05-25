package group

import (
	"github.com/gofiber/fiber/v2"
	"synapsis-challenge/internal/adapter/inbound/rest/handler"
)

func NewShoppingCartRequest(app *fiber.App, handler *handler.Handler) {
	shoppingCart := app.Group("/shopping-cart")
	shoppingCart.Use(handler.VerifyAuth)
	shoppingCart.Post("/", handler.AddShoppingCard)
}
