package dto

import "time"

// CreateOrderRequest represents the request payload for creating an order
type CreateOrderRequest struct {
	UserID    uint `json:"user_id" validate:"required"`
	ProductID uint `json:"product_id" validate:"required"`
}

// OrderResponse represents the response payload for order operations
type OrderResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// OrderWithDetailsResponse represents order with full user and product details
type OrderWithDetailsResponse struct {
	ID        uint             `json:"id"`
	UserID    uint             `json:"user_id"`
	ProductID uint             `json:"product_id"`
	User      *UserResponse    `json:"user,omitempty"`
	Product   *ProductResponse `json:"product,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// UserResponse represents user data from user service
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ProductResponse represents product data from product service
type ProductResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
