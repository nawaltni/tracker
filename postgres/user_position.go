package postgres

import (
	"fmt"

	"github.com/nawaltni/tracker/domain"
)

// UserPositionRepository is the GORM implementation of the UserPositionRepository.
type UserPositionRepository struct {
	client *Client
}

// NewUserPositionRepository creates a new instance of a GORM-based user position repository.
func NewUserPositionRepository(client *Client) *UserPositionRepository {
	return &UserPositionRepository{client: client}
}

// Insert adds a new UserPosition to the database.
func (r *UserPositionRepository) Insert(userPosition *domain.UserPosition) error {
	model := ToModelUserPosition(userPosition)
	err := r.client.db.Create(&model).Error
	if err != nil {
		return fmt.Errorf("error inserting user position: %w", err)
	}
	*userPosition = *ToDomainUserPosition(model)
	return nil
}

// GetUserPosition retrieves a user's most recent position from the database.
func (r *UserPositionRepository) GetUserPosition(userID string) (*domain.UserPosition, error) {
	var userPosition domain.UserPosition
	err := r.client.db.Where("user_id = ?", userID).Order("created_at DESC").First(&userPosition).Error
	return &userPosition, err
}

// GetUsersPositionByCoordinates retrieves a list of users' positions close to the given coordinates.
func (r *UserPositionRepository) GetUsersPositionByCoordinates(lat float64, lon float64) ([]domain.UserPosition, error) {
	var userPositions []domain.UserPosition
	// This will require raw SQL to utilize PostGIS functions.
	// You need to adjust the SQL query to your needs (e.g., distance).
	err := r.client.db.Raw(`
	SELECT * FROM user_positions
	WHERE ST_DWithin(location, ST_SetSRID(ST_Point(?, ?), 4326)::geography, 5000)
	ORDER BY timestamp DESC
	`, lon, lat).Scan(&userPositions).Error

	return userPositions, err
}
