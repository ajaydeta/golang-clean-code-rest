package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"synapsis-challenge/cmd"
	"synapsis-challenge/config"
	ireg "synapsis-challenge/internal/adapter/inbound/registry"
	"synapsis-challenge/internal/adapter/inbound/rest/group"
	"synapsis-challenge/internal/adapter/inbound/rest/handler"
	oreg "synapsis-challenge/internal/adapter/outbound/registry"
)

func main() {
	godotenv.Load()
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		TimeZone:   "Asia/Jakarta",
		Format:     "${pid} [${ip}]:${port} ${locals:requestid} ${status} - ${method} ${path}  ${latency}\n",
		TimeFormat: "02-Jan-2006",
	}))

	db := config.InitMySQL()
	rdb := config.InitRedis()

	repoReg := oreg.NewRepositoryRegistry(rdb, db)
	serviceReg := ireg.NewServiceRegistry(repoReg)
	h := handler.New(serviceReg)

	group.NewCustomerRequest(app, h)
	group.NewProductRequest(app, h)
	group.NewShoppingCartRequest(app, h)
	group.NewTransactionRequest(app, h)

	err := cmd.MigrateAndSeed(db)
	if err != nil {
		panic(err)
		return
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/stack", func(c *fiber.Ctx) error {
		return c.JSON(app.Stack())
	})

	app.Listen(":3000")
}
