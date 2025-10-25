#!/bin/bash

# Microservices Test Script
# This script demonstrates the functionality of our microservices app

echo "🚀 Testing Microservices App"
echo "=============================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to check if service is running
check_service() {
    local service_name=$1
    local port=$2
    local url="http://localhost:$port/health"
    
    echo -n "Checking $service_name on port $port... "
    
    if curl -s "$url" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ Running${NC}"
        return 0
    else
        echo -e "${RED}✗ Not running${NC}"
        return 1
    fi
}

# Check if all services are running
echo "🔍 Checking services..."
check_service "User Service" 8080
check_service "Product Service" 8081
check_service "Order Service" 8082

echo ""
echo "📋 Testing API Endpoints"
echo "========================"

# Test User Service
echo -e "${YELLOW}1. Testing User Service${NC}"
echo "Creating a user..."
USER_RESPONSE=$(curl -s -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Test User", "email": "test@example.com"}')
echo "Response: $USER_RESPONSE"

echo "Getting all users..."
curl -s http://localhost:8080/users | jq '.' 2>/dev/null || echo "Response received (install jq for pretty formatting)"

echo ""

# Test Product Service
echo -e "${YELLOW}2. Testing Product Service${NC}"
echo "Creating a product..."
PRODUCT_RESPONSE=$(curl -s -X POST http://localhost:8081/products \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Product", "description": "A test product", "category": "Test", "price": 99.99}')
echo "Response: $PRODUCT_RESPONSE"

echo "Getting all products..."
curl -s http://localhost:8081/products | jq '.' 2>/dev/null || echo "Response received (install jq for pretty formatting)"

echo ""

# Test Order Service (Inter-service communication)
echo -e "${YELLOW}3. Testing Order Service (Inter-service communication)${NC}"
echo "Creating an order..."
ORDER_RESPONSE=$(curl -s -X POST http://localhost:8082/orders \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "product_id": 1}')
echo "Response: $ORDER_RESPONSE"

echo "Getting order with full details..."
curl -s http://localhost:8082/orders?id=1 | jq '.' 2>/dev/null || echo "Response received (install jq for pretty formatting)"

echo ""
echo -e "${GREEN}✅ All tests completed!${NC}"
echo ""
echo "💡 Tips:"
echo "- Install 'jq' for pretty JSON formatting: brew install jq"
echo "- Use 'docker-compose logs -f' to see service logs"
echo "- Check individual service health: curl http://localhost:8080/health"
