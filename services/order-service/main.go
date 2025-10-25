package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// User represents a user from the user service
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Product represents a product from the product service
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

// Order represents an order combining user and product data
type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	User      *User     `json:"user,omitempty"`
	Product   *Product  `json:"product,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// OrderService handles order operations and inter-service communication
type OrderService struct {
	userServiceURL    string
	productServiceURL string
	orders            map[int]*Order
	nextID            int
}

// NewOrderService creates a new order service
func NewOrderService(userServiceURL, productServiceURL string) *OrderService {
	return &OrderService{
		userServiceURL:    userServiceURL,
		productServiceURL: productServiceURL,
		orders:            make(map[int]*Order),
		nextID:            1,
	}
}

// fetchUser fetches user data from user service
func (os *OrderService) fetchUser(userID int) (*User, error) {
	url := fmt.Sprintf("%s/users?id=%d", os.userServiceURL, userID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user service returned status %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %v", err)
	}

	return &user, nil
}

// fetchProduct fetches product data from product service
func (os *OrderService) fetchProduct(productID int) (*Product, error) {
	url := fmt.Sprintf("%s/products?id=%d", os.productServiceURL, productID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product service returned status %d", resp.StatusCode)
	}

	var product Product
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, fmt.Errorf("failed to decode product: %v", err)
	}

	return &product, nil
}

// CreateOrder creates a new order by fetching data from both services
func (os *OrderService) CreateOrder(userID, productID int) (*Order, error) {
	// Fetch user data from user service
	user, err := os.fetchUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}

	// Fetch product data from product service
	product, err := os.fetchProduct(productID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %v", err)
	}

	// Create order
	order := &Order{
		ID:        os.nextID,
		UserID:    userID,
		ProductID: productID,
		User:      user,
		Product:   product,
		CreatedAt: time.Now(),
	}

	os.orders[os.nextID] = order
	os.nextID++

	return order, nil
}

// GetOrder retrieves an order with full user and product details
func (os *OrderService) GetOrder(orderID int) (*Order, error) {
	order, exists := os.orders[orderID]
	if !exists {
		return nil, fmt.Errorf("order not found")
	}

	// Fetch fresh data from services
	user, err := os.fetchUser(order.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}

	product, err := os.fetchProduct(order.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %v", err)
	}

	order.User = user
	order.Product = product

	return order, nil
}

// GetAllOrders retrieves all orders
func (os *OrderService) GetAllOrders() []*Order {
	orders := make([]*Order, 0, len(os.orders))
	for _, order := range os.orders {
		orders = append(orders, order)
	}
	return orders
}

// HTTP handlers
func (os *OrderService) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID    int `json:"user_id"`
		ProductID int `json:"product_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.UserID <= 0 || req.ProductID <= 0 {
		http.Error(w, "Valid user_id and product_id are required", http.StatusBadRequest)
		return
	}

	order, err := os.CreateOrder(req.UserID, req.ProductID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (os *OrderService) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderIDStr := r.URL.Query().Get("id")
	if orderIDStr == "" {
		// Return all orders
		orders := os.GetAllOrders()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orders)
		return
	}

	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := os.GetOrder(orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func main() {
	// Get service URLs from environment variables or use defaults
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://localhost:8080"
	}

	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://localhost:8081"
	}

	// Initialize order service with service URLs
	orderService := NewOrderService(userServiceURL, productServiceURL)

	// Set up routes
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			orderService.handleCreateOrder(w, r)
		case http.MethodGet:
			orderService.handleGetOrder(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Order Service is healthy")
	})

	fmt.Println("Order Service starting on port 8082...")
	fmt.Println("Make sure User Service (port 8080) and Product Service (port 8081) are running!")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
