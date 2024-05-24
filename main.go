package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"synapsis-challenge/config"
	ireg "synapsis-challenge/internal/adapter/inbound/registry"
	"synapsis-challenge/internal/adapter/inbound/rest/group"
	"synapsis-challenge/internal/adapter/inbound/rest/handler"
	oreg "synapsis-challenge/internal/adapter/outbound/registry"
)

func main() {
	godotenv.Load()
	app := fiber.New()

	db := config.InitMySQL()
	rdb := config.InitRedis()

	repoReg := oreg.NewRepositoryRegistry(rdb, db)
	serviceReg := ireg.NewServiceRegistry(repoReg)
	h := handler.New(serviceReg)

	group.NewCustomerRequest(app, h)
	group.NewProductRequest(app, h)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":3000")
}
