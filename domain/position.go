package domain

import (
	"context"
	"time"
)

// UserPosition represents the geographical position of a user at a given time.
// If the user is within a known place, the PlaceID will be populated with the corresponding identifier.
type UserPosition struct {
	UserID        string     `json:"user_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	Reference     string     `json:"uuid"`
	Location      GeoPoint   `json:"location"`
	PlaceID       *string    `json:"place_id,omitempty"`    // Optional, associated when within a Place
	PlaceName     *string    `json:"place_name,omitempty"`  // Optional, name of the Place if within one
	CheckedIn     *time.Time `json:"checked_in,omitempty"`  // Optional, time when user checked into a Place
	CheckedOut    *time.Time `json:"checked_out,omitempty"` // Optional, time when user checked out of a Place
	PhoneMeta     PhoneMeta  `json:"phone_meta"`            // Metadata about the phone reporting the position
	BackendUserID string     `json:"backend_user_id,omitempty"`
	Role          string     `json:"role,omitempty"`
}

// CurrentUserPosition represents the current position of a user.
// It has a dedicated uuid
type CurrentUserPosition struct {
	UUID string
	UserPosition
}

// UserPositionRepository defines the interface for the user position storage.
type UserPositionRepository interface {
	Insert(ctx context.Context, userPosition *UserPosition) error
	GetUserCurrentPosition(ctx context.Context, userID string) (*UserPosition, error)
	GetUserCurrentPositionByBackendID(ctx context.Context, id string) (*UserPosition, error)
	GetUsersCurrentPositionByCoordinates(ctx context.Context, lat float32, lon float32, distance int) ([]UserPosition, error)
	// GetUsersCurrentPositionByDate retrieves a list of users' positions for a given date.
	GetUsersCurrentPositionByDate(ctx context.Context, date time.Time) ([]UserPosition, error)
	// GetUsersCurrentPositionSince retrieves a list of users' positions since a given time.
	GetUsersCurrentPositionSince(ctx context.Context, t time.Time) ([]UserPosition, error)
}

// UserPositionService defines the interface for the user position service.
type UserPositionService interface {
	RecordPosition(ctx context.Context, position UserPosition) error
	GetUserCurrentPosition(ctx context.Context, userID string) (*UserPosition, error)
	GetUsersCurrentPositionByCoordinates(ctx context.Context, lat float32, lon float32, distance int) ([]UserPosition, error)
	GetUsersCurrentPositionByDate(ctx context.Context, date time.Time) ([]UserPosition, error)
	GetUsersCurrentPositionsSince(ctx context.Context, t time.Time) ([]UserPosition, error)
	GetUserPositionHistorySince(ctx context.Context, userID string, t time.Time, limit int) ([]UserPosition, error)
}
