package handler

import (
	"github.com/gofiber/fiber/v2"
	productDto "synapsis-challenge/internal/adapter/inbound/rest/handler/dto/product"
	"synapsis-challenge/shared"
)

func (h *Handler) ListProduct(c *fiber.Ctx) error {
	resp := shared.NewJSONResponse()
	respDto := productDto.DTO{}
	svc := h.GetServiceRegistry().GetProductService()

	filter := respDto.GetListProductFilter(c)
	products, err := svc.FindAll(c.UserContext(), filter)
	if err != nil {
		resp.SetMessage("internal server error").SetReason(err)
		return resp.Send(c)
	}

	resp.SetData(respDto.ToResponseList(products, filter.Filter)).APIStatusSuccess()
	return resp.Send(c)
}

func (h *Handler) GetCountProduct(c *fiber.Ctx) error {
	resp := shared.NewJSONResponse()
	respDto := productDto.DTO{}
	svc := h.GetServiceRegistry().GetProductService()

	filter := respDto.GetListProductFilter(c)
	count, err := svc.CountAll(c.UserContext(), filter)
	if err != nil {
		resp.SetMessage("internal server error").SetReason(err)
		return resp.Send(c)
	}

	resp.SetCount(count).APIStatusSuccess()
	return resp.Send(c)
}

func (h *Handler) GetProduct(c *fiber.Ctx) error {
	resp := shared.NewJSONResponse()
	svc := h.GetServiceRegistry().GetProductService()
	respDto := productDto.DTO{}

	product, err := svc.FindId(c.UserContext(), c.Params("id"))
	if err != nil {
		resp.SetMessage("internal server error").SetReason(err)
		return resp.Send(c)
	}

	if product == nil {
		return resp.APIStatusNotFound().Send(c)
	}

	return resp.SetData(respDto.ToResponse(product)).APIStatusSuccess().Send(c)
}
