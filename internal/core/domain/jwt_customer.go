package domain

type JWTCustomer struct {
	ID        string `json:"customer_id"`
	Subject   string `json:"sub"`
	ExpiresAt int64  `json:"exp"`
}
