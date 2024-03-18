package domain

import (
	"context"
	"time"
)

// BigqueryClient is the interface for the bigquery client.
type BigqueryClient interface {
	RecordUserPosition(ctx context.Context, userPosition UserPosition) error
	GetUserPositionsSince(ctx context.Context, userID string, t time.Time) ([]UserPosition, error)
}
