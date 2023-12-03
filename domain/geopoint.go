package domain

import "fmt"

// GeoPoint represents a geographical coordinate with latitude and longitude.
type GeoPoint struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

// String returns the string representation of a GeoPoint.
func (g GeoPoint) String() string {
	return fmt.Sprintf("%f,%f", g.Latitude, g.Longitude)
}
