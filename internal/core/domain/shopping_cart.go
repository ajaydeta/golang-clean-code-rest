package domain

import "time"

type ShoppingCart struct {
	ID         string
	CustomerID string
	ProductID  string
	Notes      string
	Qty        float64
	CreatedAt  time.Time
}
