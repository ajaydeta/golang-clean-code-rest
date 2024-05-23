package domain

import "time"

type (
	Product struct {
		ID              string
		Name            string
		Price           float64
		CreatedAt       time.Time
		ProductCategory []ProductCategory
	}

	ProductCategory struct {
		ProductId  string
		CategoryId string
		Category   *Category
	}

	Category struct {
		ID        string
		Name      string
		CreatedAt time.Time
	}

	ProductFilter struct {
		Filter
		CategoryID []string
	}
)
