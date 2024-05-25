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

	SignInReq struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	RefreshTokenReq struct {
		Token string `json:"token" validate:"required"`
	}

	SignInResp struct {
		ID           string `json:"id"`
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
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

func (c *SignInReq) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(c)
}

func (c *RefreshTokenReq) Validate() error {
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

func (d DTO) SignInTransformIn(ctx *fiber.Ctx) (*domain.Customer, error) {
	c := new(SignInReq)

	if err := ctx.BodyParser(c); err != nil {
		return nil, err
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return &domain.Customer{
		Email:    c.Email,
		Password: c.Password,
	}, nil
}

func (d DTO) RefreshTokenTransformIn(ctx *fiber.Ctx) (string, error) {
	c := new(RefreshTokenReq)

	if err := ctx.BodyParser(c); err != nil {
		return "", err
	}

	if err := c.Validate(); err != nil {
		return "", err
	}

	return c.Token, nil
}

func (d DTO) ToSignInResp(r *domain.SignIn) *SignInResp {
	return &SignInResp{
		ID:           r.Customer.ID,
		RefreshToken: r.RefreshToken,
		AccessToken:  r.AccessToken,
	}
}
