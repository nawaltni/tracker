package services

import (
	"github.com/nawaltni/tracker/domain"
)

// UserPositionService is the service for UserPosition
type UserPositionService struct {
	repo domain.UserPositionRepository
}

// NewUserPositionService creates a new UserPositionService
func NewUserPositionService(repo domain.UserPositionRepository) (domain.UserPositionService, error) {
	return &UserPositionService{repo: repo}, nil
}
