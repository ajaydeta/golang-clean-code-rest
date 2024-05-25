package group

import (
	"github.com/gofiber/fiber/v2"
	"synapsis-challenge/internal/adapter/inbound/rest/handler"
)

func NewTransactionRequest(app *fiber.App, handler *handler.Handler) {
	transactionCart := app.Group("/transaction")
	transactionCart.Use(handler.VerifyAuth)
	transactionCart.Post("/", handler.CreateTransaction)
	transactionCart.Get("/:id", handler.GetTransaction)
	transactionCart.Post("/payoff/:id", handler.PayoffTransaction)
}
