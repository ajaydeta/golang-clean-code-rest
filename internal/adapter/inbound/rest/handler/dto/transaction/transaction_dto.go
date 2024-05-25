package transaction

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"synapsis-challenge/internal/core/domain"
	"time"
)

type (
	DTO struct{}

	CreateRequest struct {
		Discount        float64  `json:"discount"`
		ShoppingCartIDs []string `json:"shopping_cart_ids" validate:"required"`
		PaymentType     string   `json:"payment_type" validate:"required"`
	}

	Response struct {
		ID                 string    `json:"id"`
		CustomerID         string    `json:"customer_id"`
		Subtotal           float64   `json:"subtotal"`
		Discount           float64   `json:"discount"`
		Total              float64   `json:"total"`
		CreatedAt          time.Time `json:"created_at"`
		TransactionItem    []Item    `json:"transaction_item"`
		TransactionPayment *Payment  `json:"transaction_payment"`
	}

	Item struct {
		ID            string    `json:"id"`
		TransactionID string    `json:"transaction_id"`
		ProductID     string    `json:"product_id"`
		Notes         string    `json:"notes"`
		Price         float64   `json:"price"`
		Qty           float64   `json:"qty"`
		Total         float64   `json:"total"`
		CreatedAt     time.Time `json:"created_at"`
		Product       *Product  `json:"product"`
	}

	Payment struct {
		ID            string    `json:"id"`
		TransactionID string    `json:"transaction_id"`
		PaymentType   string    `json:"payment_type"`
		Paid          int       `json:"paid"`
		CreatedAt     time.Time `json:"created_at"`
	}

	Product struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Price     float64   `json:"price"`
		CreatedAt time.Time `json:"created_at"`
	}
)

func (c *CreateRequest) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}

func (d DTO) CreateTransformIn(ctx *fiber.Ctx) (*domain.TransactionCreateRequest, error) {
	c := new(CreateRequest)

	if err := ctx.BodyParser(c); err != nil {
		return nil, err
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	result := new(domain.TransactionCreateRequest)
	copier.CopyWithOption(result, c, copier.Option{DeepCopy: true})

	return result, nil
}

func (d DTO) ToCreateResp(r *domain.Transaction) *Response {
	result := new(Response)
	copier.CopyWithOption(result, r, copier.Option{DeepCopy: true})

	return result
}

func (d DTO) ToPaymentResp(r *domain.TransactionPayment) *Payment {
	result := new(Payment)
	copier.CopyWithOption(result, r, copier.Option{DeepCopy: true})

	return result
}
