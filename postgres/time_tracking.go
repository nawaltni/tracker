package postgres

import (
	"context"
	"time"

	"github.com/nawaltni/tracker/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// TimeTrackingRepository is the GORM implementation of the TimeTrackingRepository.
type TimeTrackingRepository struct {
	client *Client
}

// NewTimeTrackingRepository creates a new instance of a GORM-based time tracking repository.
func NewTimeTrackingRepository(client *Client) *TimeTrackingRepository {
	return &TimeTrackingRepository{client: client}
}

// CheckIn creates a new time tracking session with the checked-in status.
func (r *TimeTrackingRepository) CheckIn(ctx context.Context, session *domain.TimeTrackingSession) error {
	mSession := ToModelTimeTrackingSession(session)

	// Using WithContext to ensure the context is respected in the operation
	if err := r.client.db.WithContext(ctx).Create(&mSession).Error; err != nil {
		return err
	}

	dSession := ToDomainTimeTrackingSession(mSession)

	*session = *dSession

	return nil
}

// CheckOut updates an existing time tracking session with the checked-out status.
func (r *TimeTrackingRepository) CheckOut(ctx context.Context, session *domain.TimeTrackingSession) error {
	mSession := ToModelTimeTrackingSession(session)

	// Save the updated session
	return r.client.db.WithContext(ctx).Save(mSession).Error
}

// Other methods (StartBreak, EndBreak, etc.) would be here...

// StartBreak adds a break to a time tracking session.
func (r *TimeTrackingRepository) StartBreak(ctx context.Context, session *domain.TimeTrackingSession, t time.Time) error {
	mSession := ToModelTimeTrackingSession(session)

	// Save the updated session
	if err := r.client.db.WithContext(ctx).Save(mSession).Error; err != nil {
		return err
	}

	*session = *ToDomainTimeTrackingSession(mSession)

	return nil
}

// EndBreak ends the most recent break in a time tracking session.
func (r *TimeTrackingRepository) EndBreak(ctx context.Context, session *domain.TimeTrackingSession, t time.Time) error {
	mSession := ToModelTimeTrackingSession(session)

	// Save the updated session
	err := r.client.db.Debug().WithContext(ctx).
		Session(&gorm.Session{FullSaveAssociations: true}). // Save the associations
		Clauses(clause.OnConflict{UpdateAll: true}).
		Save(mSession).Error
	if err != nil {
		return err
	}

	*session = *ToDomainTimeTrackingSession(mSession)

	return nil
}

// GetTimeTrackingSession retrieves a time tracking session based on sessionID and userID.
func (r *TimeTrackingRepository) GetTimeTrackingSession(ctx context.Context, sessionID, userID string) (*domain.TimeTrackingSession, error) {
	var session TimeTrackingSession

	// Find the session by sessionID and userID
	if err := r.client.db.WithContext(ctx).Preload("Breaks").
		Where("session_id = ? AND user_id = ?", sessionID, userID).
		First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, err // or a custom error for not found
		}
		return nil, err
	}

	// Convert the model to a domain entity
	dSession := ToDomainTimeTrackingSession(session)

	return dSession, nil
}
