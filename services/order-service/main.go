package main

import (
	"fmt"
	"log"
	"net/http"
	"order-service/database"
	"order-service/handlers"
	"order-service/services"
)

func main() {
	// Connect to database
	database.ConnectDB()
	database.MigrateDB()

	// Initialize services
	orderService := services.NewOrderService(database.DB)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Set up routes
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			orderHandler.CreateOrder(w, r)
		case http.MethodGet:
			orderHandler.GetOrder(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", orderHandler.Health)

	fmt.Println("Order Service starting on port 8082...")
	fmt.Println("Make sure User Service (port 8080) and Product Service (port 8081) are running!")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
