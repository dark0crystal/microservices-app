package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Product represents a product in our system
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
}

// ProductService handles product operations
type ProductService struct {
	products map[int]*Product
	nextID   int
	mutex    sync.RWMutex
}

// NewProductService creates a new product service
func NewProductService() *ProductService {
	return &ProductService{
		products: make(map[int]*Product),
		nextID:   1,
	}
}

// CreateProduct creates a new product
func (ps *ProductService) CreateProduct(name, description, category string, price float64) *Product {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	product := &Product{
		ID:          ps.nextID,
		Name:        name,
		Description: description,
		Price:       price,
		Category:    category,
		CreatedAt:   time.Now(),
	}

	ps.products[ps.nextID] = product
	ps.nextID++

	return product
}

// GetProduct retrieves a product by ID
func (ps *ProductService) GetProduct(id int) (*Product, bool) {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	product, exists := ps.products[id]
	return product, exists
}

// GetAllProducts retrieves all products
func (ps *ProductService) GetAllProducts() []*Product {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	products := make([]*Product, 0, len(ps.products))
	for _, product := range ps.products {
		products = append(products, product)
	}

	return products
}

// GetProductsByCategory retrieves products by category
func (ps *ProductService) GetProductsByCategory(category string) []*Product {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	var products []*Product
	for _, product := range ps.products {
		if product.Category == category {
			products = append(products, product)
		}
	}

	return products
}

// UpdateProduct updates an existing product
func (ps *ProductService) UpdateProduct(id int, name, description, category string, price float64) (*Product, bool) {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	product, exists := ps.products[id]
	if !exists {
		return nil, false
	}

	product.Name = name
	product.Description = description
	product.Category = category
	product.Price = price

	return product, true
}

// DeleteProduct deletes a product by ID
func (ps *ProductService) DeleteProduct(id int) bool {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	_, exists := ps.products[id]
	if !exists {
		return false
	}

	delete(ps.products, id)
	return true
}

// HTTP handlers
func (ps *ProductService) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		Price       float64 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Category == "" || req.Price <= 0 {
		http.Error(w, "Name, category, and valid price are required", http.StatusBadRequest)
		return
	}

	product := ps.CreateProduct(req.Name, req.Description, req.Category, req.Price)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (ps *ProductService) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	if category != "" {
		// Return products by category
		products := ps.GetProductsByCategory(category)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// Return all products
		products := ps.GetAllProducts()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, exists := ps.GetProduct(id)
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (ps *ProductService) handleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		Price       float64 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Category == "" || req.Price <= 0 {
		http.Error(w, "Name, category, and valid price are required", http.StatusBadRequest)
		return
	}

	product, exists := ps.UpdateProduct(id, req.Name, req.Description, req.Category, req.Price)
	if !exists {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (ps *ProductService) handleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	success := ps.DeleteProduct(id)
	if !success {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	productService := NewProductService()

	// Add some sample data
	productService.CreateProduct("Laptop", "High-performance laptop", "Electronics", 999.99)
	productService.CreateProduct("Coffee Mug", "Ceramic coffee mug", "Kitchen", 12.99)
	productService.CreateProduct("Book", "Programming book", "Books", 29.99)

	// Set up routes
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			productService.handleCreateProduct(w, r)
		case http.MethodGet:
			productService.handleGetProduct(w, r)
		case http.MethodPut:
			productService.handleUpdateProduct(w, r)
		case http.MethodDelete:
			productService.handleDeleteProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Product Service is healthy")
	})

	fmt.Println("Product Service starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
