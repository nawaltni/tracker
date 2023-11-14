package domain

// GeoPoint represents a geographical coordinate with latitude and longitude.
type GeoPoint struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}
