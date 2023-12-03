package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/nawaltni/tracker/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
// There are
func (r *UserPositionRepository) Insert(ctx context.Context, userPosition *domain.UserPosition) error {
	model := ToModelUserPosition(userPosition)
	err := r.client.db.
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(&model).Error
	if err != nil {
		return fmt.Errorf("error inserting user position: %w", err)
	}
	*userPosition = *ToDomainUserPosition(model)
	return nil
}

// GetUserPosition retrieves a user's most recent position from the database.
func (r *UserPositionRepository) GetUserPosition(ctx context.Context, userID string) (*domain.UserPosition, error) {
	var userPosition UserPosition
	err := r.client.db.Preload("PhoneMeta").Where("user_id = ?", userID).Order("created_at DESC").First(&userPosition).Error
	if err != nil {
		if IsNotFoundError(err) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("error retrieving user position: %w", err)
	}
	pos := ToDomainUserPosition(userPosition)
	return pos, err
}

func IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// GetUsersPositionByCoordinates retrieves a list of users' positions close to the given coordinates.
func (r *UserPositionRepository) GetUsersPositionByCoordinates(ctx context.Context, lat float32, lon float32, distance int) ([]domain.UserPosition, error) {
	var userPositions []UserPosition
	// This will require raw SQL to utilize PostGIS functions.
	// You need to adjust the SQL query to your needs (e.g., distance).
	err := r.client.db.Raw(`
	SELECT * FROM user_positions
	WHERE ST_DWithin(location, ST_SetSRID(ST_Point(?, ?), 4326)::geography, ?)
	ORDER BY created_at DESC
	`, lat, lon, distance).Scan(&userPositions).Error

	if len(userPositions) == 0 {
		return nil, nil
	}

	userPositionsDomain := make([]domain.UserPosition, len(userPositions))

	for i, userPosition := range userPositions {
		userPositionsDomain[i] = *ToDomainUserPosition(userPosition)
	}

	return userPositionsDomain, err
}
