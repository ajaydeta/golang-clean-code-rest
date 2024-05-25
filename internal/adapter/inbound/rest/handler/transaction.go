package handler

import (
	"github.com/gofiber/fiber/v2"
	errors "github.com/rotisserie/eris"
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

func (h *Handler) PayoffTransaction(c *fiber.Ctx) error {
	resp := shared.NewJSONResponse()
	svc := h.GetServiceRegistry().GetTransactionService()
	respDto := transactionDto.DTO{}
	id := c.Params("id")

	payment, err := svc.PayoffTransaction(c.Context(), id)
	if err != nil {

		switch {
		case errors.Is(err, shared.ErrAlreadyPaid):
			resp.SetMessage("Payment Already Paid").SetCode(fiber.StatusBadRequest)
		case errors.Is(err, shared.ErrNotFound):
			resp.SetMessage("Payment Not Found").SetCode(fiber.StatusBadRequest)
		default:
			resp.SetMessage("internal server error").SetReason(err)
		}

		return resp.Send(c)
	}

	return resp.SetData(respDto.ToPaymentResp(payment)).APIStatusSuccess().Send(c)
}

func (h *Handler) GetTransaction(c *fiber.Ctx) error {
	resp := shared.NewJSONResponse()
	svc := h.GetServiceRegistry().GetTransactionService()
	respDto := transactionDto.DTO{}
	id := c.Params("id")

	transaction, err := svc.FindId(c.Context(), id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return resp.APIStatusNotFound().Send(c)
		}

		resp.SetMessage("internal server error").SetReason(err)
		return resp.Send(c)
	}

	return resp.SetData(respDto.ToCreateResp(transaction)).APIStatusSuccess().Send(c)
}
