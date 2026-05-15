package model

import "time"

type ProductReview struct {
	ID         int64
	ProductID  string
	CustomerID string
	OrderID    string
	Rating     int
	Comment    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
