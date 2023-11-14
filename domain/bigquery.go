package domain

import "context"

// TrackerClientBigquery is the interface for the bigquery client.
type TrackerClientBigquery interface {
	RecordUserPosition(ctx context.Context, userPosition UserPosition) error
}
