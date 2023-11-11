package domain

import (
	"context"
	"time"
)

// UserPosition represents the geographical position of a user at a given time.
// If the user is within a known place, the PlaceID will be populated with the corresponding identifier.
type UserPosition struct {
	UserID     string        `json:"user_id"`
	Location   GeoPoint      `json:"location"`
	CreatedAt  time.Time     `json:"created_at"`
	PlaceID    *string       `json:"place_id,omitempty"`    // Optional, associated when within a Place
	PlaceName  *string       `json:"place_name,omitempty"`  // Optional, name of the Place if within one
	CheckedIn  *time.Time    `json:"checked_in,omitempty"`  // Optional, time when user checked into a Place
	CheckedOut *time.Time    `json:"checked_out,omitempty"` // Optional, time when user checked out of a Place
	Metadata   PhoneMetadata `json:"metadata"`              // Metadata about the phone reporting the position
}

// UserPositionRepository defines the interface for the user position storage.
type UserPositionRepository interface {
	Insert(userPosition *UserPosition) error
	GetUserPosition(userID string) (*UserPosition, error)
	GetUsersPositionByCoordinates(lat float64, lon float64, distance int) ([]UserPosition, error)
}

// UserPositionService defines the interface for the user position service.
type UserPositionService interface {
	RecordPosition(ctx context.Context, userID string, location GeoPoint, timestamp time.Time, clientID string, metadata PhoneMetadata) error
	GetUserPosition(ctx context.Context, userID string) (*UserPosition, error)
	GetUserPositionList(ctx context.Context, userID string, clientID string, startTime, endTime time.Time) ([]UserPosition, error)
}
