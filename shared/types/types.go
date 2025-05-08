package types

import "time"

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateOrderRequest struct {
	Items []OrderItem `json:"items"`
}

// general response
type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ProductResponse struct {
	Error   bool    `json:"error"`
	Message string  `json:"message"`
	Data    Product `json:"data"`
}

type ProductsResponse struct {
	Error   bool      `json:"error"`
	Message string    `json:"message"`
	Data    []Product `json:"data"`
}

type OrderResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    Order  `json:"data"`
}

type OrdersResponse struct {
	Error   bool    `json:"error"`
	Message string  `json:"message"`
	Data    []Order `json:"data"`
}

// models types
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

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
