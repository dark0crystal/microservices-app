package main

import (
	"fmt"
	"log"
	"net/http"
	"product-service/database"
	"product-service/handlers"
	"product-service/services"
)

func main() {
	// Connect to database
	database.ConnectDB()
	database.MigrateDB()

	// Initialize services
	productService := services.NewProductService(database.DB)
	productHandler := handlers.NewProductHandler(productService)

	// Set up routes
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			productHandler.CreateProduct(w, r)
		case http.MethodGet:
			productHandler.GetProduct(w, r)
		case http.MethodPut:
			productHandler.UpdateProduct(w, r)
		case http.MethodDelete:
			productHandler.DeleteProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", productHandler.Health)

	fmt.Println("Product Service starting on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
