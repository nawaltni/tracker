package services

import (
	"context"
	"time"

	"github.com/nawaltni/tracker/domain"
)

// UserPositionService is the service for UserPosition
type UserPositionService struct {
	repo domain.UserPositionRepository
}

// NewUserPositionService creates a new UserPositionService
func NewUserPositionService(repo domain.UserPositionRepository) (*UserPositionService, error) {
	return &UserPositionService{repo: repo}, nil
}

// RecordPosition handles the logic for recording a user's position
func (s *UserPositionService) RecordPosition(ctx context.Context, userID string, location domain.GeoPoint, timestamp time.Time, clientID string, metadata domain.PhoneMetadata) error {
	// Create a UserPosition domain object
	userPosition := domain.UserPosition{
		UserID:    userID,
		Location:  location,
		CreatedAt: timestamp,
		Metadata:  metadata,
	}
	// Here is where we call the places service to see if the location matches any places.
	// We will also call the GetUserPosition to know the previous position of the user.
	// With this 2 calls we can know if the user is entering or leaving a place.

	return s.repo.Insert(ctx, &userPosition)
}

// GetUserPosition retrieves the current position of a user
func (s *UserPositionService) GetUserPosition(ctx context.Context, userID string) (*domain.UserPosition, error) {
	return s.repo.GetUserPosition(ctx, userID)
}

// GetUsersPositionByCoordinates retrieves a list of users' positions close to the given coordinates.
func (s *UserPositionService) GetUsersPositionByCoordinates(ctx context.Context, lat, long float32, distance int) ([]domain.UserPosition, error) {
	return s.repo.GetUsersPositionByCoordinates(ctx, lat, long, distance)
}
