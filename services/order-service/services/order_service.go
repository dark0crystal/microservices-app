package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"order-service/dto"
	"order-service/models"
	"os"

	"gorm.io/gorm"
)

// OrderService handles order business logic
type OrderService struct {
	db *gorm.DB
}

// NewOrderService creates a new order service
func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

// CreateOrder creates a new order by fetching data from both services
func (s *OrderService) CreateOrder(req dto.CreateOrderRequest) (*dto.OrderWithDetailsResponse, error) {
	// Fetch user data from user service
	user, err := s.fetchUser(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}

	// Fetch product data from product service
	product, err := s.fetchProduct(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %v", err)
	}

	// Create order in database
	order := models.Order{
		UserID:    req.UserID,
		ProductID: req.ProductID,
	}

	if err := s.db.Create(&order).Error; err != nil {
		return nil, err
	}

	// Return order with details
	return &dto.OrderWithDetailsResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		ProductID: order.ProductID,
		User:      user,
		Product:   product,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}, nil
}

// GetOrder retrieves an order with full user and product details
func (s *OrderService) GetOrder(orderID uint) (*dto.OrderWithDetailsResponse, error) {
	var order models.Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// Fetch fresh data from services
	user, err := s.fetchUser(order.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}

	product, err := s.fetchProduct(order.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %v", err)
	}

	return &dto.OrderWithDetailsResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		ProductID: order.ProductID,
		User:      user,
		Product:   product,
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}, nil
}

// GetAllOrders retrieves all orders
func (s *OrderService) GetAllOrders() ([]dto.OrderResponse, error) {
	var orders []models.Order
	if err := s.db.Find(&orders).Error; err != nil {
		return nil, err
	}

	var responses []dto.OrderResponse
	for _, order := range orders {
		responses = append(responses, dto.OrderResponse{
			ID:        order.ID,
			UserID:    order.UserID,
			ProductID: order.ProductID,
			CreatedAt: order.CreatedAt,
			UpdatedAt: order.UpdatedAt,
		})
	}

	return responses, nil
}

// fetchUser fetches user data from user service
func (s *OrderService) fetchUser(userID uint) (*dto.UserResponse, error) {
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://localhost:8080"
	}

	url := fmt.Sprintf("%s/users?id=%d", userServiceURL, userID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user service returned status %d", resp.StatusCode)
	}

	var user dto.UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %v", err)
	}

	return &user, nil
}

// fetchProduct fetches product data from product service
func (s *OrderService) fetchProduct(productID uint) (*dto.ProductResponse, error) {
	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://localhost:8081"
	}

	url := fmt.Sprintf("%s/products?id=%d", productServiceURL, productID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product service returned status %d", resp.StatusCode)
	}

	var product dto.ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode product: %v", err)
	}

	return &product, nil
}
