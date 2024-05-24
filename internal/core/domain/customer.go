package domain

import "time"

type (
	Customer struct {
		ID        string
		Name      string
		Email     string
		Password  string
		CreatedAt time.Time
	}
)
