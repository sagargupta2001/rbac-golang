package domain

import "time"

// User represents a user in the system
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Don't expose this
	CreatedAt    time.Time `json:"created_at"`
}

// Role represents a user role
type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Permission represents an action a role can perform
type Permission struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Product represents a resource to be protected
type Product struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	CreatedByUserID int64     `json:"created_by_user"`
	CreatedAt     time.Time `json:"created_at"`
}

// --- Request/Response Payloads ---

// RegisterRequest is the payload for user registration
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest is the payload for user login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse is the payload for a successful login
type LoginResponse struct {
	Token string `json:"token"`
}

// CreateProductRequest is the payload for creating a product
type CreateProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// --- GraphQL Client Structs ---

// GraphQLRequest represents the JSON body we send in a POST request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// GraphQLCountryResponse is the struct we will unmarshal the JSON response into.
// The struct tags (e.g., `json:"data"`) MUST match the JSON response.
type GraphQLCountryResponse struct {
	Data struct {
		Country struct {
			Name    string `json:"name"`
			Capital string `json:"capital"`
			Emoji   string `json:"emoji"`
			Currency  string            `json:"currency"`
		} `json:"country"`
	} `json:"data"`
}

// CountryDetails is a cleaner struct we'll return from our API
type CountryDetails struct {
	Name    string `json:"name"`
	Capital string `json:"capital"`
	Emoji   string `json:"emoji"`
	Currency  string            `json:"currency"`
}