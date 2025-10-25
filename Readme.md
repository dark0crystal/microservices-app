# Microservices App with Go

A simple microservices application built with Go without frameworks, demonstrating inter-service communication and basic CRUD operations.

## Architecture

This application consists of three microservices:

1. **User Service** (Port 8080) - Manages user data
2. **Product Service** (Port 8081) - Manages product data  
3. **Order Service** (Port 8082) - Demonstrates inter-service communication by combining user and product data

## Features

- Pure Go implementation without external frameworks
- RESTful API endpoints
- In-memory data storage
- Inter-service communication
- Docker containerization
- Health check endpoints
- Concurrent-safe operations with mutex

## API Endpoints

### User Service (Port 8080)

- `GET /users` - Get all users
- `GET /users?id={id}` - Get user by ID
- `POST /users` - Create a new user
- `PUT /users?id={id}` - Update user
- `DELETE /users?id={id}` - Delete user
- `GET /health` - Health check

### Product Service (Port 8081)

- `GET /products` - Get all products
- `GET /products?id={id}` - Get product by ID
- `GET /products?category={category}` - Get products by category
- `POST /products` - Create a new product
- `PUT /products?id={id}` - Update product
- `DELETE /products?id={id}` - Delete product
- `GET /health` - Health check

### Order Service (Port 8082)

- `GET /orders` - Get all orders
- `GET /orders?id={id}` - Get order by ID (with full user and product details)
- `POST /orders` - Create a new order
- `GET /health` - Health check

## Quick Start

### Option 1: Using Docker Compose (Recommended)

1. Clone the repository
2. Run all services with Docker Compose:
```bash
docker-compose up --build
```

This will start all three services and make them available at:
- User Service: http://localhost:8080
- Product Service: http://localhost:8081
- Order Service: http://localhost:8082

### Option 2: Running Services Individually

1. **Start User Service:**
```bash
cd services/user-service
go run main.go
```

2. **Start Product Service:**
```bash
cd services/product-service
go run main.go
```

3. **Start Order Service:**
```bash
cd services/order-service
go run main.go
```

## Example Usage

### Create a User
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

### Create a Product
```bash
curl -X POST http://localhost:8081/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Laptop", "description": "High-performance laptop", "category": "Electronics", "price": 999.99}'
```

### Create an Order (Inter-service communication)
```bash
curl -X POST http://localhost:8082/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "product_id": 1}'
```

### Get Order with Full Details
```bash
curl http://localhost:8082/orders?id=1
```

## Project Structure

```
microservices-app/
├── services/
│   ├── user-service/
│   │   ├── main.go
│   │   ├── go.mod
│   │   └── Dockerfile
│   ├── product-service/
│   │   ├── main.go
│   │   ├── go.mod
│   │   └── Dockerfile
│   └── order-service/
│       ├── main.go
│       ├── go.mod
│       └── Dockerfile
├── docker-compose.yml
└── README.md
```

## Development

Each service is a standalone Go application with its own `go.mod` file. The services communicate via HTTP REST APIs.

### Key Design Decisions

1. **No External Frameworks**: Uses only Go standard library for HTTP handling
2. **In-Memory Storage**: Simple map-based storage for demonstration (can be replaced with databases)
3. **Concurrent Safety**: Uses `sync.RWMutex` for thread-safe operations
4. **Environment Configuration**: Order service uses environment variables for service URLs
5. **Health Checks**: Each service provides a health check endpoint

## Next Steps

This is a basic implementation that can be extended with:

- Database integration (PostgreSQL, MongoDB, etc.)
- Service discovery
- API Gateway
- Message queues (RabbitMQ, Kafka)
- Monitoring and logging
- Authentication and authorization
- Circuit breakers
- Load balancing
- Configuration management