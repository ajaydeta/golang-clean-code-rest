package main

import (
	"github.com/gofiber/fiber/v2"
	ireg "synapsis-challenge/internal/adapter/inbound/registry"
	"synapsis-challenge/internal/adapter/inbound/rest/group"
	"synapsis-challenge/internal/adapter/inbound/rest/handler"
	oreg "synapsis-challenge/internal/adapter/outbound/registry"
)

func main() {
	app := fiber.New()

	repoReg := oreg.NewRepositoryRegistry()
	serviceReg := ireg.NewServiceRegistry(repoReg)
	h := handler.New(serviceReg)

	group.NewCustomerRequest(app, h)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}
