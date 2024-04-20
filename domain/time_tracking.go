package domain

import (
	"context"
	"time"
)

const (
	// TimeTrackingStatusUnspecified represents an unspecified time tracking status.
	TimeTrackingStatusUnspecified TimeTrackingStatus = 0
	// TimeTrackingStatusCheckedIn represents a checked in time tracking status.
	TimeTrackingStatusCheckedIn TimeTrackingStatus = 1
	// TimeTrackingStatusOnBreak represents an on break time tracking status.
	TimeTrackingStatusOnBreak TimeTrackingStatus = 2
	// TimeTrackingStatusCheckedOut represents a checked out time tracking status.
	TimeTrackingStatusCheckedOut TimeTrackingStatus = 3
)

// TimeTrackingStatus represents the status of a time tracking session.
type TimeTrackingStatus int

func (s TimeTrackingStatus) String() string {
	switch s {
	case TimeTrackingStatusCheckedIn:
		return "checked_in"
	case TimeTrackingStatusOnBreak:
		return "on_break"
	case TimeTrackingStatusCheckedOut:
		return "checked_out"
	default:
		return "unspecified"
	}
}

// TimeTrackingSession represents a time tracking session which is the data
// collected from an user when he/she is working for a single day.
type TimeTrackingSession struct {
	SessionID         string              `json:"session_id"`
	UserID            string              `json:"user_id"`
	Status            TimeTrackingStatus  `json:"status"`
	CheckedInTime     time.Time           `json:"checked_in_time"`
	CheckedOutTime    *time.Time          `json:"checked_out_time"`
	Breaks            []TimeTrackingBreak `json:"breaks" gorm:"foreignkey:SessionID;references:SessionID"`
	TotalWorkTime     time.Duration       `json:"total_work_time,omitempty"`
	TotalBreakTime    time.Duration       `json:"total_break_time,omitempty"`
	LastKnownLocation GeoPoint            `json:"last_known_position,omitempty"`
}

// TimeTrackingBreak represents a break taken by an user during a time tracking session.
type TimeTrackingBreak struct {
	BreakID   string     `json:"break_id,omitempty"`
	SessionID string     `json:"session_id,omitempty"`
	StartTime time.Time  `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}

// TimeTrackingRepository represents a repository for managing time tracking sessions.
type TimeTrackingRepository interface {
	CheckIn(ctx context.Context, session *TimeTrackingSession) error
	CheckOut(ctx context.Context, session *TimeTrackingSession) error
	StartBreak(ctx context.Context, session *TimeTrackingSession, t time.Time) error
	EndBreak(ctx context.Context, session *TimeTrackingSession, t time.Time) error
	GetTimeTrackingSession(ctx context.Context, sessionID, userID string) (*TimeTrackingSession, error)
}

// TimeTrackingService represents a service for managing time tracking sessions.
type TimeTrackingService interface {
	CheckIn(ctx context.Context, session *TimeTrackingSession) error
	CheckOut(ctx context.Context, session *TimeTrackingSession) error
	StartBreak(ctx context.Context, session *TimeTrackingSession) error
	EndBreak(ctx context.Context, session *TimeTrackingSession) error
	GetTimeTrackingSession(ctx context.Context, sessionID, userID string) (*TimeTrackingSession, error)
}
