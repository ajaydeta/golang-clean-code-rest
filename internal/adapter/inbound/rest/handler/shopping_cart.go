package handler

import (
	"github.com/gofiber/fiber/v2"
	errors "github.com/rotisserie/eris"
	shoppingCartDto "synapsis-challenge/internal/adapter/inbound/rest/handler/dto/shopping_cart"
	"synapsis-challenge/shared"
)

func (h *Handler) AddShoppingCard(c *fiber.Ctx) error {
	resp := shared.NewJSONResponse()
	svc := h.GetServiceRegistry().GetShoppingCartService()

	dto := shoppingCartDto.DTO{}

	req, err := dto.AddAddShoppingCartRequestTransformIn(c)
	if err != nil {
		resp.SetReason(err).APIStatusBadRequest()
		return resp.Send(c)
	}

	id, err := svc.Add(c.Context(), req)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			resp.SetMessage("Product Not Found").APIStatusNotFound()
			return resp.Send(c)
		}

		resp.SetMessage("internal server error").SetReason(err)
		return resp.Send(c)
	}

	return resp.SetData(shoppingCartDto.AddShoppingCartResponse{ID: id}).APIStatusSuccess().Send(c)
}
