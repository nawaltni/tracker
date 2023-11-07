package domain

// GeoPoint represents a geographical coordinate with latitude and longitude.
type GeoPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}