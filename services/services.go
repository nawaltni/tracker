package services

import (
	"github.com/nawaltni/tracker/config"
	"github.com/nawaltni/tracker/domain"
	"github.com/nawaltni/tracker/postgres"
)

// Services contains all services of the application
type Services struct {
	UserPositionService domain.UserPositionService
	TimeTrackingService domain.TimeTrackingService

	Config config.Configuration
}

// NewServices creates a new Services instance
func NewServices(
	config config.Configuration, repos *postgres.Repositories,
	placesClient domain.PlacesClientGRPC, bigqueryClient domain.BigqueryClient,
) (*Services, error) {
	userPositionService, err := NewUserPositionService(repos.UserPosition, placesClient, bigqueryClient)
	if err != nil {
		return nil, err
	}

	timeTrackingService := NewTimeTrackingService(repos.TimeTracking)
	return &Services{
		Config:              config,
		UserPositionService: userPositionService,
		TimeTrackingService: timeTrackingService,
	}, nil
}
