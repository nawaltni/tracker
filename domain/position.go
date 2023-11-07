package domain

import "time"

// UserPosition represents the geographical position of a user at a given time.
// If the user is within a known place, the PlaceID will be populated with the corresponding identifier.
type UserPosition struct {
	UserID       string        `json:"user_id"`
	Location     GeoPoint      `json:"location"`
	CreatedAt    time.Time     `json:"created_at"`
	PlaceID      *string       `json:"place_id,omitempty"`       // Optional, associated when within a Place
	PlaceName    *string       `json:"place_name,omitempty"`     // Optional, name of the Place if within one
	CheckedInAt  *time.Time    `json:"checked_in_at,omitempty"`  // Optional, time when user checked into a Place
	CheckedOutAt *time.Time    `json:"checked_out_at,omitempty"` // Optional, time when user checked out of a Place
	Metadata     PhoneMetadata `json:"metadata"`                 // Metadata about the phone reporting the position
}

// UserPositionRepository defines the interface for the user position storage.
type UserPositionRepository interface {
	Insert(userPosition *UserPosition) error
	GetUserPosition(userID string) (*UserPosition, error)
	GetUsersPositionByCoordinates(lat float64, lon float64) ([]UserPosition, error)
}
