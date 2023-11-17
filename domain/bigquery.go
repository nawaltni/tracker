package domain

import "context"

// BigqueryClient is the interface for the bigquery client.
type BigqueryClient interface {
	RecordUserPosition(ctx context.Context, userPosition UserPosition) error
}
