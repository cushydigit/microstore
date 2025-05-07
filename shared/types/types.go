package types

import "time"

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// general response
type response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

type Response struct {
	response
	Data any `json:"data,omitempty"`
}

type ProductResponse struct {
	response
	Data Product `json:"data"`
}

type ProductsReponse struct {
	response
	Data []Product `json:"data"`
}

type OrderResponse struct {
	response
	Data Order `json:"data"`
}

type OrdersResponse struct {
	response
	Data []Order `json:"data"`
}

// models types
type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

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
	Created_at time.Time   `json:"created_at"`
}
