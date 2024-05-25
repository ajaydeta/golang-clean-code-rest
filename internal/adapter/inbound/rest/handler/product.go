package handler

import (
	"github.com/gofiber/fiber/v2"
	errors "github.com/rotisserie/eris"
	productDto "synapsis-challenge/internal/adapter/inbound/rest/handler/dto/product"
	"synapsis-challenge/shared"
)

func (h *Handler) ListProduct(c *fiber.Ctx) error {
	resp := shared.NewJSONResponse()
	respDto := productDto.DTO{}
	svc := h.GetServiceRegistry().GetProductService()

	filter := respDto.GetListProductFilter(c)
	products, err := svc.FindAll(c.Context(), filter)
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
	count, err := svc.CountAll(c.Context(), filter)
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

	product, err := svc.FindId(c.Context(), c.Params("id"))
	if err != nil {
		if errors.Is(err, shared.ErrNotFound) {
			return resp.APIStatusNotFound().Send(c)
		}

		resp.SetMessage("internal server error").SetReason(err)
		return resp.Send(c)
	}

	return resp.SetData(respDto.ToResponse(product)).APIStatusSuccess().Send(c)
}
