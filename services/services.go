package services

import (
	"github.com/nawaltni/tracker/config"
	"github.com/nawaltni/tracker/domain"
	"github.com/nawaltni/tracker/postgres"
)

// Services contains all services of the application
type Services struct {
	UserPositionService domain.UserPositionService
	Config        config.Configuration
}

// NewServices creates a new Services instance
func NewServices(config config.Configuration, repos *postgres.Repositories) (*Services, error) {
	userPositionService, err := NewUserPositionService(repos.Places)
	if err != nil {
		return nil, err
	}
	return &Services{
		Config:        config,
		UserPositionService: userPositionService,
	}, nil
}
