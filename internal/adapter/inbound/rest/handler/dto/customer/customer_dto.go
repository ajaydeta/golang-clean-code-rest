package customer

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"synapsis-challenge/internal/core/domain"
)

type (
	RegisterCustomerReq struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	RegisterCustomerResp struct {
		ID string `json:"id"`
	}

	DTO struct{}
)

func (c *RegisterCustomerReq) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}

func (d DTO) CreateTransformIn(ctx *fiber.Ctx) (*domain.Customer, error) {
	c := new(RegisterCustomerReq)

	if err := ctx.BodyParser(c); err != nil {
		return nil, err
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return &domain.Customer{
		Name:     cases.Title(language.Indonesian, cases.NoLower).String(strings.ToLower(c.Name)),
		Email:    c.Email,
		Password: c.Password,
	}, nil

}
