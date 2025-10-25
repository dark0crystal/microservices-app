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

// User represents a user in our system
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserService handles user operations
type UserService struct {
	users  map[int]*User
	nextID int
	mutex  sync.RWMutex
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		users:  make(map[int]*User),
		nextID: 1,
	}
}

// CreateUser creates a new user
func (us *UserService) CreateUser(name, email string) *User {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	user := &User{
		ID:        us.nextID,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	us.users[us.nextID] = user
	us.nextID++

	return user
}

// GetUser retrieves a user by ID
func (us *UserService) GetUser(id int) (*User, bool) {
	us.mutex.RLock()
	defer us.mutex.RUnlock()

	user, exists := us.users[id]
	return user, exists
}

// GetAllUsers retrieves all users
func (us *UserService) GetAllUsers() []*User {
	us.mutex.RLock()
	defer us.mutex.RUnlock()

	users := make([]*User, 0, len(us.users))
	for _, user := range us.users {
		users = append(users, user)
	}

	return users
}

// UpdateUser updates an existing user
func (us *UserService) UpdateUser(id int, name, email string) (*User, bool) {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	user, exists := us.users[id]
	if !exists {
		return nil, false
	}

	user.Name = name
	user.Email = email

	return user, true
}

// DeleteUser deletes a user by ID
func (us *UserService) DeleteUser(id int) bool {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	_, exists := us.users[id]
	if !exists {
		return false
	}

	delete(us.users, id)
	return true
}

// HTTP handlers
func (us *UserService) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	user := us.CreateUser(req.Name, req.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (us *UserService) handleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		// Return all users
		users := us.GetAllUsers()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, exists := us.GetUser(id)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (us *UserService) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	user, exists := us.UpdateUser(id, req.Name, req.Email)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (us *UserService) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	success := us.DeleteUser(id)
	if !success {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func main() {
	userService := NewUserService()

	// Add some sample data
	userService.CreateUser("John Doe", "john@example.com")
	userService.CreateUser("Jane Smith", "jane@example.com")

	// Set up routes
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			userService.handleCreateUser(w, r)
		case http.MethodGet:
			userService.handleGetUser(w, r)
		case http.MethodPut:
			userService.handleUpdateUser(w, r)
		case http.MethodDelete:
			userService.handleDeleteUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "User Service is healthy")
	})

	fmt.Println("User Service starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
