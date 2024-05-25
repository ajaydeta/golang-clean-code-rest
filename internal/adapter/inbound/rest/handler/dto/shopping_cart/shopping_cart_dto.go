package shopping_cart

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"synapsis-challenge/internal/adapter/inbound/rest/handler/dto"
	"synapsis-challenge/internal/core/domain"
	"synapsis-challenge/shared"
	"time"
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

	Response struct {
		ID         string    `json:"id"`
		CustomerID string    `json:"customer_id"`
		ProductID  string    `json:"product_id"`
		Notes      string    `json:"notes"`
		Qty        float64   `json:"qty"`
		CreatedAt  time.Time `json:"created_at"`
		Product    *Product  `json:"product"`
	}

	Product struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Price     float64   `json:"price"`
		CreatedAt time.Time `json:"created_at"`
	}

	ResponseList struct {
		Response   []Response     `json:"shopping_cart"`
		Pagination dto.Pagination `json:"pagination"`
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

func (d DTO) ToResponseList(shoppingCards []domain.ShoppingCart, filter domain.Filter) *ResponseList {
	resp := new(ResponseList)
	resp.Response = make([]Response, len(shoppingCards))

	for i, v := range shoppingCards {
		copier.CopyWithOption(&resp.Response[i], v, copier.Option{DeepCopy: true})
	}

	page, perPage := shared.GetPageAndPerPage(filter.Limit, filter.Offset)
	resp.Pagination = dto.Pagination{
		Page:    page,
		PerPage: perPage,
	}
	return resp
}
