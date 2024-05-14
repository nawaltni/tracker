package domain

import (
	"context"
	"time"
)

// BigqueryClient is the interface for the bigquery client.
type BigqueryClient interface {
	RecordUserPosition(ctx context.Context, userPosition UserPosition) error
	GetUserPositionsHistorySince(ctx context.Context, userID string, t time.Time, limit int) ([]UserPosition, error)
}
