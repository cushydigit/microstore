package models

import "time"

type OrderItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type Order struct {
	ID         int64       `json:"id"`
	UserID     int         `json:"user_id"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price"`
	Status     string      `json:"status"` // e.g pending, confirmed, shipped
	Created_at time.Time `json:"created_at"`
}
