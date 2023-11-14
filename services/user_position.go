package services

import (
	"context"
	"fmt"
	"time"

	"github.com/nawaltni/tracker/domain"
)

// UserPositionService is the service for UserPosition
type UserPositionService struct {
	repo                  domain.UserPositionRepository
	placesClientGRPC      domain.PlacesClientGRPC
	trackerClientBigquery domain.TrackerClientBigquery
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

	// Call places grpc service to see if the location matches any places
	res, err := s.placesClientGRPC.IsWithinPlace(ctx, location.Latitude, location.Longitude)
	if err != nil {
		return err
	}

	// If the location matches a place, we will add the place to the userPosition object
	if res.IsWithin {
		userPosition.PlaceID = &res.PlaceID
		userPosition.PlaceName = &res.PlaceName
	}

	// Now we will call the GetUserPosition to know the previous position of the user.
	knownPosition, err := s.GetUserPosition(ctx, userID)
	if err != nil && err != domain.ErrNotFound {
		return fmt.Errorf("error getting user position: %w", err)
	}

	// If the user has no previous position, we will set the knownPosition to the street
	if knownPosition == nil {
		sp := StreetPosition()
		knownPosition = &sp
	}

	// Call the CalculateUserPosition function to calculate the user's position
	// and determine if the user is entering or leaving a place
	s.CalculateUserPosition(&userPosition, *knownPosition)

	// Call the tracker bigquery service to record the user's position
	err = s.trackerClientBigquery.RecordUserPosition(ctx, userPosition)
	if err != nil {
		return fmt.Errorf("error recording position in BigQuery: %w", err)
	}

	return s.repo.Insert(ctx, &userPosition)
}

// StreetPosition returns a UserPosition object for the street
func StreetPosition() domain.UserPosition {
	id := "0000-0000-0000-0000"
	name := "Street"
	return domain.UserPosition{
		PlaceID:   &id,
		PlaceName: &name,
	}
}

func (s *UserPositionService) CalculateUserPosition(userPosition *domain.UserPosition, knownPosition domain.UserPosition) {
	isKnownPlaceEmpty := knownPosition.PlaceID == nil || *knownPosition.PlaceID == "" || *knownPosition.PlaceID == "0000-0000-0000-0000"
	isUserPlaceEmpty := userPosition.PlaceID == nil || *userPosition.PlaceID == "" || *userPosition.PlaceID == "0000-0000-0000-0000"

	switch {
	case isKnownPlaceEmpty && !isUserPlaceEmpty:
		// User is entering a place
		userPosition.CheckedIn = &userPosition.CreatedAt

	case !isKnownPlaceEmpty && isUserPlaceEmpty:
		// User is leaving a place
		userPosition.CheckedOut = &userPosition.CreatedAt
	}
}

// Helper function to return a pointer to the time value
func timePointer(t time.Time) *time.Time {
	return &t
}

func stringPointer(s string) *string {
	return &s
}

// GetUserPosition retrieves the current position of a user
func (s *UserPositionService) GetUserPosition(ctx context.Context, userID string) (*domain.UserPosition, error) {
	return s.repo.GetUserPosition(ctx, userID)
}

// GetUsersPositionByCoordinates retrieves a list of users' positions close to the given coordinates.
func (s *UserPositionService) GetUsersPositionByCoordinates(ctx context.Context, lat, long float32, distance int) ([]domain.UserPosition, error) {
	return s.repo.GetUsersPositionByCoordinates(ctx, lat, long, distance)
}
