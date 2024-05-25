package handler

import (
	"github.com/gofiber/fiber/v2"
	errors "github.com/rotisserie/eris"
	"synapsis-challenge/internal/adapter/inbound/rest/handler/dto"
	shoppingCartDto "synapsis-challenge/internal/adapter/inbound/rest/handler/dto/shopping_cart"
	"synapsis-challenge/shared"
)

func (h *Handler) AddShoppingCard(c *fiber.Ctx) error {
	resp := shared.NewJSONResponse()
	svc := h.GetServiceRegistry().GetShoppingCartService()

	scDto := shoppingCartDto.DTO{}

	req, err := scDto.AddAddShoppingCartRequestTransformIn(c)
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

func (h *Handler) ListShoppingCart(c *fiber.Ctx) error {
	resp := shared.NewJSONResponse()
	respDto := shoppingCartDto.DTO{}
	svc := h.GetServiceRegistry().GetShoppingCartService()

	filter := dto.GetBaseFilter(c)
	shoppingCarts, err := svc.FindAll(c.Context(), filter)
	if err != nil {
		resp.SetMessage("internal server error").SetReason(err)
		return resp.Send(c)
	}

	resp.SetData(respDto.ToResponseList(shoppingCarts, filter)).APIStatusSuccess()
	return resp.Send(c)
}

func (h *Handler) GetCountShoppingCart(c *fiber.Ctx) error {
	resp := shared.NewJSONResponse()
	svc := h.GetServiceRegistry().GetShoppingCartService()

	filter := dto.GetBaseFilter(c)
	count, err := svc.CountAll(c.Context(), filter)
	if err != nil {
		resp.SetMessage("internal server error").SetReason(err)
		return resp.Send(c)
	}

	resp.SetCount(count).APIStatusSuccess()
	return resp.Send(c)
}

func (h *Handler) DeleteShoppingCart(c *fiber.Ctx) error {
	resp := shared.NewJSONResponse()
	svc := h.GetServiceRegistry().GetShoppingCartService()

	id := c.Params("id")

	err := svc.Delete(c.Context(), id)
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			resp.SetMessage("Shopping Cart Data Not Found").APIStatusNotFound()
			return resp.Send(c)
		}

		resp.SetMessage("internal server error").SetReason(err)
		return resp.Send(c)
	}

	resp.APIStatusSuccess()
	return resp.Send(c)
}
