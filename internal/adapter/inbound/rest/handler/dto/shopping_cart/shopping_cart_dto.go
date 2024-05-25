package shopping_cart

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"synapsis-challenge/internal/core/domain"
)

type (
	AddShoppingCartRequest struct {
		ProductID string  `json:"product_id" validate:"required"`
		Notes     string  `json:"notes"`
		Qty       float64 `json:"qty" validate:"gt=0"`
	}

	AddShoppingCartResponse struct {
		ID string `json:"id"`
	}

	DTO struct{}
)

func (c *AddShoppingCartRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}

func (d DTO) AddAddShoppingCartRequestTransformIn(ctx *fiber.Ctx) (*domain.ShoppingCart, error) {
	c := new(AddShoppingCartRequest)

	if err := ctx.BodyParser(c); err != nil {
		return nil, err
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return &domain.ShoppingCart{
		ProductID: c.ProductID,
		Notes:     c.Notes,
		Qty:       c.Qty,
	}, nil

}
