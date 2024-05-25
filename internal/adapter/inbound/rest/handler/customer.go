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

	id, err := svc.RegisterCustomer(c.Context(), req)
	if err != nil {
		if errors.Is(err, shared.ErrAlreadyExist) {
			resp.SetMessage("Email already exists").SetCode(fiber.StatusBadRequest)
			return resp.Send(c)
		}

		resp.SetMessage("internal server error").SetReason(err)
		return resp.Send(c)
	}

	return resp.SetData(customerDto.RegisterCustomerResp{ID: id}).APIStatusSuccess().Send(c)
}

func (h *Handler) SignIn(c *fiber.Ctx) error {
	resp := shared.NewJSONResponse()
	svc := h.GetServiceRegistry().GetCustomerService()

	dto := customerDto.DTO{}
	req, err := dto.SignInTransformIn(c)
	if err != nil {
		resp.SetReason(err).APIStatusBadRequest()
		return resp.Send(c)
	}

	data, err := svc.SignIn(c.Context(), req)
	if err != nil {

		switch {
		case errors.Is(err, shared.ErrInvalidPassword) || errors.Is(err, shared.ErrNotFound):
			resp.SetMessage("Credential invalid").SetCode(fiber.StatusBadRequest)
		case errors.Is(err, shared.ErrAlreadyExist):
			resp.SetMessage("Email already login").SetCode(fiber.StatusBadRequest)
		default:
			resp.SetMessage("internal server error").SetReason(err)
		}

		return resp.Send(c)
	}

	return resp.SetData(dto.ToSignInResp(data)).APIStatusSuccess().Send(c)
}
