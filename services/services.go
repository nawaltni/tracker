package services

import (
	"github.com/nawaltni/tracker/config"
	"github.com/nawaltni/tracker/domain"
	"github.com/nawaltni/tracker/postgres"
)

// Services contains all services of the application
type Services struct {
	UserPositionService domain.UserPositionService
	Config              config.Config
}

// NewServices creates a new Services instance
func NewServices(
	config config.Config, repos *postgres.Repositories,
	placesClient domain.PlacesClientGRPC,
	authClient domain.AuthClientGRPC,
	bigqueryClient domain.BigqueryClient,
) (*Services, error) {
	userPositionService, err := NewUserPositionService(repos.UserPosition, placesClient, authClient, bigqueryClient)
	if err != nil {
		return nil, err
	}
	return &Services{
		Config:              config,
		UserPositionService: userPositionService,
	}, nil
}
