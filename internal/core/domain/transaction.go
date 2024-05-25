package domain

import "time"

type (
	Transaction struct {
		ID                 string
		CustomerID         string
		Subtotal           float64
		Discount           float64
		Total              float64
		CreatedAt          time.Time
		TransactionItem    []TransactionItem
		TransactionPayment *TransactionPayment
	}

	TransactionItem struct {
		ID            string
		TransactionID string
		ProductID     string
		Notes         string
		Price         float64
		Qty           float64
		Total         float64
		CreatedAt     time.Time
		Product       *Product
	}

	TransactionPayment struct {
		ID            string
		TransactionID string
		PaymentType   string
		Paid          int
		CreatedAt     time.Time
	}

	TransactionCreateRequest struct {
		Discount        float64
		ShoppingCartIDs []string
		PaymentType     string
	}
)
