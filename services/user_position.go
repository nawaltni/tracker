package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nawaltni/tracker/domain"
)

// UserPositionService is the service for UserPosition
type UserPositionService struct {
	repo             domain.UserPositionRepository
	placesClientGRPC domain.PlacesClientGRPC
	bigqueryClient   domain.BigqueryClient
}

// NewUserPositionService creates a new UserPositionService
func NewUserPositionService(
	repo domain.UserPositionRepository,
	placesClient domain.PlacesClientGRPC,
	bigqueryClient domain.BigqueryClient,
) (*UserPositionService, error) {
	return &UserPositionService{
		repo:             repo,
		placesClientGRPC: placesClient,
		bigqueryClient:   bigqueryClient,
	}, nil
}

// RecordPosition handles the logic for recording a user's position
func (s *UserPositionService) RecordPosition(ctx context.Context, userPosition domain.UserPosition) error {
	// Call places grpc service to see if the location matches any places
	res, err := s.placesClientGRPC.IsWithinPlace(ctx, userPosition.Location.Latitude, userPosition.Location.Longitude)
	if err != nil {
		return err
	}

	// If the location matches a place, we will add the place to the userPosition object
	if res.IsWithin {
		userPosition.PlaceID = &res.PlaceID
		userPosition.PlaceName = &res.PlaceName
	}

	var knownPosition *domain.UserPosition

	if !IsValidUUID(userPosition.UserID) {
		knownPosition, err = s.GetUserCurrentPositionByBackendID(ctx, userPosition.UserID)
		if err != nil && err != domain.ErrNotFound {
			return fmt.Errorf("error getting user position: %w", err)
		}

		fmt.Printf("knownPosition: %+v\n", knownPosition)

		userPosition.BackendUserID = userPosition.UserID
		if knownPosition == nil {
			code, err := uuid.NewV7()
			if err != nil {
				return fmt.Errorf("error generating uuid: %w", err)
			}
			userPosition.UserID = code.String()
		} else {
			userPosition.UserID = knownPosition.UserID
		}
	} else {
		// Now we will call the GetCurrentUserPosition to know the previous position of the user.
		knownPosition, err = s.GetUserCurrentPosition(ctx, userPosition.UserID)
		if err != nil && err != domain.ErrNotFound {
			return fmt.Errorf("error getting user position: %w", err)
		}

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
	err = s.bigqueryClient.RecordUserPosition(ctx, userPosition)
	if err != nil {
		return fmt.Errorf("error recording position in BigQuery: %w", err)
	}

	return s.repo.Insert(ctx, &userPosition)
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
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

// GetUserCurrentPosition retrieves the current position of a user.
func (s *UserPositionService) GetUserCurrentPosition(ctx context.Context, userID string) (*domain.UserPosition, error) {
	return s.repo.GetUserCurrentPosition(ctx, userID)
}

// GetUserCurrentPositionByBackendID retrieves the current position of a user by reference.
func (s *UserPositionService) GetUserCurrentPositionByBackendID(ctx context.Context, reference string) (*domain.UserPosition, error) {
	return s.repo.GetUserCurrentPositionByBackendID(ctx, reference)
}

// GetUsersCurrentPositionByCoordinates retrieves the current position of all users within a given distance from a set of coordinates.
func (s *UserPositionService) GetUsersCurrentPositionByCoordinates(ctx context.Context, lat, long float32, distance int) ([]domain.UserPosition, error) {
	return s.repo.GetUsersCurrentPositionByCoordinates(ctx, lat, long, distance)
}

// GetUsersCurrentPositionByDate retrieves a list of users' positions for a given date.
func (s *UserPositionService) GetUsersCurrentPositionByDate(ctx context.Context, date time.Time) ([]domain.UserPosition, error) {
	return s.repo.GetUsersCurrentPositionByDate(ctx, date)
}

// GetUsersCurrentPositionsSince retrieves a list of users' positions since a given time.
func (s *UserPositionService) GetUsersCurrentPositionsSince(ctx context.Context, t time.Time) ([]domain.UserPosition, error) {
	return s.repo.GetUsersCurrentPositionSince(ctx, t)
}

// // GetUserPositionsSince retrieves a user's positions since a given time.
// func (s *UserPositionService) GetUserPositionsSince(ctx context.Context, userID string, t time.Time) ([]domain.UserPosition, error) {
// 	return s.bigqueryClient.GetUserPositionsSince(ctx, userID, t)
// }
