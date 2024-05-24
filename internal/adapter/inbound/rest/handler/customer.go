package handler

import (
	"github.com/gofiber/fiber/v2"
	errors "github.com/rotisserie/eris"
	customerDto "synapsis-challenge/internal/adapter/inbound/rest/handler/dto/customer"
	"synapsis-challenge/shared"
)

func (h *Handler) Register(c *fiber.Ctx) error {

	resp := shared.NewJSONResponse()
	svc := h.GetServiceRegistry().GetCustomerService()

	dto := customerDto.DTO{}
	req, err := dto.CreateTransformIn(c)
	if err != nil {
		resp.SetReason(err).APIStatusBadRequest()
		return resp.Send(c)
	}

	id, err := svc.RegisterCustomer(c.UserContext(), req)
	if err != nil {
		if errors.Is(err, shared.ErrAlreadyExist) {
			resp.SetMessage("Email already exists").SetCode(fiber.StatusBadRequest)
			return resp.Send(c)
		}

		resp.SetReason(err).APIStatusBadRequest()
		return resp.Send(c)
	}

	return resp.SetData(customerDto.RegisterCustomerResp{ID: id}).APIStatusSuccess().Send(c)
}
