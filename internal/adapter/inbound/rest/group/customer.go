package group

import (
	"github.com/gofiber/fiber/v2"
	"synapsis-challenge/internal/adapter/inbound/rest/handler"
)

func NewCustomerRequest(app *fiber.App, handler *handler.Handler) {
	app.Post("/register", handler.Register)
	app.Post("/sign-in", handler.SignIn)
}
