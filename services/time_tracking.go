package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/nawaltni/tracker/domain"
)

// TimeTrackingService is the service for TimeTracking of an user
type TimeTrackingService struct {
	repository domain.TimeTrackingRepository
}

// NewTimeTrackingService creates a new TimeTrackingService
func NewTimeTrackingService(repository domain.TimeTrackingRepository) *TimeTrackingService {
	return &TimeTrackingService{repository: repository}
}

// CheckIn handles the logic for checking in an user
func (s *TimeTrackingService) CheckIn(ctx context.Context, session *domain.TimeTrackingSession) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	session.SessionID = id.String()
	session.Status = domain.TimeTrackingStatusCheckedIn

	if session.CheckedInTime.IsZero() {
		session.CheckedInTime = time.Now()
	}

	err = s.repository.CheckIn(ctx, session)
	if err != nil {
		return err
	}

	return nil
}

// CheckOut handles the logic for checking out an user
func (s *TimeTrackingService) CheckOut(ctx context.Context, session *domain.TimeTrackingSession) error {
	checkOutTime := time.Now()
	if session.CheckedOutTime != nil {
		checkOutTime = *session.CheckedOutTime
	}

	// Retrieve the existing session
	eSession, err := s.repository.GetTimeTrackingSession(ctx, session.SessionID, session.UserID)
	if err != nil {
		return err
	}

	// Update the session with checkout information
	eSession.CheckedOutTime = &checkOutTime
	if session.LastKnownLocation.Latitude != 0 && session.LastKnownLocation.Longitude != 0 {
		eSession.LastKnownLocation = session.LastKnownLocation
	}

	eSession.Status = domain.TimeTrackingStatusCheckedOut

	// Save the updated session using the repository's CheckOut method
	err = s.repository.CheckOut(ctx, eSession)
	if err != nil {
		return err
	}

	// Update the session with the updated session
	*session = *eSession

	return nil
}

// StartBreak handles the logic for starting a break for an user
func (s *TimeTrackingService) StartBreak(ctx context.Context, session *domain.TimeTrackingSession) error {
	breakStartTime := time.Now()

	eSession, err := s.repository.GetTimeTrackingSession(ctx, session.SessionID, session.UserID)
	if err != nil {
		return err
	}

	eSession.Status = domain.TimeTrackingStatusOnBreak
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	// Create a new break
	newBreak := domain.TimeTrackingBreak{
		BreakID:   id.String(),
		SessionID: session.SessionID,
		StartTime: breakStartTime,
		// EndTime is not set yet as the break has just started
	}

	// Add the new break to the session
	eSession.Breaks = append(eSession.Breaks, newBreak)

	err = s.repository.StartBreak(ctx, eSession, breakStartTime)
	if err != nil {
		return err
	}

	*session = *eSession
	return nil
}

// EndBreak handles the logic for ending a break for an user
func (s *TimeTrackingService) EndBreak(ctx context.Context, session *domain.TimeTrackingSession) error {
	breakEndTime := time.Now()

	eSession, err := s.repository.GetTimeTrackingSession(ctx, session.SessionID, session.UserID)
	if err != nil {
		return err
	}

	if len(eSession.Breaks) == 0 {
		return errors.New("no breaks to end")
	}

	// Assuming the last break is the one to end
	lastBreakIndex := len(eSession.Breaks) - 1
	if eSession.Breaks[lastBreakIndex].EndTime != nil {
		return errors.New("break has already ended")
	}
	eSession.Breaks[lastBreakIndex].EndTime = &breakEndTime

	eSession.Status = domain.TimeTrackingStatusCheckedIn
	err = s.repository.EndBreak(ctx, eSession, breakEndTime)
	if err != nil {
		return err
	}

	*session = *eSession
	return nil
}

// GetTimeTrackingSession handles the logic for getting a time tracking session for an user
func (s *TimeTrackingService) GetTimeTrackingSession(ctx context.Context, sessionID, userID string) (*domain.TimeTrackingSession, error) {
	return s.repository.GetTimeTrackingSession(ctx, sessionID, userID)
}
