package domain

import (
	"context"
	"time"
)

// User is the representation of a user in the domain
type User struct {
	ID         string    `json:"id"`
	BackendID  int       `json:"backend_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	FirebaseID string    `json:"firebase_id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Address    string    `json:"address"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Password   string    `json:"password"`
	Role       string    `json:"role"`
}

// AuthClientGRPC is the interface that defines the methods to interact with the auth gRPC service
type AuthClientGRPC interface {
	GetUserByBackendID(ctx context.Context, id int) (*User, error)
}
