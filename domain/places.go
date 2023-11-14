package domain

import "context"

type WithinPlaceResponse struct {
	IsWithin  bool
	PlaceID   string
	PlaceName string
}

type PlacesClientGRPC interface {
	IsWithinPlace(ctx context.Context, lat float32, lon float32) (*WithinPlaceResponse, error)
}
