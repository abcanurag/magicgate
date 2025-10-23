package database

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Don't expose password hash in JSON
	CreatedAt    time.Time `json:"created_at"`
}

// Key represents a cryptographic key associated with a user
type Key struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	KeyMaterial []byte    `json:"-"` // Don't expose raw key material in JSON
	CreatedAt   time.Time `json:"created_at"`
}

// KeyResponse is used for API responses to avoid exposing raw key material
type KeyResponse struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Secret represents a secret associated with a user and a key
type Secret struct {
	ID        int       `json:"id"`
	KeyID     int       `json:"key_id"`
	Data      []byte    `json:"-"` // Don't expose raw secret data in JSON
	CreatedAt time.Time `json:"created_at"`
}