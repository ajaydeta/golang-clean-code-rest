package product

import (
	"github.com/gofiber/fiber/v2"
	"synapsis-challenge/internal/adapter/inbound/rest/handler/dto"
	"synapsis-challenge/internal/core/domain"
	"synapsis-challenge/shared"
	"time"
)

type (
	Response struct {
		ID        string     `json:"id"`
		Name      string     `json:"name"`
		Price     float64    `json:"price"`
		CreatedAt time.Time  `json:"created_at"`
		Category  []Category `json:"category"`
	}

	Category struct {
		ID        string
		Name      string
		CreatedAt time.Time
	}

	ResponseList struct {
		Response   []Response     `json:"products"`
		Pagination dto.Pagination `json:"pagination"`
	}

	DTO struct{}
)

func (d *DTO) ToResponseList(products []domain.Product, filter domain.Filter) *ResponseList {
	resp := new(ResponseList)
	resp.Response = make([]Response, len(products))

	for i, v := range products {
		resp.Response[i] = d.ToResponse(&v)
	}

	page, perPage := shared.GetPageAndPerPage(filter.Limit, filter.Offset)
	resp.Pagination = dto.Pagination{
		Page:    page,
		PerPage: perPage,
	}
	return resp
}

func (d *DTO) ToResponse(v *domain.Product) Response {
	rTemp := Response{
		ID:        v.ID,
		Name:      v.Name,
		Price:     v.Price,
		CreatedAt: v.CreatedAt,
	}

	if len(v.ProductCategory) > 0 {
		rTemp.Category = make([]Category, len(v.ProductCategory))

		for j, c := range v.ProductCategory {
			rTemp.Category[j] = Category{
				ID:        c.Category.ID,
				Name:      c.Category.Name,
				CreatedAt: c.Category.CreatedAt,
			}
		}
	}

	return rTemp
}

func (d *DTO) GetListProductFilter(c *fiber.Ctx) domain.ProductFilter {
	return domain.ProductFilter{
		Filter:     dto.GetBaseFilter(c),
		CategoryID: shared.SplitStringBySeparator(c.Query("categoryId"), ","),
	}
}
