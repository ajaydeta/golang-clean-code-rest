package handler

import (
	"github.com/gofiber/fiber/v2"
	transactionDto "synapsis-challenge/internal/adapter/inbound/rest/handler/dto/transaction"
	"synapsis-challenge/shared"
)

func (h *Handler) CreateTransaction(c *fiber.Ctx) error {
	resp := shared.NewJSONResponse()
	respDto := transactionDto.DTO{}
	svc := h.GetServiceRegistry().GetTransactionService()

	req, err := respDto.CreateTransformIn(c)
	if err != nil {
		resp.SetReason(err).APIStatusBadRequest()
		return resp.Send(c)
	}

	transaction, err := svc.CreateTransaction(c.Context(), req)
	if err != nil {
		resp.SetMessage("internal server error").SetReason(err)
		return resp.Send(c)
	}

	return resp.SetData(respDto.ToCreateResp(transaction)).APIStatusSuccess().Send(c)
}
